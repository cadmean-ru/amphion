// This package is used to manage and load fonts, layout and draw text.
package atext

import "github.com/cadmean-ru/amphion/common/a"

// Text contains the layered out text.
type Text struct {
	lines  []*Line
	width  int
	height int
}

func (t *Text) append(line *Line) {
	t.lines = append(t.lines, line)
}

// ForEachChar iterates over each character in this text.
func (t *Text) ForEachChar(delegate func(c *Char)) {
	for _, l := range t.lines {
		for _, c := range l.chars {
			delegate(c)
		}
	}
}

// GetSize returns the calculated text size.
func (t *Text) GetSize() a.IntVector2 {
	return a.NewIntVector2(t.width, t.height)
}
