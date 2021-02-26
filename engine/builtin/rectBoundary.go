package builtin

import (
	"github.com/cadmean-ru/amphion/common/a"
	"github.com/cadmean-ru/amphion/engine"
)

type RectBoundary struct {
	engine.ComponentImpl
}

func (r *RectBoundary) GetName() string {
	return engine.NameOfComponent(r)
}

func (r *RectBoundary) IsPointInside(point a.Vector3) bool {
	return r.SceneObject.Transform.GetGlobalRect().IsPointInside(point)
}

func (r *RectBoundary) IsPointInside2D(point a.Vector3) bool {
	return r.SceneObject.Transform.GetGlobalRect().IsPointInside2D(point)
}

func NewRectBoundary() *RectBoundary {
	return &RectBoundary{}
}