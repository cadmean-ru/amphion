package engine

import (
	"github.com/cadmean-ru/amphion/common"
	"github.com/cadmean-ru/amphion/common/a"
	"github.com/cadmean-ru/amphion/rendering"
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
	Component
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