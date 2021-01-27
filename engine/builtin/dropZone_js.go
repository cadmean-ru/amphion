// +build js

package builtin

import (
	"github.com/cadmean-ru/amphion/engine"
	"syscall/js"
)

type FileDropZone struct {
	engine.ComponentImpl
	div    js.Value
	style  js.Value
	OnDrop engine.EventHandler `state:"onDrop"`
}

func (z *FileDropZone) OnInit(ctx engine.InitContext) {
	z.ComponentImpl.OnInit(ctx)
	document := js.Global().Get("document")
	z.div = document.Call("createElement", "div")
	z.style = z.div.Get("style")
	z.style.Set("position", "absolute")
	z.div.Set("ondrop", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		ev := args[0]
		ev.Call("preventDefault")

		files := ev.Get("dataTransfer").Get("files")
		file := files.Index(0)

		reader := js.Global().Get("FileReader").New()
		reader.Call("readAsArrayBuffer", file)

		reader.Set("onload", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			buffer := reader.Get("result")
			bytesLength := buffer.Get("byteLength").Int()
			uint8array := js.Global().Get("Uint8Array").New(buffer)
			bytes := make([]byte, bytesLength)
			js.CopyBytesToGo(bytes, uint8array)

			data := engine.InputFileData{
				Name: file.Get("name").String(),
				Data: bytes,
				Mime: file.Get("type").String(),
			}

			if z.OnDrop != nil {
				z.OnDrop(engine.NewAmphionEvent(z, engine.EventDropFile, data))
			}
			return nil
		}))

		return nil
	}))
	z.div.Set("ondragover", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		ev := args[0]
		ev.Call("preventDefault")
		return nil
	}))
	z.div.Set("id", "test")
	document.Get("body").Call("appendChild", z.div)
}

func (z *FileDropZone) OnStart() {
	z.ComponentImpl.OnStart()
	z.style.Set("display", "block")
}

func (z *FileDropZone) OnDraw(_ engine.DrawingContext) {
	t := z.SceneObject.Transform.ToRenderingTransform()
	z.style.Set("left", t.Position.X)
	z.style.Set("top", t.Position.Y)
	z.style.Set("width", t.Size.X)
	z.style.Set("height", t.Size.Y)
}

func (z *FileDropZone) OnStop() {
	z.ComponentImpl.OnStop()
	z.style.Set("display", "none")
}

func (z *FileDropZone) ForceRedraw() {
	z.Engine.GetMessageDispatcher().DispatchDown(z.SceneObject, engine.NewMessage(z, engine.MessageRedraw, nil), engine.MessageMaxDepth)
}

func (z *FileDropZone) GetName() string {
	return engine.NameOfComponent(z)
}

func NewFileDropZone(onDrop engine.EventHandler) *FileDropZone {
	return &FileDropZone{
		OnDrop: onDrop,
	}
}