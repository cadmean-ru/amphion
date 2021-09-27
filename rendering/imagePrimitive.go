package rendering

type ImagePrimitive struct {
	Transform Transform
	Bitmaps   []*Bitmap
	Index     int
}

func (p *ImagePrimitive) GetType() byte {
	return PrimitiveImage
}

func (p *ImagePrimitive) GetTransform() Transform {
	return p.Transform
}

func (p *ImagePrimitive) SetTransform(t Transform) {
	p.Transform = t
}

func NewImagePrimitive(bitmaps []*Bitmap, index int) *ImagePrimitive {
	return &ImagePrimitive{
		Transform: NewTransform(),
		Bitmaps:   bitmaps,
		Index:     index,
	}
}
