package rendering

import (
	"github.com/cadmean-ru/amphion/common/a"
)

const imagePrimitiveBytesSize = primitiveBytesSize + 4

type ImagePrimitive struct {
	Transform Transform
	ResIndex  a.Int
}

func (p *ImagePrimitive) GetType() a.Byte {
	return PrimitiveImage
}

func (p *ImagePrimitive) GetTransform() Transform {
	return p.Transform
}

func (p *ImagePrimitive) BuildPrimitive() *Primitive {
	pr := NewPrimitive(PrimitiveImage)
	pr.AddAttribute(NewAttribute(AttributeTransform, p.Transform))
	pr.AddAttribute(NewAttribute(AttributeResIndex, p.ResIndex))
	return pr
}

func NewImagePrimitive(index a.Int) *ImagePrimitive {
	return &ImagePrimitive{
		Transform: NewTransform(),
		ResIndex:  index,
	}
}
