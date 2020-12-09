package builtin

import (
	"github.com/cadmean-ru/amphion/common"
	"github.com/cadmean-ru/amphion/engine"
)

type RectBoundary struct {
	obj *engine.SceneObject
}

func (r *RectBoundary) GetName() string {
	return "RectBoundary"
}

func (r *RectBoundary) OnInit(ctx engine.InitContext) {
	r.obj = ctx.GetSceneObject()
}

func (r *RectBoundary) OnStart() {

}

func (r *RectBoundary) OnStop() {

}

func (r *RectBoundary) IsPointInside(point common.Vector3) bool {
	return r.obj.Transform.GetGlobalRect().IsPointInside(point)
}

func (r *RectBoundary) IsPointInside2D(point common.Vector3) bool {
	return r.obj.Transform.GetGlobalRect().IsPointInside2D(point)
}

func NewRectBoundary() *RectBoundary {
	return &RectBoundary{}
}