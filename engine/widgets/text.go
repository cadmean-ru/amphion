package widgets

import (
	"github.com/cadmean-ru/amphion/common/a"
	"github.com/cadmean-ru/amphion/engine"
	"github.com/cadmean-ru/amphion/engine/builtin"
)

type TextOptions struct {
	Options
	Text       string
	TextColor  a.Color
	Font       string
	FontSize   byte
	FontWeight byte
	HTextAlign a.TextAlign
	VTextAlign a.TextAlign
	SingleLine bool
}

func Text(options TextOptions) (*engine.SceneObject, *builtin.TextView) {
	obj := engine.NewSceneObject("text")
	text := builtin.NewTextView(options.Text)
	text.TextColor = options.TextColor
	text.Font = options.Font
	text.FontSize = options.FontSize
	text.FontWeight = options.FontWeight
	text.HTextAlign = options.HTextAlign
	text.VTextAlign = options.VTextAlign
	text.SingleLine = options.SingleLine
	obj.AddComponent(text)
	return obj, text
}
