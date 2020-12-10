package builtin

import (
	"github.com/cadmean-ru/amphion/common"
	"github.com/cadmean-ru/amphion/engine"
	"github.com/cadmean-ru/amphion/rendering"
)

type TextView struct {
	ViewImpl
	Appearance     rendering.Appearance
	TextAppearance rendering.TextAppearance
	text           common.AString
}

func (t *TextView) GetName() string {
	return "TextView"
}

func (t *TextView) OnDraw(ctx engine.DrawingContext) {
	pr := rendering.NewTextPrimitive(t.text)
	pr.Transform = transformToRenderingTransform(t.obj.Transform)
	pr.Appearance = t.Appearance
	pr.TextAppearance = t.TextAppearance
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
		text:    common.AString(text),
	}
}
