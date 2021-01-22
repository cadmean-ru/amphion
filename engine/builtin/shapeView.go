package builtin

import (
	"github.com/cadmean-ru/amphion/common/a"
	"github.com/cadmean-ru/amphion/engine"
	"github.com/cadmean-ru/amphion/rendering"
)

type ShapeType byte

const (
	ShapeRectangle ShapeType = rendering.PrimitiveRectangle
	ShapeEllipse   ShapeType = rendering.PrimitiveEllipse
	ShapeTriangle  ShapeType = rendering.PrimitiveTriangle
	ShapeLine      ShapeType = rendering.PrimitiveLine
	ShapePoint     ShapeType = rendering.PrimitivePoint
)

// Displays a basic shape: rectangle, ellipse, triangle, line, point.
type ShapeView struct {
	engine.ViewImpl
	FillColor    a.Color `state:"fillColor"`
	StrokeColor  a.Color `state:"strokeColor"`
	StrokeWeight byte    `state:"strokeWeight"`
	CornerRadius byte    `state:"cornerRadius"`
	pType        byte    `state:"pType"`
}

func (c *ShapeView) OnDraw(ctx engine.DrawingContext) {
	pr := rendering.NewGeometryPrimitive(c.pType)
	pr.Transform = transformToRenderingTransform(c.SceneObject.Transform)
	pr.Appearance = rendering.Appearance{
		FillColor:    c.FillColor,
		StrokeColor:  c.StrokeColor,
		StrokeWeight: c.StrokeWeight,
		CornerRadius: c.CornerRadius,
	}
	ctx.GetRenderer().SetPrimitive(c.PrimitiveId, pr, c.ShouldRedraw())
	c.Redraw = false
}

func (c *ShapeView) GetName() string {
	return engine.NameOfComponent(c)
}

func NewShapeView(pType ShapeType) *ShapeView {
	if pType < 0 || pType > 5 {
		panic("Invalid primitive type")
	}

	return &ShapeView{
		pType:        byte(pType),
		FillColor:    a.WhiteColor(),
		StrokeColor:  a.BlackColor(),
		StrokeWeight: 1,
	}
}
