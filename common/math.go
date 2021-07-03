package common

// ClampFloat32 limits the x value between min and max values.
func ClampFloat32(x, min, max float32) float32 {
	if x < min {
		return min
	} else if x > max {
		return max
	} else {
		return x
	}
}

// ClampFloat64 limits the x value between min and max values.
func ClampFloat64(x, min, max float64) float64 {
	if x < min {
		return min
	} else if x > max {
		return max
	} else {
		return x
	}
}

// ClampInt limits the x value between min and max values.
func ClampInt(x, min, max int) int {
	if x < min {
		return min
	} else if x > max {
		return max
	} else {
		return x
	}
}

//MaxFloat32 returns the maximum of two floats.
func MaxFloat32(a, b float32) float32 {
	if b > a {
		return b
	}

	return a
}

//MinFloat32 returns the minimum of two floats.
func MinFloat32(a, b float32) float32 {
	if b < a {
		return b
	}

	return a
}