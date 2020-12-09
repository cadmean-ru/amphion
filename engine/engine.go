package engine

import (
	"errors"
	"fmt"
	"github.com/cadmean-ru/amphion/common"
	"github.com/cadmean-ru/amphion/rendering"
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
	currentApp         *App
	stopChan           chan bool
	eventChan          chan AmphionEvent
	updateRoutine      *updateRoutine
	eventBinder        *EventBinder
	globalContext      *GlobalContext
	forceRedraw        bool
	messageDispatcher  *MessageDispatcher
	currentComponent   Component
	closeSceneCallback func()
	tasksRoutine       *TasksRoutine
	resourceManager    *ResourceManager
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

func Initialize(platform common.Platform) *AmphionEngine {
	if instance != nil {
		return instance
	}

	instance = &AmphionEngine{
		platform:        platform,
		logger:          GetLoggerForPlatform(platform),
		renderer:        rendering.NewRenderer(),
		idgen:           common.NewIdGenerator(),
		state:           StateStopped,
		stopChan:        make(chan bool),
		eventChan:       make(chan AmphionEvent, 100),
		updateRoutine:   newUpdateRoutine(),
		eventBinder:     newEventBinder(),
		globalContext:   getGlobalContext(),
		tasksRoutine:    newTasksRoutine(),
		resourceManager: newResourceManager(),
	}
	return instance
}

func GetInstance() *AmphionEngine {
	return instance
}

func (engine *AmphionEngine) Start() {
	engine.started = true
	engine.renderer.Prepare()
	prepareInterop()
	engine.registerInternalEvenHandlers()
	engine.logger.Info(engine, "Amphion started")
	engine.state = StateStarted
	go engine.eventLoop()
	engine.tasksRoutine.start()
	//if engine.tryLoadApp() {
	//	engine.LoadScene(engine.currentApp.Scenes[0])
	//}
}

func (engine *AmphionEngine) Stop() {
	engine.eventChan<-NewAmphionEvent(engine, EventStop, nil)
}

func (engine *AmphionEngine) WaitForStop() {
	<-engine.stopChan
}

func (engine *AmphionEngine) GetRenderer() rendering.Renderer {
	return engine.renderer
}

func (engine *AmphionEngine) GetLogger() *Logger {
	return engine.logger
}

func (engine *AmphionEngine) GetCurrentScene() *SceneObject {
	return engine.currentScene
}

func (engine *AmphionEngine) GetGlobalContext() *GlobalContext {
	return engine.globalContext
}

func (engine *AmphionEngine) LoadScene(scene string) {
	engine.RunTask(NewTaskBuilder().Run(func() (interface{}, error) {
		return loadScene(scene)
	}).Than(func(res interface{}) {
		engine.loadedScene = res.(*SceneObject)
	}).Err(func(err error) {
		engine.logger.Error(engine, fmt.Sprintf("Error loading scene: %e", err))
	}).Build())
}

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

	engine.logger.Info(engine, "Scene showing")

	return nil
}

func (engine *AmphionEngine) CloseScene(callback func()) {
	engine.closeSceneCallback = callback
	engine.eventChan<-NewAmphionEvent(engine, EventCloseScene, nil)
}

func (engine *AmphionEngine) configureScene(scene *SceneObject) {
	screenInfo := engine.globalContext.screenInfo
	scene.Transform.Size.X = float64(screenInfo.GetWidth())
	scene.Transform.Size.Y = float64(screenInfo.GetHeight())
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

func (engine *AmphionEngine) RequestUpdate() {
	engine.updateRoutine.requestUpdate()
}

func (engine *AmphionEngine) RequestRendering() {
	engine.updateRoutine.requestRendering()
}

func (engine *AmphionEngine) IsForcedToRedraw() bool {
	return engine.forceRedraw
}

func (engine *AmphionEngine) handleFrontEndCallback(callback frontEndCallback) {
	switch callback.code {
	case FrontCallbackClick:
		coords := strings.Split(callback.data, ";")
		if len(coords) != 2 {
			panic("Invalid click callback Data")
		}
		x, err := strconv.ParseFloat(coords[0], 64)
		if err != nil {
			panic("Invalid click callback Data")
		}
		y, err := strconv.ParseFloat(coords[1], 64)
		if err != nil {
			panic("Invalid click callback Data")
		}
		event := NewAmphionEvent(engine, EventClick, common.NewVector3(x, y, 0))
		engine.eventChan<-event
		break
	case FrontCallbackContextChange:
		engine.globalContext = getGlobalContext()
		engine.configureScene(engine.currentScene)
		engine.RequestRendering()
		engine.forceRedraw = true
	}
}

func (engine *AmphionEngine) BindEventHandler(code int, handler EventHandler) {
	engine.eventBinder.Bind(code, handler)
}

func (engine *AmphionEngine) handleFrontEndInterrupt(msg string) {
	commencePanic("Kernel panic", msg)
	panic(msg)
}

func (engine *AmphionEngine) recover() {
	if err := recover(); err != nil {
		engine.logger.Error(engine, "Fatal error.")
		engine.logger.Error(engine, fmt.Sprintf("Current state: %s", engine.GetStateString()))
		if engine.currentComponent != nil {
			engine.logger.Error(engine, fmt.Sprintf("Error in component %s", engine.currentComponent.GetName()))
		}
		commencePanic("Kernel panic", fmt.Sprintf("%v", err))
		panic(err)
	}
}

func (engine *AmphionEngine) IsRenderingRequested() bool {
	return engine.updateRoutine.renderingRequested
}

func (engine *AmphionEngine) IsUpdateRequested() bool {
	return engine.updateRoutine.updateRequested
}

func (engine *AmphionEngine) OnMessage(_ Message) bool {
	return true
}

func (engine *AmphionEngine) GetMessageDispatcher() *MessageDispatcher {
	return engine.messageDispatcher
}

func (engine *AmphionEngine) GetState() byte {
	return engine.state
}

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

func (engine *AmphionEngine) handleClickEvent(event AmphionEvent) bool {
	clickPos := event.Data.(common.Vector3)
	engine.currentScene.ForEachObject(func(o *SceneObject) {
		if o.IsRendering() && o.HasBoundary() && o.IsPointInsideBoundaries2D(clickPos) {
			engine.messageDispatcher.Dispatch(NewMessage(o, MessageBuiltinEvent, NewAmphionEvent(o, EventClick, clickPos)))
		}
	})
	return true
}

func (engine *AmphionEngine) handleCloseSceneEvent(_ AmphionEvent) bool {
	engine.logger.Info(engine, "Closing scene")
	engine.updateRoutine.stop()
	engine.updateRoutine.waitForStop()
	frontEndCloseScene()
	engine.currentScene = nil
	engine.state = StateStarted
	engine.renderer.Clear()
	engine.logger.Info(engine, "Scene closed")
	if engine.closeSceneCallback != nil {
		engine.closeSceneCallback()
	}
	return false
}

func (engine *AmphionEngine) registerInternalEvenHandlers() {
	engine.BindEventHandler(EventClick, engine.handleClickEvent)
	engine.BindEventHandler(EventCloseScene, engine.handleCloseSceneEvent)
}

func (engine *AmphionEngine) GetTasksRoutine() *TasksRoutine {
	return engine.tasksRoutine
}

func (engine *AmphionEngine) RunTask(task Task) {
	engine.tasksRoutine.RunTask(task)
}

func (engine *AmphionEngine) GetResourceManager() *ResourceManager {
	return engine.resourceManager
}

func (engine *AmphionEngine) GetName() string {
	return "Amphion Engine"
}

func (engine *AmphionEngine) tryLoadApp() bool {
	if data, err := loadAppData(); err == nil {
		if app, err := DecodeApp(data); err == nil {
			engine.currentApp = app
			return true
		} else {
			engine.logger.Warning(engine, "Failed to decode app")
		}
	} else {
		engine.logger.Warning(engine, "Failed to load app")
	}

	return false
}