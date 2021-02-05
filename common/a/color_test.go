package a

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseHexColor(t *testing.T) {
	ass := assert.New(t)

	hex := "#2c68a8"
	expectedColor := NewColor(44, 104, 168)
	color := ParseHexColor(hex)

	fmt.Printf("%+v\n", color)

	ass.Equal(expectedColor, color)
}

func TestColor_GetHex(t *testing.T) {
	ass := assert.New(t)

	expectedHex := "#2c68a8"
	color := NewColor(44, 104, 168)
	hex := color.GetHex()

	fmt.Printf("%+v\n", hex)

	ass.Equal(expectedHex, hex)
}