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
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"math"
)

const SleepTimeS = float64(1) / 120

type Frontend struct {
	window           *glfw.Window
	wSize            a.IntVector3
	handler          frontend.CallbackHandler
	rendererDelegate *OpenGLRenderer
	renderer         *rendering.ARenderer
	initialized      bool
	context          frontend.Context
	msgChan          chan frontend.Message
	inputMan         *InputManager
	resMan           *ResourceManager
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

	if f.window, err = glfw.CreateWindow(f.wSize.X, f.wSize.Y, "Amphion", nil, nil); err != nil {
		panic(err)
	}

	f.window.SetMouseButtonCallback(f.mouseButtonCallback)
	f.window.SetKeyCallback(f.keyCallback)
	f.window.SetFramebufferSizeCallback(f.frameBufferSizeCallback)
	f.window.SetFocusCallback(f.focusCallback)
	f.window.SetRefreshCallback(f.windowRefreshCallback)
	f.window.SetCursorPosCallback(f.cursorPosCallback)
	f.window.SetScrollCallback(f.scrollCallback)

	f.context.ScreenInfo = common.NewScreenInfo(f.wSize.X, f.wSize.Y)

	f.rendererDelegate.window = f.window
	f.rendererDelegate.wSize = f.wSize

	f.inputMan.window = f.window

	f.renderer.Prepare()

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
				case frontend.MessageExec:
					if msg.Data != nil {
						if action, ok := msg.Data.(func()); ok {
							action()
						}
					}
				}
			} else {

			}
		default:

		}

		glfw.WaitEventsTimeout(SleepTimeS)
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
	f.rendererDelegate.wSize = f.wSize
	f.context.ScreenInfo = common.NewScreenInfo(width, height)
	f.rendererDelegate.handleWindowResize(width, height)
	f.handler(frontend.NewCallback(frontend.CallbackContextChange, ""))
	f.renderer.PerformRendering()
	fmt.Printf("New window size: %d %d\n", width, height)
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

func (f *Frontend) windowRefreshCallback(_ *glfw.Window) {
	f.handler(frontend.NewCallback(frontend.CallbackContextChange, ""))
}

func (f *Frontend) cursorPosCallback(_ *glfw.Window, _ float64, _ float64) {
	f.handler(frontend.NewCallback(frontend.CallbackMouseMove, ""))
}

func (f *Frontend) scrollCallback(_ *glfw.Window, xoff float64, yoff float64) {
	f.handler(frontend.NewCallback(frontend.CallbackMouseScroll, fmt.Sprintf("%f:%f", xoff, yoff)))
}

func (f *Frontend) Stop() {
	f.renderer.Stop()
}

func (f *Frontend) Reset() {

}

func (f *Frontend) SetCallback(handler frontend.CallbackHandler) {
	f.handler = handler
}

func (f *Frontend) GetInputManager() frontend.InputManager {
	return f.inputMan
}

func (f *Frontend) GetRenderer() *rendering.ARenderer {
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
	f.msgChan <- message
}

func (f *Frontend) GetResourceManager() frontend.ResourceManager {
	return f.resMan
}

func (f *Frontend) GetApp() *frontend.App {
	data, err := ioutil.ReadFile("app.yaml")
	if err != nil {
		return nil
	}

	app := frontend.App{}
	err = yaml.Unmarshal(data, &app)
	if err != nil {
		return nil
	}

	return &app
}

func (f *Frontend) GetLaunchArgs() a.SiMap {
	args := make(a.SiMap)
	return args
}
func NewFrontend() *Frontend {
	f := &Frontend{
		wSize: a.NewIntVector3(500, 500, 0),
		inputMan: &InputManager{},
		msgChan:  make(chan frontend.Message, 10),
		resMan:   newResourceManager(),
	}
	f.rendererDelegate = &OpenGLRenderer{
		front:      f,
	}
	f.renderer = rendering.NewARenderer(f.rendererDelegate)
	return f
}
