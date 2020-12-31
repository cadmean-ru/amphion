package a

import (
	"fmt"
	"github.com/cadmean-ru/amphion/common/require"
	"math"
)

// Represents a point in 3D space
type Vector3 struct {
	X, Y, Z float32
}

func (v Vector3) SetXYZ(x, y, z float32) {
	v.X = x
	v.Y = y
	v.Z = z
}

func NewVector3(x, y, z float32) Vector3 {
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

func (v *Vector3) FromMap(siMap SiMap) {
	v.X = require.Float32(siMap["x"])
	v.Y = require.Float32(siMap["y"])
	v.Z = require.Float32(siMap["z"])
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
		X: int(math.Round(float64(v.X))),
		Y: int(math.Round(float64(v.Y))),
		Z: int(math.Round(float64(v.Z))),
	}
}

// Transforms vector ro normalized device coordinates vector
func (v Vector3) Ndc(screen Vector3) Vector3 {
	xs := screen.X
	ys := screen.Y
	x0 := screen.X / 2
	y0 := screen.Y / 2
	x := v.X
	y := v.Y
	newX := (2*(x-x0))/xs
	newY := (-2*(y-y0))/ys

	return Vector3{newX, newY, 0}
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
	arr := make([]byte, 12)
	_ = CopyByteArray(Float32ToByteArray(v.X), arr, 0, 4)
	_ = CopyByteArray(Float32ToByteArray(v.Y), arr, 4, 4)
	_ = CopyByteArray(Float32ToByteArray(v.Z), arr, 8, 4)
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

func NewVector3FromMap(siMap SiMap) Vector3 {
	var v Vector3
	v.FromMap(siMap)
	return v
}