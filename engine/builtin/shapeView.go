package builtin

import (
	"github.com/cadmean-ru/amphion/common/a"
	"github.com/cadmean-ru/amphion/engine"
	"github.com/cadmean-ru/amphion/rendering"
)

// Displays a basic shape: rectangle, ellipse, triangle, line, point.
type ShapeView struct {
	ViewImpl
	FillColor    a.Color `state:"fillColor"`
	StrokeColor  a.Color `state:"strokeColor"`
	StrokeWeight byte    `state:"strokeWeight"`
	CornerRadius byte    `state:"cornerRadius"`
	pType        byte    `state:"pType"`
}

func (c *ShapeView) OnDraw(ctx engine.DrawingContext) {
	pr := rendering.NewGeometryPrimitive(c.pType)
	pr.Transform = transformToRenderingTransform(c.obj.Transform)
	pr.Appearance = rendering.Appearance{
		FillColor:    c.FillColor,
		StrokeColor:  c.StrokeColor,
		StrokeWeight: c.StrokeWeight,
		CornerRadius: c.CornerRadius,
	}
	ctx.GetRenderer().SetPrimitive(c.pId, pr, c.ShouldRedraw())
	c.redraw = false
}

func (c *ShapeView) GetName() string {
	return engine.NameOfComponent(c)
}

func NewShapeView(pType byte) *ShapeView {
	if pType < 0 || pType > 5 {
		panic("Invalid primitive type")
	}

	return &ShapeView{
		pType:        pType,
		FillColor:    a.WhiteColor(),
		StrokeColor:  a.BlackColor(),
		StrokeWeight: 1,
	}
}
