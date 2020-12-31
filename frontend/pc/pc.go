// +build windows linux darwin
// +build !android

package pc

import (
	"fmt"
	"github.com/cadmean-ru/amphion/common"
	"github.com/cadmean-ru/amphion/common/a"
	"github.com/cadmean-ru/amphion/frontend"
	"github.com/cadmean-ru/amphion/rendering"
	"github.com/go-gl/glfw/v3.3/glfw"
	"math"
)

const SleepTimeMs = 16

type Frontend struct {
	window      *glfw.Window
	wSize       a.IntVector3
	handler     frontend.CallbackHandler
	renderer    *OpenGLRenderer
	initialized bool
	context     frontend.Context
	msgChan     chan frontend.Message
	inputMan    *InputManager
}

func (f *Frontend) Init() {
	var err error

	if err = glfw.Init(); err != nil {
		panic(err)
	}

	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)

	if f.window, err = glfw.CreateWindow(int(f.wSize.X), int(f.wSize.Y), "Amphion", nil, nil); err != nil {
		panic(err)
	}

	f.window.SetMouseButtonCallback(f.mouseButtonCallback)
	f.window.SetKeyCallback(f.keyCallback)
	f.window.SetFramebufferSizeCallback(f.frameBufferSizeCallback)
	f.window.SetFocusCallback(f.focusCallback)

	f.context.ScreenInfo = common.NewScreenInfo(int(f.wSize.X), int(f.wSize.Y))

	f.renderer.window = f.window
	f.renderer.wSize = f.wSize
	f.renderer.Prepare()

	f.inputMan.window = f.window

	f.initialized = true
}

func (f *Frontend) Run() {
	for !f.window.ShouldClose() {
		select {
		case msg, ok := <-f.msgChan:
			if ok {
				switch msg.Code {
				case frontend.MessageRender:
					f.renderer.PerformRendering()
				}
			} else {

			}
		default:

		}

		//glfw.PollEvents()
		//
		//time.Sleep(SleepTimeMs)

		glfw.WaitEventsTimeout(SleepTimeMs)
	}

	glfw.Terminate()
}

func (f *Frontend) mouseButtonCallback(w *glfw.Window, button glfw.MouseButton, action glfw.Action, _ glfw.ModifierKey) {
	var callback frontend.Callback

	mouseX, mouseY := w.GetCursorPos()
	data := fmt.Sprintf("%d;%d", int(math.Floor(mouseX)), int(math.Floor(mouseY)))

	switch button {
	case glfw.MouseButton1:
		switch action {
		case glfw.Press:
			callback = frontend.NewCallback(frontend.CallbackMouseDown, data)
		case glfw.Release:
			callback = frontend.NewCallback(frontend.CallbackMouseUp, data)
		}
	}

	f.handler(callback)
}

func (f *Frontend) keyCallback(_ *glfw.Window, key glfw.Key, scancode int, action glfw.Action, _ glfw.ModifierKey) {
	keyName := glfw.GetKeyName(key, scancode)
	data := fmt.Sprintf("%s\n", keyName)
	var code int

	switch action {
	case glfw.Press:
		code = frontend.CallbackKeyDown
	}

	f.handler(frontend.NewCallback(code, data))
}

func (f *Frontend) frameBufferSizeCallback(_ *glfw.Window, width int, height int) {
	f.wSize = a.NewIntVector3(width, height, 0)
	f.renderer.wSize = f.wSize
	f.context.ScreenInfo = common.NewScreenInfo(width, height)
	f.renderer.handleWindowResize()
	f.handler(frontend.NewCallback(frontend.CallbackContextChange, ""))
}

func (f *Frontend) focusCallback(_ *glfw.Window, focused bool) {
	var code int
	if focused {
		code = frontend.CallbackAppShow
	} else {
		code = frontend.CallbackAppHide
	}
	f.handler(frontend.NewCallback(code, ""))
}

func (f *Frontend) Stop() {
	glfw.Terminate()
}

func (f *Frontend) Reset() {

}

func (f *Frontend) SetCallback(handler frontend.CallbackHandler) {
	f.handler = handler
}

func (f *Frontend) GetInputManager() frontend.InputManager {
	return f.inputMan
}

func (f *Frontend) GetRenderer() rendering.Renderer {
	if !f.initialized {
		return nil
	}
	return f.renderer
}

func (f *Frontend) GetPlatform() common.Platform {
	return common.PlatformFromString("pc")
}

func (f *Frontend) GetContext() frontend.Context {
	return f.context
}

func (f *Frontend) CommencePanic(reason, msg string) {
	panic(fmt.Sprintf("%s: %s", reason, msg))
}

func (f *Frontend) ReceiveMessage(message frontend.Message) {
	f.msgChan<-message
}

func NewFrontend() *Frontend {
	return &Frontend{
		wSize: a.NewIntVector3(500, 500, 0),
		renderer: &OpenGLRenderer{
			primitives: make(map[int64]*glContainer),
			fonts:      make(map[string]*glFont),
		},
		inputMan: &InputManager{},
		msgChan: make(chan frontend.Message, 10),
	}
}
