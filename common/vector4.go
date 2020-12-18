package common

type Vector4 struct {
	X, Y, Z, W float32
}

func NewVector4(x, y, z, w float32) Vector4 {
	return Vector4{
		X: x,
		Y: y,
		Z: z,
		W: w,
	}
}