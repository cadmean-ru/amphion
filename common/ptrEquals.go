package common

import "reflect"

func PtrEquals(v1, v2 any) bool {
	p1 := reflect.ValueOf(v1).Pointer()
	p2 := reflect.ValueOf(v2).Pointer()
	return p1 == p2
}
