package builtin

import (
	"github.com/cadmean-ru/amphion/common/a"
	"github.com/cadmean-ru/amphion/engine"
	"github.com/cadmean-ru/amphion/rendering"
)

type BoundaryView struct {
	engine.ViewImpl
}

func (r *BoundaryView) OnDraw(ctx engine.DrawingContext) {
	pr := rendering.NewGeometryPrimitive(rendering.PrimitiveRectangle)
	pr.Transform = r.SceneObject.Transform.ToRenderingTransform()
	pr.Transform.Position.Z = 100
	pr.Appearance.FillColor = a.Transparent()
	pr.Appearance.StrokeColor = a.Pink()
	ctx.GetRenderingNode().SetPrimitive(r.PrimitiveId, pr)
	r.ShouldRedraw = false
}

func NewBoundaryView() *BoundaryView {
	return &BoundaryView{}
}