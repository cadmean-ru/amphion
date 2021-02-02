//+build js

package builtin

import (
	"github.com/cadmean-ru/amphion/common/a"
	"github.com/cadmean-ru/amphion/engine"
	"github.com/cadmean-ru/amphion/native/web"
	"syscall/js"
)

func (n *NativeInputView) onInitWeb(_ engine.InitContext) {
	html := web.CreateHtmlElement("input")
	n.nativeView = html
	html.SetProperty("oninput", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if n.onTextChange != nil {
			n.onTextChange(n.getTextNative())
		}
		return nil
	}))
	n.setTextWeb(n.text)
	n.setHintWeb(n.hint)
}

func (n *NativeInputView) onStartWeb() {
	html := n.nativeView.(*web.HtmlElement)
	web.InstantiateHtml(html)
}

func (n *NativeInputView) onStopWeb() {
	html := n.nativeView.(*web.HtmlElement)
	web.RemoveHtml(html)
}

func (n *NativeInputView) onDrawWeb(_ engine.DrawingContext) {
	html := n.nativeView.(*web.HtmlElement)
	t := n.SceneObject.Transform.ToRenderingTransform()
	html.SetPosition(t.Position)
	html.SetSize(a.NewIntVector2(t.Size.X, t.Size.Y))
}

func (n *NativeInputView) setTextWeb(t string) {
	n.text = t
	html := n.nativeView.(*web.HtmlElement)
	html.SetProperty("value", n.text)
}

func (n *NativeInputView) getTextWeb() string {
	html := n.nativeView.(*web.HtmlElement)
	n.text = html.GetStringProperty("value")
	return n.text
}

func (n *NativeInputView) setHintWeb(h string) {
	html := n.nativeView.(*web.HtmlElement)
	html.SetProperty("placeholder", h)
}

func (n *NativeInputView) getHintWeb() string {
	html := n.nativeView.(*web.HtmlElement)
	return html.GetStringProperty("placeholder")
}