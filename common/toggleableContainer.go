package common

type Toggleable interface {
	SetEnabled(enabled bool)
	IsEnabled() bool
}

type ToggleableContainer interface {
	Toggleable
	Container
}

type ToggleableContainerImpl struct {
	enabled bool
	value   interface{}
}

func (t *ToggleableContainerImpl) SetEnabled(enabled bool) {
	t.enabled = enabled
}

func (t *ToggleableContainerImpl) IsEnabled() bool {
	return t.enabled
}

func (t *ToggleableContainerImpl) SetValue(value interface{}) {
	t.value = value
}

func (t *ToggleableContainerImpl) GetValue() interface{} {
	return t.value
}
