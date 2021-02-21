package atext

import (
	_ "embed"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
)

// Font represents a true type font.
type Font struct {
	name string
	ttf  *truetype.Font
}

// Get Name returns the name of the font.
func (f *Font) GetName() string {
	return f.name
}

// NewFace creates a new font face with the specifies size.
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
		glyphs:     make(map[rune]*Glyph),
	}
}

// ParseFont parses a true type font from the data.
func ParseFont(fontData []byte) (*Font, error) {
	f, err := truetype.Parse(fontData)
	if err != nil {
		return nil, err
	}

	return &Font{
		name: f.Name(1),
		ttf:  f,
	}, nil
}

//go:embed Roboto.ttf
var DefaultFontData []byte