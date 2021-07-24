package gpu

type Buffer interface {
	GetSize() int
	PutBytes(data []byte)
}
