package shape

import (
	"github.com/cadmean-ru/amphion/common"
	"github.com/cadmean-ru/amphion/common/a"
	"math"
)

type TriangleShape struct {
	rect *common.RectBoundary
}

func (t *TriangleShape) IsPointInside(point a.Vector3) bool {
	return t.rect.IsPointInside(point) && t.IsPointInside2D(point)
}

func (t *TriangleShape) IsPointInside2D(point a.Vector3) bool {
	pos := t.rect.Min()
	size := t.rect.Size()
	a1 := t.rect.X.GetLength()
	h := t.rect.Y.GetLength()
	s := a1 * h / 2
	pointA := a.NewVector3(pos.X, pos.Y + size.Y, 0)
	pointB := a.NewVector3(pos.X + size.X / 2, pos.Y, 0)
	pointC := a.NewVector3(pos.X + size.X, pos.Y + size.Y, 0)
	s1 := t.area(point, pointA, pointB)
	s2 := t.area(point, pointB, pointC)
	s3 := t.area(point, pointA, pointC)
	return math.Abs(float64(s-s1-s2-s3)) <= 0.0001
}

func (t *TriangleShape) Kind() Kind {
	return Triangle
}

func (t *TriangleShape) area(a, b, c a.Vector3) float32 {
	ab := t.length(a, b)
	bc := t.length(b, c)
	ac := t.length(a, c)
	p := (ab + bc + ac) / 2
	return float32(math.Sqrt(float64(p * (p - ab) * (p - bc) * (p - ac))))
}

func (t *TriangleShape) length(a, b a.Vector3) float32 {
	x := math.Abs(float64(b.X - a.X))
	y := math.Abs(float64(b.Y - a.Y))
	return float32(math.Sqrt(x*x + y*y))
}

func NewTriangleShape(boundary *common.RectBoundary) *TriangleShape {
	return &TriangleShape{rect: boundary}
}

