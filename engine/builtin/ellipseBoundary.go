package builtin

import (
	"github.com/cadmean-ru/amphion/common/a"
	"github.com/cadmean-ru/amphion/common/shape"
	"github.com/cadmean-ru/amphion/engine"
)

type EllipseBoundary struct {
	engine.ComponentImpl
}

func (s *EllipseBoundary) IsPointInside(point a.Vector3) bool {
	e := shape.NewEllipseShape(s.SceneObject.Transform.GlobalRect())
	return e.IsPointInside(point)
}

func (s *EllipseBoundary) IsPointInside2D(point a.Vector3) bool {
	e := shape.NewEllipseShape(s.SceneObject.Transform.GlobalRect())
	return e.IsPointInside2D(point)
}

func (s *EllipseBoundary) IsSolid() bool {
	return true
}

func NewCircleBoundary() *EllipseBoundary {
	return &EllipseBoundary{}
}