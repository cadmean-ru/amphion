package builtin

import (
	"github.com/cadmean-ru/amphion/common/a"
	"github.com/cadmean-ru/amphion/common/shape"
	"github.com/cadmean-ru/amphion/engine"
	"github.com/cadmean-ru/amphion/rendering"
)

type ClipArea struct {
	engine.ComponentImpl
	Shape shape.Kind `state:"shape"`
}

func (c *ClipArea) OnUpdate(_ engine.UpdateContext) {
    c.SceneObject.GetRenderingNode().SetClipArea2D(rendering.NewClipArea2D(c.Shape, c.SceneObject.Transform.GlobalRect()))
}

func (c *ClipArea) OnStop() {
	c.SceneObject.GetRenderingNode().RemoveClipArea2D()
}

func (c *ClipArea) IsPointInside(point a.Vector3) bool {
	s := shape.New(c.Shape, c.SceneObject.Transform.GlobalRect())
	return s.IsPointInside(point)
}

func (c *ClipArea) IsPointInside2D(point a.Vector3) bool {
	s := shape.New(c.Shape, c.SceneObject.Transform.GlobalRect())
	return s.IsPointInside2D(point)
}

func (c *ClipArea) IsSolid() bool {
	return false
}

func NewClipArea(shape shape.Kind) *ClipArea {
	return &ClipArea{
		Shape: shape,
	}
}