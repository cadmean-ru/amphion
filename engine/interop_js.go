// +build js

package engine

import (
	"github.com/cadmean-ru/amphion/common"
	"syscall/js"
)

func prepareInterop() {
	js.Global().Set("frontEndCallback", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		e := args[0]
		code := e.Get("code").Int()
		jsEvent := e.Get("data").String()
		event := frontEndCallback{
			code: code,
			data: jsEvent,
		}
		instance.handleFrontEndCallback(event)
		return nil
	}))

	js.Global().Set("frontEndInterrupt", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		msg := args[0].String()
		instance.handleFrontEndInterrupt(msg)
		return nil
	}))
}

func getGlobalContext() *GlobalContext {
	contextJs := js.Global().Get("getGlobalContext").Invoke()
	contextNative := newGlobalContext()

	screenJs := contextJs.Get("screen")
	screenW := screenJs.Get("width").Int()
	screenH := screenJs.Get("height").Int()

	contextNative.screenInfo = common.newScreenInfo(screenW, screenH)

	host := contextJs.Get("host").String()
	contextNative.domain = host
	contextNative.host = host

	port := contextJs.Get("port").String()
	contextNative.port = port

	return contextNative
}

func frontEndCloseScene() {
	js.Global().Get("closeScene").Invoke()
}

func commencePanic(reason, message string) {
	js.Global().Get("commencePanic").Invoke(reason, message)
}

func loadAppData() ([]byte, error) {
	json := js.Global().Get("getAppData").Invoke().String()
	return []byte(json), nil
}

