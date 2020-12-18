package engine

import (
	"github.com/cadmean-ru/amphion/common"
	"github.com/cadmean-ru/amphion/rendering"
)

type Component interface {
	NamedObject
	OnInit(ctx InitContext)
	OnStart()
	OnStop()
}

type UpdatingComponent interface {
	Component
	OnUpdate(ctx UpdateContext)
}

type ViewComponent interface {
	OnDraw(ctx DrawingContext)
	ForceRedraw()
}

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

type UpdateContext struct {
	DeltaTime float32
}

func newUpdateContext(dTime float32) UpdateContext {
	return UpdateContext{
		DeltaTime: dTime,
	}
}

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

type BoundaryComponent interface {
	Component
	common.Boundary
}

type StatefulComponent interface {
	SaveInstanceState() common.SiMap
	LoadState(state common.SiMap)
}