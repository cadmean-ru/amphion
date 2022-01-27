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
		rect *common.Rect
		options LayoutOptions
		nLines int
	} {
		{
			"Hello bruh",
			common.NewRect2D(0, 100, 0, 100, 0),
			LayoutOptions{},
			1,
		},
		{
			"Hello\nbruh",
			common.NewRect2D(0, 100, 0, 100, 0),
			LayoutOptions{},
			2,
		},
		{
			"Hello\nbruh",
			common.NewRect2D(0, 100, 0, 100, 0),
			LayoutOptions{SingleLine: true},
			1,
		},
	}

	for i, testCase := range testCases {
		fmt.Printf("Test case %d\n", i)
		runes := []rune(testCase.input)
		text := LayoutRunes(face, runes, testCase.rect, testCase.options)

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