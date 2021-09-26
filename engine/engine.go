// Package engine provides functionality for the app engine.
package engine

import (
	"errors"
	"fmt"
	"github.com/cadmean-ru/amphion/common"
	"github.com/cadmean-ru/amphion/common/a"
	"github.com/cadmean-ru/amphion/common/dispatch"
	"github.com/cadmean-ru/amphion/frontend"
	"github.com/cadmean-ru/amphion/rendering"
	"sort"
	"sync"
)

var instance *AmphionEngine

type AmphionEngine struct {
	platform common.Platform
	logger   *Logger
	renderer *rendering.ARenderer
	idgen    *common.IdGenerator
	started  bool
	state    byte

	loadedScene        *SceneObject
	currentScene       *SceneObject
	sceneContext       *SceneContext
	currentApp         *frontend.App
	appContext         *AppContext
	stopChan           chan bool

	updateRoutine      *updateRoutine

	globalContext      frontend.Context
	forceRedraw        bool
	messageDispatcher  *MessageDispatcher
	componentsManager  *ComponentsManager
	currentComponent   Component
	closeSceneCallback func()
	tasksRoutine       *TasksRoutine
	front              frontend.Frontend
	suspend            bool
	inputManager       *InputManager
	startingWg         *sync.WaitGroup
	callbackHandler    dispatch.MessageDispatcher

	*FeaturesManager
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

// Initialize initializes a new instance of Amphion Engine and configures it to run with the specified frontend.
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
		stopChan:          make(chan bool, 1),
		updateRoutine:     newUpdateRoutine(),
		tasksRoutine:      newTasksRoutine(),
		componentsManager: newComponentsManager(),
		front:             front,
		inputManager:      newInputManager(),
		startingWg:        &sync.WaitGroup{},
		FeaturesManager:   newFeaturesManager(),
		callbackHandler:   newFrontendCallbackHandler(),
	}
	instance.startingWg.Add(2)
	instance.renderer = front.GetRenderer()
	instance.renderer.SetPreparedCallback(func() {
		instance.logger.Info(instance, "Rendering prepared")
		instance.startingWg.Done()
	})
	instance.globalContext = instance.front.GetContext()
	instance.front.SetEngineDispatcher(instance.callbackHandler)
	return instance
}

// GetInstance returns pointer to the engine instance.
func GetInstance() *AmphionEngine {
	return instance
}

// Start starts the engine.
// Must be called, before any interaction with the engine.
func (engine *AmphionEngine) Start() {
	engine.startingWg.Wait()
	engine.started = true
	engine.registerInternalEventHandlers()
	engine.state = StateStarted
	engine.tasksRoutine.start()
	engine.logger.Info(engine, "Amphion started")
}

// Stop closes the current SceneObject if any, and stops the engine.
func (engine *AmphionEngine) Stop() {
	engine.updateRoutine.enqueueEventAndRequestUpdate(NewAmphionEvent(engine, EventStop, nil))
}

// WaitForStop blocks the calling goroutine until the engine is stopped.
func (engine *AmphionEngine) WaitForStop() {
	<-engine.stopChan
}

// GetRenderer returns the renderer.
func (engine *AmphionEngine) GetRenderer() *rendering.ARenderer {
	return engine.renderer
}

// GetLogger returns the logger.
func (engine *AmphionEngine) GetLogger() *Logger {
	return engine.logger
}

// GetCurrentScene returns the currently displaying SceneObject object.
func (engine *AmphionEngine) GetCurrentScene() *SceneObject {
	return engine.currentScene
}

// GetCurrentApp returns the current loaded app or nil if no app is loaded.
func (engine *AmphionEngine) GetCurrentApp() *frontend.App {
	return engine.currentApp
}

// GetGlobalContext returns the global application context.
// See frontend.Context.
func (engine *AmphionEngine) GetGlobalContext() frontend.Context {
	return engine.globalContext
}

// LoadScene loads SceneObject from a resource file asynchronously.
// If show is true, after loading also shows this SceneObject.
func (engine *AmphionEngine) LoadScene(scene a.ResId, show bool) {
	if engine.state != StateStarted {
		panic("Invalid engine state")
	}

	engine.RunTask(NewTaskBuilder().Run(func() (interface{}, error) {
		return engine.GetResourceManager().ReadFile(scene)
	}).Then(func(res interface{}) {
		data := res.([]byte)
		so := &SceneObject{}
		err := so.DecodeFromYaml(data)
		if err != nil {
			engine.logger.Info(engine, fmt.Sprintf("Failed to decode SceneObject: %s", err.Error()))
			return
		}

		engine.loadedScene = so

		if show {
			engine.SwapScenes()
		}
	}).Err(func(err error) {
		engine.logger.Error(engine, fmt.Sprintf("Error loading SceneObject with id %d: %v", scene, err))
	}).Build())
}

// SwapScenes hides the current showing SceneObject and shows the currently loaded SceneObject (using LoadScene).
// The previously showing SceneObject will be properly stopped and the deleted.
// So, calling SwapScenes again will not swap the two scenes back.
// If no SceneObject is loaded, will not do anything.
func (engine *AmphionEngine) SwapScenes() {
	if engine.loadedScene == nil {
		return
	}

	engine.CloseScene(func() {
		_ = engine.ShowScene(engine.loadedScene)
		engine.loadedScene = nil
	})
}

// ShowScene shows the specified SceneObject object.
// Returns an error, if the engine is not yet ready or if another SceneObject is already showing.
func (engine *AmphionEngine) ShowScene(scene *SceneObject) error {
	if !engine.started {
		engine.logger.Error(engine, "Cannot show SceneObject. Engine not started!")
		return errors.New("engine not started")
	}

	if engine.currentScene != nil {
		engine.logger.Error(engine, "Cannot show SceneObject. A SceneObject is already showing!")
		return errors.New("a SceneObject is already loaded")
	}

	engine.logger.Info(engine, fmt.Sprintf("Starting SceneObject %s", scene.name))

	engine.sceneContext = makeSceneContext()

	newScene := engine.prepareScene(scene)
	engine.configureScene(newScene)
	engine.messageDispatcher = newMessageDispatcherForScene(newScene)
	engine.currentScene = newScene
	engine.renderer.SetRoot(engine.currentScene.renderingNode)

	engine.logger.Info(engine, "Starting Loop")
	engine.updateRoutine.start()

	// Perform first update
	engine.updateRoutine.requestRendering()

	// Update frontend window title. As any UI action must be executed on frontend thread.
	engine.SetWindowTitle(engine.currentScene.name)

	engine.logger.Info(engine, "Scene showing")

	return nil
}

// CloseScene closes the currently showing SceneObject asynchronously.
// It will call the provided callback function as soon as the SceneObject was closed.
// If no SceneObject is showing calls the callback function immediately.
func (engine *AmphionEngine) CloseScene(callback func()) {
	if engine.currentScene == nil {
		callback()
		return
	}

	engine.closeSceneCallback = callback
	engine.updateRoutine.stop()
}

func (engine *AmphionEngine) prepareScene(scene *SceneObject) *SceneObject {
	var scene2 *SceneObject
	if scene.initialized {
		scene2 = scene.Copy(scene.name)
	} else {
		scene2 = scene
	}

	return scene2
}

func (engine *AmphionEngine) configureScene(scene *SceneObject) {
	if scene == nil {
		return
	}

	screenInfo := engine.globalContext.ScreenInfo
	size := a.NewVector3(float32(screenInfo.GetWidth()), float32(screenInfo.GetHeight()), 0)
	scene.Transform.size = size
	scene.Transform.actualSize = size
}

// RequestUpdate tells the engine to schedule an update as soon as possible.
func (engine *AmphionEngine) RequestUpdate() {
	engine.updateRoutine.requestUpdate()
}

// RequestRendering tells the engine to schedule rendering in the next update cycle.
// It will also request an update, if it was not requested already.
func (engine *AmphionEngine) RequestRendering() {
	engine.updateRoutine.requestRendering()
}

//ForceAllViewsRedraw will request all view in the SceneObject to redraw on the next rendering cycle.
//It will not request rendering, you will need to call RequestRendering after that.
func (engine *AmphionEngine) ForceAllViewsRedraw() {
	engine.forceRedraw = true
}

//IsForcedToRedraw checks if all views redraw was requested in the next rendering cycle by calling ForceAllViewsRedraw.
func (engine *AmphionEngine) IsForcedToRedraw() bool {
	return engine.forceRedraw
}

// BindEventHandler binds an event handler for the specified event code.
// The handler will be invoked in the event Loop goroutine, when the event with the specified code is raised.
func (engine *AmphionEngine) BindEventHandler(code int, handler EventHandler) {
	engine.updateRoutine.eventBinder.Bind(code, handler)
}

// UnbindEventHandler unbinds the event handler for the specified event code.
// The handler will no longer be invoked, when the event with the specified code is raised.
func (engine *AmphionEngine) UnbindEventHandler(code int, handler EventHandler) {
	engine.updateRoutine.eventBinder.Unbind(code, handler)
}

// RaiseEvent raises a new event.
func (engine *AmphionEngine) RaiseEvent(event AmphionEvent) {
	engine.updateRoutine.enqueueEventAndRequestUpdate(event)
}

// LoadPrefab synchronously loads prefab from file.
func (engine *AmphionEngine) LoadPrefab(resId a.ResId) (*SceneObject, error) {
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
		engine.logger.Error(engine, "")
		engine.logger.Error(engine, "--- BEGINNING OF CRASH ---")
		engine.logger.Error(engine, "")
		engine.logger.Error(engine, "Fatal error.")
		engine.logger.Error(engine, fmt.Sprintf("Current state: %s", engine.GetStateString()))
		reason := "Kernel panic"
		if engine.currentComponent != nil {
			reason = fmt.Sprintf("Error in component %s", NameOfComponent(engine.currentComponent))
			engine.logger.Error(engine, reason)
		}
		engine.front.CommencePanic(reason, fmt.Sprintf("%v", err))
		panic(err)
	}
}

// IsRenderingRequested returns whether rendering is requested for the next update cycle.
func (engine *AmphionEngine) IsRenderingRequested() bool {
	return engine.updateRoutine.renderingRequested
}

// IsUpdateRequested returns whether an update is requested for the next frame.
func (engine *AmphionEngine) IsUpdateRequested() bool {
	return engine.updateRoutine.updateRequested
}

func (engine *AmphionEngine) OnMessage(_ Message) bool {
	return true
}

func (engine *AmphionEngine) GetMessageDispatcher() *MessageDispatcher {
	return engine.messageDispatcher
}

// GetState returns the current engine state.
func (engine *AmphionEngine) GetState() byte {
	return engine.state
}

// GetStateString returns the current engine state as string.
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
		engine.logger.Info(engine, "Unable to stop. Scene is showing. Closing SceneObject and retrying.")
		engine.CloseScene(func() {
			engine.handleStop()
		})
		return
	}

	if engine.appContext != nil {
		engine.appContext.onAppStopping()
	}

	engine.logger.Info(engine, "Stopping")

	engine.state = StateStopped

	engine.updateRoutine.close()
	engine.renderer.Stop()
	engine.tasksRoutine.stop()

	engine.logger.Info(engine, "Amphion stopped")

	engine.stopChan<-true
	close(engine.stopChan)

	engine.front.GetMessageDispatcher().SendMessage(dispatch.NewMessage(frontend.MessageEngineStopped))
}

func (engine *AmphionEngine) canStop() bool {
	return engine.currentScene == nil
}

func (engine *AmphionEngine) handleMouseMove(mousePos a.IntVector2) {
	if engine.currentScene == nil {
		return
	}

	candidates := make([]*SceneObject, 0, 1)
	engine.currentScene.ForEachObject(func(o *SceneObject) {
		if o.HasView() && o.HasBoundary() && o.IsPointInsideBoundaries2D(a.NewVector3(float32(mousePos.X), float32(mousePos.Y), 0)) {
			candidates = append(candidates, o)
		}
	})

	if len(candidates) > 0 {
		sort.Slice(candidates, func(i, j int) bool {
			return candidates[i].Transform.GlobalPosition().Z > candidates[j].Transform.GlobalPosition().Z
		})
		o := candidates[0]

		if o == engine.sceneContext.hoveredObject {
			return
		}

		if engine.sceneContext.hoveredObject != nil {
			engine.messageDispatcher.DispatchDirectly(
				engine.sceneContext.hoveredObject,
				NewMessage(
					engine.sceneContext.hoveredObject,
					MessageBuiltinEvent,
					NewAmphionEvent(engine.sceneContext.hoveredObject, EventMouseOut, nil),
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

		engine.sceneContext.hoveredObject = o
	} else {
		if engine.sceneContext.hoveredObject != nil {
			engine.messageDispatcher.DispatchDirectly(
				engine.sceneContext.hoveredObject,
				NewMessage(
					engine.sceneContext.hoveredObject,
					MessageBuiltinEvent,
					NewAmphionEvent(engine.sceneContext.hoveredObject, EventMouseOut, nil),
				),
			)

			engine.sceneContext.hoveredObject = nil
		}
	}
}

func (engine *AmphionEngine) handleSceneClose() {
	instance.renderer.Clear()
	engine.currentScene = nil
	engine.sceneContext = nil
	engine.state = StateStarted
	engine.logger.Info(engine, "Scene closed")
	if engine.closeSceneCallback != nil {
		engine.closeSceneCallback()
	}
}

func (engine *AmphionEngine) registerInternalEventHandlers() {
	//TODO: remove maybe?
}

func (engine *AmphionEngine) GetTasksRoutine() *TasksRoutine {
	return engine.tasksRoutine
}

// RunTask runs the given task in the background goroutine.
func (engine *AmphionEngine) RunTask(task *Task) {
	engine.tasksRoutine.RunTask(task)
}

// GetResourceManager returns the current resource manager.
func (engine *AmphionEngine) GetResourceManager() frontend.ResourceManager {
	return engine.front.GetResourceManager()
}

// GetFrontend returns the current frontend.
func (engine *AmphionEngine) GetFrontend() frontend.Frontend {
	return engine.front
}

// GetInputManager returns the current input manager.
func (engine *AmphionEngine) GetInputManager() *InputManager {
	return engine.inputManager
}

//GetComponentsManager returns the current ComponentsManager.
func (engine *AmphionEngine) GetComponentsManager() *ComponentsManager {
	return engine.componentsManager
}

func (engine *AmphionEngine) GetName() string {
	return "Amphion Engine"
}

// LoadApp loads app data from well-known source and shows the main SceneObject.
func (engine *AmphionEngine) LoadApp(delegate ...AppDelegate) {
	if engine.state != StateStarted {
		panic("Invalid engine state")
	}
	if engine.appContext != nil {
		panic("app already loaded")
	}

	engine.RunTask(NewTaskBuilder().Run(func() (interface{}, error) {
		app := engine.front.GetApp()
		return app, nil
	}).Then(func(res interface{}) {
		app := res.(*frontend.App)
		if app != nil {
			engine.currentApp = app
			engine.appContext = makeAppContext(app)
			if len(delegate) > 0 {
				engine.appContext.delegate = delegate[0]
			}

			engine.appContext.onAppLoaded()

			args := engine.front.GetLaunchArgs()
			path := args.GetString("path")

			if path == "" {
				path = "/"
			}

			err := Navigate(path, nil)
			if err != nil {
				engine.logger.Warning(engine, fmt.Sprintf("Failed to navigate to main SceneObject: %s", err.Error()))
			}
		} else {
			engine.logger.Warning(engine, "No app info found in well-known location!")
		}
	}).Build())
}

// GetAppContext returns the current app's context.
func (engine *AmphionEngine) GetAppContext() *AppContext {
	return engine.appContext
}

// GetSceneContext returns the current SceneObject's context.
func (engine *AmphionEngine) GetSceneContext() *SceneContext {
	return engine.sceneContext
}

// ExecuteOnFrontendThread executes the specified action on frontend thread.
// Can be used to execute UI related functions from another goroutine.
func (engine *AmphionEngine) ExecuteOnFrontendThread(action func()) {
	engine.front.GetWorkDispatcher().Execute(dispatch.NewWorkItemFunc(action))
}

// SetWindowTitle updates app's window title.
// On web sets the tab's title.
func (engine *AmphionEngine) SetWindowTitle(title string) {
	engine.front.GetMessageDispatcher().SendMessage(dispatch.NewMessageWithStringData(frontend.MessageTitle, title))
}

//GetFeaturesManager returns the current FeaturesManager.
func (engine *AmphionEngine) GetFeaturesManager() *FeaturesManager {
	return engine.FeaturesManager
}

func (engine *AmphionEngine) rebuildMessageTree() {
	if engine.currentScene == nil {
		return
	}
	engine.messageDispatcher = newMessageDispatcherForScene(engine.currentScene)
}