package common

import (
	"fmt"
	"math"
)

// Represents a range of values between min and max inclusive
type FloatRange struct {
	Min, Max float64
}

func NewFloatRange(min, max float64) FloatRange {
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
func (r FloatRange) GetLength() float64 {
	return math.Abs(r.Max - r.Min)
}

// Checks if specific value falls inside the range
func (r FloatRange) IsValueInside(value float64) bool {
	return value >= r.Min && value <= r.Max
}

func (r FloatRange) ToString() string {
	return fmt.Sprintf("[%f %f]", r.Min, r.Max)
}