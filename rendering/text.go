package rendering

import "github.com/cadmean-ru/amphion/common/a"

const textPrimitiveMinBytesSize = 1 + transformBytesSize + appearanceBytesSize + textAppearanceBytesSize + 4

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

func NewTextPrimitive(text string) *TextPrimitive {
	return &TextPrimitive{
		Transform:      NewTransform(),
		Appearance:     DefaultAppearance(),
		TextAppearance: DefaultTextAppearance(),
		Text:           text,
	}
}
