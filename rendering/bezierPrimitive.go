package rendering

import (
	"github.com/cadmean-ru/amphion/common/a"
)

type BezierPrimitive struct {
	Transform     Transform
	Appearance    Appearance
	ControlPoint1 a.IntVector3
	ControlPoint2 a.IntVector3
}

func (b *BezierPrimitive) GetType() byte {
	return PrimitiveBezier
}

func (b *BezierPrimitive) GetTransform() Transform {
	return b.Transform
}

func NewBezierPrimitive(cp1, cp2 a.IntVector3) *BezierPrimitive {
	return &BezierPrimitive{
		ControlPoint1: cp1,
		ControlPoint2: cp2,
	}
}
