package require

import "reflect"

func Int(num interface{}) int {
	switch num.(type) {
	case byte:
		return int(num.(byte))
	case int:
		return num.(int)
	case int32:
		return int(num.(int32))
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

func Int32(num interface{}) int32 {
	switch num.(type) {
	case byte:
		return int32(num.(byte))
	case int:
		return int32(num.(int))
	case int32:
		return num.(int32)
	case int64:
		return int32(num.(int64))
	case float32:
		return int32(num.(float32))
	case float64:
		return int32(num.(float64))
	default:
		return 0
	}
}

func Int64(num interface{}) int64 {
	switch num.(type) {
	case byte:
		return int64(num.(byte))
	case int:
		return int64(num.(int))
	case int32:
		return int64(num.(int32))
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

func Float32(num interface{}) float32 {
	switch num.(type) {
	case byte:
		return float32(num.(byte))
	case int:
		return float32(num.(int))
	case int32:
		return float32(num.(int32))
	case int64:
		return float32(num.(int64))
	case float32:
		return num.(float32)
	case float64:
		return float32(num.(float64))
	default:
		return 0
	}
}

func Float64(num interface{}) float64 {
	switch num.(type) {
	case byte:
		return float64(num.(byte))
	case int:
		return float64(num.(int))
	case int32:
		return float64(num.(int32))
	case int64:
		return float64(num.(int64))
	case float32:
		return float64(num.(float32))
	case float64:
		return num.(float64)
	default:
		return 0
	}
}

func Byte(num interface{}) byte {
	switch num.(type) {
	case byte:
		return num.(byte)
	case int:
		return byte(num.(int))
	case int32:
		return byte(num.(int32))
	case int64:
		return byte(num.(int64))
	case float32:
		return byte(num.(float32))
	case float64:
		return byte(num.(float64))
	default:
		return 0
	}
}

func Number(value interface{}, kind reflect.Kind) interface{} {
	switch kind {
	case reflect.Uint8:
		return Byte(value)
	case reflect.Int:
		return Int(value)
	case reflect.Int32:
		return Int32(value)
	case reflect.Int64:
		return Int64(value)
	case reflect.Float32:
		return Float32(value)
	case reflect.Float64:
		return Float64(value)
	default:
		return 0
	}
}