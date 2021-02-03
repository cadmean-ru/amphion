//+build js

package web

import (
	"fmt"
	"github.com/cadmean-ru/amphion/common/a"
	"github.com/cadmean-ru/amphion/frontend"
	"syscall/js"
)

type InputManager struct {
	mousePos a.IntVector2
}

func (m *InputManager) init(f *Frontend) {
	js.Global().Get("document").Set("onmousemove", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		e := args[0]
		x := e.Get("pageX").Int()
		y := e.Get("pageY").Int()
		m.mousePos = a.IntVector2{X: x, Y: y}
		f.handler(frontend.NewCallback(frontend.CallbackMouseMove, ""))
		return nil
	}))

	js.Global().Get("document").Set("onmousedown", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		f.handler(frontend.NewCallback(frontend.CallbackMouseDown, fmt.Sprintf("%d;%d", m.mousePos.X, m.mousePos.Y)))
		return nil
	}))

	js.Global().Get("document").Set("onmouseup", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		f.handler(frontend.NewCallback(frontend.CallbackMouseUp, fmt.Sprintf("%d;%d", m.mousePos.X, m.mousePos.Y)))
		return nil
	}))

	js.Global().Get("document").Set("onkeydown", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		e := args[0]
		key := e.Get("key").String()
		code := e.Get("code").String()

		f.handler(frontend.NewCallback(frontend.CallbackKeyDown, fmt.Sprintf("%s\n%s", key, code)))
		return nil
	}))
}

func (m *InputManager) GetMousePosition() a.IntVector2 {
	return m.mousePos
}
