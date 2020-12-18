package common

import "fmt"

// Represents a point in 3D space only in integer values
type IntVector3 struct {
	X, Y, Z int
}

func (v IntVector3) SetXYZ(x, y, z int) {
	v.X = x
	v.Y = y
	v.Z = z
}

func NewIntVector3(x, y, z int) IntVector3 {
	return IntVector3{
		X: x,
		Y: y,
		Z: z,
	}
}

func (v IntVector3) ToMap() SiMap {
	return map[string]interface{}{
		"x": v.X,
		"y": v.Y,
		"z": v.Z,
	}
}

func (v IntVector3) ToString() string {
	return fmt.Sprintf("(%d, %d, %d)", v.X, v.Y, v.Z)
}

// Returns a new vector - the sum of two vectors
func (v IntVector3) Add(v2 IntVector3) IntVector3 {
	return IntVector3{
		X: v.X + v2.X,
		Y: v.Y + v2.Y,
		Z: v.Z + v2.Z,
	}
}

// Returns a new vector - v-v2
func (v IntVector3) Sub(v2 IntVector3) IntVector3 {
	return IntVector3{
		X: v.X - v2.X,
		Y: v.Y - v2.Y,
		Z: v.Z - v2.Z,
	}
}

// Returns new vector - multiplication of two vectors
func (v IntVector3) Multiply(v2 IntVector3) IntVector3 {
	return IntVector3{
		X: v.X * v2.X,
		Y: v.Y * v2.Y,
		Z: v.Z * v2.Z,
	}
}

func (v IntVector3) ToFloat() Vector3 {
	return NewVector3(float32(v.X), float32(v.Y), float32(v.Z))
}

// Transforms vector ro normalized device coordinates vector
func (v IntVector3) Ndc(screen IntVector3) Vector3 {
	xs := float32(screen.X)
	ys := float32(screen.Y)
	x0 := float32(screen.X) / 2
	y0 := float32(screen.Y) / 2
	x := float32(v.X)
	y := float32(v.Y)
	newX := (2*(x-x0))/xs
	newY := (-2*(y-y0))/ys

	return Vector3{newX, newY, 0}
}

// Checks if the vector is the same as other vector
func (v IntVector3) Equals(other interface{}) bool {
	switch other.(type) {
	case IntVector3:
		v2 := other.(IntVector3)
		return v.X == v2.X && v.Y == v2.Y && v.Z == v2.Z
	default:
		return false
	}
}

func (v IntVector3) EncodeToByteArray() []byte {
	arr := make([]byte, 12)
	_ = CopyByteArray(IntToByteArray(v.X), arr, 0, 4)
	_ = CopyByteArray(IntToByteArray(v.Y), arr, 4, 4)
	_ = CopyByteArray(IntToByteArray(v.Z), arr, 8, 4)
	return arr
}

func ZeroIntVector() IntVector3 {
	return IntVector3{}
}

func OneIntVector() IntVector3 {
	return IntVector3{
		X: 1,
		Y: 1,
		Z: 1,
	}
}
