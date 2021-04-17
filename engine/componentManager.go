package engine

import (
	"github.com/cadmean-ru/amphion/common/a"
	"github.com/cadmean-ru/amphion/common/require"
	"reflect"
	"runtime"
)

//ComponentsManager keeps track  of types of components and event handlers. that are present in the application.
//It also allows getting and setting the state of a component.
type ComponentsManager struct {
	typesMap map[string]reflect.Type
	handlers map[string]EventHandler
}

//RegisterComponentType registers component, so an instance of it can be later created using MakeComponent.
func (m *ComponentsManager) RegisterComponentType(component Component) {
	cName := component.GetName()
	cType := reflect.TypeOf(component)
	m.typesMap[cName] = cType
}

//RegisterEventHandler registers an event handler, so that it can be serialized and deserialized as a part of component's state.
func (m *ComponentsManager) RegisterEventHandler(handler EventHandler) {
	for _, h := range m.handlers {
		if reflect.ValueOf(handler).Pointer() == reflect.ValueOf(h).Pointer() {
			return
		}
	}

	m.handlers[getFunctionName(handler)] = handler
}

//GetEventHandler retrieves the previously registered event handler by it's name.
func (m *ComponentsManager) GetEventHandler(name string) EventHandler {
	return m.handlers[name]
}

// MakeComponent creates an instance of a component with the specified name.
// Returns the new instance of component or nil if component with that name was not registered.
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

// GetComponentState retrieves component's state and returns it as string-interface map.
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

// SetComponentState sets the component's state from the given state map.
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
		newValue = m.setReflectNumberValue(vf, value)
	case reflect.String:
		newValue = m.setReflectStringValue(vf, value)
	case reflect.Struct:
		newValue = m.setReflectStructValue(vf, value)
	case reflect.Ptr:
		newValue = m.setReflectPtrValue(vf, value)
	case reflect.Slice:
		newValue = m.setReflectSliceValue(vf, value)
	case reflect.Func:
		newValue = m.setReflectFuncValue(vf, value)
	}

	if vf.CanSet() {
		vf.Set(newValue)
	}
}

func (m *ComponentsManager) setReflectNumberValue(vf reflect.Value, value interface{}) reflect.Value {
	var newValue reflect.Value

	if vf.Type().AssignableTo(reflect.TypeOf(a.ResId(0))) {
		newValue = reflect.ValueOf(a.ResId(require.Int(value)))
	} else if vf.Type().AssignableTo(reflect.TypeOf(a.TextAlign(0))) {
		if IsSpecialValueString(value) {
			newValue = reflect.ValueOf(GetSpecialValueFromString(value))
		} else {
			newValue = reflect.ValueOf(a.TextAlign(require.Int(value)))
		}
	} else {
		if IsSpecialValueString(value) {
			newValue = reflect.ValueOf(require.Number(GetSpecialValueFromString(value), vf.Kind()))
		} else {
			newValue = reflect.ValueOf(require.Number(value, vf.Kind()))
		}
	}

	return newValue
}

func (m *ComponentsManager) setReflectStringValue(_ reflect.Value, value interface{}) reflect.Value {
	return reflect.ValueOf(require.String(value))
}

func (m *ComponentsManager) setReflectStructValue(vf reflect.Value, value interface{}) reflect.Value {
	structValue := reflect.New(vf.Type())

	if structValue.Type().Implements(reflect.TypeOf((*a.Unstringable)(nil)).Elem()) {
		structValue.Interface().(a.Unstringable).FromString(require.String(value))
	} else if structValue.Type().Implements(reflect.TypeOf((*a.Unmappable)(nil)).Elem()) {
		structValue.Interface().(a.Unmappable).FromMap(a.RequireSiMap(value))
	}

	return reflect.Indirect(structValue)
}

func (m *ComponentsManager) setReflectPtrValue(vf reflect.Value, value interface{}) reflect.Value {
	ptr := reflect.New(vf.Type().Elem())
	m.setReflectValue(reflect.Indirect(ptr), value)
	return ptr
}

func (m *ComponentsManager) setReflectSliceValue(vf reflect.Value, value interface{}) reflect.Value {
	vv := reflect.ValueOf(value)

	if vv.Kind() != reflect.Slice {
		return reflect.MakeSlice(vf.Type(), 0, 0)
	}

	arrValue := reflect.MakeSlice(vf.Type(), vv.Len(), vv.Len())

	for i := 0; i < vv.Len(); i++ {
		elemValue := reflect.Indirect(reflect.New(vf.Type().Elem()))
		m.setReflectValue(elemValue, vv.Index(i).Interface())
		arrValue.Index(i).Set(elemValue)
	}

	return arrValue
}

func (m *ComponentsManager) setReflectFuncValue(vf reflect.Value, value interface{}) reflect.Value {
	if hName, ok := value.(string); vf.Type().AssignableTo(reflect.TypeOf(eh)) && ok && hName != "" {
		if h := m.GetEventHandler(hName); h != nil {
			return reflect.ValueOf(h)
		}
	}

	return reflect.Zero(vf.Type())
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

//getStructTypeAndValue retrieves the reflect.Type and reflect.Value of a struct.
//If the given value is a pointer to struct, returns the type and value of struct it points to.
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

func eh(event AmphionEvent) bool {
	return true
}