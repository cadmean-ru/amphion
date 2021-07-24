package gpu

type MemoryLayout struct {
	value interface{}
	stride uintptr
}

func (m *MemoryLayout) Calculate() {

}

func (m *MemoryLayout) Stride() uint64 {
	return uint64(m.stride)
}

func NewMemoryLayout(value interface{}) *MemoryLayout {
	return &MemoryLayout{
		value: value,
	}
}
