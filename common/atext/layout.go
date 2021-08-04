package atext

import (
	"github.com/cadmean-ru/amphion/common"
	"github.com/cadmean-ru/amphion/common/a"
)

// LayoutRunes splits the given text into lines, aligns text and calculates the coordinates of each rune in the slice.
func LayoutRunes(face *Face, runes []rune, bounds *common.RectBoundary, options LayoutOptions) *Text {
	//text := &Text{
	//	lines:       make([]*Line, 0, 10),
	//	initialRect: bounds,
	//}
	//
	//allChars := make([]*Char, len(runes))
	//
	//for i, r := range runes {
	//	if g := face.GetGlyph(r); g != nil {
	//		allChars[i] = &Char{
	//			rune:  r,
	//			glyph: g,
	//		}
	//	} else {
	//		allChars[i] = &Char{
	//			rune:  r,
	//			glyph: nil,
	//		}
	//	}
	//}
	//
	//x := int(bounds.X.Min)
	//xMin := int(bounds.X.Min)
	//xMax := int(bounds.X.Max)
	//y := int(bounds.Y.Min)
	//yMax := int(bounds.Y.Max)
	//currentLine := newLine(face)
	//
	//for i, char := range allChars {
	//	if char.rune == ' ' {
	//		x += face.GetSize() / 4
	//		currentLine.append(char)
	//		continue
	//	}
	//
	//	if char.glyph == nil {
	//		continue
	//	}
	//
	//	gs := char.glyph.GetSize()
	//
	//	if y+gs.Y > yMax {
	//		//cannot fit in more text
	//		break
	//	}
	//
	//	if char.rune == '\n' || x+char.glyph.GetAdvance() > xMax {
	//		//move to new line
	//		x = xMin
	//		y += face.GetLineHeight()
	//		text.append(currentLine)
	//
	//		if currentLine.width > text.width {
	//			text.width = currentLine.width
	//		}
	//		text.height += face.GetLineHeight()
	//
	//		if char.rune == '\n' {
	//			currentLine.append(char)
	//
	//			if y+gs.Y > yMax {
	//				//cannot fit in more text
	//				break
	//			}
	//		}
	//
	//		currentLine = newLine(face)
	//
	//		if char.rune == '\n' {
	//			continue
	//		}
	//	}
	//
	//	char.pos = a.NewIntVector2(x+char.glyph.GetBearing().X, y+face.GetAscent()-char.glyph.GetAscent())
	//
	//	var kern int
	//	if i < len(allChars)-1 {
	//		kern = face.GetKerning(char.rune, allChars[i+1].rune)
	//	}
	//
	//	bruh := char.GetGlyph().GetAdvance() + kern
	//	currentLine.width += bruh
	//	x += bruh
	//
	//	currentLine.append(char)
	//}
	//
	//text.append(currentLine)
	//if currentLine.width > text.width {
	//	text.width = currentLine.width
	//}
	//text.height += face.GetLineHeight()
	//
	//var yOffset int
	//if options.VTextAlign == a.TextAlignCenter {
	//	yOffset = (int(bounds.Y.GetLength()) - text.height) / 2
	//} else if options.VTextAlign == a.TextAlignBottom {
	//	yOffset = int(bounds.Y.GetLength()) - text.height
	//}
	//
	//for _, line := range text.lines {
	//	lineSize := line.GetSize()
	//
	//	var xOffset int
	//	if options.HTextAlign == a.TextAlignCenter {
	//		xOffset = (int(bounds.X.GetLength()) - lineSize.X) / 2
	//	} else if options.HTextAlign == a.TextAlignRight {
	//		xOffset = int(bounds.X.GetLength()) - lineSize.X
	//	}
	//
	//	for _, c := range line.chars {
	//		c.pos = c.pos.Add(a.NewIntVector2(xOffset, yOffset))
	//	}
	//}
	//
	//return text

	return newLayoutManager(face, runes, bounds, options).layout()
}

func LayoutStringCompat(face *Face, text string, minX, maxX, minY, maxY, minZ, maxZ float32, vTextAlign, hTextAlign int) *Text {
	runes := []rune(text)
	bounds := common.NewRectBoundary(minX, maxX, minY, maxY, minZ, maxZ)
	return LayoutRunes(face, runes, bounds, LayoutOptions{
		VTextAlign: a.TextAlign(vTextAlign),
		HTextAlign: a.TextAlign(hTextAlign),
	})
}

type LayoutOptions struct {
	VTextAlign a.TextAlign
	HTextAlign a.TextAlign
}
