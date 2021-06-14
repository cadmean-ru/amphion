package builtin

import (
	"github.com/cadmean-ru/amphion/common/a"
	"github.com/cadmean-ru/amphion/engine"
	"github.com/cadmean-ru/amphion/rendering"
)

// Component for displaying text
type TextView struct {
	engine.ViewImpl
	TextColor  a.Color     `state:"textColor"`
	Font       string      `state:"font"`
	FontSize   byte        `state:"fontSize"`
	FontWeight byte        `state:"fontWeight"`
	Text       string      `state:"text"`
	HTextAlign a.TextAlign `state:"hTextAlign"`
	VTextAlign a.TextAlign `state:"vTextAlign"`
}

func (t *TextView) GetName() string {
	return engine.NameOfComponent(t)
}

func (t *TextView) OnDraw(ctx engine.DrawingContext) {
	pr := rendering.NewTextPrimitive(t.Text)
	pr.Transform = t.SceneObject.Transform.ToRenderingTransform()
	pr.Appearance = rendering.Appearance{
		FillColor:    t.TextColor,
		StrokeWeight: t.FontWeight,
	}
	pr.TextAppearance = rendering.TextAppearance{
		Font:     t.Font,
		FontSize: t.FontSize,
	}
	pr.HTextAlign = t.HTextAlign
	pr.VTextAlign = t.VTextAlign
	ctx.GetRenderingNode().SetPrimitive(t.PrimitiveId, pr)
	t.ShouldRedraw = false
}

// Sets the text equal to the specified value, forcing the view to redraw and requesting rendering.
func (t *TextView) SetText(text string) {
	t.Text = text
	t.ShouldRedraw = true
	engine.RequestRendering()
}

// Returns the current text value.
func (t *TextView) GetText() string {
	return t.Text
}

// Sets the current text color.
func (t *TextView) SetTextColor(color interface{}) {
	switch color.(type) {
	case a.Color:
		t.TextColor = color.(a.Color)
	case string:
		t.TextColor = a.ParseHexColor(color.(string))
	default:
		t.TextColor = a.BlackColor()
	}

	t.ShouldRedraw = true
	engine.RequestRendering()
}

// Sets the current text size.
func (t *TextView) SetFontSize(fontSize byte) {
	t.FontSize = fontSize
	t.ShouldRedraw = true
	engine.RequestRendering()
}

// Sets the current horizontal text alignment.
func (t *TextView) SetHTextAlign(align a.TextAlign) {
	t.HTextAlign = align
	t.ShouldRedraw = true
	engine.RequestRendering()
}

// Sets the current horizontal text alignment.
func (t *TextView) SetVTextAlign(align a.TextAlign) {
	t.VTextAlign = align
	t.ShouldRedraw = true
	engine.RequestRendering()
}

func NewTextView(text string) *TextView {
	return &TextView{
		TextColor:  a.BlackColor(),
		FontSize:   16,
		FontWeight: 0,
		Text:       text,
	}
}
