package common

type SiMap map[string]interface{}

type Mappable interface {
	ToMap() SiMap
}

type Unmappable interface {
	FromMap(siMap SiMap)
}

//func BuildMappableArray(arr interface{}) []map[string]interface{} {
//	//switch arr.(type) {
//	//case []interface{}:
//	//	fmt.Println("fuu")
//	//	mArr := arr.([]interface{})
//	//	newArr := make([]map[string]interface{}, len(mArr))
//	//	for i, el := range mArr {
//	//		switch el.(type) {
//	//		case Mappable:
//	//			newArr[i] = el.(Mappable).ToMap()
//	//		default:
//	//			newArr[i] = nil
//	//		}
//	//	}
//	//	return newArr
//	//default:
//	//	return nil
//	//}
//
//	if mArr, ok := arr.([]Mappable); ok {
//		fmt.Println("fuu")
//		newArr := make([]map[string]interface{}, len(mArr))
//		for i, el := range mArr {
//			newArr[i] = el.ToMap()
//		}
//		return newArr
//	}
//
//	return nil
//}