package builtin

import (
	"github.com/cadmean-ru/amphion/common/a"
	"github.com/cadmean-ru/amphion/engine"
	"github.com/cadmean-ru/amphion/native"
	"github.com/cadmean-ru/amphion/native/web"
	"syscall/js"
)

// Component displays a native input widget for current platform.
type NativeInputView struct {
	engine.ComponentImpl
	initFeature                    *nativeInputViewInit
	startFeature                   *nativeInputViewStart
	stopFeature                    *nativeInputViewStop
	drawFeature                    *nativeInputViewDraw
	getTextFeature                 *nativeInputViewGetText
	setTextFeature                 *nativeInputViewSetText
	onTextChange                   func(text string)
}

func (n *NativeInputView) OnInit(ctx engine.InitContext) {
	n.ComponentImpl.OnInit(ctx)
	n.initFeature = &nativeInputViewInit{view: n}
	native.Invoke(n.initFeature)
	n.startFeature = &nativeInputViewStart{init: n.initFeature}
	n.stopFeature = &nativeInputViewStop{init: n.initFeature}
	n.drawFeature = &nativeInputViewDraw{
		init: n.initFeature,
		obj:  n.SceneObject,
	}
	n.getTextFeature = &nativeInputViewGetText{init: n.initFeature}
	n.setTextFeature = &nativeInputViewSetText{init: n.initFeature}
}

func (n *NativeInputView) OnStart() {
	native.Invoke(n.startFeature)
}

func (n *NativeInputView) OnDraw(_ engine.DrawingContext) {
	native.Invoke(n.drawFeature)
}

func (n *NativeInputView) ForceRedraw() {

}

func (n *NativeInputView) OnStop() {
	native.Invoke(n.stopFeature)
}

func (n *NativeInputView) GetName() string {
	return engine.NameOfComponent(n)
}

func (n *NativeInputView) GetText() string {
	native.Invoke(n.getTextFeature)
	return n.getTextFeature.text
}

func (n *NativeInputView) SetText(text string) {
	n.setTextFeature.text = text
	native.Invoke(n.setTextFeature)
	n.Engine.RequestRendering()
}

func (n *NativeInputView) SetOnTextChangeListener(listener func(text string)) {
	n.onTextChange = listener
}

// Init
type nativeInputViewInit struct {
	native.FeatureImpl
	html *web.HtmlElement
	view *NativeInputView
}

func (n *nativeInputViewInit) OnWeb() {
	n.html = web.CreateHtmlElement("input")
	n.html.SetProperty("oninput", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if n.view.onTextChange != nil {
			native.Invoke(n.view.getTextFeature)
			n.view.onTextChange(n.view.getTextFeature.text)
		}
		return nil
	}))
}

// Start
type nativeInputViewStart struct {
	native.FeatureImpl
	init *nativeInputViewInit
}

func (n *nativeInputViewStart) OnWeb() {
	web.InstantiateHtml(n.init.html)
}

// Stop
type nativeInputViewStop struct {
	native.FeatureImpl
	init *nativeInputViewInit
}

func (n *nativeInputViewStop) OnWeb() {
	web.RemoveHtml(n.init.html)
}

// Draw
type nativeInputViewDraw struct {
	native.FeatureImpl
	init *nativeInputViewInit
	obj  *engine.SceneObject
}

func (n *nativeInputViewDraw) OnWeb() {
	t := transformToRenderingTransform(n.obj.Transform)
	n.init.html.SetPosition(t.Position)
	n.init.html.SetSize(a.NewIntVector2(t.Size.X, t.Size.Y))
}

// Set text
type nativeInputViewSetText struct {
	native.FeatureImpl
	init *nativeInputViewInit
	text string
}

func (n *nativeInputViewSetText) OnWeb() {
	n.init.html.SetProperty("value", n.text)
}

// Get text
type nativeInputViewGetText struct {
	native.FeatureImpl
	init *nativeInputViewInit
	text string
}

func (n *nativeInputViewGetText) OnWeb() {
	n.text = n.init.html.GetStringProperty("value")
}

// Creates a new NativeInputView. Returns pointer to the instance.
func NewNativeInputView() *NativeInputView {
	return &NativeInputView{}
}