package engine

import "github.com/cadmean-ru/amphion/common/dispatch"

type sceneLifecycleManager struct {
	newSceneObjects   *dispatch.MessageQueue
	startSceneObjects *dispatch.MessageQueue
	stopSceneObjects  *dispatch.MessageQueue
	componentsToStop  *dispatch.MessageQueue
}

func (m *sceneLifecycleManager) start() {
	m.newSceneObjects = dispatch.NewMessageQueue(MaxSceneObjects / 2)
	m.startSceneObjects = dispatch.NewMessageQueue(MaxSceneObjects / 2)
	m.stopSceneObjects = dispatch.NewMessageQueue(MaxSceneObjects / 2)
	m.componentsToStop = dispatch.NewMessageQueue(MaxSceneObjects / 2)
}

func (m *sceneLifecycleManager) initSceneObject(object *SceneObject) {
	m.newSceneObjects.Enqueue(dispatch.NewMessageWithAnyData(0, object))
	m.startSceneObject(object)
}

func (m *sceneLifecycleManager) startSceneObject(object *SceneObject) {
	m.startSceneObjects.Enqueue(dispatch.NewMessageWithAnyData(0, object))
}

func (m *sceneLifecycleManager) stopSceneObject(object *SceneObject) {
	m.stopSceneObjects.Enqueue(dispatch.NewMessageWithAnyData(0, object))
}

func (m *sceneLifecycleManager) stopComponent(c *ComponentContainer) {
	m.componentsToStop.Enqueue(dispatch.NewMessageWithAnyData(0, c))
}

func (m *sceneLifecycleManager) handleSceneObjectsLifecycle() {
	m.newSceneObjects.LockMainChannel()
	m.startSceneObjects.LockMainChannel()
	m.stopSceneObjects.LockMainChannel()
	m.componentsToStop.LockMainChannel()

	for !m.newSceneObjects.IsEmpty() {
		o := m.newSceneObjects.Dequeue().AnyData.(*SceneObject)
		m.onInitSceneObject(o, newInitContext(o))
	}

	for !m.startSceneObjects.IsEmpty() {
		o := m.startSceneObjects.Dequeue().AnyData.(*SceneObject)
		m.onStartSceneObject(o)
	}

	for !m.stopSceneObjects.IsEmpty() {
		o := m.stopSceneObjects.Dequeue().AnyData.(*SceneObject)
		m.onStopSceneObject(o)
	}

	for !m.componentsToStop.IsEmpty() {
		c := m.componentsToStop.Dequeue().AnyData.(*ComponentContainer)
		c.stop()
	}

	m.newSceneObjects.UnlockMainChannel()
	m.startSceneObjects.UnlockMainChannel()
	m.stopSceneObjects.UnlockMainChannel()
	m.componentsToStop.UnlockMainChannel()
}

func (m *sceneLifecycleManager) loopInit(obj *SceneObject) {
	temp := make([]*SceneObject, len(obj.children))
	copy(temp, obj.children)
	for _, c := range temp {
		m.loopInit(c)
	}

	m.onInitSceneObject(obj, newInitContext(obj))
}

func (m *sceneLifecycleManager) onInitSceneObject(o *SceneObject, ctx InitContext) {
	for _, c := range o.components {
		if c.initialized {
			continue
		}
		instance.currentComponent = c.component
		c.component.OnInit(ctx)
		c.initialized = true
	}

	o.initialized = true
	instance.currentComponent = nil
}

func (m *sceneLifecycleManager) loopStart(obj *SceneObject) {
	if !obj.enabled {
		return
	}

	temp := make([]*SceneObject, len(obj.children))
	copy(temp, obj.children)
	for _, c := range temp {
		m.loopStart(c)
	}

	m.onStartSceneObject(obj)
}

func (m *sceneLifecycleManager) onStartSceneObject(o *SceneObject) {
	for _, c := range o.components {
		if c.IsDirty() || c.started {
			continue
		}
		instance.currentComponent = c.component
		c.component.OnStart()
		c.started = true
	}

	o.started = true
	instance.currentComponent = nil
}


func (m *sceneLifecycleManager) loopUpdate(obj *SceneObject, ctx UpdateContext) {
	if !obj.enabled {
		return
	}

	for _, c := range obj.children {
		m.loopUpdate(c, ctx)
	}

	m.onUpdateSceneObject(obj, ctx)
}

func (m *sceneLifecycleManager) onUpdateSceneObject(o *SceneObject, ctx UpdateContext) {
	for _, c := range o.updatingComponents {
		if c.IsDirty() || !c.started {
			continue
		}

		instance.currentComponent = c.component
		c.component.(UpdatingComponent).OnUpdate(ctx)
	}

	instance.currentComponent = nil
}

func (m *sceneLifecycleManager) loopLayout(obj *SceneObject) {
	obj.TraversePostOrder(func(object *SceneObject) {
		if !object.HasLayout() {
			return
		}

		m.onLayoutSceneObject(object)
	})
}

func (m *sceneLifecycleManager) onLayoutSceneObject(o *SceneObject) {
	if o.layout == nil || o.layout.IsDirty() {
		return
	}

	instance.currentComponent = o.layout.component
	o.layout.component.(Layout).LayoutChildren()
	instance.currentComponent = nil
}

func (m *sceneLifecycleManager) loopRender(obj *SceneObject) {
	obj.TraversePreOrder(func(object *SceneObject) bool {
		m.onDrawSceneObject(object, newDrawingContext(object))
		return true
	})
}

func (m *sceneLifecycleManager) onDrawSceneObject(o *SceneObject, ctx DrawingContext) {
	if o.HasView() {
		if o.view.IsDirty() || !o.view.started {
			return
		}

		view := o.view.component.(ViewComponent)
		instance.currentComponent = o.view.component

		if !view.ShouldDraw() {
			return
		}

		view.OnDraw(ctx)
	}
	instance.currentComponent = nil
}

func (m *sceneLifecycleManager) loopStop(obj *SceneObject) {
	for _, c := range obj.children {
		m.loopStop(c)
	}

	m.onStopSceneObject(obj)
}

func (m *sceneLifecycleManager) onStopSceneObject(o *SceneObject) {
	for _, c := range o.components {
		if !c.enabled || !c.started {
			continue
		}
		instance.currentComponent = c.component
		c.component.OnStop()
		c.started = false
	}
	o.started = false
	instance.currentComponent = nil
}