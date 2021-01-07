package rendering

const imagePrimitiveBytesSize = primitiveBytesSize + 4

type ImagePrimitive struct {
	Transform Transform
	ResIndex  int
}

func (p *ImagePrimitive) GetType() byte {
	return PrimitiveImage
}

func (p *ImagePrimitive) GetTransform() Transform {
	return p.Transform
}

func NewImagePrimitive(index int) *ImagePrimitive {
	return &ImagePrimitive{
		Transform: NewTransform(),
		ResIndex:  index,
	}
}
