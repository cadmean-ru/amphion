// +build windows linux darwin
// +build !android

package pc

import (
	"github.com/cadmean-ru/amphion/common/a"
	"github.com/go-gl/glfw/v3.3/glfw"
)

type InputManager struct {
	window *glfw.Window
}

func (m *InputManager) GetMousePosition() a.IntVector3 {
	x, y := m.window.GetCursorPos()
	return a.NewIntVector3(int(x), int(y), 0)
}