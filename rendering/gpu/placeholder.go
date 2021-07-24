package gpu

import (
	"github.com/cadmean-ru/amphion/common"
)

type Placeholder struct {

}

func (p Placeholder) AllocateBuffer(size int) Buffer {
	return nil
}

func (p Placeholder) SetClippingArea(rect *common.RectBoundary) {

}

func (p Placeholder) ClearClippingArea() {

}

