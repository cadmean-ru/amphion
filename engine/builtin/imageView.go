package builtin

import (
	"github.com/cadmean-ru/amphion/engine"
	"github.com/cadmean-ru/amphion/rendering"
)

// Displays an image given it's resource index
type ImageView struct {
	engine.ViewImpl
	ResIndex int    `state:"resIndex"`
	ImageUrl string `state:"imageUrl"`
}

func (v *ImageView) OnDraw(ctx engine.DrawingContext) {
	var url string

	if v.ResIndex > -1 {
		url = v.Engine.GetResourceManager().FullPathOf(v.ResIndex)
	} else {
		url = v.ImageUrl
	}

	pr := rendering.NewImagePrimitive(url)
	pr.Transform = v.SceneObject.Transform.ToRenderingTransform()
	ctx.GetRenderer().SetPrimitive(v.PrimitiveId, pr, v.ShouldRedraw())
	v.Redraw = false
}

// Sets the resource index equal to the specified value, forcing the view to redraw and requesting rendering.
func (v *ImageView) SetResIndex(i int) {
	v.ResIndex = i
	v.ImageUrl = ""
	v.Redraw = true
	engine.RequestRendering()
}

// Sets the image url equal to the specified value, forcing the view to redraw and requesting rendering.
func (v *ImageView) SetImageUrl(url string) {
	v.ResIndex = -1
	v.ImageUrl = url
	v.Redraw = true
	engine.RequestRendering()
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
