package builtin

import (
	"github.com/cadmean-ru/amphion/common/a"
	"github.com/cadmean-ru/amphion/engine"
)

type RectBoundary struct {
	engine.ComponentImpl
}

func (r *RectBoundary) IsPointInside(point a.Vector3) bool {
	return r.SceneObject.Transform.GlobalRect().IsPointInside(point)
}

func (r *RectBoundary) IsPointInside2D(point a.Vector3) bool {
	return r.SceneObject.Transform.GlobalRect().IsPointInside2D(point)
}

func (r *RectBoundary) IsSolid() bool {
	return true
}

func NewRectBoundary() *RectBoundary {
	return &RectBoundary{}
}