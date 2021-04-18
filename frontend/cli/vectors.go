package cli

import "github.com/cadmean-ru/amphion/common/a"

type Vector2 struct {
	X, Y float32
}

func NewVector2(x, y float32) *Vector2 {
	return &Vector2{
		X: x,
		Y: y,
	}
}

type Vector3 struct {
	X, Y, Z float32
}

func NewVector3(x, y, z float32) *Vector3 {
	return &Vector3{
		X: x,
		Y: y,
		Z: z,
	}
}

func NewVector3FromAVector3(vector3 a.Vector3) *Vector3 {
	return NewVector3(vector3.X, vector3.Y, vector3.Z)
}

type Vector4 struct {
	X, Y, Z, W float32
}

func NewVector4(x, y, z, w float32) *Vector4 {
	return &Vector4{
		X: x,
		Y: y,
		Z: z,
		W: w,
	}
}

func NewVector4FromAVector4(vector4 a.Vector4) *Vector4 {
	return NewVector4(vector4.X, vector4.Y, vector4.Z, vector4.W)
}

func Vector3Ndc(v *Vector3, size *Vector3) *Vector3 {
	ndc := a.NewVector3(v.X, v.Y, v.Z).Ndc(a.NewVector3(size.X, size.Y, size.Z))
	return NewVector3(ndc.X, ndc.Y, ndc.Z)
}
