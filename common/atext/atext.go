// This package is used to manage and load fonts, layout and draw text.
package atext

import (
	"github.com/cadmean-ru/amphion/common/a"
)

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

//GetAllChars returns all characters in the text.
func (t *Text) GetAllChars() []*Char {
	count := t.GetCharsCount()

	allChars := make([]*Char, count)
	count = 0

	for _, l := range t.lines {
		for _, c := range l.chars {
			allChars[count] = c
			count++
		}
	}

	return allChars
}

//GetCharCount returns the total number of characters in the text.
func (t *Text) GetCharsCount() int {
	count := 0
	for _, l := range t.lines {
		count += len(l.chars)
	}
	return count
}

//GetCharAt returns the character at the specified position in text.
func (t *Text) GetCharAt(index int) *Char {
	i := 0

	if index < 0 || index > t.GetCharsCount() {
		return nil
	}

	for _, l := range t.lines {
		for _, c := range l.chars {
			if i == index {
				return c
			}
			i++
		}
	}

	return nil
}

func (t *Text) GetLinesCount() int {
	return len(t.lines)
}

func (t *Text) GetLineAt(index int) *Line {
	return t.lines[index]
}

// GetSize returns the calculated text size.
func (t *Text) GetSize() a.IntVector2 {
	return a.NewIntVector2(t.width, t.height)
}
