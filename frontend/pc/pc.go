// +build windows linux darwin
// +build !android
// +build !ios

package pc

import (
	"fmt"
	"github.com/cadmean-ru/amphion/common"
	"github.com/cadmean-ru/amphion/common/a"
	"github.com/cadmean-ru/amphion/common/dispatch"
	"github.com/cadmean-ru/amphion/frontend"
	"github.com/cadmean-ru/amphion/rendering"
	"github.com/go-gl/glfw/v3.3/glfw"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"math"
	"time"
)

const SleepTimeS = 1.0 / 60.0

type Frontend struct {
	window           *glfw.Window
	wSize            a.IntVector3
	rendererDelegate *OpenGLRenderer
	renderer         *rendering.ARenderer
	initialized      bool
	context          frontend.Context
	msgChan          *dispatch.MessageQueue
	resMan           *ResourceManager
	app              *frontend.App
	disp             dispatch.MessageDispatcher
}

func (f *Frontend) Init() {
	var err error

	if err = glfw.Init(); err != nil {
		panic(err)
	}

	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
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

	f.renderer.Prepare()

	f.initialized = true
}

func (f *Frontend) Run() {
	f.disp.SendMessage(dispatch.NewMessage(frontend.CallbackReady))
	fmt.Printf("Frontend target frame time: %fms\n", SleepTimeS*1000)

	//defer profile.Start(profile.ProfilePath(".")).Stop()

	for !f.window.ShouldClose() {
		msgTime := f.processMessages()

		timeToWait := 0.0
		if SleepTimeS - msgTime.Seconds() > 0 {
			timeToWait = SleepTimeS - msgTime.Seconds()
		}

		//if timeToWait == 0 && f.app != nil && f.app.Debug {
		//	fmt.Println("Warning! The frontend is skipping frames!")
		//	fmt.Printf("Message processing took: %dms\n", msgTime.Milliseconds())
		//}

		glfw.WaitEventsTimeout(timeToWait)
	}

	// TODO: notify the engine to stop correctly

	glfw.Terminate()
}

func (f *Frontend) processMessages() time.Duration {
	f.msgChan.LockMainChannel()

	var total time.Duration = 0

	for !f.msgChan.IsEmpty() {
		total += f.processMessage(f.msgChan.Dequeue())
	}

	f.msgChan.UnlockMainChannel()

	return total
}

func (f *Frontend) processMessage(msg *dispatch.Message) time.Duration {
	processingStart := time.Now()

	switch msg.What {
	case frontend.MessageExec, dispatch.MessageWorkExec:
		if msg.AnyData != nil {
			if work, ok := msg.AnyData.(dispatch.WorkItem); ok {
				work.Execute()
			}
		}
	case frontend.MessageTitle:
		f.window.SetTitle(msg.StrData)
	}

	processingTime := time.Since(processingStart)
	//if processingTime.Seconds() > SleepTimeS && f.app != nil && f.app.Debug {
	//	fmt.Printf("Warning! message processing took too long: %dms code: %d\n", processingTime.Milliseconds(), msg.Code)
	//}

	return processingTime
}

func (f *Frontend) mouseButtonCallback(w *glfw.Window, button glfw.MouseButton, action glfw.Action, _ glfw.ModifierKey) {
	var callback *dispatch.Message

	mouseX, mouseY := w.GetCursorPos()
	data := fmt.Sprintf("%d;%d", int(math.Floor(mouseX)), int(math.Floor(mouseY)))

	switch button {
	case glfw.MouseButton1:
		switch action {
		case glfw.Press:
			callback = dispatch.NewMessageWithStringData(frontend.CallbackMouseDown, data)
		case glfw.Release:
			callback = dispatch.NewMessageWithStringData(frontend.CallbackMouseUp, data)
		}
	}

	f.disp.SendMessage(callback)
}

func (f *Frontend) keyCallback(_ *glfw.Window, key glfw.Key, scancode int, action glfw.Action, _ glfw.ModifierKey) {
	keyName := glfw.GetKeyName(key, scancode)
	data := fmt.Sprintf("%s\n", keyName)
	var code int

	switch action {
	case glfw.Press:
		code = frontend.CallbackKeyDown
	}

	f.disp.SendMessage(dispatch.NewMessageWithStringData(code, data))
}

func (f *Frontend) frameBufferSizeCallback(_ *glfw.Window, width int, height int) {
	w, h := f.window.GetSize()
	f.wSize = a.NewIntVector3(w, h, 0)
	f.rendererDelegate.wSize = f.wSize
	f.context.ScreenInfo = common.NewScreenInfo(w, h)
	f.rendererDelegate.handleWindowResize(w, h)
	f.disp.SendMessage(dispatch.NewMessage(frontend.CallbackContextChange))
}

func (f *Frontend) focusCallback(_ *glfw.Window, focused bool) {
	var code int
	if focused {
		code = frontend.CallbackAppShow
	} else {
		code = frontend.CallbackAppHide
	}
	f.disp.SendMessage(dispatch.NewMessage(code))
}

func (f *Frontend) windowRefreshCallback(_ *glfw.Window) {
	f.disp.SendMessage(dispatch.NewMessage(frontend.CallbackContextChange))
}

func (f *Frontend) cursorPosCallback(_ *glfw.Window, x float64, y float64) {
	f.disp.SendMessage(dispatch.NewMessageWithStringData(frontend.CallbackMouseMove, fmt.Sprintf("%d;%d", int(x), int(y))))
}

func (f *Frontend) scrollCallback(_ *glfw.Window, xoff float64, yoff float64) {
	f.disp.SendMessage(dispatch.NewMessageWithStringData(frontend.CallbackMouseScroll, fmt.Sprintf("%f:%f", xoff, yoff)))
}

func (f *Frontend) Stop() {
	f.renderer.Stop()
}

func (f *Frontend) SetCallback(handler frontend.CallbackHandler) {

}

func (f *Frontend) SetEngineDispatcher(disp dispatch.MessageDispatcher) {
	f.disp = disp
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

}

func (f *Frontend) Execute(item dispatch.WorkItem) {
	f.msgChan.Enqueue(dispatch.NewMessageWithAnyData(dispatch.MessageWorkExec, item))
}

func (f *Frontend) SendMessage(message *dispatch.Message) {
	f.msgChan.Enqueue(message)
}

func (f *Frontend) GetMessageDispatcher() dispatch.MessageDispatcher {
	return f
}

func (f *Frontend) GetWorkDispatcher() dispatch.WorkDispatcher {
	return f
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

	f.app = &app
	return f.app
}

func (f *Frontend) GetLaunchArgs() a.SiMap {
	args := make(a.SiMap)
	return args
}

func NewFrontend() *Frontend {
	f := &Frontend{
		wSize: a.NewIntVector3(500, 500, 0),
		msgChan:  dispatch.NewMessageQueue(1000),
		resMan:   newResourceManager(),
	}
	f.rendererDelegate = &OpenGLRenderer{
		front:      f,
	}
	f.renderer = rendering.NewARenderer(f.rendererDelegate, f)
	return f
}
