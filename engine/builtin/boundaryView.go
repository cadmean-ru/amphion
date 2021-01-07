package builtin

import (
	"github.com/cadmean-ru/amphion/common/a"
	"github.com/cadmean-ru/amphion/engine"
	"github.com/cadmean-ru/amphion/rendering"
)

type BoundaryView struct {
	ViewImpl
}

func (r *BoundaryView) GetName() string {
	return "BoundaryView"
}

func (r *BoundaryView) OnDraw(ctx engine.DrawingContext) {
	pr := rendering.NewGeometryPrimitive(rendering.PrimitiveRectangle)
	pr.Transform = transformToRenderingTransform(r.obj.Transform)
	pr.Transform.Position.Z = 100
	pr.Appearance.FillColor = a.TransparentColor()
	pr.Appearance.StrokeColor = a.PinkColor()
	ctx.GetRenderer().SetPrimitive(r.pId, pr, r.ShouldRedraw())
	r.redraw = false
}

func NewBoundaryView() *BoundaryView {
	return &BoundaryView{}
}