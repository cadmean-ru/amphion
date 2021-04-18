//+build js

// Package web provides implementation of web frontend.
package web

import (
	"fmt"
	"github.com/cadmean-ru/amphion/common"
	"github.com/cadmean-ru/amphion/common/a"
	"github.com/cadmean-ru/amphion/common/require"
	"github.com/cadmean-ru/amphion/engine"
	"github.com/cadmean-ru/amphion/frontend"
	"github.com/cadmean-ru/amphion/rendering"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net/http"
	"syscall/js"
)

type Frontend struct {
	renderer         *rendering.ARenderer
	handler          frontend.CallbackHandler
	context          frontend.Context
	msgChan          chan frontend.Message
	resManager       *ResourceManager
	rendererDelegate *P5RendererDelegate
}

func (f *Frontend) Init() {
	f.renderer.SetManagementMode(rendering.FrontendManaged)
	f.renderer.Prepare()

	js.Global().Get("document").Set("onmousemove", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		e := args[0]
		x := e.Get("pageX").Int()
		y := e.Get("pageY").Int()
		f.handler(frontend.NewCallback(frontend.CallbackMouseMove, fmt.Sprintf("%d;%d", x, y)))
		return nil
	}))

	js.Global().Get("document").Set("onmousedown", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		e := args[0]
		x := e.Get("pageX").Int()
		y := e.Get("pageY").Int()
		f.handler(frontend.NewCallback(frontend.CallbackMouseDown, fmt.Sprintf("%d;%d", x, y)))
		return nil
	}))

	js.Global().Get("document").Set("onmouseup", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		e := args[0]
		x := e.Get("pageX").Int()
		y := e.Get("pageY").Int()
		f.handler(frontend.NewCallback(frontend.CallbackMouseUp, fmt.Sprintf("%d;%d", x, y)))
		return nil
	}))

	js.Global().Get("document").Set("onkeydown", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		e := args[0]
		key := e.Get("key").String()
		code := e.Get("code").String()

		f.handler(frontend.NewCallback(frontend.CallbackKeyDown, fmt.Sprintf("%s\n%s", key, code)))
		return nil
	}))

	js.Global().Get("addEventListener").Invoke("resize", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		ws := getWindowSize()
		f.rendererDelegate.p5.resizeCanvas(ws.X, ws.Y)
		f.context.ScreenInfo = common.NewScreenInfo(ws.X, ws.Y)
		f.handler(frontend.NewCallback(frontend.CallbackContextChange, ""))
		return nil
	}))

	js.Global().Get("addEventListener").Invoke("blur", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		f.handler(frontend.NewCallback(frontend.CallbackAppHide, ""))
		return nil
	}))

	js.Global().Get("addEventListener").Invoke("focus", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		f.handler(frontend.NewCallback(frontend.CallbackAppShow, ""))
		return nil
	}))

	js.Global().Get("document").Set("onpaste", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		e := args[0]
		e.Call("preventDefault")

		//item := e.Get("clipboardData").Get("items").Index(0)
		//
		//kind := item.Get("kind").String()
		//mime := item.Get("type").String()
		//data := engine.NewClipboardData()

		fmt.Println("pasted")

		return nil
	}))

	js.Global().Set("dumpCurrentScene", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if json, err := engine.GetCurrentScene().DumpToJson(); err == nil {
			return js.Global().Get("JSON").Call("parse", string(json))
		}
		return nil
	}))

	ws := getWindowSize()
	f.context.ScreenInfo = common.NewScreenInfo(ws.X, ws.Y)
}

func (f *Frontend) Run() {
	f.handler(frontend.NewCallback(frontend.CallbackReady, ""))

	for msg := range f.msgChan {
		if msg.Code == frontend.MessageExit {
			break
		}

		f.handleMessage(msg)
	}
	close(f.msgChan)
}

func (f *Frontend) SetCallback(handler frontend.CallbackHandler) {
	f.handler = handler
}

func (f *Frontend) GetRenderer() *rendering.ARenderer {
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

func (f *Frontend) GetApp() *frontend.App {
	resp, err := http.Get("app.yaml")
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		data, err := ioutil.ReadAll(resp.Body)
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

	return nil
}

func (f *Frontend) GetLaunchArgs() a.SiMap {
	args := make(a.SiMap)

	launchArgsJs := js.Global().Get("launchArgs")
	if launchArgsJs.IsUndefined() {
		return args
	}

	pathJs := launchArgsJs.Get("path")
	if !pathJs.IsUndefined() {
		args["path"] = pathJs.String()
	}

	return args
}

func (f *Frontend) handleMessage(msg frontend.Message) {
	switch msg.Code {
	case frontend.MessageRender:
		f.renderer.PerformRendering()
	case frontend.MessageExec:
		if msg.Data != nil {
			if action, ok := msg.Data.(func()); ok {
				action()
			}
		}
	case frontend.MessageTitle:
		setWindowTitle(require.String(msg.Data, "No title"))
	case frontend.MessageNavigate:
		path := require.String(msg.Data)
		if path != "" {
			setWindowLocation(path)
		}
	}
}

func NewFrontend() *Frontend {
	f := &Frontend{
		msgChan:          make(chan frontend.Message, 10),
		resManager:       newResourceManager(),
		rendererDelegate: newP5RendererDelegate(),
	}
	f.renderer = rendering.NewARenderer(f.rendererDelegate)
	f.rendererDelegate.aRenderer = f.renderer
	return f
}

func getWindowSize() a.IntVector2 {
	w := js.Global().Get("innerWidth").Int()
	h := js.Global().Get("innerHeight").Int()

	return a.IntVector2{w, h}
}

func setWindowTitle(title string) {
	js.Global().Get("document").Set("title", title)
}

func setWindowLocation(path string) {
	app := engine.GetInstance().GetCurrentApp()
	if app != nil && app.Debug {
		path += "?connectDebugger=4200"
	}
	js.Global().Get("history").Call("replaceState", js.Value{}, "Loading", path)
}
