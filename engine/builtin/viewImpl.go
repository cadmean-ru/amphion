package builtin

import (
	"github.com/cadmean-ru/amphion/common"
	"github.com/cadmean-ru/amphion/engine"
)

type ClickHandler func(vector3 common.IntVector3) bool

type ViewImpl struct {
	redraw  bool
	ctx     engine.InitContext
	pId     int64
	eng     *engine.AmphionEngine
	obj     *engine.SceneObject
	onClick ClickHandler
}

func (v *ViewImpl) OnInit(ctx engine.InitContext) {
	v.ctx = ctx
	v.eng = ctx.GetEngine()
	v.obj = ctx.GetSceneObject()
	v.redraw = true
	v.pId = -1
	//fmt.Printf("Init view: %s\n", v.obj.GetName())
}

func (v *ViewImpl) OnStart() {
	if v.pId != -1 {
		return
	}
	v.pId = v.ctx.GetRenderer().AddPrimitive()
	v.ForceRedraw()
	//fmt.Printf("Start view: %d %s\n", v.pId, v.obj.GetName())
}

func (v *ViewImpl) OnStop() {
	v.ctx.GetRenderer().RemovePrimitive(v.pId)
	v.ForceRedraw()
	v.eng.RequestRendering()
	//fmt.Printf("Stop view: %d %s\n", v.pId, v.obj.GetName())
	v.pId = -1
}

func (v *ViewImpl) OnDraw(_ engine.DrawingContext) {

}

func (v *ViewImpl) ForceRedraw() {
	v.redraw = true
	v.eng.GetMessageDispatcher().DispatchDown(v.obj, engine.NewMessage(v.obj, engine.MessageRedraw, nil))
}

func (v *ViewImpl) ShouldRedraw() bool {
	return v.redraw || v.eng.IsForcedToRedraw()
}

func (v *ViewImpl) OnMessage(message engine.Message) bool {
	if message.Code == engine.MessageRedraw {
		v.redraw = true
		return true
	}

	return true
}


func (v *ViewImpl) GetName() string {
	return "ViewImpl"
}