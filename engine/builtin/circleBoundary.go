package builtin

import (
	"github.com/cadmean-ru/amphion/common"
	"github.com/cadmean-ru/amphion/engine"
)

type CircleBoundary struct {
	obj *engine.SceneObject
}

func (r *CircleBoundary) GetName() string {
	return "CircleBoundary"
}

func (r *CircleBoundary) OnInit(ctx engine.InitContext) {
	r.obj = ctx.GetSceneObject()
}

func (r *CircleBoundary) OnStart() {

}

func (r *CircleBoundary) OnStop() {

}

func (r *CircleBoundary) IsPointInside(_ common.Vector3) bool {
	return false
}

func (r *CircleBoundary) IsPointInside2D(point common.Vector3) bool {
	rect := r.obj.Transform.GetGlobalRect()
	pos := r.obj.Transform.GetGlobalTopLeftPosition()
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