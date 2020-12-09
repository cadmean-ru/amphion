package common

import "fmt"

// Represents the boundaries of an object, like collider in unity
type Boundary interface {
	IsPointInside(point Vector3) bool
	IsPointInside2D(point Vector3) bool
}

// Represents a boundary in 3D space
type RectBoundary struct {
	X, Y, Z FloatRange
}

// Checks if specific point is inside the boundary
func (b RectBoundary) IsPointInside(v Vector3) bool {
	return b.X.IsValueInside(v.X) && b.Y.IsValueInside(v.Y) && b.Z.IsValueInside(v.Z)
}

// Checks if specific point is inside the boundary ignoring z position
func (b RectBoundary) IsPointInside2D(v Vector3) bool {
	return b.X.IsValueInside(v.X) && b.Y.IsValueInside(v.Y)
}

func (b RectBoundary) ToString() string {
	return fmt.Sprintf("(%s %s %s)", b.X.ToString(), b.Y.ToString(), b.Z.ToString())
}

func NewRectBoundary(minX, maxX, minY, maxY, minZ, maxZ float64) RectBoundary {
	return RectBoundary{
		X: NewFloatRange(minX, maxX),
		Y: NewFloatRange(minY, maxY),
		Z: NewFloatRange(minZ, maxZ),
	}
}