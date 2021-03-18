package atext

import "github.com/cadmean-ru/amphion/common/a"

// Char is a representation of a character of text on the screen.
type Char struct {
	rune  rune
	glyph *Glyph
	pos   a.IntVector2
}

// GetRune returns the rune the char is associated with.
func (c *Char) GetRune() rune {
	return c.rune
}

// GetGlyph returns the Glyph for this character. Can return nil if the character is not associated with any Glyph.
func (c *Char) GetGlyph() *Glyph {
	return c.glyph
}

// GetPosition returns the position of the character on the screen, where it should be drawn.
func (c *Char) GetPosition() a.IntVector2 {
	return c.pos
}

func (c *Char) GetX() int {
	return c.pos.X
}

func (c *Char) GetY() int {
	return c.pos.Y
}

// IsVisible tells if the character has a visual representation.
func (c *Char) IsVisible() bool {
	return c.rune != ' ' && c.rune != '\n'
}