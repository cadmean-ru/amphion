package rendering

import "github.com/cadmean-ru/amphion/common/a"

type TextPrimitive struct {
	Transform      Transform
	Appearance     Appearance
	TextAppearance TextAppearance
	Text           string
	HTextAlign     a.TextAlign
	VTextAlign     a.TextAlign
}

func (p *TextPrimitive) GetType() byte {
	return PrimitiveText
}

func (p *TextPrimitive) GetTransform() Transform {
	return p.Transform
}

func (p *TextPrimitive) SetTransform(t Transform) {
	p.Transform = t
}

func NewTextPrimitive(text string) *TextPrimitive {
	return &TextPrimitive{
		Transform:      NewTransform(),
		Appearance:     DefaultAppearance(),
		TextAppearance: DefaultTextAppearance(),
		Text:           text,
	}
}
