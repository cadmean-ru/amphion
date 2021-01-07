package a

// Deprecated: use corresponding go type instead
type Byte byte

func (b Byte) EncodeToByteArray() []byte {
	arr := make([]byte, 1)
	arr[0] = byte(b)
	return arr
}

func (b Byte) GetName() string {
	return "Byte"
}

// Deprecated: use corresponding go type instead
type Int int32

func (i Int) EncodeToByteArray() []byte {
	return IntToByteArray(int32(i))
}

func (i Int) GetName() string {
	return "Int"
}

// Deprecated: use corresponding go type instead
type Long int64

func (l Long) EncodeToByteArray() []byte {
	return Int64ToByteArray(int64(l))
}

func (l Long) GetName() string {
	return "Long"
}

// Deprecated: use corresponding go type instead
type Float float32

func (f Float) EncodeToByteArray() []byte {
	return Float64ToByteArray(float64(f))
}

func (f Float) GetName() string {
	return "Float"
}

// Deprecated: use corresponding go type instead
type Double float32

func (d Double) EncodeToByteArray() []byte {
	return Float64ToByteArray(float64(d))
}

func (d Double) GetName() string {
	return "Double"
}

// Deprecated: use corresponding go type instead
type String string

func (s String) EncodeToByteArray() []byte {
	sbytes := []byte(s)
	data := make([]byte, len(sbytes) + 4)
	_ = CopyByteArray(Int(len(sbytes)).EncodeToByteArray(), data, 0, 4)
	_ = CopyByteArray(sbytes, data, 4, len(sbytes))
	return data
}

func (s String) GetName() string {
	return "String"
}