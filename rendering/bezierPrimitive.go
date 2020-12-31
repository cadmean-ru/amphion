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

func (b *BezierPrimitive) GetType() a.Byte {
	return PrimitiveBezier
}

func (b *BezierPrimitive) GetTransform() Transform {
	return b.Transform
}

func (b *BezierPrimitive) BuildPrimitive() *Primitive {
	pr := NewPrimitive(PrimitiveBezier)
	pr.AddAttribute(NewAttribute(AttributeTransform, b.Transform))
	pr.AddAttribute(NewAttribute(AttributeFillColor, b.Appearance.FillColor))
	pr.AddAttribute(NewAttribute(AttributeStrokeColor, b.Appearance.StrokeColor))
	pr.AddAttribute(NewAttribute(AttributeStrokeWeight, b.Appearance.StrokeWeight))
	pr.AddAttribute(NewAttribute(AttributePoint, b.Transform.Position.Add(b.ControlPoint1)))
	pr.AddAttribute(NewAttribute(AttributePoint, b.Transform.Position.Add(b.ControlPoint2)))
	return pr
}

func NewBezierPrimitive(cp1, cp2 a.IntVector3) *BezierPrimitive {
	return &BezierPrimitive{
		ControlPoint1: cp1,
		ControlPoint2: cp2,
	}
}
