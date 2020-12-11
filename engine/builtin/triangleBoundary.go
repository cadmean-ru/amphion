package builtin

import (
	"github.com/cadmean-ru/amphion/common"
	"github.com/cadmean-ru/amphion/engine"
	"math"
)

type TriangleBoundary struct {
	obj *engine.SceneObject
}

func (r *TriangleBoundary) GetName() string {
	return "TriangleBoundary"
}

func (r *TriangleBoundary) OnInit(ctx engine.InitContext) {
	r.obj = ctx.GetSceneObject()
}

func (r *TriangleBoundary) OnStart() {

}

func (r *TriangleBoundary) OnStop() {

}

func (r *TriangleBoundary) IsPointInside(_ common.Vector3) bool {
	return false
}

func (r *TriangleBoundary) IsPointInside2D(point common.Vector3) bool {
	rect := r.obj.Transform.GetGlobalRect()
	pos := r.obj.Transform.GetGlobalTopLeftPosition()
	size := r.obj.Transform.Size
	a := rect.X.GetLength()
	h := rect.Y.GetLength()
	s := a * h / 2
	pointA := common.NewVector3(pos.X, pos.Y + size.Y, 0)
	pointB := common.NewVector3(pos.X + size.X / 2, pos.Y, 0)
	pointC := common.NewVector3(pos.X + size.X, pos.Y + size.Y, 0)
	s1 := area(point, pointA, pointB)
	s2 := area(point, pointB, pointC)
	s3 := area(point, pointA, pointC)
	return math.Abs(s - s1 - s2 - s3) <= 0.0001
}

func area(a, b, c common.Vector3) float64 {
	ab := length(a, b)
	bc := length(b, c)
	ac := length(a, c)
	p := (ab + bc + ac) / 2
	return math.Sqrt(p * (p - ab) * (p - bc) * (p - ac))
}

func length(a, b common.Vector3) float64 {
	x := math.Abs(b.X - a.X)
	y := math.Abs(b.Y - a.Y)
	return math.Sqrt(x * x + y * y)
}

func NewTriangleBoundary() *TriangleBoundary {
	return &TriangleBoundary{}
}