package common

import (
	"github.com/cadmean-ru/amphion/common/a"
)

// Boundary represents the boundaries of an object, like collider in unity
type Boundary interface {
	IsPointInside(point a.Vector3) bool
	IsPointInside2D(point a.Vector3) bool
}
