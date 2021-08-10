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
	newSceneObjects    *dispatch.MessageQueue
	startSceneObjects  *dispatch.MessageQueue
	stopSceneObjects   *dispatch.MessageQueue
	componentsToStop   *dispatch.MessageQueue
	eventBinder        *EventBinder
	dt                 time.Duration
	lastFrameTime      time.Time
	updateTime         time.Duration
	renderingTime      time.Duration
	updateWasRequested bool
	renderingWasRequested bool
}

//region Internal API

func (r *updateRoutine) start() {
	if r.running {
		return
	}

	r.updateChan = dispatch.NewMessageQueue(10)
	r.eventQueue = dispatch.NewMessageQueue(100)
	r.newSceneObjects = dispatch.NewMessageQueue(MaxSceneObjects/2)
	r.startSceneObjects = dispatch.NewMessageQueue(MaxSceneObjects/2)
	r.stopSceneObjects = dispatch.NewMessageQueue(MaxSceneObjects/2)
	r.componentsToStop =  dispatch.NewMessageQueue(MaxSceneObjects/2)

	go r.Loop()
}

func (r *updateRoutine) requestUpdate() {
	if !r.running || r.updateRequested {
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
	LogDebug("Request to init %s", object.GetName())
	r.newSceneObjects.Enqueue(dispatch.NewMessageWithAnyData(0, object))
	r.startSceneObject(object)
}

func (r *updateRoutine) startSceneObject(object *SceneObject) {
	LogDebug("Request to start %s", object.GetName())
	r.startSceneObjects.Enqueue(dispatch.NewMessageWithAnyData(0, object))
}

func (r *updateRoutine) stopSceneObject(object *SceneObject) {
	r.stopSceneObjects.Enqueue(dispatch.NewMessageWithAnyData(0, object))
}

func (r *updateRoutine) stopComponent(c *ComponentContainer) {
	r.componentsToStop.Enqueue(dispatch.NewMessageWithAnyData(0, c))
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

		if instance.suspend {
			instance.state = StateStarted
			r.resetFlags()
			continue
		}

		r.handleFlags()

		r.handleEvents()

		r.handleSceneObjectsLifecycle()

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

func (r *updateRoutine) handleFlags() {
	r.updateWasRequested = r.updateRequested
	r.renderingWasRequested = r.renderingRequested
	r.updateRequested = false
	r.renderingRequested = false
}

func (r *updateRoutine) resetFlags() {
	r.updateWasRequested = false
	r.renderingWasRequested = false
	r.updateRequested = false
	r.renderingRequested = false
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

func (r *updateRoutine) handleSceneObjectsLifecycle() {
	r.newSceneObjects.LockMainChannel()
	r.startSceneObjects.LockMainChannel()
	r.stopSceneObjects.LockMainChannel()
	r.componentsToStop.LockMainChannel()

	for !r.newSceneObjects.IsEmpty() {
		o := r.newSceneObjects.Dequeue().AnyData.(*SceneObject)
		o.init(newInitContext(o))
	}

	for !r.startSceneObjects.IsEmpty() {
		o := r.startSceneObjects.Dequeue().AnyData.(*SceneObject)
		o.start()
	}

	for !r.stopSceneObjects.IsEmpty() {
		o := r.stopSceneObjects.Dequeue().AnyData.(*SceneObject)
		o.stop()
	}

	for !r.componentsToStop.IsEmpty() {
		c := r.componentsToStop.Dequeue().AnyData.(*ComponentContainer)
		c.stop()
	}

	r.newSceneObjects.UnlockMainChannel()
	r.startSceneObjects.UnlockMainChannel()
	r.stopSceneObjects.UnlockMainChannel()
	r.componentsToStop.UnlockMainChannel()
}

func (r *updateRoutine) performUpdateIfNeeded() {
	if !r.updateWasRequested {
		return
	}

	r.updateWasRequested = false
	instance.state = StateUpdating

	ctx := newUpdateContext(float32(r.dt.Seconds()))

	// Calling OnUpdate for all objects in scene
	r.loopUpdate(instance.currentScene, ctx)
}

func (r *updateRoutine) performRenderingIdNeeded() {
	if !r.renderingWasRequested {
		return
	}

	r.renderingWasRequested = false
	instance.state = StateRendering

	instance.currentScene.Traverse(func(object *SceneObject) bool {
		ctx := newDrawingContext(object)
		object.draw(ctx)

		return true
	})

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

	r.running = false

	r.newSceneObjects.Close()
	r.startSceneObjects.Close()
	r.stopSceneObjects.Close()
	r.componentsToStop.Close()
	r.eventQueue.Close()
	r.updateChan.Close()

	instance.logger.Info(r, "Stopped")

	instance.handleSceneClose()
}

//endregion

func (r *updateRoutine) loopInit(obj *SceneObject) {
	obj.init(newInitContext(obj))
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
		eventBinder:       newEventBinder(),
	}
}
