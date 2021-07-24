package common

import (
	"fmt"
	"github.com/cadmean-ru/amphion/common/a"
)

// Represents the boundaries of an object, like collider in unity
type Boundary interface {
	IsPointInside(point a.Vector3) bool
	IsPointInside2D(point a.Vector3) bool
}

// Represents a boundary in 3D space
type RectBoundary struct {
	X, Y, Z FloatRange
}

// Checks if specific point is inside the boundary
func (b *RectBoundary) IsPointInside(v a.Vector3) bool {
	return b.X.IsValueInside(v.X) && b.Y.IsValueInside(v.Y) && b.Z.IsValueInside(v.Z)
}

// Checks if specific point is inside the boundary ignoring z position
func (b *RectBoundary) IsPointInside2D(v a.Vector3) bool {
	return b.X.IsValueInside(v.X) && b.Y.IsValueInside(v.Y)
}

//IsRectInside checks if another rect is fully inside this rect.
func (b *RectBoundary) IsRectInside(rect *RectBoundary) bool {
	return b.X.IsRangeInside(rect.X) && b.Y.IsRangeInside(rect.Y) && b.Z.IsRangeInside(rect.Z)
}

//Move shifts all coordinates of the rect by the given vector.
func (b *RectBoundary) Move(by a.Vector3) {
	b.X.Move(by.X)
	b.Y.Move(by.Y)
	b.Z.Move(by.Z)
}

func (b *RectBoundary) GetMin() a.Vector3 {
	return a.NewVector3(b.X.Min, b.Y.Min, b.Z.Min)
}

func (b *RectBoundary) GetMax() a.Vector3 {
	return a.NewVector3(b.X.Max, b.Y.Max, b.Z.Max)
}

func (b *RectBoundary) GetSize() a.Vector3 {
	return a.NewVector3(b.X.GetLength(), b.Y.GetLength(), b.Z.GetLength())
}

func (b *RectBoundary) ToString() string {
	return fmt.Sprintf("(%s %s %s)", b.X.ToString(), b.Y.ToString(), b.Z.ToString())
}

func NewRectBoundary(minX, maxX, minY, maxY, minZ, maxZ float32) *RectBoundary {
	return &RectBoundary{
		X: NewFloatRange(minX, maxX),
		Y: NewFloatRange(minY, maxY),
		Z: NewFloatRange(minZ, maxZ),
	}
}

func NewRectBoundarySized(position a.Vector3, size a.Vector3) *RectBoundary {
	return &RectBoundary{
		X: NewFloatRange(position.X, position.X + size.X),
		Y: NewFloatRange(position.Y, position.Y + size.Y),
		Z: NewFloatRange(position.Z, position.Z + size.Z),
	}
}