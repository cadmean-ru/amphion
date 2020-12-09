package builtin

import (
	"github.com/cadmean-ru/amphion/common"
	"github.com/cadmean-ru/amphion/engine"
	"github.com/cadmean-ru/amphion/rendering"
)

type BoundaryView struct {
	obj       *engine.SceneObject
	engine    *engine.AmphionEngine
	renderer  rendering.Renderer
	pId       int64
	redraw    bool
}

func (r *BoundaryView) GetName() string {
	return "BoundaryView"
}

func (r *BoundaryView) OnInit(ctx engine.InitContext) {
	r.obj = ctx.GetSceneObject()
	r.engine = ctx.GetEngine()
	r.renderer = ctx.GetRenderer()
}

func (r *BoundaryView) OnStart() {
	r.pId = r.renderer.AddPrimitive()
}

func (r *BoundaryView) OnStop() {

}

func (r *BoundaryView) OnDraw(ctx engine.DrawingContext) {
	pr := rendering.NewGeometryPrimitive(rendering.PrimitiveRectangle)
	pr.Transform = transformToRenderingTransform(r.obj.Transform)
	pr.Transform.Position.Z = 100
	pr.Appearance.FillColor = common.TransparentColor()
	pr.Appearance.StrokeColor = common.PinkColor()
	ctx.GetRenderer().SetPrimitive(r.pId, pr, r.redraw || r.engine.IsForcedToRedraw())
	r.redraw = false
}

func (r *BoundaryView) ForceRedraw() {
	r.redraw = true
}

func (r *BoundaryView) OnMessage(message engine.Message) bool {
	if message.Code == engine.MessageRedraw {
		r.redraw = true
	}

	return true
}

func NewBoundaryView() *BoundaryView {
	return &BoundaryView{
		redraw: true,
	}
}