package atext

import (
	"github.com/cadmean-ru/amphion/common/a"
	"strings"
)

// Line represents a  single line of characters in the text.
type Line struct {
	text  *Text
	index int
	chars []*Char
	width int
	x, y  int
}

func (l *Line) append(c *Char, xOffset int) {
	l.chars = append(l.chars, c)
	l.width += xOffset
	c.line = l
}

//GetCharsCount returns the number of characters in this line.
func (l *Line) GetCharsCount() int {
	return len(l.chars)
}

//GetCharAt returns the char at the given index.
//Returns nil if index is out of bounds.
func (l *Line) GetCharAt(index int) *Char {
	return l.chars[index]
}

func (l *Line) GetChars() []*Char {
	charsCopy := make([]*Char, len(l.chars))
	copy(charsCopy, l.chars)
	return charsCopy
}

// GetSize returns the size of the line in pixels.
func (l *Line) GetSize() a.IntVector2 {
	return a.NewIntVector2(l.width, l.text.face.GetLineHeight())
}

//GetHeight return the height of this line in pixels.
func (l *Line) GetHeight() int {
	return l.text.face.GetLineHeight()
}

//GetWidth returns the width of all characters in this line in pixels.
func (l *Line) GetWidth() int {
	return l.width
}

//GetPosition returns the top-left position of the line.
func (l *Line) GetPosition() a.IntVector2 {
	return a.NewIntVector2(l.x, l.y)
}

//GetX returns the x coordinate of the top-left position.
func (l *Line) GetX() int {
	return l.x
}

//GetY returns the y coordinate of the top-left position.
func (l *Line) GetY() int {
	return l.y
}

//GetIndex returns the index of the line in text.
func (l *Line) GetIndex() int {
	return l.index
}

//IsEmpty checks if this line has no characters.
func (l *Line) IsEmpty() bool {
	return len(l.chars) == 0
}

func (l *Line) String() string {
	sb := strings.Builder{}

	for _, c := range l.chars {
		sb.WriteRune(c.rune)
	}

	return sb.String()
}

func newLine(t *Text, i, x, y int) *Line {
	return &Line{
		text:  t,
		index: i,
		chars: make([]*Char, 0, 50),
		x:     x,
		y:     y,
	}
}
