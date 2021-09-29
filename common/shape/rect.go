package shape

import (
	"github.com/cadmean-ru/amphion/common"
	"github.com/cadmean-ru/amphion/common/a"
)

type RectangleShape struct {
	rect *common.RectBoundary
}

func (r *RectangleShape) IsPointInside(point a.Vector3) bool {
	return r.rect.IsPointInside(point)
}

func (r *RectangleShape) IsPointInside2D(point a.Vector3) bool {
	return r.rect.IsPointInside2D(point)
}

func (r *RectangleShape) Kind() Kind {
	return Rectangle
}

func NewRectangleShape(boundary *common.RectBoundary) *RectangleShape {
	return &RectangleShape{rect: boundary}
}
