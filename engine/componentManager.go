package engine

import (
	"reflect"
)

// Manages types of components that are present in the application.
type ComponentsManager struct {
	typesMap map[string]reflect.Type
}

// Register component, so an instance of it can be later created using MakeComponent.
func (m *ComponentsManager) RegisterComponentType(component Component) {
	cName := component.GetName()
	cType := reflect.TypeOf(component)
	m.typesMap[cName] = cType
}

// Creates an instance of a component with the specified name.
// Returns the new instance of component ot nil if component with the name was not registered.
func (m *ComponentsManager) MakeComponent(name string) Component {
	var t reflect.Type
	var ok bool
	if t, ok = m.typesMap[name]; !ok {
		return nil
	}

	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	c := reflect.New(t)
	return c.Interface().(Component)
}

func (m *ComponentsManager) GetName() string {
	return "ComponentsManager"
}

func newComponentsManager() *ComponentsManager {
	return &ComponentsManager{
		typesMap: make(map[string]reflect.Type),
	}
}
