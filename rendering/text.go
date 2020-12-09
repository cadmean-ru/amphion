package rendering

import "github.com/cadmean-ru/amphion/common"

const textPrimitiveMinBytesSize = 1 + transformBytesSize + appearanceBytesSize + textAppearanceBytesSize + 4

type TextPrimitive struct {
	Transform      Transform
	Appearance     Appearance
	Text           common.AString
	TextAppearance TextAppearance
}

func (p *TextPrimitive) BuildPrimitive() *Primitive {
	pr := NewPrimitive(PrimitiveText)
	pr.AddAttribute(NewAttribute(AttributeTransform, p.Transform))
	pr.AddAttribute(NewAttribute(AttributeFillColor, p.Appearance.FillColor))
	pr.AddAttribute(NewAttribute(AttributeStrokeColor, p.Appearance.StrokeColor))
	pr.AddAttribute(NewAttribute(AttributeStrokeWeight, p.Appearance.StrokeWeight))
	pr.AddAttribute(NewAttribute(AttributeText, p.Text))
	pr.AddAttribute(NewAttribute(AttributeFontSize, p.TextAppearance.FontSize))
	return pr
}

func NewTextPrimitive(text common.AString) *TextPrimitive {
	return &TextPrimitive{
		Transform:      NewTransform(),
		Appearance:     DefaultAppearance(),
		TextAppearance: DefaultTextAppearance(),
		Text:           text,
	}
}
