//+build js

package web

import (
	"github.com/cadmean-ru/amphion/common"
	"github.com/cadmean-ru/amphion/common/a"
	"github.com/cadmean-ru/amphion/frontend"
	"github.com/cadmean-ru/amphion/rendering"
	"syscall/js"
)

type Frontend struct {
	renderer   *P5Renderer
	handler    frontend.CallbackHandler
	input      *InputManager
	context    frontend.Context
	msgChan    chan frontend.Message
	resManager *ResourceManager
}

func (f *Frontend) Init() {
	f.input.init(f)
	f.renderer.Prepare()

	js.Global().Get("addEventListener").Invoke("resize", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		ws := getWindowSize()
		f.renderer.p5.resizeCanvas(ws.X, ws.Y)
		f.context.ScreenInfo = common.NewScreenInfo(ws.X, ws.Y)
		f.handler(frontend.NewCallback(frontend.CallbackContextChange, ""))
		return nil
	}))

	js.Global().Get("addEventListener").Invoke("visibilitychange", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		var code int
		if js.Global().Get("document").Get("visibilityState").String() == "visible" {
			code = frontend.CallbackAppShow
		} else {
			code = frontend.CallbackAppHide
		}
		f.handler(frontend.NewCallback(code, ""))
		return nil
	}))

	ws := getWindowSize()
	f.context.ScreenInfo = common.NewScreenInfo(ws.X, ws.Y)
}

func (f *Frontend) Run() {
	for msg := range f.msgChan {
		if msg.Code == frontend.MessageExit {
			break
		}

		switch msg.Code {
		case frontend.MessageRender:
			f.renderer.PerformRendering()
		}
	}
	close(f.msgChan)
}

func (f *Frontend) Reset() {

}

func (f *Frontend) SetCallback(handler frontend.CallbackHandler) {
	f.handler = handler
}

func (f *Frontend) GetInputManager() frontend.InputManager {
	return f.input
}

func (f *Frontend) GetRenderer() rendering.Renderer {
	return f.renderer
}

func (f *Frontend) GetContext() frontend.Context {
	return f.context
}

func (f *Frontend) GetPlatform() common.Platform {
	return common.PlatformFromString("web")
}

func (f *Frontend) CommencePanic(reason, msg string) {
	js.Global().Get("commencePanic").Invoke(reason, msg)
}

func (f *Frontend) ReceiveMessage(message frontend.Message) {
	f.msgChan <- message
}

func (f *Frontend) GetResourceManager() frontend.ResourceManager {
	return f.resManager
}

func NewFrontend() *Frontend {
	f := &Frontend{
		input:      &InputManager{},
		msgChan:    make(chan frontend.Message, 10),
		resManager: newResourceManager(),
	}
	f.renderer = newP5Renderer(f)
	return f
}

func getWindowSize() a.IntVector2 {
	w := js.Global().Get("innerWidth").Int()
	h := js.Global().Get("innerHeight").Int()

	return a.IntVector2{w, h}
}
