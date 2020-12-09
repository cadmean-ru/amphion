package builtin

import (
	"github.com/cadmean-ru/amphion/common"
	"github.com/cadmean-ru/amphion/engine"
	"github.com/cadmean-ru/amphion/rendering"
)

type ImageView struct {
	ViewImpl
	resIndex common.AInt
}

func (v *ImageView) OnDraw(ctx engine.DrawingContext) {
	pr := rendering.NewImagePrimitive(v.resIndex)
	pr.Transform = transformToRenderingTransform(v.obj.Transform)
	ctx.GetRenderer().SetPrimitive(v.pId, pr, v.redraw || v.ctx.GetEngine().IsForcedToRedraw())
}

func (v *ImageView) GetName() string {
	return "ImageView"
}

func NewImageView(index common.AInt) *ImageView {
	return &ImageView{
		resIndex: index,
	}
}
