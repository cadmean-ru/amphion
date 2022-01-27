package common

import (
	"fmt"
	"github.com/cadmean-ru/amphion/common/a"
)

// Rect represents a boundary in 3D space
type Rect struct {
	X, Y, Z FloatRange
}

// IsPointInside checks if specific point is inside the boundary
func (b *Rect) IsPointInside(v a.Vector3) bool {
	return b.X.IsValueInside(v.X) && b.Y.IsValueInside(v.Y) && b.Z.IsValueInside(v.Z)
}

// IsPointInside2D checks if specific point is inside the boundary ignoring z position
func (b *Rect) IsPointInside2D(v a.Vector3) bool {
	return b.X.IsValueInside(v.X) && b.Y.IsValueInside(v.Y)
}

//IsRectInside checks if another rect is fully inside this rect.
func (b *Rect) IsRectInside(rect *Rect) bool {
	return b.X.IsRangeInside(rect.X) && b.Y.IsRangeInside(rect.Y) && b.Z.IsRangeInside(rect.Z)
}

//Move shifts all coordinates of the rect by the given vector.
func (b *Rect) Move(by a.Vector3) {
	b.X.Move(by.X)
	b.Y.Move(by.Y)
	b.Z.Move(by.Z)
}

//Shrink shrinks all dimensions of the rect by the given vector.
func (b *Rect) Shrink(by a.Vector3) {
	b.X.Shrink(by.X)
	b.Y.Shrink(by.Y)
	b.Z.Shrink(by.Z)
}

//ShrinkMax shrinks max value of all dimensions of the rect by the given vector.
func (b *Rect) ShrinkMax(by a.Vector3) {
	b.X.ShrinkMax(by.X)
	b.Y.ShrinkMax(by.Y)
	b.Z.ShrinkMax(by.Z)
}

//Min returns the min (top-left) point of the rect.
func (b *Rect) Min() a.Vector3 {
	return a.NewVector3(b.X.Min, b.Y.Min, b.Z.Min)
}

//Max returns the max (bottom-right) point of the rect.
func (b *Rect) Max() a.Vector3 {
	return a.NewVector3(b.X.Max, b.Y.Max, b.Z.Max)
}

//Size returns the size of the rect.
func (b *Rect) Size() a.Vector3 {
	return a.NewVector3(b.X.GetLength(), b.Y.GetLength(), b.Z.GetLength())
}

func (b *Rect) String() string {
	return fmt.Sprintf("(%s %s %s)", b.X.String(), b.Y.String(), b.Z.String())
}

func NewRect(minX, maxX, minY, maxY, minZ, maxZ float32) *Rect {
	return &Rect{
		X: NewFloatRange(minX, maxX),
		Y: NewFloatRange(minY, maxY),
		Z: NewFloatRange(minZ, maxZ),
	}
}

func NewRect2D(minX, maxX, minY, maxY, z float32) *Rect {
	return &Rect{
		X: NewFloatRange(minX, maxX),
		Y: NewFloatRange(minY, maxY),
		Z: NewFloatRange(z, z),
	}
}

func NewRectFromPositionAndSize(position a.Vector3, size a.Vector3) *Rect {
	return NewRect(position.X, position.X+size.X, position.Y, position.Y+size.Y, position.Z, position.Z+size.Z)
}

func NewRectFromMinMaxPositions(min a.Vector3, max a.Vector3) *Rect {
	return NewRect(min.X, max.X, min.Y, max.Y, min.Z, max.Z)
}

func NewZeroRect() *Rect {
	return NewRect(0, 0, 0, 0, 0, 0)
}
