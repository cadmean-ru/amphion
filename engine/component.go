package engine

import (
	"github.com/cadmean-ru/amphion/common"
	"github.com/cadmean-ru/amphion/common/a"
	"github.com/cadmean-ru/amphion/rendering"
)

//Component is a basic component interface.
//A component is a piece of functionality, that can be attached to SceneObject objects.
type Component interface {
	// OnInit is called only once when the component is first created.
	OnInit(ctx InitContext)

	// OnStart is called every time the component is being enabled.
	// If the SceneObject object is enabled on component attachment this method will also be called.
	OnStart()

	// OnStop is called when the component is being disabled.
	OnStop()
}

// UpdatingComponent is an interface for components that can receive updates.
type UpdatingComponent interface {
	Component
	OnUpdate(ctx UpdateContext)
	OnLateUpdate(ctx UpdateContext)
}

// ViewComponent is an interface for views.
type ViewComponent interface {
	Component
	Measurable
	OnDraw(ctx DrawingContext)
	ShouldDraw() bool
	Redraw()
}

// InitContext contains necessary objects for component initialization.
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

func (c InitContext) GetRenderer() *rendering.ARenderer {
	return c.engine.renderer
}

func (c InitContext) GetLogger() *Logger {
	return c.engine.logger
}

func (c InitContext) GetRenderingNode() *rendering.Node {
	return c.sceneObject.renderingNode
}

func newInitContext(object *SceneObject) InitContext {
	return InitContext{
		engine:      instance,
		sceneObject: object,
	}
}

type DrawingContext struct {
	sceneObject *SceneObject
}

func (c DrawingContext) GetRenderer() *rendering.ARenderer {
	return instance.renderer
}

func (c DrawingContext) GetRenderingNode() *rendering.Node {
	return c.sceneObject.renderingNode
}

func newDrawingContext(object *SceneObject) DrawingContext {
	return DrawingContext{
		sceneObject: object,
	}
}

// UpdateContext contains info about current update cycle
type UpdateContext struct {
	DeltaTime float32
}

func newUpdateContext(dTime float32) UpdateContext {
	return UpdateContext{
		DeltaTime: dTime,
	}
}

// BoundaryComponent is a component, that determines the bounding box of an object in the SceneObject. Used for mouse interactions.
type BoundaryComponent interface {
	Component
	common.Boundary
	IsSolid() bool
}

// StatefulComponent is an interface for components, that can persist state.
type StatefulComponent interface {
	Component
	GetInstanceState() a.SiMap
	SetInstanceState(state a.SiMap)
}

// IsStatefulComponent checks if the given component has state.
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
