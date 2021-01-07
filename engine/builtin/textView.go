package builtin

import (
	"github.com/cadmean-ru/amphion/common/a"
	"github.com/cadmean-ru/amphion/engine"
	"github.com/cadmean-ru/amphion/rendering"
)

// Component for displaying text
type TextView struct {
	ViewImpl
	TextColor  a.Color  `state:"textColor"`
	Font       string   `state:"font"`
	FontSize   byte     `state:"fontSize"`
	FontWeight byte     `state:"fontWeight"`
	Text       string   `state:"text"`
}

func (t *TextView) GetName() string {
	return engine.NameOfComponent(t)
}

func (t *TextView) OnDraw(ctx engine.DrawingContext) {
	pr := rendering.NewTextPrimitive(t.Text)
	pr.Transform = transformToRenderingTransform(t.obj.Transform)
	pr.Appearance = rendering.Appearance{
		FillColor:    t.TextColor,
		StrokeWeight: t.FontWeight,
	}
	pr.TextAppearance = rendering.TextAppearance{
		Font:     t.Font,
		FontSize: t.FontSize,
	}
	ctx.GetRenderer().SetPrimitive(t.pId, pr, t.ShouldRedraw())
	t.redraw = false
}

func (t *TextView) SetText(text string) {
	t.Text = text
	t.redraw = true
}

func (t *TextView) GetText() string {
	return t.Text
}

func NewTextView(text string) *TextView {
	return &TextView{
		TextColor:  a.BlackColor(),
		FontSize:   16,
		FontWeight: 0,
		Text:       text,
	}
}
