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
	initFeature   *nativeInputViewInit
	onTextChange  func(text string)
	onInitNative  func(ctx engine.InitContext)
	onStartNative func()
	onStopNative  func()
	onDrawNative  func(ctx engine.DrawingContext)
	setTextNative func(t string)
	getTextNative func() string
	setHintNative func(h string)
	getHintNative func() string
	html          *web.HtmlElement
	text          string
	hint          string
}

func (n *NativeInputView) OnInit(ctx engine.InitContext) {
	n.ComponentImpl.OnInit(ctx)
	n.initFeature = &nativeInputViewInit{view: n}
	native.Invoke(n.initFeature)
	n.onInitNative(ctx)
}

func (n *NativeInputView) OnStart() {
	n.onStartNative()
}

func (n *NativeInputView) OnDraw(ctx engine.DrawingContext) {
	n.onDrawNative(ctx)
}

func (n *NativeInputView) ForceRedraw() {

}

func (n *NativeInputView) OnStop() {
	n.onStopNative()
}

func (n *NativeInputView) GetName() string {
	return engine.NameOfComponent(n)
}

// Returns the current text value of the input view.
func (n *NativeInputView) GetText() string {
	return n.text
}

// Updates the text.
func (n *NativeInputView) SetText(text string) {
	n.setTextNative(text)
	n.Engine.RequestRendering()
}

// Sets the callback that is invoked when the text of the input view is changed.
func (n *NativeInputView) SetOnTextChangeListener(listener func(text string)) {
	n.onTextChange = listener
}

// Updates the hint text of the input view.
func (n *NativeInputView) SetHint(hint string) {
	n.hint = hint
	n.setHintNative(hint)
	n.Engine.RequestRendering()
}

// Returns the current hint value of the input view.
func (n *NativeInputView) GetHint() string {
	n.hint = n.getHintNative()
	return n.hint
}

// Init
type nativeInputViewInit struct {
	native.FeatureImpl
	html *web.HtmlElement
	view *NativeInputView
}

func (n *nativeInputViewInit) OnWeb() {
	n.view.onInitNative = n.view.onInitWeb
	n.view.onStartNative = n.view.onStartWeb
	n.view.onStopNative = n.view.onStopWeb
	n.view.onDrawNative = n.view.onDrawWeb
	n.view.setTextNative = n.view.setTextWeb
	n.view.getTextNative = n.view.getTextWeb
	n.view.setHintNative = n.view.setHintWeb
	n.view.getHintNative = n.view.getHintWeb
}

func (n *NativeInputView) onInitWeb(_ engine.InitContext) {
	n.html = web.CreateHtmlElement("input")
	n.html.SetProperty("oninput", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if n.onTextChange != nil {
			n.onTextChange(n.getTextNative())
		}
		return nil
	}))
	n.setTextWeb(n.text)
	n.setHintWeb(n.hint)
}

func (n *NativeInputView) onStartWeb() {
	web.InstantiateHtml(n.html)
}

func (n *NativeInputView) onStopWeb() {
	web.RemoveHtml(n.html)
}

func (n *NativeInputView) onDrawWeb(_ engine.DrawingContext) {
	t := transformToRenderingTransform(n.SceneObject.Transform)
	n.html.SetPosition(t.Position)
	n.html.SetSize(a.NewIntVector2(t.Size.X, t.Size.Y))
}

func (n *NativeInputView) setTextWeb(t string) {
	n.text = t
	n.html.SetProperty("value", n.text)
}

func (n *NativeInputView) getTextWeb() string {
	n.text = n.html.GetStringProperty("value")
	return n.text
}

func (n *NativeInputView) setHintWeb(h string) {
	n.html.SetProperty("placeholder", h)
}

func (n *NativeInputView) getHintWeb() string {
	return n.html.GetStringProperty("placeholder")
}

// Creates a new NativeInputView. Returns pointer to the instance.
// This function takes a set of parameters to initialize the input view. All of them are optional.
// The parameters are in the following order:
// 0 - initial text
// 1 - initial hint
// Further values are ignored.
func NewNativeInputView(values ...string) *NativeInputView {
	var initText, initHint string

	if len(values) > 0 {
		initText = values[0]
	}
	if len(values) > 1 {
		initHint = values[1]
	}

	return &NativeInputView{
		text: initText,
		hint: initHint,
	}
}