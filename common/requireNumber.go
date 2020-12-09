package common

func RequireInt(num interface{}) int {
	switch num.(type) {
	case int:
		return num.(int)
	case int64:
		return int(num.(int64))
	case float32:
		return int(num.(float32))
	case float64:
		return int(num.(float64))
	default:
		return 0
	}
}

func RequireInt64(num interface{}) int64 {
	switch num.(type) {
	case int:
		return int64(num.(int))
	case int64:
		return num.(int64)
	case float32:
		return int64(num.(float32))
	case float64:
		return int64(num.(float64))
	default:
		return 0
	}
}