package builtin

import (
	"github.com/cadmean-ru/amphion/common/a"
	"github.com/cadmean-ru/amphion/engine"
	"github.com/cadmean-ru/amphion/rendering"
)

type BezierView struct {
	engine.ViewImpl
	rendering.Appearance
	ControlPoint1, ControlPoint2 a.Vector3
}

func (b *BezierView) OnDraw(ctx engine.DrawingContext) {
	bezier := rendering.NewBezierPrimitive(b.ControlPoint1.Round(), b.ControlPoint2.Round())
	bezier.Transform = transformToRenderingTransform(b.SceneObject.Transform)
	bezier.Appearance = b.Appearance
	ctx.GetRenderer().SetPrimitive(b.PrimitiveId, bezier, b.ShouldRedraw())
}

func (b *BezierView) GetName() string {
	return "BezierView"
}

func NewBezierView(cp1, cp2 a.Vector3) *BezierView {
	return &BezierView{
		ControlPoint1: cp1,
		ControlPoint2: cp2,
		Appearance: rendering.Appearance{
			FillColor:    a.TransparentColor(),
			StrokeColor:  a.BlackColor(),
			StrokeWeight: 1,
		},
	}
}
