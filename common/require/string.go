package require

import (
	"fmt"
	"reflect"
)

// Returns the given interface{} as string.
func String(i interface{}, defaultString ...string) string {
	def := ""
	if len(defaultString) > 0 {
		def = defaultString[0]
	}

	if i == nil {
		return def
	}

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