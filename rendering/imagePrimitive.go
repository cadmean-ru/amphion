package rendering

import "github.com/cadmean-ru/amphion/common"

const imagePrimitiveBytesSize = primitiveBytesSize + 4

type ImagePrimitive struct {
	Transform Transform
	ResIndex  common.AInt
}

func (p *ImagePrimitive) GetType() common.AByte {
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

func NewImagePrimitive(index common.AInt) *ImagePrimitive {
	return &ImagePrimitive{
		Transform: NewTransform(),
		ResIndex:  index,
	}
}
