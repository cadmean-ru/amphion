package rendering

import "github.com/cadmean-ru/amphion/common/a"

type IPrimitive interface {
	GetType() a.Byte
	GetTransform() Transform
}
