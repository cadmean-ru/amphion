package engine

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

func (v *ViewImpl) OnStop() {
	v.Context.GetRenderingNode().RemovePrimitive(v.PrimitiveId)
	v.ShouldRedraw = true
	v.PrimitiveId = -1
}

func (v *ViewImpl) OnDraw(_ DrawingContext) {

}

func (v *ViewImpl) Redraw() {
	v.ShouldRedraw = true
	v.Engine.GetMessageDispatcher().DispatchDown(v.SceneObject, NewMessage(v.SceneObject, MessageRedraw, nil), MessageMaxDepth)
}

func (v *ViewImpl) ShouldDraw() bool {
	return v.ShouldRedraw || v.Engine.IsForcedToRedraw()
}

func (v *ViewImpl) OnMessage(message Message) bool {
	if message.Code == MessageRedraw {
		v.ShouldRedraw = true
		return true
	}

	return true
}