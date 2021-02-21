package atext

import (
	"github.com/cadmean-ru/amphion/common/a"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
	"image"
	"image/draw"
)

// Face represents a font face, i.e. a font variant of the specific size.
type Face struct {
	font       *Font
	face       font.Face
	size       int
	capHeight  int
	lineHeight int
	xHeight    int
	ascent     int
	descent    int
	glyphs     map[rune]*Glyph
}

//GetFont returns the Font, that this face is associates with.
func (f *Face) GetFont() *Font {
	return f.font
}

//GetSize returns the size of this font face.
func (f *Face) GetSize() int {
	return f.size
}

//GetCapHeight returns the height of capital letter in this face.
func (f *Face) GetCapHeight() int {
	return f.capHeight
}

//GetLineHeight returns the line height of this face.
func (f *Face) GetLineHeight() int {
	return f.lineHeight
}

//GetXHeight returns the xHeight of this face.
func (f *Face) GetXHeight() int {
	return f.xHeight
}

//GetAscent returns the ascent of this face.
func (f *Face) GetAscent() int {
	return f.ascent
}

//GetDescent returns the descent of this face.
func (f *Face) GetDescent() int {
	return f.descent
}

//GetKerning returns the kerning between the two given runes.
func (f *Face) GetKerning(r1, r2 rune) int {
	return int(f.face.Kern(r1, r2)) >> 6
}

//GetGlyph creates a glyph i.e. a visual representation for the given rune.
func (f *Face) GetGlyph(r rune) *Glyph {
	if g, ok := f.glyphs[r]; ok {
		return g
	}

	bounds, adv, ok := f.face.GlyphBounds(r)

	if !ok {
		return nil
	}

	width := int(bounds.Max.X - bounds.Min.X) >> 6
	height := int(bounds.Max.Y - bounds.Min.Y) >> 6

	//if glyph has no dimensions set to a max value
	if width == 0 || height == 0 {
		bounds = f.font.ttf.Bounds(fixed.Int26_6(f.size))
		width = int((bounds.Max.X - bounds.Min.X) >> 6)
		height = int((bounds.Max.Y - bounds.Min.Y) >> 6)

		//above can sometimes yield 0 for font smaller than 48pt, 1 is minimum
		if width == 0 || height == 0 {
			width = 1
			height = 1
		}
	}

	ascent := int(-bounds.Min.Y) >> 6
	descent := int(bounds.Max.Y) >> 6

	dr, mask, maskp, _, ok := f.face.Glyph(fixed.Point26_6{}, r)

	if dr.Size().X == 0 || dr.Size().Y == 0 {
		return nil
	}

	img := image.NewGray(dr)
	draw.Draw(img, dr, mask, maskp, draw.Src)

	g := Glyph{
		face:    f,
		rune:    r,
		pixels:  img.Pix,
		size:    a.NewIntVector2(width, height),
		bearing: a.NewIntVector2(int(bounds.Min.X) >> 6, descent),
		advance: int(adv) >> 6,
		ascent:  ascent,
		descent: descent,
	}

	f.glyphs[r] = &g

	return &g
}