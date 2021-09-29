package shape

import (
	"github.com/cadmean-ru/amphion/common/a"
)

type EmptyShape struct {

}

func (e EmptyShape) IsPointInside(_ a.Vector3) bool {
	return false
}

func (e EmptyShape) IsPointInside2D(_ a.Vector3) bool {
	return false
}

func (e EmptyShape) Kind() Kind {
	return Empty
}

