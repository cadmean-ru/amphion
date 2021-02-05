package engine

import (
	"github.com/cadmean-ru/amphion/common/a"
	"github.com/cadmean-ru/amphion/common/require"
	"reflect"
	"runtime"
)

// Manages types of components that are present in the application.
type ComponentsManager struct {
	typesMap map[string]reflect.Type
	handlers map[string]EventHandler
}

// Register component, so an instance of it can be later created using MakeComponent.
func (m *ComponentsManager) RegisterComponentType(component Component) {
	cName := component.GetName()
	cType := reflect.TypeOf(component)
	m.typesMap[cName] = cType
}

// Registers an event handler, so that it can be serialized and deserialized as a part of component's state.
func (m *ComponentsManager) RegisterEventHandler(handler EventHandler) {
	for _, h := range m.handlers {
		if reflect.ValueOf(handler).Pointer() == reflect.ValueOf(h).Pointer() {
			return
		}
	}

	m.handlers[getFunctionName(handler)] = handler
}

// Retrieves the previously registered event handler by it's name.
func (m *ComponentsManager) GetEventHandler(name string) EventHandler {
	return m.handlers[name]
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

// Retrieves component's state and returns it as string-interface map.
func (m *ComponentsManager) GetComponentState(component Component) a.SiMap {
	if sc, ok := component.(StatefulComponent); ok {
		return sc.GetInstanceState()
	}

	t, v := getStructTypeAndValue(component)

	state := make(a.SiMap)

	fCount := t.NumField()
	for i := 0; i < fCount; i++ {
		sf := t.Field(i)
		vf := v.Field(i)
		var key string
		if sf.Tag == "state" {
			key = sf.Name
		} else if sf.Tag.Get("state") != "" {
			key = sf.Tag.Get("state")
		} else {
			continue
		}
		value := vf.Interface()
		if s, ok := value.(a.Stringable); ok {
			if s != nil {
				value = s.ToString()
			} else {
				value = ""
			}
		} else if m1, ok := value.(a.Mappable); ok {
			if m1 != nil {
				value = m1.ToMap()
			} else {
				value = a.SiMap{}
			}
		} else if eh, ok := value.(EventHandler); ok {
			if eh != nil {
				value = getFunctionName(eh)
			} else {
				value = ""
			}
		}
		state[key] = value
	}

	return state
}

// Sets the component's state to the given state map.
func (m *ComponentsManager) SetComponentState(component Component, state a.SiMap) {
	if sc, ok := component.(StatefulComponent); ok {
		sc.SetInstanceState(state)
		return
	}

	t, v := getStructTypeAndValue(component)

	for key, value := range state {
		fCount := t.NumField()
		for i := 0; i < fCount; i++ {
			sf := t.Field(i)
			vf := v.Field(i)
			if sf.Tag == "state" && sf.Name == key {
				m.setReflectValue(vf, value)
			} else if sf.Tag.Get("state") == key {
				m.setReflectValue(vf, value)
			} else {
				continue
			}
		}
	}
}

// Sets the reflect.Value vf (field) of a struct equal to the specified value trying to convert it to the field's type.
func (m *ComponentsManager) setReflectValue(vf reflect.Value, value interface{}) {
	var newValue reflect.Value

	switch vf.Kind() {
	case reflect.Int, reflect.Int32, reflect.Int64, reflect.Uint8, reflect.Float32, reflect.Float64:
		newValue = reflect.ValueOf(require.Number(value, vf.Kind()))
	case reflect.String:
		newValue = reflect.ValueOf(require.String(value))
	case reflect.Struct:
		structValue := reflect.New(vf.Type())

		if structValue.Type().Implements(reflect.TypeOf((*a.Unstringable)(nil)).Elem()) {
			structValue.Interface().(a.Unstringable).FromString(require.String(value))
		}

		if structValue.Type().Implements(reflect.TypeOf((*a.Unmappable)(nil)).Elem()) {
			structValue.Interface().(a.Unmappable).FromMap(a.RequireSiMap(value))
		}

		newValue = reflect.Indirect(structValue)
	case reflect.Ptr:
		m.setReflectValue(reflect.Indirect(vf), value)
	case reflect.Slice:
		if arr, ok := value.([]interface{}); ok {
			arrValue := reflect.MakeSlice(vf.Type(), len(arr), len(arr))

			for i, v := range arr {
				elemValue := reflect.New(vf.Type().Elem())
				m.setReflectValue(elemValue, v)
				arrValue.Index(i).Set(reflect.Indirect(elemValue))
			}

			newValue = arrValue
		} else {
			newValue = reflect.MakeSlice(vf.Type(), 0, 0)
		}
	case reflect.Func:
		if hName, ok := value.(string); ok && hName != "" {
			if h := m.GetEventHandler(hName); h != nil {
				newValue = reflect.ValueOf(h)
			}
		}
	}

	if vf.CanSet() {
		vf.Set(newValue)
	}
}

func (m *ComponentsManager) GetName() string {
	return "ComponentsManager"
}

func newComponentsManager() *ComponentsManager {
	return &ComponentsManager{
		typesMap: make(map[string]reflect.Type),
		handlers: make(map[string]EventHandler),
	}
}

func getStructTypeAndValue(i interface{}) (reflect.Type, reflect.Value) {
	t := reflect.TypeOf(i)
	v := reflect.ValueOf(i)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		v = reflect.Indirect(v)
	}
	return t, v
}

func getFunctionName(i interface{}) string {
	v := reflect.ValueOf(i)
	if v.Kind() != reflect.Func {
		panic("not a func")
	}
	return runtime.FuncForPC(v.Pointer()).Name()
}