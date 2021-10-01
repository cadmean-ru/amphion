package engine

import (
	"github.com/cadmean-ru/amphion/common/a"
	"github.com/cadmean-ru/amphion/common/dispatch"
)

// ViewImpl is a basic view component implementation.
type ViewImpl struct {
	ShouldRedraw bool
	Context      InitContext
	PrimitiveId  int
	Engine       *AmphionEngine
	SceneObject  *SceneObject
}

func (v *ViewImpl) OnInit(ctx InitContext) {
	v.Context = ctx
	v.Engine = ctx.GetEngine()
	v.SceneObject = ctx.GetSceneObject()
	v.ShouldRedraw = true
	v.PrimitiveId = -1
}

func (v *ViewImpl) OnStart() {
	if v.PrimitiveId != -1 {
		return
	}
	v.PrimitiveId = v.Context.GetRenderingNode().AddPrimitive()
	v.Redraw()
}

func (v *ViewImpl) OnUpdate(_ UpdateContext) {

}

func (v *ViewImpl) OnLateUpdate(_ UpdateContext) {

}

func (v *ViewImpl) OnDraw(_ DrawingContext) {

}

func (v *ViewImpl) OnStop() {
	v.Context.GetRenderingNode().RemovePrimitive(v.PrimitiveId)
	v.ShouldRedraw = true
	v.PrimitiveId = -1
}

func (v *ViewImpl) Redraw() {
	v.ShouldRedraw = true
	v.Engine.GetMessageDispatcher().DispatchDown(v.SceneObject, dispatch.NewMessageFrom(v.SceneObject, MessageRedraw), MessageMaxDepth)
}

func (v *ViewImpl) ShouldDraw() bool {
	return v.ShouldRedraw || v.Engine.IsForcedToRedraw()
}

func (v *ViewImpl) MeasureContents() a.Vector3 {
	 return a.ZeroVector()
}

func (v *ViewImpl) OnMessage(message *dispatch.Message) bool {
	if message.What == MessageRedraw {
		v.ShouldRedraw = true
		return true
	}

	return true
}