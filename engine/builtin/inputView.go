package builtin

import (
	"github.com/cadmean-ru/amphion/engine"
	"strings"
)

type InputView struct {
	TextView
	AllowMultiline bool
	cursorPrId     int64
	cursorPos      int
}

func (v *InputView) OnStart() {
	v.TextView.OnStart()
	//v.cursorPrId = v.eng.GetRenderer().AddPrimitive()
	v.eng.BindEventHandler(engine.EventKeyDown, v.handleKeyDown)
}

func (v *InputView) handleKeyDown(event engine.AmphionEvent) bool {
	if !v.obj.IsFocused() {
		return true
	}

	keyEvent := event.Data.(engine.KeyEvent)

	if isPrintableKeyCode(keyEvent.Code) {
		v.SetText(string(v.text) + keyEvent.Key)
		v.ForceRedraw()
	} else if v.AllowMultiline && keyEvent.Code == "Enter" {
		v.SetText(string(v.text) + "\n")
		v.ForceRedraw()
	} else if keyEvent.Code == "Backspace" {
		v.SetText(string(v.text)[:len(v.text)-1])
		v.ForceRedraw()
	}

	if v.redraw {
		v.eng.RequestRendering()
	}

	v.cursorPos = len(v.text)

	return true
}

func isPrintableKeyCode(code string) bool {
	return strings.HasPrefix(code, "Key") ||
		strings.HasPrefix(code, "Numpad") ||
		strings.HasPrefix(code, "Digit") ||
		strings.HasPrefix(code, "Bracket") ||
		code == "Space" || code == "Period" || code == "Comma" ||
		code == "Slash" || code == "Backslash" ||
		code == "Backquote"
}

//func (v *InputView) OnDraw(ctx engine.DrawingContext) {
//	line := rendering.NewGeometryPrimitive(rendering.PrimitiveLine)
//	line.Transform.Size = common.NewIntVector3(0, int(v.TextAppearance.FontSize), 0)
//	line.Appearance.StrokeColor = common.BlackColor()
//	line.Appearance.StrokeWeight = 1
//	p := v.obj.Transform.GetGlobalTopLeftPosition().Round()
//	p = p.Add(common.NewIntVector3(1, 3, 0).Multiply(common.NewIntVector3(v.cursorPos * 10, 1, 1)))
//	line.Transform.Position = p
//	ctx.GetRenderer().SetPrimitive(v.cursorPrId, line, v.redraw || v.eng.IsForcedToRedraw())
//	v.TextView.OnDraw(ctx)
//}

func (v *InputView) OnStop() {
	v.TextView.OnStop()
	v.eng.UnbindEventHandler(engine.EventKeyDown, v.handleKeyDown)
}

func (v *InputView) GetName() string {
	return engine.NameOfComponent(v)
}

func NewInputView() *InputView {
	return &InputView{
		TextView: TextView{
			text: "\n",
		},
	}
}