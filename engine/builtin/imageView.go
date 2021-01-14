package builtin

import (
	"github.com/cadmean-ru/amphion/engine"
	"github.com/cadmean-ru/amphion/rendering"
)

// Displays an image given it's resource index
type ImageView struct {
	ViewImpl
	ResIndex int    `state:"resIndex"`
	ImageUrl string `state:"imageUrl"`
}

func (v *ImageView) OnDraw(ctx engine.DrawingContext) {
	var url string

	if v.ResIndex > -1 {
		url = v.eng.GetResourceManager().FullPathOf(v.ResIndex)
	} else {
		url = v.ImageUrl
	}

	pr := rendering.NewImagePrimitive(url)
	pr.Transform = transformToRenderingTransform(v.obj.Transform)
	ctx.GetRenderer().SetPrimitive(v.pId, pr, v.ShouldRedraw())
	v.redraw = false
}

func (v *ImageView) GetName() string {
	return engine.NameOfComponent(v)
}

func NewImageView(index int) *ImageView {
	return &ImageView{
		ResIndex: index,
	}
}

func NewImageViewWithUrl(url string) *ImageView {
	return &ImageView{
		ResIndex: -1,
		ImageUrl: url,
	}
}
