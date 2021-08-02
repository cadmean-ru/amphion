package builtin

import (
	"github.com/cadmean-ru/amphion/common/a"
	"github.com/cadmean-ru/amphion/engine"
	"github.com/cadmean-ru/amphion/rendering"
)

// TextView component displays the given text
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

// SetText sets the text equal to the specified value, forcing the view to redraw and requesting rendering.
func (t *TextView) SetText(text string) {
	t.Text = text
	t.ShouldRedraw = true
	engine.RequestRendering()
}

// GetText returns the current text value.
func (t *TextView) GetText() string {
	return t.Text
}

// SetTextColor sets the current text color.
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

// SetFontSize sets the current text size.
func (t *TextView) SetFontSize(fontSize byte) {
	t.FontSize = fontSize
	t.ShouldRedraw = true
	engine.RequestRendering()
}

// SetHTextAlign sets the current horizontal text alignment.
func (t *TextView) SetHTextAlign(align a.TextAlign) {
	t.HTextAlign = align
	t.ShouldRedraw = true
	engine.RequestRendering()
}

// SetVTextAlign sets the current horizontal text alignment.
func (t *TextView) SetVTextAlign(align a.TextAlign) {
	t.VTextAlign = align
	t.ShouldRedraw = true
	engine.RequestRendering()
}

//NewTextView creates a new TextView with the given text.
func NewTextView(text string) *TextView {
	return &TextView{
		TextColor:  a.BlackColor(),
		FontSize:   16,
		FontWeight: 0,
		Text:       text,
	}
}
