package require

//Bool converts the given value to bool.
//If the value's type is bool returns this value.
//Otherwise, returns true if the value is not nil.
func Bool(value interface{}) bool {
	switch value.(type) {
	case bool:
		return value.(bool)
	default:
		return value != nil
	}
}
