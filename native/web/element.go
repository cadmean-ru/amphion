//+build js

package web

import (
	"github.com/cadmean-ru/amphion/common/a"
	"syscall/js"
)

// HtmlElement used create and manage HTML elements in DOM.
type HtmlElement struct {
	jsValue js.Value
	tagName string
	style   js.Value
}

func (e *HtmlElement) GetTagName() string {
	return e.tagName
}

func (e *HtmlElement) GetId() string {
	return e.jsValue.Get("id").String()
}

func (e *HtmlElement) SetId(id string) {
	e.jsValue.Set("id", id)
}

func (e *HtmlElement) GetClass() string {
	return e.jsValue.Get("className").String()
}

func (e *HtmlElement) SetClass(class string) {
	e.jsValue.Set("className", class)
}

func (e *HtmlElement) GetStyles() map[string]string {
	return nil
}

func (e *HtmlElement) SetStyle(styleName, value string) {
	e.style.Set(styleName, value)
}

func (e *HtmlElement) GetPosition() a.IntVector3 {
	x := e.style.Get("left").Int()
	y := e.style.Get("top").Int()
	z := e.style.Get("z-index").Int()
	return a.NewIntVector3(x, y, z)
}

func (e *HtmlElement) SetPosition(pos a.IntVector3) {
	e.style.Set("left", pos.X)
	e.style.Set("top", pos.Y)
	e.style.Set("z-index", pos.Z)
}

func (e *HtmlElement) GetSize() a.IntVector2 {
	w := e.style.Get("width").Int()
	h := e.style.Get("height").Int()
	return a.IntVector2{X: w, Y: h}
}

func (e *HtmlElement) SetSize(size a.IntVector2) {
	e.style.Set("width", size.X)
	e.style.Set("height", size.Y)
}

func (e *HtmlElement) SetEventListener(eventName string, listener func(args ...interface{})) {
	e.jsValue.Set(eventName, js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		listener(args)
		return nil
	}))
}

func (e *HtmlElement) SetProperty(property string, value interface{}) {
	e.jsValue.Set(property, value)
}

func (e *HtmlElement) GetStringProperty(property string) string {
	return e.jsValue.Get(property).String()
}

func (e *HtmlElement) GetIntProperty(property string) int {
	return e.jsValue.Get(property).Int()
}

func (e *HtmlElement) GetFloatProperty(property string) float64 {
	return e.jsValue.Get(property).Float()
}

func CreateHtmlElement(tagName string) *HtmlElement {
	value := js.Global().Get("document").Call("createElement", tagName)
	style := value.Get("style")
	el := &HtmlElement{
		jsValue: value,
		style:   style,
	}
	el.SetStyle("position", "absolute")
	return el
}

func InstantiateHtml(el *HtmlElement) {
	js.Global().Get("document").Get("body").Call("appendChild", el.jsValue)
}

func RemoveHtml(el *HtmlElement) {
	el.jsValue.Call("remove")
}