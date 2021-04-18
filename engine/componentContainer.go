package engine

// Wrap around container interface, that can be enabled or disabled
type ComponentContainer struct {
	enabled     bool
	sceneObject *SceneObject
	component   Component
	initialized bool
	started     bool
}

func (c *ComponentContainer) SetEnabled(enabled bool) {
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
	return !c.initialized || !c.enabled
}

func NewComponentContainer(sceneObject *SceneObject, component Component) *ComponentContainer {
	return &ComponentContainer{
		enabled:     true,
		component:   component,
		sceneObject: sceneObject,
	}
}