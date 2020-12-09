package common

type AByte byte

func (b AByte) EncodeToByteArray() []byte {
	arr := make([]byte, 1)
	arr[0] = byte(b)
	return arr
}

func (b AByte) GetName() string {
	return "AByte"
}

type AInt int

func (i AInt) EncodeToByteArray() []byte {
	return IntToByteArray(int(i))
}

func (i AInt) GetName() string {
	return "AInt"
}

type ALong int64

func (l ALong) EncodeToByteArray() []byte {
	return Int64ToByteArray(int64(l))
}

func (l ALong) GetName() string {
	return "ALong"
}

type AFloat float32

func (f AFloat) EncodeToByteArray() []byte {
	return Float64ToByteArray(float64(f))
}

func (f AFloat) GetName() string {
	return "AFloat"
}

type ADouble float32

func (d ADouble) EncodeToByteArray() []byte {
	return Float64ToByteArray(float64(d))
}

func (d ADouble) GetName() string {
	return "ADouble"
}

type AString string

func (s AString) EncodeToByteArray() []byte {
	sbytes := []byte(s)
	data := make([]byte, len(sbytes) + 4)
	_ = CopyByteArray(AInt(len(sbytes)).EncodeToByteArray(), data, 0, 4)
	_ = CopyByteArray(sbytes, data, 4, len(sbytes))
	return data
}

func (s AString) GetName() string {
	return "AString"
}