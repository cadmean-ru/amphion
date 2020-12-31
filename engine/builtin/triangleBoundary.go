package builtin

import (
	"github.com/cadmean-ru/amphion/common/a"
	"github.com/cadmean-ru/amphion/engine"
	"math"
)

type TriangleBoundary struct {
	obj *engine.SceneObject
}

func (r *TriangleBoundary) GetName() string {
	return engine.NameOfComponent(r)
}

func (r *TriangleBoundary) OnInit(ctx engine.InitContext) {
	r.obj = ctx.GetSceneObject()
}

func (r *TriangleBoundary) OnStart() {

}

func (r *TriangleBoundary) OnStop() {

}

func (r *TriangleBoundary) IsPointInside(_ a.Vector3) bool {
	return false
}

func (r *TriangleBoundary) IsPointInside2D(point a.Vector3) bool {
	rect := r.obj.Transform.GetGlobalRect()
	pos := r.obj.Transform.GetGlobalTopLeftPosition()
	size := r.obj.Transform.Size
	a1 := rect.X.GetLength()
	h := rect.Y.GetLength()
	s := a1 * h / 2
	pointA := a.NewVector3(pos.X, pos.Y + size.Y, 0)
	pointB := a.NewVector3(pos.X + size.X / 2, pos.Y, 0)
	pointC := a.NewVector3(pos.X + size.X, pos.Y + size.Y, 0)
	s1 := area(point, pointA, pointB)
	s2 := area(point, pointB, pointC)
	s3 := area(point, pointA, pointC)
	return math.Abs(float64(s-s1-s2-s3)) <= 0.0001
}

func area(a, b, c a.Vector3) float32 {
	ab := length(a, b)
	bc := length(b, c)
	ac := length(a, c)
	p := (ab + bc + ac) / 2
	return float32(math.Sqrt(float64(p * (p - ab) * (p - bc) * (p - ac))))
}

func length(a, b a.Vector3) float32 {
	x := math.Abs(float64(b.X - a.X))
	y := math.Abs(float64(b.Y - a.Y))
	return float32(math.Sqrt(x*x + y*y))
}

func NewTriangleBoundary() *TriangleBoundary {
	return &TriangleBoundary{}
}