package atext

import (
	"github.com/cadmean-ru/amphion/common"
	"github.com/cadmean-ru/amphion/common/a"
)

// LayoutRunes splits the given text into lines, aligns text and calculates the coordinates of each rune in the slice.
func LayoutRunes(face *Face, runes []rune, bounds *common.Rect, options LayoutOptions) *Text {
	return newLayoutManager(face, runes, bounds, options).layout()
}

func LayoutStringCompat(face *Face, text string, minX, maxX, minY, maxY, minZ, maxZ float32, vTextAlign, hTextAlign int) *Text {
	runes := []rune(text)
	bounds := common.NewRect(minX, maxX, minY, maxY, minZ, maxZ)
	return LayoutRunes(face, runes, bounds, LayoutOptions{
		VTextAlign: a.TextAlign(vTextAlign),
		HTextAlign: a.TextAlign(hTextAlign),
	})
}

type LayoutOptions struct {
	VTextAlign a.TextAlign
	HTextAlign a.TextAlign
	SingleLine bool
}
