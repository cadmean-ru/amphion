package rendering

const imagePrimitiveBytesSize = primitiveBytesSize + 4

type ImagePrimitive struct {
	Transform Transform
	ImageUrl  string
}

func (p *ImagePrimitive) GetType() byte {
	return PrimitiveImage
}

func (p *ImagePrimitive) GetTransform() Transform {
	return p.Transform
}

func NewImagePrimitive(url string) *ImagePrimitive {
	return &ImagePrimitive{
		Transform: NewTransform(),
		ImageUrl:  url,
	}
}
