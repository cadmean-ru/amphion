package engine

// Basic view component implementation.
type ViewImpl struct {
	Redraw      bool
	Context     InitContext
	PrimitiveId int
	Engine      *AmphionEngine
	SceneObject *SceneObject
}

func (v *ViewImpl) OnInit(ctx InitContext) {
	v.Context = ctx
	v.Engine = ctx.GetEngine()
	v.SceneObject = ctx.GetSceneObject()
	v.Redraw = true
	v.PrimitiveId = -1
}

func (v *ViewImpl) OnStart() {
	if v.PrimitiveId != -1 {
		return
	}
	v.PrimitiveId = v.Context.GetRenderer().AddPrimitive()
	v.ForceRedraw()
}

func (v *ViewImpl) OnStop() {
	v.Context.GetRenderer().RemovePrimitive(v.PrimitiveId)
	v.Redraw = true
	v.PrimitiveId = -1
}

func (v *ViewImpl) OnDraw(_ DrawingContext) {

}

func (v *ViewImpl) ForceRedraw() {
	v.Redraw = true
	v.Engine.GetMessageDispatcher().DispatchDown(v.SceneObject, NewMessage(v.SceneObject, MessageRedraw, nil), MessageMaxDepth)
}

func (v *ViewImpl) ShouldRedraw() bool {
	return v.Redraw || v.Engine.IsForcedToRedraw()
}

func (v *ViewImpl) OnMessage(message Message) bool {
	if message.Code == MessageRedraw {
		v.Redraw = true
		return true
	}

	return true
}


func (v *ViewImpl) GetName() string {
	return "ViewImpl"
}