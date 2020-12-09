package engine

import "time"

type updateRoutine struct {
	running            bool
	updateChan         chan bool
	updateRequested    bool
	renderingRequested bool
}

func (r *updateRoutine) start() {
	if r.running {
		return
	}

	go r.loop()
}

func (r *updateRoutine) requestUpdate() {
	if r.updateRequested {
		return
	}

	r.updateRequested = true
	r.updateChan<-true
}

func (r *updateRoutine) requestRendering() {
	if r.renderingRequested {
		return
	}

	r.renderingRequested = true
	r.requestUpdate()
}

func (r *updateRoutine) stop() {
	r.updateChan<-false
}

func (r *updateRoutine) waitForStop() {
	if !r.running {
		return
	}

	for r.running {
		instance.logger.Info(r, "Waiting for update loop to stop")
		time.Sleep(100)
	}
}

func (r *updateRoutine) close() {
	if r.running {
		instance.logger.Error(r, "Cannot close update loop before stopping")
		panic("Cannot close update loop before stopping")
	}
	close(r.updateChan)
}

func (r *updateRoutine) loop() {
	instance.logger.Info(r, "Starting")

	r.running = true

	defer instance.recover()

	// Initialize all components
	r.loopInit(instance.currentScene)

	// Calling OnStart for all objects in scene
	r.loopStart(instance.currentScene)

	lastFrameTime := time.Now()

	// Updating every frame or wait for update chan
	for b := range r.updateChan {
		//engine.logger.Info("Update loop", "Here")

		if !b {
			instance.logger.Info(r, "Stopping")
			break
		}

		elapsed := time.Since(lastFrameTime)
		lastFrameTime = time.Now()

		if r.updateRequested {
			//engine.logger.Info("Update loop", "Updating components")

			r.updateRequested = false
			instance.state = StateUpdating

			ctx := newUpdateContext(elapsed.Seconds())

			// Calling OnUpdate for all objects in scene
			r.loopUpdate(instance.currentScene, ctx)
		}

		if r.renderingRequested {
			//engine.logger.Info("Update loop", "Rendering components")

			r.renderingRequested = false
			instance.state = StateRendering

			ctx := newRenderingContext(instance.renderer)

			// Render objects
			r.loopRender(instance.currentScene, ctx)
			instance.renderer.PerformRendering()

			instance.forceRedraw = false
		}

		instance.state = StateStarted

		// Wait until next frame
		timeToSleep := TargetFrameTime - elapsed.Milliseconds()
		time.Sleep(time.Duration(timeToSleep) * time.Millisecond)
	}

	r.loopStop(instance.currentScene)

	r.running = false
}

func (r *updateRoutine) loopInit(obj *SceneObject) {
	obj.init(newInitContext(instance, obj))
	for _, c := range obj.children {
		r.loopInit(c)
	}
}

func (r *updateRoutine) loopStart(obj *SceneObject) {
	if !obj.enabled {
		return
	}
	obj.start()
	for _, c := range obj.children {
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

func (r *updateRoutine) GetName() string {
	return "Update routine"
}

func newUpdateRoutine() *updateRoutine {
	return &updateRoutine{
		running:    false,
		updateChan: make(chan bool, 10),
	}
}
