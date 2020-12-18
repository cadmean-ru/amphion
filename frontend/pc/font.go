// +build windows linux darwin
// +build !android

package pc

import (
	"fmt"
	"github.com/4ydx/gltext"
	v41 "github.com/4ydx/gltext/v4.1"
	"github.com/cadmean-ru/amphion/common"
	"github.com/flopp/go-findfont"
	"golang.org/x/image/math/fixed"
	"os"
)

type glCharacter struct {
	textureId uint32
	char      rune
	size      common.IntVector3
	bearing   common.IntVector3
	advance   uint32
}

type glFont struct {
	name       string
	characters []*glCharacter
	japanFont  *v41.Font
}

func loadFont(name string) (*glFont, error) {
	fontPath, err := findfont.Find(fmt.Sprintf("%s.ttf", name))
	if err != nil {
		return nil, err
	}

	fontData, err := os.Open(fontPath)
	if err != nil {
		return nil, err
	}
	defer fontData.Close()

	//f, err := truetype.Parse(fontData)
	//if err != nil {
	//	return nil, err
	//}
	//
	//fupe := fixed.Int26_6(f.FUnitsPerEm())
	//
	//face := truetype.NewFace(f, &truetype.Options{
	//	Size:              18,
	//	DPI:               72,
	//	Hinting:           font.HintingNone,
	//})
	//
	//for _, r := range "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890 " {
	//	//dr, mask, maskp, adv, ok := face.Glyph(fixed.Point26_6{}, r)
	//
	//	i := f.Index(r)
	//	hm := f.HMetric(fupe, i)
	//	g := &truetype.GlyphBuf{}
	//
	//	err = g.Load(f, fupe, i, font.HintingNone)
	//	if err != nil {
	//		return nil, err
	//	}
	//
	//	c := glCharacter{
	//		char: r,
	//	}
	//
	//	gl.GenTextures(1, &c.textureId)
	//	gl.BindTexture(gl.TEXTURE_2D, c.textureId)
	//	gl.TexImage2D(
	//		gl.TEXTURE_2D,
	//		0,
	//		gl.RED,
	//		int32(dr.Size().X),
	//		int32(dr.Size().Y),
	//		0,
	//		gl.RED,
	//		gl.UNSIGNED_BYTE,
	//		gl.Ptr([]byte {}),
	//	)
	//}

	runeRanges := make(gltext.RuneRanges, 0)
	runeRanges = append(runeRanges, gltext.RuneRange{Low: 32, High: 128})
	runeRanges = append(runeRanges, gltext.RuneRange{Low: 0x3000, High: 0x3030})
	runeRanges = append(runeRanges, gltext.RuneRange{Low: 0x3040, High: 0x309f})
	runeRanges = append(runeRanges, gltext.RuneRange{Low: 0x30a0, High: 0x30ff})
	runeRanges = append(runeRanges, gltext.RuneRange{Low: 0x4e00, High: 0x9faf})
	runeRanges = append(runeRanges, gltext.RuneRange{Low: 0xff00, High: 0xffef})

	scale := fixed.Int26_6(32)
	runesPerRow := fixed.Int26_6(128)
	config, err := gltext.NewTruetypeFontConfig(fontData, scale, runeRanges, runesPerRow, 5)
	if err != nil {
		return nil, err
	}

	font := &glFont{}

	font.japanFont, err = v41.NewFont(config)
	if err != nil {
		return nil, err
	}

	return font, nil
}
