package engine

// ComponentContainer is a wrap around container interface, that can be enabled or disabled
type ComponentContainer struct {
	enabled       bool
	sceneObject   *SceneObject
	component     Component
	initialized   bool
	started       bool
	willBeRemoved bool
}

func (c *ComponentContainer) SetEnabled(enabled bool) {
	if c.willBeRemoved || c.enabled == enabled {
		return
	}

	c.enabled = enabled

	if c.sceneObject.inCurrentScene {
		instance.RequestRendering()
	}
}

func (c *ComponentContainer) IsEnabled() bool {
	return c.enabled
}

func (c *ComponentContainer) GetComponent() Component {
	return c.component
}

func (c *ComponentContainer) IsDirty() bool {
	return !c.initialized || !c.enabled || c.willBeRemoved
}

func (c *ComponentContainer) stop() {
	if !c.enabled || !c.started {
		return
	}

	instance.updateRoutine.currentComponent = c
	c.component.OnStop()
	instance.updateRoutine.currentComponent = nil
}

func NewComponentContainer(sceneObject *SceneObject, component Component) *ComponentContainer {
	return &ComponentContainer{
		enabled:     true,
		component:   component,
		sceneObject: sceneObject,
	}
}
