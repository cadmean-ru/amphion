package rendering

import (
	"github.com/cadmean-ru/amphion/common"
	"github.com/cadmean-ru/amphion/common/a"
)

const (
	PrimitiveEmpty     = 0
	PrimitivePoint     = 1
	PrimitiveLine      = 2
	PrimitiveRectangle = 3
	PrimitiveEllipse   = 4
	PrimitiveTriangle  = 5
	PrimitiveText      = 6
	PrimitiveImage     = 7
	PrimitiveBezier    = 8
)

type Primitive struct {
	Type       a.Byte
	Attributes []Attribute
}

func (p *Primitive) AddAttribute(attr Attribute) {
	p.Attributes = append(p.Attributes, attr)
}

func (p *Primitive) EncodeToByteArray() []byte {
	length := 1
	for _, attr := range p.Attributes {
		length += attr.GetLength()
	}
	data := make([]byte, length)
	data[0] = byte(p.Type)
	counter := 1
	for _, attr := range p.Attributes {
		_ = common.CopyByteArray(attr.EncodeToByteArray(), data, counter, attr.GetLength())
		counter += attr.GetLength()
	}
	return data
}

func NewPrimitive(pType a.Byte) *Primitive {
	return &Primitive{
		Type:       pType,
		Attributes: make([]Attribute, 0, 1),
	}
}

type PrimitiveBuilder interface {
	BuildPrimitive() *Primitive
}