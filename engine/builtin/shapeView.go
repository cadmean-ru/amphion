package builtin

import (
	"github.com/cadmean-ru/amphion/common"
	"github.com/cadmean-ru/amphion/engine"
	"github.com/cadmean-ru/amphion/rendering"
)

type ShapeView struct {
	ViewImpl
	rendering.Appearance
	pType     common.AByte
}

func (c *ShapeView) OnDraw(ctx engine.DrawingContext) {
	pr := rendering.NewGeometryPrimitive(c.pType)
	pr.Transform = transformToRenderingTransform(c.obj.Transform)
	pr.Appearance = c.Appearance
	ctx.GetRenderer().SetPrimitive(c.pId, pr, c.ShouldRedraw())
	c.redraw = false
}

func (c *ShapeView) GetName() string {
	return "ShapeView"
}

func NewShapeView(pType common.AByte) *ShapeView {
	if pType < 0 || pType > 5 {
		panic("Invalid primitive type")
	}

	return &ShapeView{
		pType:      pType,
		Appearance: rendering.DefaultAppearance(),
	}
}