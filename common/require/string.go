package require

import (
	"fmt"
	"reflect"
)

func String(i interface{}) string {
	switch i.(type) {
	case string:
		return i.(string)
	case []byte:
		return string(i.([]byte))
	default:
		return fmt.Sprintf("%v", i)
	}
}

var SupportedNumberTypes = []reflect.Kind { reflect.Uint8, reflect.Int32, reflect.Int64, reflect.Int, reflect.Float32, reflect.Float64 }

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