package atext

import "github.com/cadmean-ru/amphion/common/a"

type Glyph struct {
	face    *Face
	rune    rune
	pixels  []byte
	size    a.IntVector2
	bearing a.IntVector2
	advance int
	ascent  int
	descent int
}

func (g *Glyph) GetFace() *Face {
	return g.face
}

func (g *Glyph) GetPixels() []byte {
	return g.pixels
}

func (g *Glyph) GetSize() a.IntVector2 {
	return g.size
}

func (g *Glyph) GetBearing() a.IntVector2 {
	return g.bearing
}

func (g *Glyph) GetAdvance() int {
	return g.advance
}

func (g *Glyph) GetAscent() int {
	return g.ascent
}

func (g *Glyph) GetDescent() int {
	return g.descent
}

func (g *Glyph) GetRune() rune {
	return g.rune
}
