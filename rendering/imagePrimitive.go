package rendering

import "github.com/cadmean-ru/amphion/common"

const imagePrimitiveBytesSize = primitiveBytesSize + 4

type ImagePrimitive struct {
	Transform Transform
	resIndex  common.AInt
}

func (p *ImagePrimitive) BuildPrimitive() *Primitive {
	pr := NewPrimitive(PrimitiveImage)
	pr.AddAttribute(NewAttribute(AttributeTransform, p.Transform))
	pr.AddAttribute(NewAttribute(AttributeResIndex, p.resIndex))
	return pr
}

func NewImagePrimitive(index common.AInt) *ImagePrimitive {
	return &ImagePrimitive{
		Transform: NewTransform(),
		resIndex:  index,
	}
}
