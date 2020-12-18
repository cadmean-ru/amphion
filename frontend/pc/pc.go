// +build windows linux darwin

package pc

import (
	"fmt"
	"github.com/cadmean-ru/amphion/common"
	"github.com/cadmean-ru/amphion/frontend/commonFrontend"
	"github.com/cadmean-ru/amphion/rendering"
	"github.com/go-gl/glfw/v3.3/glfw"
	"time"
)

const SleepTimeMs = 16

type Frontend struct {
	window      *glfw.Window
	wSize       common.IntVector3
	handler     commonFrontend.CallbackHandler
	renderer    *OpenGLRenderer
	initialized bool
	context     commonFrontend.Context
	msgChan     chan commonFrontend.Message
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

	f.context.ScreenInfo = common.NewScreenInfo(f.wSize.X, f.wSize.Y)

	f.renderer.window = f.window
	f.renderer.wSize = f.wSize

	f.renderer.Prepare()

	f.initialized = true
}

func (f *Frontend) Run() {
	for !f.window.ShouldClose() {
		glfw.PollEvents()

		select {
		case msg, ok := <-f.msgChan:
			if ok {
				switch msg.Code {
				case commonFrontend.MessageRender:
					f.renderer.PerformRendering()
				}
			} else {

			}
		default:

		}

		time.Sleep(SleepTimeMs)
	}

	glfw.Terminate()
}

func (f *Frontend) Stop() {
	//f.renderer.Stop()

	glfw.Terminate()
}

func (f *Frontend) Reset() {

}

func (f *Frontend) SetCallback(handler commonFrontend.CallbackHandler) {
	f.handler = handler
}

func (f *Frontend) GetInputManager() commonFrontend.InputManager {
	return nil
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

func (f *Frontend) GetContext() commonFrontend.Context {
	return f.context
}

func (f *Frontend) CommencePanic(reason, msg string) {
	panic(fmt.Sprintf("%s: %s", reason, msg))
}

func (f *Frontend) ReceiveMessage(message commonFrontend.Message) {
	f.msgChan<-message
}

func NewFrontend() *Frontend {
	return &Frontend{
		wSize:    common.NewIntVector3(500, 500, 0),
		renderer: &OpenGLRenderer{
			primitives: make(map[int64]*glContainer),
		},
		msgChan: make(chan commonFrontend.Message, 10),
	}
}
