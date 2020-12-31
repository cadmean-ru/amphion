package engine

import (
	"github.com/cadmean-ru/amphion/common"
	"github.com/cadmean-ru/amphion/common/a"
	"github.com/cadmean-ru/amphion/common/require"
	"github.com/cadmean-ru/amphion/rendering"
	"reflect"
)

// Basic component interface. A component is a piece of functionality, that can be attached to scene objects.
type Component interface {
	NamedObject

	// This method is called only once when the component is first created.
	OnInit(ctx InitContext)

	// This method is called every time the component is being enabled.
	// If the scene object is enabled on component attachment this method will also be called.
	OnStart()

	// This method is called when the component is being disabled.
	OnStop()
}

// Interface for components that can receive updates.
type UpdatingComponent interface {
	Component
	OnUpdate(ctx UpdateContext)
}

// Interface for views.
type ViewComponent interface {
	OnDraw(ctx DrawingContext)
	ForceRedraw()
}

// Contains necessary objects for component initialization
type InitContext struct {
	engine      *AmphionEngine
	sceneObject *SceneObject
}

func (c InitContext) GetEngine() *AmphionEngine {
	return c.engine
}

func (c InitContext) GetSceneObject() *SceneObject {
	return c.sceneObject
}

func (c InitContext) GetRenderer() rendering.Renderer {
	return c.engine.renderer
}

func (c InitContext) GetLogger() *Logger {
	return c.engine.logger
}

func newInitContext(engine *AmphionEngine, object *SceneObject) InitContext {
	return InitContext{
		engine:      engine,
		sceneObject: object,
	}
}

// Contains info about current update cycle
type UpdateContext struct {
	DeltaTime float32
}

func newUpdateContext(dTime float32) UpdateContext {
	return UpdateContext{
		DeltaTime: dTime,
	}
}

// Contains renderer
type DrawingContext struct {
	renderer rendering.Renderer
}

func (c DrawingContext) GetRenderer() rendering.Renderer {
	return c.renderer
}

func newRenderingContext(renderer rendering.Renderer) DrawingContext {
	return DrawingContext{
		renderer: renderer,
	}
}

// A component, that determines the bounding box of an object in the scene. Used for mouse interactions.
type BoundaryComponent interface {
	Component
	common.Boundary
}

// Interface for components, that can persist state.
type StatefulComponent interface {
	GetInstanceState() a.SiMap
	SetInstanceState(state a.SiMap)
}

// Checks if the given component has state.
// A component becomes stateful if is implements StatefulComponent interface or contains fields with state tag.
func IsStatefulComponent(component Component) bool {
	if _, ok := component.(StatefulComponent); ok {
		return true
	}

	t, _ := getStructTypeAndValue(component)

	fCount := t.NumField()
	for i := 0; i < fCount; i++ {
		f := t.Field(i)
		if f.Tag == "state" || f.Tag.Get("state") != "" {
			return true
		}
	}

	return false
}

// Retrieves component's state and returns it as string-interface map.
func GetComponentState(component Component) a.SiMap {
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
		if m, ok := value.(a.Mappable); ok {
			value = m.ToMap()
		}
		state[key] = value
	}

	return state
}

// Sets the component's state to the given state map.
func SetComponentState(component Component, state a.SiMap) {
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
				setReflectValue(vf, value)
			} else if sf.Tag.Get("state") == key {
				setReflectValue(vf, value)
			} else {
				continue
			}
		}
	}
}

func setReflectValue(vf reflect.Value, value interface{}) {
	var newValue reflect.Value

	switch vf.Kind() {
	case reflect.Int, reflect.Int32, reflect.Int64, reflect.Uint8, reflect.Float32, reflect.Float64:
		newValue = reflect.ValueOf(require.Number(value, vf.Kind()))
	case reflect.String:
		newValue = reflect.ValueOf(require.String(value))
	case reflect.Struct:
		structValue := reflect.New(vf.Type())
		if structValue.Type().Implements(reflect.TypeOf((*a.Unmappable)(nil)).Elem()) {
			structValue.Interface().(a.Unmappable).FromMap(a.RequireSiMap(value))
		}
		newValue = reflect.Indirect(structValue)
	case reflect.Ptr:
		setReflectValue(reflect.Indirect(vf), value)
	case reflect.Slice:
		if arr, ok := value.([]interface{}); ok {
			arrValue := reflect.MakeSlice(vf.Type(), len(arr), len(arr))

			for i, v := range arr {
				elemValue := reflect.New(vf.Type().Elem())
				setReflectValue(elemValue, v)
				arrValue.Index(i).Set(reflect.Indirect(elemValue))
			}

			newValue = arrValue
		} else {
			newValue = reflect.MakeSlice(vf.Type(), 0, 0)
		}
	}

	if vf.CanSet() {
		vf.Set(newValue)
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