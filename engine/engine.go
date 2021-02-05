// This package provides functionality for the app engine.
package engine

import (
	"errors"
	"fmt"
	"github.com/cadmean-ru/amphion/common"
	"github.com/cadmean-ru/amphion/common/a"
	"github.com/cadmean-ru/amphion/frontend"
	"github.com/cadmean-ru/amphion/rendering"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

var instance *AmphionEngine

type AmphionEngine struct {
	platform common.Platform
	logger   *Logger
	renderer rendering.Renderer
	idgen    *common.IdGenerator
	started  bool
	state    byte

	loadedScene        *SceneObject
	currentScene       *SceneObject
	currentApp         *frontend.App
	appContext         *AppContext
	stopChan           chan bool
	eventChan          chan AmphionEvent
	updateRoutine      *updateRoutine
	eventBinder        *EventBinder
	globalContext      frontend.Context
	forceRedraw        bool
	messageDispatcher  *MessageDispatcher
	componentsManager  *ComponentsManager
	currentComponent   Component
	closeSceneCallback func()
	tasksRoutine       *TasksRoutine
	focusedObject      *SceneObject
	hoveredObject      *SceneObject
	front              frontend.Frontend
	suspend            bool
}

const (
	TargetFrameTime = 16
)

const (
	StateStopped = 0
	StateStarted = 1
	StateUpdating = 2
	StateRendering = 3
)

// Initializes a new instance of Amphion Engine and configures it to run with the specified frontend.
// Returns pointer to the created engine instance.
// The engine is a singleton, so calling Initialize more than once will have no effect.
func Initialize(front frontend.Frontend) *AmphionEngine {
	if instance != nil {
		return instance
	}

	instance = &AmphionEngine{
		platform:          front.GetPlatform(),
		logger:            GetLoggerForPlatform(front.GetPlatform()),
		idgen:             common.NewIdGenerator(),
		state:             StateStopped,
		stopChan:          make(chan bool),
		eventChan:         make(chan AmphionEvent, 100),
		updateRoutine:     newUpdateRoutine(),
		eventBinder:       newEventBinder(),
		tasksRoutine:      newTasksRoutine(),
		componentsManager: newComponentsManager(),
		front:             front,
	}
	instance.renderer = instance.front.GetRenderer()
	instance.globalContext = instance.front.GetContext()
	instance.front.SetCallback(instance.handleFrontEndCallback)
	return instance
}

// Returns pointer to the engine instance.
func GetInstance() *AmphionEngine {
	return instance
}

// Starts the engine.
// Must be called, before any interaction with the engine.
func (engine *AmphionEngine) Start() {
	engine.started = true
	engine.registerInternalEvenHandlers()
	engine.logger.Info(engine, "Amphion started")
	engine.state = StateStarted
	go engine.eventLoop()
	engine.tasksRoutine.start()
}

// Closes the current scene if any, and stops the engine.
func (engine *AmphionEngine) Stop() {
	engine.eventChan<-NewAmphionEvent(engine, EventStop, nil)
}

// Blocks the calling goroutine until the engine is stopped.
func (engine *AmphionEngine) WaitForStop() {
	<-engine.stopChan
}

// Returns the renderer.
func (engine *AmphionEngine) GetRenderer() rendering.Renderer {
	return engine.renderer
}

// Returns the logger.
func (engine *AmphionEngine) GetLogger() *Logger {
	return engine.logger
}

// Returns the currently displaying scene object.
func (engine *AmphionEngine) GetCurrentScene() *SceneObject {
	return engine.currentScene
}

// Returns the current loaded app or nil if no app is loaded.
func (engine *AmphionEngine) GetCurrentApp() *frontend.App {
	return engine.currentApp
}

// Returns the global application context.
// See frontend.Context.
func (engine *AmphionEngine) GetGlobalContext() frontend.Context {
	return engine.globalContext
}

// Loads scene from a resource file asynchronously.
// If show is true, after loading also shows this scene.
func (engine *AmphionEngine) LoadScene(scene int, show bool) {
	engine.RunTask(NewTaskBuilder().Run(func() (interface{}, error) {
		return engine.GetResourceManager().ReadFile(scene)
	}).Then(func(res interface{}) {
		data := res.([]byte)
		so := &SceneObject{}
		err := so.DecodeFromYaml(data)
		if err != nil {
			engine.logger.Info(engine, fmt.Sprintf("Failed to decode scene: %s", err.Error()))
			return
		}

		engine.loadedScene = so

		if show {
			engine.SwapScenes()
		}
	}).Err(func(err error) {
		engine.logger.Error(engine, fmt.Sprintf("Error loading scene with id %d: %e", scene, err))
	}).Build())
}

// Hides the current showing scene and shows the currently loaded scene (using LoadScene).
// The previously showing scene will be properly stopped and the deleted.
// So, calling SwapScenes again wont swap the two scenes back.
// If no scene is loaded, wont do anything.
func (engine *AmphionEngine) SwapScenes() {
	if engine.loadedScene == nil {
		return
	}

	engine.CloseScene(func() {
		_ = engine.ShowScene(engine.loadedScene)
		engine.loadedScene = nil
	})
}

// Shows the specified scene object.
// Returns an error, if the engine is not yet ready or if another scene is already showing.
func (engine *AmphionEngine) ShowScene(scene *SceneObject) error {
	if !engine.started {
		engine.logger.Error(engine, "Cannot show scene. Engine not started!")
		return errors.New("engine not started")
	}

	if engine.currentScene != nil {
		engine.logger.Error(engine, "Cannot show scene. A scene is already showing!")
		return errors.New("a scene is already loaded")
	}

	engine.logger.Info(engine, fmt.Sprintf("Starting scene %s", scene.name))

	engine.configureScene(scene)
	engine.messageDispatcher = newMessageDispatcherForScene(scene)
	engine.currentScene = scene

	engine.logger.Info(engine, "Starting loop")
	engine.updateRoutine.start()

	// Perform first update
	engine.updateRoutine.requestRendering()

	// Update frontend window title. As any UI action must be executed on frontend thread.
	engine.ExecuteOnFrontendThread(func() {
		engine.front.SetWindowTitle(scene.name)
	})

	engine.logger.Info(engine, "Scene showing")

	return nil
}

// Closes the currently showing scene asynchronously.
// It will call the provided callback function as soon as the scene was closed.
// If no scene is showing calls the callback function immediately.
func (engine *AmphionEngine) CloseScene(callback func()) {
	if engine.currentScene == nil {
		callback()
		return
	}

	engine.closeSceneCallback = callback
	engine.eventChan<-NewAmphionEvent(engine, EventCloseScene, nil)
}

func (engine *AmphionEngine) configureScene(scene *SceneObject) {
	screenInfo := engine.globalContext.ScreenInfo
	scene.Transform.Size.X = float32(screenInfo.GetWidth())
	scene.Transform.Size.Y = float32(screenInfo.GetHeight())
}

func (engine *AmphionEngine) eventLoop() {
	engine.logger.Info(engine, "Event loop starting")

	defer engine.recover()

	for event := range engine.eventChan {
		if event.Code == EventStop {
			if engine.canStop() {
				engine.logger.Info(nil, "Stopping")
				break
			} else {
				engine.handleStop()
			}
		}

		engine.eventBinder.InvokeHandlers(event)
	}

	engine.handleStop()
}

// Tells the engine to schedule an update as soon as possible.
func (engine *AmphionEngine) RequestUpdate() {
	engine.updateRoutine.requestUpdate()
}

// Tells the engine to schedule rendering in the next update cycle.
// It will also request an update, if it was not requested already.
func (engine *AmphionEngine) RequestRendering() {
	engine.updateRoutine.requestRendering()
}

func (engine *AmphionEngine) IsForcedToRedraw() bool {
	return engine.forceRedraw
}

func (engine *AmphionEngine) handleFrontEndCallback(callback frontend.Callback) {
	switch callback.Code {
	case frontend.CallbackMouseDown:
		coords := strings.Split(callback.Data, ";")
		if len(coords) != 2 {
			panic("Invalid click callback Data")
		}
		x, err := strconv.ParseInt(coords[0], 10, 32)
		if err != nil {
			panic("Invalid click callback Data")
		}
		y, err := strconv.ParseInt(coords[1], 10, 32)
		if err != nil {
			panic("Invalid click callback Data")
		}
		engine.handleClickEvent(a.IntVector2{X: int(x), Y: int(y)})
	case frontend.CallbackMouseUp:
		coords := strings.Split(callback.Data, ";")
		if len(coords) != 2 {
			panic("Invalid click callback Data")
		}
		x, err := strconv.ParseInt(coords[0], 10, 32)
		if err != nil {
			panic("Invalid click callback Data")
		}
		y, err := strconv.ParseInt(coords[1], 10, 32)
		if err != nil {
			panic("Invalid click callback Data")
		}
		event := NewAmphionEvent(engine, EventMouseUp, a.NewIntVector3(int(x), int(y), 0))
		engine.eventChan<-event
	case frontend.CallbackContextChange:
		engine.globalContext = engine.front.GetContext()
		engine.configureScene(engine.currentScene)
		engine.forceRedraw = true
		engine.RequestRendering()
	case frontend.CallbackKeyDown:
		tokens := strings.Split(callback.Data, "\n")
		if len(tokens) != 2 {
			panic("Invalid key down callback Data")
		}
		event := NewAmphionEvent(engine, EventKeyDown, KeyEvent{
			Key:  tokens[0],
			Code: tokens[1],
		})
		engine.eventChan<-event
	case frontend.CallbackAppHide:
		engine.eventChan<-NewAmphionEvent(engine, EventAppHide, nil)
		engine.suspend = true
	case frontend.CallbackAppShow:
		engine.suspend = false
		engine.eventChan<-NewAmphionEvent(engine, EventAppShow, nil)
		engine.RequestRendering()
	case frontend.CallbackMouseMove:
		event := NewAmphionEvent(engine, EventMouseMove, nil)
		engine.eventChan <- event
	}
}

// Binds an event handler for the specified event code.
// The handler will be invoked in the event loop goroutine, when the event with the specified code is raised.
func (engine *AmphionEngine) BindEventHandler(code int, handler EventHandler) {
	engine.eventBinder.Bind(code, handler)
}

// Unbinds the event handler for the specified event code.
// The handler will no longer be invoked, when the event with the specified code is raised.
func (engine *AmphionEngine) UnbindEventHandler(code int, handler EventHandler) {
	engine.eventBinder.Unbind(code, handler)
}

// Raises a new event.
func (engine *AmphionEngine) RaiseEvent(event AmphionEvent) {
	engine.eventChan<-event
}

// Synchronously loads prefab from file.
func (engine *AmphionEngine) LoadPrefab(resId int) (*SceneObject, error) {
	prefab := &SceneObject{}
	data, err := engine.GetResourceManager().ReadFile(resId)
	if err != nil {
		return nil, err
	}

	err = prefab.DecodeFromYaml(data)
	if err != nil {
		return nil, err
	}

	prefab.id = engine.idgen.NextId()

	return prefab, nil
}

func (engine *AmphionEngine) handleFrontEndInterrupt(msg string) {
	engine.front.CommencePanic("Kernel panic", msg)
	panic(msg)
}

func (engine *AmphionEngine) recover() {
	if err := recover(); err != nil {
		engine.logger.Error(engine, "Fatal error.")
		engine.logger.Error(engine, fmt.Sprintf("Current state: %s", engine.GetStateString()))
		if engine.currentComponent != nil {
			engine.logger.Error(engine, fmt.Sprintf("Error in component %s", engine.currentComponent.GetName()))
		}
		engine.front.CommencePanic("Kernel panic", fmt.Sprintf("%v", err))
		panic(err)
	}
}

// Returns whether rendering is requested for the next update cycle.
func (engine *AmphionEngine) IsRenderingRequested() bool {
	return engine.updateRoutine.renderingRequested
}

// Returns whether an update is requested for the next frame.
func (engine *AmphionEngine) IsUpdateRequested() bool {
	return engine.updateRoutine.updateRequested
}

func (engine *AmphionEngine) OnMessage(_ Message) bool {
	return true
}

func (engine *AmphionEngine) GetMessageDispatcher() *MessageDispatcher {
	return engine.messageDispatcher
}

// Returns the current engine state.
func (engine *AmphionEngine) GetState() byte {
	return engine.state
}

// Returns the current engine state as string.
func (engine *AmphionEngine) GetStateString() string {
	switch engine.state {
	case StateStarted:
		return "Started"
	case StateStopped:
		return "Stopped"
	case StateUpdating:
		return "Updating"
	case StateRendering:
		return "Rendering"
	default:
		return "Unknown"
	}
}

func (engine *AmphionEngine) handleStop() {
	if engine.currentScene != nil {
		engine.logger.Info(engine, "Unable to stop. Scene is showing. Closing scene and retrying.")
		engine.CloseScene(engine.closeSceneCallback)
		engine.Stop()
		return
	}

	engine.logger.Info(engine, "Stopping")

	engine.state = StateStopped

	close(engine.eventChan)
	engine.updateRoutine.close()

	engine.logger.Info(engine, "Amphion stopped")

	engine.stopChan<-true
	close(engine.stopChan)
}

func (engine *AmphionEngine) canStop() bool {
	return engine.currentScene == nil
}

func (engine *AmphionEngine) handleClickEvent(clickPos a.IntVector2) {
	if engine.currentScene == nil {
		return
	}

	candidates := make([]*SceneObject, 0, 1)

	engine.currentScene.ForEachObject(func(o *SceneObject) {
		if o.IsRendering() && o.HasBoundary() && o.IsPointInsideBoundaries2D(a.NewVector3(float32(clickPos.X), float32(clickPos.Y), 0)) {
			candidates = append(candidates, o)
		}
	})

	if len(candidates) > 0 {
		sort.Slice(candidates, func(i, j int) bool {
			return candidates[i].Transform.GetGlobalPosition().Z > candidates[j].Transform.GetGlobalPosition().Z
		})
		o := candidates[0]
		engine.messageDispatcher.DispatchDirectly(o, NewMessage(o, MessageBuiltinEvent, NewAmphionEvent(o, EventMouseDown, clickPos)))
		engine.focusedObject = o

		event := NewAmphionEvent(engine, EventMouseDown, MouseEventData{
			MousePosition: clickPos,
			SceneObject:   o,
		})
		engine.eventChan<-event
	} else {
		engine.focusedObject = nil
		event := NewAmphionEvent(engine, EventMouseDown, MouseEventData{
			MousePosition: clickPos,
			SceneObject:   nil,
		})
		engine.eventChan<-event
	}
}

func (engine *AmphionEngine) handleMouseMove(_ AmphionEvent) bool {
	if engine.currentScene == nil {
		return false
	}

	mousePos := engine.GetInputManager().GetMousePosition()

	candidates := make([]*SceneObject, 0, 1)
	engine.currentScene.ForEachObject(func(o *SceneObject) {
		if o.IsRendering() && o.HasBoundary() && o.IsPointInsideBoundaries2D(a.NewVector3(float32(mousePos.X), float32(mousePos.Y), 0)) {
			candidates = append(candidates, o)
		}
	})

	if len(candidates) > 0 {
		sort.Slice(candidates, func(i, j int) bool {
			return candidates[i].Transform.GetGlobalPosition().Z > candidates[j].Transform.GetGlobalPosition().Z
		})
		o := candidates[0]

		if o == engine.hoveredObject {
			return true
		}

		if engine.hoveredObject != nil {
			engine.messageDispatcher.DispatchDirectly(
				engine.hoveredObject,
				NewMessage(
					engine.hoveredObject,
					MessageBuiltinEvent,
					NewAmphionEvent(engine.hoveredObject, EventMouseOut, nil),
				),
			)
		}

		engine.messageDispatcher.DispatchDirectly(
			o,
			NewMessage(
				o,
				MessageBuiltinEvent,
				NewAmphionEvent(o, EventMouseIn, nil),
			),
		)

		engine.hoveredObject = o
	} else {
		if engine.hoveredObject != nil {
			engine.messageDispatcher.DispatchDirectly(
				engine.hoveredObject,
				NewMessage(
					engine.hoveredObject,
					MessageBuiltinEvent,
					NewAmphionEvent(engine.hoveredObject, EventMouseOut, nil),
				),
			)

			engine.hoveredObject = nil
		}
	}

	return true
}

func (engine *AmphionEngine) handleCloseSceneEvent(_ AmphionEvent) bool {
	engine.logger.Info(engine, "Closing scene")
	engine.updateRoutine.stop()
	engine.updateRoutine.waitForStop()
	engine.front.Reset()
	engine.currentScene = nil
	engine.state = StateStarted
	engine.logger.Info(engine, "Scene closed")
	if engine.closeSceneCallback != nil {
		engine.closeSceneCallback()
	}
	return false
}

func (engine *AmphionEngine) registerInternalEvenHandlers() {
	engine.BindEventHandler(EventCloseScene, engine.handleCloseSceneEvent)
	engine.BindEventHandler(EventMouseMove, engine.handleMouseMove)
}

func (engine *AmphionEngine) GetTasksRoutine() *TasksRoutine {
	return engine.tasksRoutine
}

// Runs the given task in the background goroutine.
func (engine *AmphionEngine) RunTask(task Task) {
	engine.tasksRoutine.RunTask(task)
}

// Returns the current resource manager.
func (engine *AmphionEngine) GetResourceManager() frontend.ResourceManager {
	return engine.front.GetResourceManager()
}

// Returns the current frontend.
func (engine *AmphionEngine) GetFrontend() frontend.Frontend {
	return engine.front
}

// Returns the current input manager.
func (engine *AmphionEngine) GetInputManager() frontend.InputManager {
	return engine.front.GetInputManager()
}

func (engine *AmphionEngine) GetComponentsManager() *ComponentsManager {
	return engine.componentsManager
}

func (engine *AmphionEngine) GetName() string {
	return "Amphion Engine"
}

// Loads app data from well-known source and shows the main scene.
func (engine *AmphionEngine) LoadApp() {
	engine.RunTask(NewTaskBuilder().Run(func() (interface{}, error) {
		app := engine.front.GetApp()
		return app, nil
	}).Then(func(res interface{}) {
		app := res.(*frontend.App)
		if app != nil {
			engine.currentApp = app
			engine.appContext = makeAppContext(app)
			engine.RaiseEvent(NewAmphionEvent(engine, EventAppLoaded, nil))
			_ = Navigate(app.MainScene, nil)
		} else {
			engine.logger.Warning(engine, "No app info found in well-known location!")
		}
	}).Build())
}

// Returns the current app context.
func (engine *AmphionEngine) GetAppContext() *AppContext {
	return engine.appContext
}

// Executes the specified action on frontend thread.
// Can be used to execute UI related functions from another goroutine.
func (engine *AmphionEngine) ExecuteOnFrontendThread(action func()) {
	engine.front.ReceiveMessage(frontend.NewFrontendMessageWithData(frontend.MessageExec, action))
}

func (engine *AmphionEngine) rebuildMessageTree() {
	if engine.currentScene == nil {
		return
	}
	engine.messageDispatcher = newMessageDispatcherForScene(engine.currentScene)
}

// Return the name of the given component suitable for serialization.
func NameOfComponent(component interface{}) string {
	t := reflect.TypeOf(component)

	if t.Kind() == reflect.Ptr {
		t = reflect.Indirect(reflect.ValueOf(component)).Type()
	}

	return t.PkgPath() + "." + t.Name()
}