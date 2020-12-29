package builtin

import (
	"github.com/cadmean-ru/amphion/common/a"
	"github.com/cadmean-ru/amphion/engine"
	"github.com/cadmean-ru/amphion/rendering"
)

// Displays an image given it's resource index
type ImageView struct {
	ViewImpl
	resIndex a.Int
}

func (v *ImageView) OnDraw(ctx engine.DrawingContext) {
	pr := rendering.NewImagePrimitive(v.resIndex)
	pr.Transform = transformToRenderingTransform(v.obj.Transform)
	ctx.GetRenderer().SetPrimitive(v.pId, pr, v.redraw || v.ctx.GetEngine().IsForcedToRedraw())
	v.redraw = false
}

func (v *ImageView) GetName() string {
	return engine.NameOfComponent(v)
}

func NewImageView(index a.Int) *ImageView {
	return &ImageView{
		resIndex: index,
	}
}
