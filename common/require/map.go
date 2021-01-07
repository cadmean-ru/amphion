package require

func SiMap(i interface{}) map[string]interface{} {
	if sim, ok := i.(map[string]interface{}); ok {
		return sim
	}

	switch i.(type) {
	case map[string]interface{}:
		return i.(map[string]interface{})
	case map[interface{}]interface{}:
		m := make(map[string]interface{})
		for k, v := range i.(map[interface{}]interface{}) {
			if sk, ok := k.(string); ok {
				m[sk] = v
			}
		}
		return m
	default:
		return map[string]interface{} {}
	}
}