package common

import (
	"fmt"
	"github.com/cadmean-ru/amphion/common/a"
)

// Boundary represents the boundaries of an object, like collider in unity
type Boundary interface {
	IsPointInside(point a.Vector3) bool
	IsPointInside2D(point a.Vector3) bool
}

// RectBoundary represents a boundary in 3D space
type RectBoundary struct {
	X, Y, Z FloatRange
}

// IsPointInside checks if specific point is inside the boundary
func (b *RectBoundary) IsPointInside(v a.Vector3) bool {
	return b.X.IsValueInside(v.X) && b.Y.IsValueInside(v.Y) && b.Z.IsValueInside(v.Z)
}

// IsPointInside2D checks if specific point is inside the boundary ignoring z position
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

//Shrink shrinks all dimensions of the rect by the given vector.
func (b *RectBoundary) Shrink(by a.Vector3) {
	b.X.Shrink(by.X)
	b.Y.Shrink(by.Y)
	b.Z.Shrink(by.Z)
}

//Min returns the min (top-left) point of the rect.
func (b *RectBoundary) Min() a.Vector3 {
	return a.NewVector3(b.X.Min, b.Y.Min, b.Z.Min)
}

//Max returns the max (bottom-right) point of the rect.
func (b *RectBoundary) Max() a.Vector3 {
	return a.NewVector3(b.X.Max, b.Y.Max, b.Z.Max)
}

//Size returns the size of the rect.
func (b *RectBoundary) Size() a.Vector3 {
	return a.NewVector3(b.X.GetLength(), b.Y.GetLength(), b.Z.GetLength())
}

func (b *RectBoundary) String() string {
	return fmt.Sprintf("(%s %s %s)", b.X.String(), b.Y.String(), b.Z.String())
}

func NewRectBoundary(minX, maxX, minY, maxY, minZ, maxZ float32) *RectBoundary {
	return &RectBoundary{
		X: NewFloatRange(minX, maxX),
		Y: NewFloatRange(minY, maxY),
		Z: NewFloatRange(minZ, maxZ),
	}
}

func NewRectBoundary2D(minX, maxX, minY, maxY, z float32) *RectBoundary {
	return &RectBoundary{
		X: NewFloatRange(minX, maxX),
		Y: NewFloatRange(minY, maxY),
		Z: NewFloatRange(z, z),
	}
}

func NewRectBoundaryFromPositionAndSize(position a.Vector3, size a.Vector3) *RectBoundary {
	return NewRectBoundary(position.X, position.X + size.X, position.Y, position.Y + size.Y, position.Z, position.Z + size.Z)
}