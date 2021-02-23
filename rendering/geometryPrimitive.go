package rendering

type GeometryPrimitive struct {
	Transform     Transform
	Appearance    Appearance
	primitiveType byte
}

func (p *GeometryPrimitive) GetType() byte {
	return p.primitiveType
}

func (p *GeometryPrimitive) GetTransform() Transform {
	return p.Transform
}

func (p *GeometryPrimitive) SetTransform(t Transform) {
	p.Transform = t
}

func NewGeometryPrimitive(pType byte) *GeometryPrimitive {
	if pType < 0 || pType > 5 {
		pType = 0
	}

	return &GeometryPrimitive{
		primitiveType: pType,
		Transform:     NewTransform(),
		Appearance:    DefaultAppearance(),
	}
}

