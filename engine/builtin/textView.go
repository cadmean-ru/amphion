package builtin

import (
	"github.com/cadmean-ru/amphion/common"
	"github.com/cadmean-ru/amphion/engine"
	"github.com/cadmean-ru/amphion/rendering"
)

type TextView struct {
	ViewImpl
	TextColor common.Color   `state:"TextColor"`
	Font      string         `state:"Font"`
	FontSize  common.AByte   `state:"FontSize"`
	text      common.AString `state:"Text"`
}

func (t *TextView) GetName() string {
	return "TextView"
}

func (t *TextView) OnDraw(ctx engine.DrawingContext) {
	pr := rendering.NewTextPrimitive(t.text)
	pr.Transform = transformToRenderingTransform(t.obj.Transform)
	pr.Appearance = rendering.Appearance{
		FillColor: t.TextColor,
	}
	pr.TextAppearance = rendering.TextAppearance{
		Font:     t.Font,
		FontSize: t.FontSize,
	}
	ctx.GetRenderer().SetPrimitive(t.pId, pr, t.ShouldRedraw())
	t.redraw = false
}

func (t *TextView) SetText(text string) {
	t.text = common.AString(text)
	t.redraw = true
}

func (t *TextView) GetText() string {
	return string(t.text)
}

func NewTextView(text string) *TextView {
	return &TextView{
		TextColor: common.BlackColor(),
		FontSize:  16,
		text:      common.AString(text),
	}
}
