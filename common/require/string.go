package require

import (
	"fmt"
)

// String returns the given interface{} as string.
//If the given value is nil returns the given default string or "".
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
	case []rune:
		return string(i.([]rune))
	default:
		return fmt.Sprintf("%v", i)
	}
}