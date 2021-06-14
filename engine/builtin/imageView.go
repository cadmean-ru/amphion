package builtin

import (
	"github.com/cadmean-ru/amphion/common/a"
	"github.com/cadmean-ru/amphion/engine"
	"github.com/cadmean-ru/amphion/rendering"
)

// Displays an image given it's resource index
type ImageView struct {
	engine.ViewImpl
	ResId    a.ResId `state:"resId"`
	ImageUrl string  `state:"imageUrl"`
}

func (v *ImageView) OnDraw(ctx engine.DrawingContext) {
	var url string

	if v.ResId > -1 {
		url = v.Engine.GetResourceManager().FullPathOf(v.ResId)
	} else {
		url = v.ImageUrl
	}

	pr := rendering.NewImagePrimitive(url)
	pr.Transform = v.SceneObject.Transform.ToRenderingTransform()
	ctx.GetRenderingNode().SetPrimitive(v.PrimitiveId, pr)
	v.ShouldRedraw = false
}

// Sets the resource index equal to the specified value, forcing the view to redraw and requesting rendering.
func (v *ImageView) SetResId(i a.ResId) {
	v.ResId = i
	v.ImageUrl = ""
	v.ShouldRedraw = true
	engine.RequestRendering()
}

// Sets the image url equal to the specified value, forcing the view to redraw and requesting rendering.
func (v *ImageView) SetImageUrl(url string) {
	v.ResId = -1
	v.ImageUrl = url
	v.ShouldRedraw = true
	engine.RequestRendering()
}

func (v *ImageView) GetName() string {
	return engine.NameOfComponent(v)
}

func NewImageView(index a.ResId) *ImageView {
	return &ImageView{
		ResId: index,
	}
}

func NewImageViewWithUrl(url string) *ImageView {
	return &ImageView{
		ResId:    -1,
		ImageUrl: url,
	}
}
