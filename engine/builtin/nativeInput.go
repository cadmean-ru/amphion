package builtin

import (
	"github.com/cadmean-ru/amphion/engine"
)

// NativeInputView component displays a native input widget for current platform.
type NativeInputView struct {
	engine.ComponentImpl
	onTextChange  func(text string)
	onInitNative  func(ctx engine.InitContext)
	onStartNative func()
	onStopNative  func()
	onDrawNative  func(ctx engine.DrawingContext)
	setTextNative func(t string)
	getTextNative func() string
	setHintNative func(h string)
	getHintNative func() string
	text          string
	hint          string
	nativeView    interface{}
}

func (n *NativeInputView) OnInit(ctx engine.InitContext) {
	n.ComponentImpl.OnInit(ctx)
	n.onInitNative = n.onInitWeb
	n.onStartNative = n.onStartWeb
	n.onStopNative = n.onStopWeb
	n.onDrawNative = n.onDrawWeb
	n.setTextNative = n.setTextWeb
	n.getTextNative = n.getTextWeb
	n.setHintNative = n.setHintWeb
	n.getHintNative = n.getHintWeb
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

// GetText returns the current text value of the input view.
func (n *NativeInputView) GetText() string {
	return n.text
}

// SetText updates the text.
func (n *NativeInputView) SetText(text string) {
	n.setTextNative(text)
	n.Engine.RequestRendering()
}

// SetOnTextChangeListener sets the callback that is invoked when the text of the input view is changed.
func (n *NativeInputView) SetOnTextChangeListener(listener func(text string)) {
	n.onTextChange = listener
}

// SetHint updates the hint text of the input view.
func (n *NativeInputView) SetHint(hint string) {
	n.hint = hint
	n.setHintNative(hint)
	n.Engine.RequestRendering()
}

// GetHint returns the current hint value of the input view.
func (n *NativeInputView) GetHint() string {
	n.hint = n.getHintNative()
	return n.hint
}

// NewNativeInputView creates a new NativeInputView. Returns pointer to the instance.
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