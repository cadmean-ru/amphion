package engine

import (
	"fmt"
	"github.com/cadmean-ru/amphion/common/dispatch"
	"time"
)

type updateRoutine struct {
	sceneLifecycleManager
	running            bool
	updateChan         *dispatch.MessageQueue
	eventQueue         *dispatch.MessageQueue
	updateRequested    bool
	renderingRequested bool
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

	r.updateRequested = true
	r.renderingRequested = true
	r.updateChan = dispatch.NewMessageQueue(10)
	r.eventQueue = dispatch.NewMessageQueue(100)

	r.sceneLifecycleManager.start()

	r.updateChan.Enqueue(dispatch.NewMessage(MessageUpdate))

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
		//instance.logger.Info(r, "Waiting for update")
		msg := r.updateChan.DequeueBlocking()

		//instance.logger.Info(r, "Performing update")

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

	// Calling OnStart for all objects in SceneObject
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

func (r *updateRoutine) performUpdateIfNeeded() {
	if !r.updateWasRequested {
		return
	}

	r.updateWasRequested = false
	instance.state = StateUpdating

	ctx := newUpdateContext(float32(r.dt.Seconds()))

	r.loopUpdate(instance.currentScene, ctx)
	r.loopLayout(instance.currentScene)
	r.loopLateUpdate(instance.currentScene, ctx)
}

func (r *updateRoutine) performRenderingIdNeeded() {
	if !r.renderingWasRequested {
		return
	}

	//instance.logger.Info(r, "Performing rendering")

	r.renderingWasRequested = false
	instance.state = StateRendering

	r.loopRender(instance.currentScene)

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
