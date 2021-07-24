package gpu

import "github.com/cadmean-ru/amphion/common"

type Gpu interface {
	Init()
	AllocateBuffer(size int) Buffer
	SetClippingArea(rect *common.RectBoundary)
	ClearClippingArea()
}
