package builtin

import (
	"github.com/cadmean-ru/amphion/engine"
	"github.com/cadmean-ru/amphion/rendering"
)

// Displays an image given it's resource index
type ImageView struct {
	ViewImpl
	resIndex int
}

func (v *ImageView) OnDraw(ctx engine.DrawingContext) {
	pr := rendering.NewImagePrimitive(v.resIndex)
	pr.Transform = transformToRenderingTransform(v.obj.Transform)
	ctx.GetRenderer().SetPrimitive(v.pId, pr, v.ShouldRedraw())
	v.redraw = false
}

func (v *ImageView) GetName() string {
	return engine.NameOfComponent(v)
}

func NewImageView(index int) *ImageView {
	return &ImageView{
		resIndex: index,
	}
}
