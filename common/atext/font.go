package atext

import (
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
)

type Font struct {
	name string
	ttf  *truetype.Font
}

func (f *Font) GetName() string {
	return f.name
}

func (f *Font) NewFace(size int) *Face {
	face := truetype.NewFace(f.ttf, &truetype.Options{
		Size:    float64(size),
		DPI:     72,
		Hinting: font.HintingFull,
	})

	metrics := face.Metrics()

	return &Face{
		font:       f,
		face:       face,
		size:       size,
		capHeight:  int(metrics.CapHeight) >> 6,
		lineHeight: int(metrics.Height) >> 6,
		xHeight:    int(metrics.XHeight) >> 6,
		ascent:     int(metrics.Ascent) >> 6,
		descent:    int(metrics.Descent) >> 6,
	}
}