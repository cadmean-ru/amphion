// Package atext is used to manage and load fonts, layout and draw text.
package atext

import (
	"github.com/cadmean-ru/amphion/common"
	"github.com/cadmean-ru/amphion/common/a"
	"strings"
)

// Text contains the layered out text. Read only.
type Text struct {
	lines       []*Line
	width       int
	height      int
	initialRect *common.RectBoundary
}

func (t *Text) append(line *Line) {
	t.lines = append(t.lines, line)

	if line.width > t.width {
		t.width = line.width
	}

	t.height += line.face.GetLineHeight()
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

//GetCharsCount returns the total number of characters in the text.
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

func (t *Text) GetInitialRect() *common.RectBoundary {
	return t.initialRect
}

func (t *Text) GetActualRect() *common.RectBoundary {
	return common.NewRectBoundaryFromPositionAndSize(t.initialRect.GetMin(), a.NewVector3(float32(t.width), float32(t.height), 0))
}

func (t *Text) String() string {
	sb := strings.Builder{}

	for _, l := range t.lines {
		sb.WriteString(l.String())
	}

	return sb.String()
}