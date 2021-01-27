package common

// Limits the x value between min and max values.
func ClampFloat32(x, min, max float32) float32 {
	if x < min {
		return min
	} else if x > max {
		return max
	} else {
		return x
	}
}

// Limits the x value between min and max values.
func ClampFloat64(x, min, max float64) float64 {
	if x < min {
		return min
	} else if x > max {
		return max
	} else {
		return x
	}
}

// Limits the x value between min and max values.
func ClampInt(x, min, max int) int {
	if x < min {
		return min
	} else if x > max {
		return max
	} else {
		return x
	}
}
