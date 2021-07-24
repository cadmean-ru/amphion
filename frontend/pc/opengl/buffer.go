package opengl

type Buffer struct {
	id   int
	size int
}

func (b *Buffer) GetSize() int {
	return b.size
}

func (b *Buffer) PutBytes(data []byte) {

}

