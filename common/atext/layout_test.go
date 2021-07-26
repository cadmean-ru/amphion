package atext

import (
	"fmt"
	"github.com/cadmean-ru/amphion/common"
	"testing"
)

func TestLayoutRunes(t *testing.T) {
	font, err := ParseFont(DefaultFontData)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("Using font %s\n", font.GetName())

	face := font.NewFace(14)

	text := LayoutRunes(face, []rune("Hello bruh"), common.NewRectBoundary(0, 30, 0, 100, 0, 0), LayoutOptions{})

	fmt.Println("The layered out chars:")
	text.ForEachChar(func(c *Char) {
		fmt.Printf("%c: %+v\n", c.rune, c)
	})
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