package atext

import (
	"github.com/cadmean-ru/amphion/common/a"
	"strings"
)

// Line represents a  single line of characters in the text.
type Line struct {
	face   *Face
	chars  []*Char
	width  int
}

func (l *Line) append(c *Char, xOffset int) {
	l.chars = append(l.chars, c)
	l.width += xOffset
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

// GetSize returns the size of the line in pixels.
func (l *Line) GetSize() a.IntVector2 {
	return a.NewIntVector2(l.width, l.face.GetLineHeight())
}

//GetHeight return the height of this line in pixels.
func (l *Line) GetHeight() int {
	return l.face.GetLineHeight()
}

//GetWidth returns the width of all characters in this line in pixels.
func (l *Line) GetWidth() int {
	return l.width
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

func newLine(f *Face) *Line {
	return &Line{
		face:  f,
		chars: make([]*Char, 0, 50),
	}
}