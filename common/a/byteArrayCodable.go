package a

import (
	"encoding/binary"
	"math"
	"unsafe"
)

type ByteArrayEncodable interface {
	EncodeToByteArray() []byte
}

type ByteArrayDecodable interface {
	DecodeFromByteArray([]byte)
}

type ByteArrayCodable interface {
	ByteArrayEncodable
	ByteArrayDecodable
}

func Int64ToByteArray(num int64) []byte {
	//size := int(unsafe.Sizeof(num))
	//arr := make([]byte, size)
	//for i := 0 ; i < size ; i++ {
	//	byt := *(*uint8)(unsafe.Pointer(uintptr(unsafe.Pointer(&num)) + uintptr(i)))
	//	arr[i] = byt
	//}
	//return arr

	bytes := make([]byte, 8)
	for i := 0; i < 8; i++ {
		bytes[i] = byte(num)
		num >>= 8
	}

	return bytes
}

func ByteArrayToInt64(arr []byte) int64 {
	val := int64(0)
	size := len(arr)
	for i := 0 ; i < size ; i++ {
		*(*uint8)(unsafe.Pointer(uintptr(unsafe.Pointer(&val)) + uintptr(i))) = arr[i]
	}
	return val
}

func IntToByteArray(num int32) []byte {
	//return Int64ToByteArray(int64(num))

	bytes := make([]byte, 4)
	for i := 0; i < 4; i++ {
		bytes[i] = byte(num)
		num >>= 8
	}

	return bytes
}

func ByteArrayToInt(arr []byte) int {
	val := 0
	size := len(arr)
	for i := 0 ; i < size ; i++ {
		*(*uint8)(unsafe.Pointer(uintptr(unsafe.Pointer(&val)) + uintptr(i))) = arr[i]
	}
	return val
}

func Float64ToByteArray(num float64) []byte {
	arr := make([]byte, 8)
	n := math.Float64bits(num)
	binary.LittleEndian.PutUint64(arr, n)
	return arr
}

func Float32ToByteArray(num float32) []byte {
	arr := make([]byte, 4)
	n := math.Float32bits(num)
	binary.LittleEndian.PutUint32(arr, n)
	return arr
}