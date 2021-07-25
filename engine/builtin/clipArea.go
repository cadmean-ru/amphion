package builtin

import (
	"github.com/cadmean-ru/amphion/common/shape"
	"github.com/cadmean-ru/amphion/engine"
	"github.com/cadmean-ru/amphion/rendering"
)

type ClipArea struct {
	engine.ComponentImpl
	shape shape.Kind
}

func (c *ClipArea) OnUpdate(ctx engine.UpdateContext) {
    c.SceneObject.GetRenderingNode().SetClipArea2D(rendering.NewClipArea2D(c.shape, c.SceneObject.Transform.GetGlobalRect()))
}

func (c *ClipArea) OnStop() {
	c.SceneObject.GetRenderingNode().RemoveClipArea2D()
}

func (c *ClipArea) GetName() string {
	return engine.NameOfComponent(c)
}

func NewClipArea(shape shape.Kind) *ClipArea {
	return &ClipArea{
		shape: shape,
	}
}