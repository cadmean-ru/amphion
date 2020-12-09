package common

import (
	"fmt"
	"math"
)

// Represents a point in 3D space
type Vector3 struct {
	X, Y, Z float64
}

func (v Vector3) SetXYZ(x, y, z float64) {
	v.X = x
	v.Y = y
	v.Z = z
}

func NewVector3(x, y, z float64) Vector3 {
	return Vector3{
		X: x,
		Y: y,
		Z: z,
	}
}

func (v Vector3) ToMap() SiMap {
	return map[string]interface{}{
		"x": v.X,
		"y": v.Y,
		"z": v.Z,
	}
}

func (v Vector3) ToString() string {
	return fmt.Sprintf("(%f, %f, %f)", v.X, v.Y, v.Z)
}

// Returns a new vector - the sum of two vectors
func (v Vector3) Add(v2 Vector3) Vector3 {
	return Vector3{
		X: v.X + v2.X,
		Y: v.Y + v2.Y,
		Z: v.Z + v2.Z,
	}
}

// Returns a new vector - v-v2
func (v Vector3) Sub(v2 Vector3) Vector3 {
	return Vector3{
		X: v.X - v2.X,
		Y: v.Y - v2.Y,
		Z: v.Z - v2.Z,
	}
}

// Returns new vector - multiplication of two vectors
func (v Vector3) Multiply(v2 Vector3) Vector3 {
	return Vector3{
		X: v.X * v2.X,
		Y: v.Y * v2.Y,
		Z: v.Z * v2.Z,
	}
}

// Rounds the valued of the vector
func (v Vector3) Round() IntVector3 {
	return IntVector3{
		X: int(math.Round(v.X)),
		Y: int(math.Round(v.Y)),
		Z: int(math.Round(v.Z)),
	}
}

// Checks if the vector is the same as other vector
func (v Vector3) Equals(other interface{}) bool {
	switch other.(type) {
	case Vector3:
		v2 := other.(Vector3)
		return v.X == v2.X && v.Y == v2.Y && v.Z == v2.Z
	default:
		return false
	}
}

func (v Vector3) EncodeToByteArray() []byte {
	arr := make([]byte, 24)
	_ = CopyByteArray(Float64ToByteArray(v.X), arr, 0, 8)
	_ = CopyByteArray(Float64ToByteArray(v.Y), arr, 8, 8)
	_ = CopyByteArray(Float64ToByteArray(v.Z), arr, 16, 8)
	return arr
}

func ZeroVector() Vector3 {
	return Vector3{}
}

func OneVector() Vector3 {
	return Vector3{
		X: 1,
		Y: 1,
		Z: 1,
	}
}
