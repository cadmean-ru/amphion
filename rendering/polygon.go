package rendering

import "github.com/cadmean-ru/amphion/common/a"

type PolygonPrimitive struct {
	Transform  Transform
	Vertices   []a.Vector3
	Indexes    []uint
	Appearance Appearance
}

func (p *PolygonPrimitive) GetType() byte {
	return PrimitivePolygon
}

func (p *PolygonPrimitive) GetTransform() Transform {
	return p.Transform
}

func (p *PolygonPrimitive) SetTransform(t Transform) {
	p.Transform = t
}

func NewPolygonPrimitive() *PolygonPrimitive {
	return &PolygonPrimitive{}
}
