package rendering

import (
	"github.com/cadmean-ru/amphion/common"
)

const primitiveBytesSize = 1 + transformBytesSize + appearanceBytesSize

type GeometryPrimitive struct {
	Transform     Transform
	Appearance    Appearance
	primitiveType common.AByte
}

func (p *GeometryPrimitive) GetType() common.AByte {
	return p.primitiveType
}

func (p *GeometryPrimitive) GetTransform() Transform {
	return p.Transform
}

func (p *GeometryPrimitive) BuildPrimitive() *Primitive {
	pr := NewPrimitive(p.primitiveType)
	pr.AddAttribute(NewAttribute(AttributeTransform, p.Transform))
	pr.AddAttribute(NewAttribute(AttributeFillColor, p.Appearance.FillColor))
	pr.AddAttribute(NewAttribute(AttributeStrokeColor, p.Appearance.StrokeColor))
	pr.AddAttribute(NewAttribute(AttributeStrokeWeight, p.Appearance.StrokeWeight))
	if p.Appearance.CornerRadius > 0 {
		pr.AddAttribute(NewAttribute(AttributeCornerRadius, p.Appearance.CornerRadius))
	}
	return pr
}

func NewGeometryPrimitive(pType common.AByte) *GeometryPrimitive {
	if pType < 0 || pType > 5 {
		pType = 0
	}

	return &GeometryPrimitive{
		primitiveType: pType,
		Transform:     NewTransform(),
		Appearance:    DefaultAppearance(),
	}
}

