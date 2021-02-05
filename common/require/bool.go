package require

func Bool(value interface{}) bool {
	switch value.(type) {
	case bool:
		return value.(bool)
	default:
		return value != nil
	}
}
