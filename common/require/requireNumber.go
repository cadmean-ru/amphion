package require

import (
	"reflect"
	"strconv"
)

//Int converts the given number or string value to int.
//If conversion is not possible returns the given default value or 0 if no default value is specified.
func Int(num interface{}, defaultValue ...int) int {
	def := 0
	if len(defaultValue) > 0 {
		def = defaultValue[0]
	}

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
	case string:
		n, err := strconv.Atoi(num.(string))
		if err != nil {
			return def
		}
		return n
	default:
		return def
	}
}

//Int32 converts the given number or string value to int32.
//If conversion is not possible returns the given default value or 0 if no default value is specified.
func Int32(num interface{}, defaultValue ...int32) int32 {
	var def int32
	if len(defaultValue) > 0 {
		def = defaultValue[0]
	}

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
	case string:
		n, err := strconv.ParseInt(num.(string), 10, 32)
		if err != nil {
			return def
		}
		return int32(n)
	default:
		return def
	}
}

//Int64 converts the given number or string value to int64.
//If conversion is not possible returns the given default value or 0 if no default value is specified.
func Int64(num interface{}, defaultValue ...int64) int64 {
	var def int64
	if len(defaultValue) > 0 {
		def = defaultValue[0]
	}

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
	case string:
		n, err := strconv.ParseInt(num.(string), 10, 64)
		if err != nil {
			return def
		}
		return n
	default:
		return def
	}
}

//Float32 converts the given number or string value to float32.
//If conversion is not possible returns the given default value or 0 if no default value is specified.
func Float32(num interface{}, defaultValue ...float32) float32 {
	var def float32
	if len(defaultValue) > 0 {
		def = defaultValue[0]
	}

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
	case string:
		n, err := strconv.ParseFloat(num.(string), 10)
		if err != nil {
			return def
		}
		return float32(n)
	default:
		return def
	}
}

//Float64 converts the given number or string value to float64.
//If conversion is not possible returns the given default value or 0 if no default value is specified.
func Float64(num interface{}, defaultValue ...float64) float64 {
	var def float64
	if len(defaultValue) > 0 {
		def = defaultValue[0]
	}

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
	case string:
		n, err := strconv.ParseFloat(num.(string), 10)
		if err != nil {
			return def
		}
		return n
	default:
		return def
	}
}

//Byte converts the given number or string value to byte.
//If conversion is not possible returns the given default value or 0 if no default value is specified.
func Byte(num interface{}, defaultValue ...byte) byte {
	var def byte
	if len(defaultValue) > 0 {
		def = defaultValue[0]
	}

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
	case string:
		n, err := strconv.ParseUint(num.(string), 10, 8)
		if err != nil {
			return def
		}
		return byte(n)
	default:
		return def
	}
}

//Number converts the given number value to the given number type.
//If conversion is not possible returns 0.
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

var SupportedNumberTypes = []reflect.Kind { reflect.Uint8, reflect.Int32, reflect.Int64, reflect.Int, reflect.Float32, reflect.Float64 }

//IsNumber checks if the given value's type is one of the supported number types: byte, int, int32, int64, float32, float64.
func IsNumber(i interface{}) bool {
	t := reflect.TypeOf(i)
	k := t.Kind()

	for _, nt := range SupportedNumberTypes {
		if k == nt {
			return true
		}
	}

	return false
}