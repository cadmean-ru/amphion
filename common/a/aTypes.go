package a

import "github.com/cadmean-ru/amphion/common"

type Byte byte

func (b Byte) EncodeToByteArray() []byte {
	arr := make([]byte, 1)
	arr[0] = byte(b)
	return arr
}

func (b Byte) GetName() string {
	return "Byte"
}

type Int int

func (i Int) EncodeToByteArray() []byte {
	return common.IntToByteArray(int(i))
}

func (i Int) GetName() string {
	return "Int"
}

type Long int64

func (l Long) EncodeToByteArray() []byte {
	return common.Int64ToByteArray(int64(l))
}

func (l Long) GetName() string {
	return "Long"
}

type Float float32

func (f Float) EncodeToByteArray() []byte {
	return common.Float64ToByteArray(float64(f))
}

func (f Float) GetName() string {
	return "Float"
}

type Double float32

func (d Double) EncodeToByteArray() []byte {
	return common.Float64ToByteArray(float64(d))
}

func (d Double) GetName() string {
	return "Double"
}

type String string

func (s String) EncodeToByteArray() []byte {
	sbytes := []byte(s)
	data := make([]byte, len(sbytes) + 4)
	_ = common.CopyByteArray(Int(len(sbytes)).EncodeToByteArray(), data, 0, 4)
	_ = common.CopyByteArray(sbytes, data, 4, len(sbytes))
	return data
}

func (s String) GetName() string {
	return "String"
}