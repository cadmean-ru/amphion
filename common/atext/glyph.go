package atext

import "github.com/cadmean-ru/amphion/common/a"

// Glyph is the visual representation of a rune.
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

// GetFace returns the font face it is associated with.
func (g *Glyph) GetFace() *Face {
	return g.face
}

// GetPixels returns the actual pixels of the glyph, that can be drawn to the screen.
func (g *Glyph) GetPixels() []byte {
	return g.pixels
}

// GetSize returns the size of the glyph.
func (g *Glyph) GetSize() a.IntVector2 {
	return g.size
}

func (g *Glyph) GetWidth() int {
	return g.size.X
}

func (g *Glyph) GetHeight() int {
	return g.size.Y
}

// GetBearing returns the x and y bearing of the glyph.
func (g *Glyph) GetBearing() a.IntVector2 {
	return g.bearing
}

// GetAdvance returns the x advance of the glyph.
func (g *Glyph) GetAdvance() int {
	return g.advance
}

// GetAscent returns the ascent of the glyph.
func (g *Glyph) GetAscent() int {
	return g.ascent
}

// GetDescent returns the descent of the glyph.
func (g *Glyph) GetDescent() int {
	return g.descent
}

// GetRune returns the rune glyph is associated with.
func (g *Glyph) GetRune() rune {
	return g.rune
}