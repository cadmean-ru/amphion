package common

type SiMap map[string]interface{}

type Mappable interface {
	ToMap() SiMap
}

type Unmappable interface {
	FromMap(siMap SiMap)
}

func RequireSiMap(i interface{}) SiMap {
	switch i.(type) {
	case SiMap:
		return i.(SiMap)
	case map[string]interface{}:
		return i.(SiMap)
	case map[interface{}]interface{}:
		m := make(SiMap)
		for k, v := range i.(map[interface{}]interface{}) {
			switch k.(type) {
			case string:
				m[k.(string)] = v
			default:
				continue
			}
		}
		return m
	default:
		return nil
	}
}