package builtin

import (
	"github.com/cadmean-ru/amphion/common/shape"
	"github.com/cadmean-ru/amphion/engine"
	"github.com/cadmean-ru/amphion/rendering"
)

type ClipArea struct {
	engine.ComponentImpl
	Shape shape.Kind `state:"shape"`
}

func (c *ClipArea) OnUpdate(_ engine.UpdateContext) {
    c.SceneObject.GetRenderingNode().SetClipArea2D(rendering.NewClipArea2D(c.Shape, c.SceneObject.Transform.GetGlobalRect()))
}

func (c *ClipArea) OnStop() {
	c.SceneObject.GetRenderingNode().RemoveClipArea2D()
}

func NewClipArea(shape shape.Kind) *ClipArea {
	return &ClipArea{
		Shape: shape,
	}
}