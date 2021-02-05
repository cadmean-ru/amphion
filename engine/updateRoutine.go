package engine

import (
	"github.com/cadmean-ru/amphion/frontend"
	"time"
)

type updateRoutine struct {
	running            bool
	updateChan         chan bool
	updateRequested    bool
	renderingRequested bool
	newSceneObjects    []*SceneObject
	startSceneObjects  []*SceneObject
	stopSceneObjects   []*SceneObject
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
	r.renderingRequested = true
	r.requestUpdate()
}

func (r *updateRoutine) stop() {
	r.updateChan<-false
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
	instance.currentScene.setInCurrentScene(true)
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

		if len(r.newSceneObjects) > 0 {
			for _, o := range r.newSceneObjects {
				//r.loopInit(o)
				o.init(newInitContext(instance, o))
			}
			r.newSceneObjects = make([]*SceneObject, 0)
		}

		if len(r.startSceneObjects) > 0 {
			for _, o := range r.startSceneObjects {
				//r.loopStart(o)
				o.start()
			}
			r.startSceneObjects = make([]*SceneObject, 0)
		}

		if len(r.stopSceneObjects) > 0 {
			for _, o := range r.stopSceneObjects {
				//r.loopStop(o)
				o.stop()
			}
			r.stopSceneObjects = make([]*SceneObject, 0)
		}

		if instance.suspend {
			instance.state = StateStarted
			r.updateRequested = false
			r.renderingRequested = false
			continue
		}

		if r.updateRequested {
			//engine.logger.Info("Update loop", "Updating components")

			r.updateRequested = false
			instance.state = StateUpdating

			ctx := newUpdateContext(float32(elapsed.Seconds()))

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
			instance.front.ReceiveMessage(frontend.NewFrontendMessage(frontend.MessageRender))

			instance.forceRedraw = false
		}

		instance.state = StateStarted

		// Wait until next frame
		timeToSleep := TargetFrameTime - elapsed.Milliseconds()
		time.Sleep(time.Duration(timeToSleep) * time.Millisecond)
	}

	r.loopStop(instance.currentScene)
	instance.currentScene.setInCurrentScene(false)

	instance.renderer.Clear()
	//instance.renderer.PerformRendering()

	r.running = false
	r.newSceneObjects = make([]*SceneObject, 0)
	r.startSceneObjects = make([]*SceneObject, 0)
	r.stopSceneObjects = make([]*SceneObject, 0)

	instance.logger.Info(r, "Stopped")
}

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

func (r *updateRoutine) GetName() string {
	return "Update routine"
}

func newUpdateRoutine() *updateRoutine {
	return &updateRoutine{
		running:           false,
		updateChan:        make(chan bool, 10),
		newSceneObjects:   make([]*SceneObject, 0),
		startSceneObjects: make([]*SceneObject, 0),
		stopSceneObjects:  make([]*SceneObject, 0),
	}
}
