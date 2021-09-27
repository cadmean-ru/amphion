package rendering

type Bitmap struct {
	Pixels []byte
	Width  int
	Height int
}

func (b *Bitmap) Dispose() {
	b.Pixels = nil
}
