package common

import (
	"fmt"
	"math"
)

// Represents a range of values between min and max inclusive
type FloatRange struct {
	Min, Max float32
}

//NewFloatRange creates a new FloatRange with the given min and max values.
func NewFloatRange(min, max float32) FloatRange {
	if max < min {
		return FloatRange{
			Min: max,
			Max: min,
		}
	}
	return FloatRange{
		Min: min,
		Max: max,
	}
}

// Gets the length of the range
func (r *FloatRange) GetLength() float32 {
	return float32(math.Abs(float64(r.Max - r.Min)))
}

// Checks if specific value falls inside the range
func (r *FloatRange) IsValueInside(value float32) bool {
	return value >= r.Min && value <= r.Max
}

// IsRangeInside checks if the given range is inside this range.
func (r *FloatRange) IsRangeInside(r2 FloatRange) bool {
	return r.IsValueInside(r2.Min) && r.IsValueInside(r2.Max)
}

//Move shifts the range's min and max values by the given value.
func (r *FloatRange) Move(d float32) {
	r.Min += d
	r.Max += d
}

func (r *FloatRange) ToString() string {
	return fmt.Sprintf("[%f %f]", r.Min, r.Max)
}