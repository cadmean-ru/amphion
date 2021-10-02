package rendering

import (
	"github.com/cadmean-ru/amphion/common/a"
)

const appearanceBytesSize = 9
const textAppearanceBytesSize = 1

type Appearance struct {
	FillColor    a.Color
	StrokeColor  a.Color
	StrokeWeight byte
	CornerRadius byte
}

func (ap Appearance) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"fillColor":    ap.FillColor.ToMap(),
		"strokeColor":  ap.StrokeColor.ToMap(),
		"strokeWeight": ap.StrokeWeight,
	}
}

func DefaultAppearance() Appearance {
	return Appearance{
		FillColor:    a.White(),
		StrokeColor:  a.Black(),
		StrokeWeight: 1,
	}
}

//TextAppearance
//Deprecated
type TextAppearance struct {
	Font     string
	FontSize byte
}

func (a TextAppearance) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"font":      a.Font,
		"fontSize":  a.FontSize,
	}
}

func DefaultTextAppearance() TextAppearance {
	return TextAppearance{
		Font:      "",
		FontSize:  14,
	}
}

//func (a TextAppearance) EncodeToByteArray() []byte {
//	arr := make([]byte, textAppearanceBytesSize)
//	arr[0] = byte(a.FontSize)
//	return arr
//}
