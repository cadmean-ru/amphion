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

//func (ap Appearance) EncodeToByteArray() []byte {
//	arr := make([]byte, appearanceBytesSize)
//	_ = a.CopyByteArray(ap.FillColor.EncodeToByteArray(), arr, 0, 4)
//	_ = a.CopyByteArray(ap.StrokeColor.EncodeToByteArray(), arr, 4, 4)
//	_ = a.CopyByteArray(ap.StrokeWeight.EncodeToByteArray(), arr, 8, 1)
//	return arr
//}

func DefaultAppearance() Appearance {
	return Appearance{
		FillColor:    a.WhiteColor(),
		StrokeColor:  a.BlackColor(),
		StrokeWeight: 1,
	}
}


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
