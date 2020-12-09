package builtin

import (
	"github.com/cadmean-ru/amphion/engine"
)

type ViewImpl struct {
	redraw bool
	ctx    engine.InitContext
	pId    int64
	eng    *engine.AmphionEngine
	obj    *engine.SceneObject
}

func (v *ViewImpl) OnInit(ctx engine.InitContext) {
	v.ctx = ctx
	v.eng = ctx.GetEngine()
	v.obj = ctx.GetSceneObject()
	v.redraw = true
}

func (v *ViewImpl) OnStart() {
	v.pId = v.ctx.GetRenderer().AddPrimitive()
}

func (v *ViewImpl) OnStop() {
	v.ctx.GetRenderer().RemovePrimitive(v.pId)
}

func (v *ViewImpl) OnDraw(_ engine.DrawingContext) {

}

func (v *ViewImpl) ForceRedraw() {
	v.redraw = true
	v.eng.GetMessageDispatcher().DispatchDown(v.obj, engine.NewMessage(v.obj, engine.MessageRedraw, nil))
}

func (v *ViewImpl) GetName() string {
	return "ViewImpl"
}