package builtin

import (
	"github.com/cadmean-ru/amphion/common/a"
	"github.com/cadmean-ru/amphion/common/atext"
	"github.com/cadmean-ru/amphion/engine"
	"github.com/cadmean-ru/amphion/rendering"
)

// TextView component displays the given text
type TextView struct {
	engine.ViewImpl
	TextColor     a.Color     `state:"textColor"`
	Font          string      `state:"font"`
	FontSize      byte        `state:"fontSize"`
	FontWeight    byte        `state:"fontWeight"`
	Text          string      `state:"text"`
	HTextAlign    a.TextAlign `state:"hTextAlign"`
	VTextAlign    a.TextAlign `state:"vTextAlign"`
	SingleLine    bool        `state:"singleLine"`
	prevTransform engine.Transform
	aText         *atext.Text
	aFont         *atext.Font
	aFace         *atext.Face
}

func (t *TextView) OnInit(ctx engine.InitContext) {
    t.ViewImpl.OnInit(ctx)

    t.aFont, _ = atext.ParseFont(atext.DefaultFontData)
    t.aFace = t.aFont.NewFace(int(t.FontSize))

    t.layoutText()
}

func (t *TextView) OnUpdate(_ engine.UpdateContext) {
	if !t.ShouldDraw() && t.prevTransform.Equals(t.SceneObject.Transform) {
		return
	}

	t.prevTransform = t.SceneObject.Transform

	t.layoutText()
}

func (t *TextView) OnDraw(ctx engine.DrawingContext) {
	pr := rendering.NewTextPrimitive(t.Text, t)
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

func (t *TextView) MeasureContents() a.Vector3 {
	return t.aText.GetSize().ToFloat3()
}

func (t *TextView) layoutText() {
	bounds := t.SceneObject.Transform.GlobalRect()
	wantedSize := t.SceneObject.Transform.WantedSize()
	if wantedSize.X == a.WrapContent {
		bounds.X.Max = atext.Unbounded
	}
	if wantedSize.Y == a.WrapContent {
		bounds.Y.Max = atext.Unbounded
	}
	t.aText = atext.LayoutRunes(t.aFace, []rune(t.Text), bounds, atext.LayoutOptions{
		VTextAlign: t.VTextAlign,
		HTextAlign: t.HTextAlign,
		SingleLine: t.SingleLine,
	})
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
	t.aFace = t.aFont.NewFace(int(fontSize))
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

//SetSingleLine sets if the text should be on single line.
func (t *TextView) SetSingleLine(singleLine bool) {
	t.SingleLine = singleLine
	t.ShouldRedraw = true
	engine.RequestRendering()
}

func (t *TextView) GetAText() *atext.Text {
	return t.aText
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
