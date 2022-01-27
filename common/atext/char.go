package atext

import (
	"github.com/cadmean-ru/amphion/common"
	"github.com/cadmean-ru/amphion/common/a"
)

// Char is a representation of a character of text on the screen.
type Char struct {
	line  *Line
	index int
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

// GetPosition returns the top-left position of the character on the screen, where it should be drawn.
func (c *Char) GetPosition() a.IntVector2 {
	return c.pos
}

//GetX return the x coordinate of the top-left position of the char.
func (c *Char) GetX() int {
	return c.pos.X
}

//GetY return the y coordinate of the top-left position of the char.
func (c *Char) GetY() int {
	return c.pos.Y
}

//GetSize reruns the size of the char in pixels.
func (c *Char) GetSize() a.IntVector2 {
	if c.glyph != nil {
		return c.glyph.GetSize()
	}

	if c.rune == ' ' {
		return a.NewIntVector2(c.line.text.face.GetSize() / 4, c.line.text.face.GetSize())
	}

	return a.NewIntVector2(0, c.line.text.face.GetSize())
}

//GetRect returns the rect of the char.
func (c *Char) GetRect() *common.Rect {
	return common.NewRectFromPositionAndSize(c.GetPosition().ToFloat3(), c.GetGlyph().GetSize().ToFloat3())
}

// IsVisible tells if the character has a visual representation.
func (c *Char) IsVisible() bool {
	return c.glyph != nil && c.rune != ' ' && c.rune != '\n'
}

//GetLine returns the line containing this char.
func (c *Char) GetLine() *Line {
	return c.line
}

//GetIndex returns the index of the char in the whole text.
func (c *Char) GetIndex() int {
	return c.index
}

func (c *Char) String() string {
	return string(c.rune)
}
