package cli

import "github.com/cadmean-ru/amphion/common/atext"

type GeometryPrimitiveData struct {
	GeometryType int
	TlPosition   *Vector3
	BrPosition   *Vector3
	FillColor    *Vector4
	StrokeColor  *Vector4
	StrokeWeight int
	CornerRadius int
}

type ImagePrimitiveData struct {
	TlPosition  *Vector3
	BrPosition  *Vector3
	Bitmaps     []*Bitmap
	Index       int
}

func (i *ImagePrimitiveData) BitmapAt(index int) *Bitmap {
	return i.Bitmaps[index]
}

func (i *ImagePrimitiveData) GetBitmapCount() int {
	return len(i.Bitmaps)
}

type TextPrimitiveData struct {
	Text       string
	TlPosition *Vector3
	Size       *Vector3
	TextColor  *Vector4
	Provider   atext.Provider
}
