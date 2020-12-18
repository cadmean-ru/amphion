// +build windows linux

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
	window   *glfw.Window
	wSize    common.IntVector3
	handler  commonFrontend.CallbackHandler
	renderer *OpenGLRenderer
	prepared bool
	context  commonFrontend.Context
}

func (f *Frontend) Init() {
	var err error

	if err = glfw.Init(); err != nil {
		panic(err)
	}
	if f.window, err = glfw.CreateWindow(f.wSize.X, f.wSize.Y, "Amphion", nil, nil); err != nil {
		panic(err)
	}

	f.context.ScreenInfo = common.NewScreenInfo(f.wSize.X, f.wSize.Y)

	f.renderer.window = f.window
	f.renderer.wSize = f.wSize

	f.prepared = true
}

func (f *Frontend) Start() {
	go f.loop()
}

func (f *Frontend) loop() {
	for !f.window.ShouldClose() {
		glfw.PollEvents()

		time.Sleep(SleepTimeMs)
	}

	glfw.Terminate()
}

func (f *Frontend) Stop() {
	f.window.SetShouldClose(true)
	f.renderer.Stop()
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
	if !f.prepared {
		return nil
	}
	return f.renderer
}

func (f *Frontend) GetContext() commonFrontend.Context {
	return f.context
}

func (f *Frontend) CommencePanic(reason, msg string) {
	panic(fmt.Sprintf("%s: %s", reason, msg))
}

func NewFrontend() *Frontend {
	return &Frontend{
		wSize:    common.NewIntVector3(500, 500, 0),
		renderer: &OpenGLRenderer{
			primitives: make(map[int64]*glContainer),
		},
	}
}
