package engine

// Wrap around container interface, that can be enabled or disabled
type ComponentContainer struct {
	enabled     bool
	sceneObject *SceneObject
	component   Component
	initialized bool
}

func (c *ComponentContainer) SetEnabled(enabled bool) {
	c.enabled = enabled
	instance.RequestRendering()
}

func (c *ComponentContainer) IsEnabled() bool {
	return c.enabled
}

func (c *ComponentContainer) GetComponent() Component {
	return c.component
}

func NewComponentContainer(sceneObject *SceneObject, component Component) *ComponentContainer {
	return &ComponentContainer{
		enabled:     true,
		component:   component,
		sceneObject: sceneObject,
	}
}