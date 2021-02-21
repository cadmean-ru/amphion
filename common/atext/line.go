package atext

import "github.com/cadmean-ru/amphion/common/a"

// Line represents a  single line of characters in the text.
type Line struct {
	face   *Face
	chars  []*Char
	width  int
}

func (l *Line) append(c *Char) {
	l.chars = append(l.chars, c)
}

// Returns the size of the line in pixels.
func (l *Line) GetSize() a.IntVector2 {
	return a.NewIntVector2(l.width, l.face.GetLineHeight())
}

func newLine(f *Face) *Line {
	return &Line{
		face:  f,
		chars: make([]*Char, 0, 50),
	}
}