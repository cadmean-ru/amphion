package engine

import (
	"fmt"
	"github.com/cadmean-ru/amphion/common/dispatch"
	"time"
)

type updateRoutine struct {
	running            bool
	updateChan         *dispatch.MessageQueue
	eventQueue         *dispatch.MessageQueue
	updateRequested    bool
	renderingRequested bool
	newSceneObjects    []*SceneObject
	startSceneObjects  []*SceneObject
	stopSceneObjects   []*SceneObject
	componentsToStop   []*ComponentContainer
	eventBinder        *EventBinder
	dt                 time.Duration
	lastFrameTime      time.Time
	updateTime         time.Duration
	renderingTime      time.Duration
}

//region Internal API

func (r *updateRoutine) start() {
	if r.running {
		return
	}

	go r.Loop()
}

func (r *updateRoutine) requestUpdate() {
	if r.updateRequested {
		return
	}

	r.updateRequested = true
	r.updateChan.Enqueue(dispatch.NewMessage(MessageUpdate))
}

func (r *updateRoutine) requestRendering() {
	r.renderingRequested = true
	r.requestUpdate()
}

func (r *updateRoutine) stop() {
	r.updateChan.Enqueue(dispatch.NewMessage(MessageUpdateStop))
}

func (r *updateRoutine) initSceneObject(object *SceneObject) {
	r.newSceneObjects = append(r.newSceneObjects, object)
	r.startSceneObject(object)
}

func (r *updateRoutine) startSceneObject(object *SceneObject) {
	r.startSceneObjects = append(r.startSceneObjects, object)
}

func (r *updateRoutine) stopSceneObject(object *SceneObject) {
	r.stopSceneObjects = append(r.stopSceneObjects, object)
}

func (r *updateRoutine) stopComponent(c *ComponentContainer) {
	r.componentsToStop = append(r.componentsToStop, c)
}

func (r *updateRoutine) waitForStop() {
	if !r.running {
		return
	}

	for r.running {
		instance.logger.Info(r, "Waiting for update Loop to stop")
		time.Sleep(100 * time.Millisecond)
	}
}

func (r *updateRoutine) enqueueEventAndRequestUpdate(event AmphionEvent) {
	r.eventQueue.Enqueue(dispatch.NewMessageWithAnyData(0, event))
	r.requestUpdate()
}

func (r *updateRoutine) close() {
	if r.running {
		instance.logger.Error(r, "Cannot close update Loop before stopping")
		panic("Cannot close update Loop before stopping")
	}
	r.updateChan.Close()
	r.eventQueue.Close()
}

//endregion

//region Loop

func (r *updateRoutine) Loop() {
	instance.logger.Info(r, "Starting")

	r.running = true

	defer instance.recover()

	r.handleStart()

	// Updating every frame or wait for update chan
	for {
		msg := r.updateChan.DequeueBlocking()

		if msg.What == MessageUpdateStop {
			instance.logger.Info(r, "Stopping")
			break
		}

		r.dt = time.Since(r.lastFrameTime)
		r.lastFrameTime = time.Now()

		r.handleSceneObjectsLifecycle()

		if instance.suspend {
			instance.state = StateStarted
			r.updateRequested = false
			r.renderingRequested = false
			continue
		}

		r.handleEvents()

		updateStart := time.Now()

		r.performUpdateIfNeeded()

		r.updateTime = time.Since(updateStart)
		renderingStart := time.Now()

		r.performRenderingIdNeeded()

		r.renderingTime = time.Since(renderingStart)

		instance.state = StateStarted

		r.waitForNextFrame()
	}

	r.handleStop()
}

func (r *updateRoutine) handleStart() {
	// Initialize all components
	instance.currentScene.setInCurrentScene(true)
	r.loopInit(instance.currentScene)

	// Calling OnStart for all objects in scene
	r.loopStart(instance.currentScene)

	r.lastFrameTime = time.Now()
}

func (r *updateRoutine) handleSceneObjectsLifecycle() {
	if len(r.newSceneObjects) > 0 {
		for _, o := range r.newSceneObjects {
			o.init(newInitContext(instance, o))
		}
		r.newSceneObjects = make([]*SceneObject, 0, 10)
	}

	if len(r.startSceneObjects) > 0 {
		for _, o := range r.startSceneObjects {
			o.start()
		}
		r.startSceneObjects = make([]*SceneObject, 0, 10)
	}

	if len(r.stopSceneObjects) > 0 {
		for _, o := range r.stopSceneObjects {
			o.stop()
		}
		r.stopSceneObjects = make([]*SceneObject, 0, 10)
	}

	if len(r.componentsToStop) > 0 {
		for _, c := range r.componentsToStop {
			c.stop()
		}
		r.componentsToStop = make([]*ComponentContainer, 0, 10)
	}
}

func (r *updateRoutine) handleEvents() {
	r.eventQueue.LockMainChannel()

	for !r.eventQueue.IsEmpty() {
		msg := r.eventQueue.Dequeue()
		event := msg.AnyData.(AmphionEvent)

		if event.Code == EventStop {
			if instance.canStop() {
				instance.logger.Info(nil, "Stopping")
				break
			} else {
				instance.handleStop()
			}
		}

		r.eventBinder.InvokeHandlers(event)
	}

	r.eventQueue.UnlockMainChannel()
}

func (r *updateRoutine) performUpdateIfNeeded() {
	if !r.updateRequested {
		return
	}

	//engine.logger.Info("Update Loop", "Updating components")

	r.updateRequested = false
	instance.state = StateUpdating

	ctx := newUpdateContext(float32(r.dt.Seconds()))

	// Calling OnUpdate for all objects in scene
	r.loopUpdate(instance.currentScene, ctx)
}

func (r *updateRoutine) performRenderingIdNeeded() {
	if !r.renderingRequested {
		return
	}

	//instance.logger.Warning(r, "Rendering components")

	r.renderingRequested = false
	instance.state = StateRendering

	ctx := newRenderingContext(instance.renderer)

	// Render objects
	r.loopRender(instance.currentScene, ctx)
	instance.renderer.PerformRendering()

	instance.forceRedraw = false
}

func (r *updateRoutine) waitForNextFrame() {
	// Wait until next frame
	timeToSleep := TargetFrameTime - time.Since(r.lastFrameTime).Milliseconds()

	if timeToSleep == 0 {
		instance.logger.Warning(r,
			fmt.Sprintf("The application is skipping frames! Update time: %d, Rendering time: %d",
				r.updateTime.Milliseconds(),
				r.renderingTime.Milliseconds()))
	}

	time.Sleep(time.Duration(timeToSleep) * time.Millisecond)
}

func (r *updateRoutine) handleStop() {
	r.loopStop(instance.currentScene)
	instance.currentScene.setInCurrentScene(false)

	instance.renderer.Clear()

	r.running = false
	r.newSceneObjects = make([]*SceneObject, 0)
	r.startSceneObjects = make([]*SceneObject, 0)
	r.stopSceneObjects = make([]*SceneObject, 0)

	instance.logger.Info(r, "Stopped")
}

//endregion

func (r *updateRoutine) loopInit(obj *SceneObject) {
	obj.init(newInitContext(instance, obj))
	temp := make([]*SceneObject, len(obj.children))
	copy(temp, obj.children)
	for _, c := range temp {
		r.loopInit(c)
	}
}

func (r *updateRoutine) loopStart(obj *SceneObject) {
	if !obj.enabled {
		return
	}
	obj.start()
	temp := make([]*SceneObject, len(obj.children))
	copy(temp, obj.children)
	for _, c := range temp {
		r.loopStart(c)
	}
}

func (r *updateRoutine) loopUpdate(obj *SceneObject, ctx UpdateContext) {
	if !obj.enabled {
		return
	}

	obj.update(ctx)

	for _, c := range obj.children {
		r.loopUpdate(c, ctx)
	}
}

func (r *updateRoutine) loopRender(obj *SceneObject, ctx DrawingContext) {
	if !obj.enabled {
		return
	}

	obj.draw(ctx)

	for _, c := range obj.children {
		r.loopRender(c, ctx)
	}
}

func (r *updateRoutine) loopStop(obj *SceneObject) {
	obj.stop()
	for _, c := range obj.children {
		r.loopStop(c)
	}
}

func (r *updateRoutine) GetMessageDispatcher() dispatch.MessageDispatcher {
	return r
}

func (r *updateRoutine) SendMessage(message *dispatch.Message) {
	r.updateChan.Enqueue(message)
}

func (r *updateRoutine) GetName() string {
	return "Update routine"
}

func newUpdateRoutine() *updateRoutine {
	return &updateRoutine{
		running:           false,
		updateChan:        dispatch.NewMessageQueue(10),
		eventQueue:        dispatch.NewMessageQueue(100),
		newSceneObjects:   make([]*SceneObject, 0, 10),
		startSceneObjects: make([]*SceneObject, 0, 10),
		stopSceneObjects:  make([]*SceneObject, 0, 10),
		componentsToStop:  make([]*ComponentContainer, 0, 10),
		eventBinder:       newEventBinder(),
	}
}
