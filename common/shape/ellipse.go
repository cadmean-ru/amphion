package shape

import (
	"github.com/cadmean-ru/amphion/common"
	"github.com/cadmean-ru/amphion/common/a"
)

type EllipseShape struct {
	rect *common.Rect
}

func (e *EllipseShape) IsPointInside(point a.Vector3) bool {
	return e.rect.IsPointInside(point) && e.IsPointInside2D(point)
}

func (e *EllipseShape) IsPointInside2D(point a.Vector3) bool {
	pos := e.rect.Min()
	a := e.rect.X.GetLength() / 2
	b := e.rect.Y.GetLength() / 2
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

func (e *EllipseShape) Kind() Kind {
	return Ellipse
}

func NewEllipseShape(boundary *common.Rect) *EllipseShape {
	return &EllipseShape{rect: boundary}
}

