package atext

import (
	"github.com/cadmean-ru/amphion/common"
	"github.com/cadmean-ru/amphion/common/a"
)

type layoutManager struct {
	text             *Text
	x, y             int
	xMin, xMax, yMax int
	currentLine      *Line
	face             *Face
	allChars         []*Char
	bounds           *common.RectBoundary
	options          LayoutOptions
	runes            []rune
}

func (m *layoutManager) layout() *Text {
	m.makeAllChars(m.runes)
	m.splitIntoLines()
	m.align()

	return m.text
}

func (m *layoutManager) makeAllChars(runes []rune) {
	m.allChars = make([]*Char, len(runes))

	for i, r := range runes {
		if g := m.face.GetGlyph(r); g != nil {
			m.allChars[i] = &Char{
				rune:  r,
				glyph: g,
				face:  m.face,
			}
		} else {
			m.allChars[i] = &Char{
				rune:  r,
				glyph: nil,
				face:  m.face,
			}
		}
	}
}

func (m *layoutManager) splitIntoLines() {
	for i, char := range m.allChars {
		if char.rune == ' ' || char.glyph == nil {
			char.pos = a.NewIntVector2(m.x, m.y+m.face.GetAscent())
			space := m.face.GetSize() / 4
			m.x += space
			m.currentLine.append(char, space)
			continue
		}

		if char.rune == '\n' || !m.charFits(char) {
			m.lineBreak()
		}

		if !m.lineFits() {
			break
		}

		char.pos = a.NewIntVector2(m.x+char.glyph.GetBearing().X, m.y+m.face.GetAscent()-char.glyph.GetAscent())

		var kern int
		if i < len(m.allChars)-1 {
			kern = m.face.GetKerning(char.rune, m.allChars[i+1].rune)
		}

		xOffset := char.GetGlyph().GetAdvance() + kern
		m.x += xOffset
		m.currentLine.append(char, xOffset)
	}

	if !m.currentLine.IsEmpty() {
		m.text.append(m.currentLine)
	}
}

func (m *layoutManager) align() {
	var yOffset int
	if m.options.VTextAlign == a.TextAlignCenter {
		yOffset = (int(m.bounds.Y.GetLength()) - m.text.height) / 2
	} else if m.options.VTextAlign == a.TextAlignBottom {
		yOffset = int(m.bounds.Y.GetLength()) - m.text.height
	}

	for _, line := range m.text.lines {
		lineSize := line.GetSize()

		var xOffset int
		if m.options.HTextAlign == a.TextAlignCenter {
			xOffset = (int(m.bounds.X.GetLength()) - lineSize.X) / 2
		} else if m.options.HTextAlign == a.TextAlignRight {
			xOffset = int(m.bounds.X.GetLength()) - lineSize.X
		}

		for _, c := range line.chars {
			c.pos = c.pos.Add(a.NewIntVector2(xOffset, yOffset))
		}
	}
}

func (m *layoutManager) lineBreak() {
	if m.options.SingleLine {
		return
	}

	m.x = m.xMin
	m.y += m.face.GetLineHeight()
	m.text.append(m.currentLine)
	m.currentLine = newLine(m.face)
}

func (m *layoutManager) lineFits() bool {
	return m.bounds.Y.Max == Unbounded || m.y < m.yMax
}

func (m *layoutManager) charFits(char *Char) bool {
	return m.bounds.X.Max == Unbounded || !char.IsVisible() || m.x+char.glyph.GetSize().X <= m.xMax
}

func newLayoutManager(face *Face, runes []rune, bounds *common.RectBoundary, options LayoutOptions) *layoutManager {
	m := layoutManager{}
	m.text = &Text{
		lines:       make([]*Line, 0, 10),
		initialRect: bounds,
	}
	m.x = int(bounds.X.Min)
	m.xMin = int(bounds.X.Min)
	m.xMax = int(bounds.X.Max)
	m.y = int(bounds.Y.Min)
	m.yMax = int(bounds.Y.Max)
	m.currentLine = newLine(face)
	m.face = face
	m.bounds = bounds
	m.runes = runes
	m.options = options
	return &m
}
