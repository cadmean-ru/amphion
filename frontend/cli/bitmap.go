package cli

import "github.com/cadmean-ru/amphion/rendering"

type Bitmap struct {
	*rendering.Bitmap
}

func (b *Bitmap) GetPixels() []byte {
	return b.Pixels
}

func (b *Bitmap) GetWidth() int {
	return b.Width
}

func (b *Bitmap) GetHeight() int {
	return b.Height
}

func (b *Bitmap) Dispose() {
	b.Bitmap.Dispose()
}

func newBitmap(original *rendering.Bitmap) *Bitmap {
	return &Bitmap{Bitmap: original}
}