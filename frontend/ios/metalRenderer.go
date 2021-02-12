package ios

import (
	"github.com/cadmean-ru/amphion/rendering"
)

type MetalRenderer struct {

}

func (m *MetalRenderer) Prepare() {
}

func (m *MetalRenderer) AddPrimitive() int {
	return 0
}

func (m *MetalRenderer) SetPrimitive(id int, primitive rendering.IPrimitive, shouldRerender bool) {
}

func (m *MetalRenderer) RemovePrimitive(id int) {
}

func (m *MetalRenderer) PerformRendering() {
}

func (m *MetalRenderer) Clear() {
}

func (m *MetalRenderer) Stop() {
}

