package builtin

import (
	"github.com/cadmean-ru/amphion/common/a"
	"github.com/cadmean-ru/amphion/engine"
)

type CircleBoundary struct {
	engine.ComponentImpl
}

func (s *CircleBoundary) IsPointInside(_ a.Vector3) bool {
	return false
}

func (s *CircleBoundary) IsPointInside2D(point a.Vector3) bool {
	rect := s.SceneObject.Transform.GlobalRect()
	pos := s.SceneObject.Transform.GlobalTopLeftPosition()
	a := rect.X.GetLength() / 2
	b := rect.Y.GetLength() / 2
	x := point.X
	y := point.Y
	xc := pos.X + a
	yc := pos.Y + b
	x2 := x - xc
	y2 := y - yc
	c := x2 / a
	d := y2 / b
	res := c * c + d * d
	return res <= 1
}

func NewCircleBoundary() *CircleBoundary {
	return &CircleBoundary{}
}