package builtin

import (
	"github.com/cadmean-ru/amphion/common/a"
	"github.com/cadmean-ru/amphion/common/shape"
	"github.com/cadmean-ru/amphion/engine"
)

type TriangleBoundary struct {
	engine.ComponentImpl
}

func (r *TriangleBoundary) IsPointInside(point a.Vector3) bool {
	triangle := shape.NewTriangleShape(r.SceneObject.Transform.GlobalRect())
	return triangle.IsPointInside(point)
}

func (r *TriangleBoundary) IsPointInside2D(point a.Vector3) bool {
	triangle := shape.NewTriangleShape(r.SceneObject.Transform.GlobalRect())
	return triangle.IsPointInside2D(point)
}

func (r *TriangleBoundary) IsSolid() bool {
	return true
}

func NewTriangleBoundary() *TriangleBoundary {
	return &TriangleBoundary{}
}