package atext

import (
	"fmt"
	"github.com/cadmean-ru/amphion/common"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLayoutRunes(t *testing.T) {
	as := assert.New(t)

	font, err := ParseFont(DefaultFontData)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("Using font %s\n", font.GetName())

	face := font.NewFace(14)

	testCases := []struct{
		input string
		rect *common.RectBoundary
		nLines int
	} {
		{
			"Hello bruh",
			common.NewRectBoundaryXY(0, 100, 0, 100),
			1,
		},
		{
			"Hello\nbruh",
			common.NewRectBoundaryXY(0, 100, 0, 100),
			2,
		},
	}

	for i, testCase := range testCases {
		fmt.Printf("Test case %d\n", i)
		runes := []rune(testCase.input)
		text := LayoutRunes(face, runes, testCase.rect, LayoutOptions{})

		fmt.Printf("Initial rect: %v\n", text.GetInitialRect())
		fmt.Printf("Actual rect: %v\n", text.GetActualRect())

		as.Equal(testCase.nLines, text.GetLinesCount())
		as.Equal(len(runes), text.GetCharsCount())

		fmt.Println()
		fmt.Println()
	}
}

func TestLayoutStringCompat(t *testing.T) {
	font, err := ParseFont(DefaultFontData)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("Using font %s\n", font.GetName())

	face := font.NewFace(14)

	text := LayoutStringCompat(face, "Hello bruh", 0, 30, 0, 100, 0, 0, 0, 0)

	fmt.Println("The layered out chars:")
	text.ForEachChar(func(c *Char) {
		fmt.Printf("%+v\n", c)
	})
}