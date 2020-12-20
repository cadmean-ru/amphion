// +build windows linux darwin
// +build !android

package pc

import (
	"errors"
	"fmt"
	"github.com/cadmean-ru/amphion/common"
	"github.com/flopp/go-findfont"
	"github.com/go-gl/gl/all-core/gl"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
	"image"
	"image/draw"
	"io/ioutil"
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
	characters map[rune]*glCharacter
}

//var ft C.FT_Library
//
//func initFreeType() {
//	if C.FT_Init_FreeType(&ft) != 0 {
//		panic("could not init freetype")
//	}
//}

func loadFont(name string) (*glFont, error) {
	fontPath, err := findfont.Find(fmt.Sprintf("%s.ttf", name))
	if err != nil {
		return nil, err
	}

	//fontData, err := os.Open(fontPath)
	//if err != nil {
	//	return nil, err
	//}
	//defer fontData.Close()

	fontData, err := ioutil.ReadFile(fontPath)

	f, err := truetype.Parse(fontData)
	if err != nil {
		return nil, err
	}

	fupe := fixed.Int26_6(f.FUnitsPerEm())

	face := truetype.NewFace(f, &truetype.Options{
		Size:              18,
		DPI:               72,
		Hinting:           font.HintingNone,
	})

	glFont := glFont{
		name:       name,
		characters: make(map[rune]*glCharacter),
	}

	for _, r := range "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890\uE00A" {
		dr, mask, maskp, _, ok := face.Glyph(fixed.Point26_6{}, r)

		if !ok {
			return nil, errors.New("failed to get glyph")
		}

		i := f.Index(r)
		hm := f.HMetric(fupe, i)
		vm := f.VMetric(fupe, i)

		bruh := image.NewGray(dr)
		draw.Draw(bruh, dr, mask, maskp, draw.Src)

		//bruh1 := make([]byte, dr.Size().X * dr.Size().Y)
		//
		//l := 0
		//for j := dr.Size().Y-1; j >= 0; j-- {
		//	for k := 0; k < dr.Size().X; k ++ {
		//		i1 := j * dr.Size().X + k
		//		i2 := l * dr.Size().X + k
		//		bruh1[i2] = bruh.Pix[i1]
		//	}
		//	l++
		//}
		//
		//cf, err := os.Create(fmt.Sprintf("%v.png", r))
		//if err == nil {
		//	err = png.Encode(cf, bruh)
		//	_ = cf.Close()
		//}

		c := glCharacter{
			textureId: 0,
			char:      r,
			size:      common.NewIntVector3(dr.Size().X, dr.Size().Y, 0),
			bearing:   common.NewIntVector3(hm.LeftSideBearing.Round(), vm.TopSideBearing.Round(), 0),
			advance:   uint32(hm.AdvanceWidth.Round()),
		}

		gl.GenTextures(1, &c.textureId)
		gl.BindTexture(gl.TEXTURE_2D, c.textureId)
		gl.TexImage2D(
			gl.TEXTURE_2D,
			0,
			gl.RED,
			int32(dr.Size().X),
			int32(dr.Size().Y),
			0,
			gl.RED,
			gl.UNSIGNED_BYTE,
			gl.Ptr(bruh.Pix),
		)

		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)

		glFont.characters[r] = &c
	}

	return &glFont, nil

	//runeRanges := make(gltext.RuneRanges, 0)
	//runeRanges = append(runeRanges, gltext.RuneRange{Low: 32, High: 128})
	//runeRanges = append(runeRanges, gltext.RuneRange{Low: 0x3000, High: 0x3030})
	//runeRanges = append(runeRanges, gltext.RuneRange{Low: 0x3040, High: 0x309f})
	//runeRanges = append(runeRanges, gltext.RuneRange{Low: 0x30a0, High: 0x30ff})
	//runeRanges = append(runeRanges, gltext.RuneRange{Low: 0x4e00, High: 0x9faf})
	//runeRanges = append(runeRanges, gltext.RuneRange{Low: 0xff00, High: 0xffef})
	//
	//scale := fixed.Int26_6(32)
	//runesPerRow := fixed.Int26_6(128)
	//config, err := gltext.NewTruetypeFontConfig(fontData, scale, runeRanges, runesPerRow, 5)
	//if err != nil {
	//	return nil, err
	//}
	//
	//font := &glFont{}
	//
	//font.japanFont, err = v41.NewFont(config)
	//if err != nil {
	//	return nil, err
	//}
	//
	//return font, nil

	//var face C.FT_Face
	//if C.FT_New_Face(&ft, fontPath, 0, &face) != 0 {
	//	return nil, errors.New("failed to load font")
	//}
	//
	//C.FT_Set_Pixel_Sizes(face, 0, 48)
	//
	//abc := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890 "
	//
	//f := glFont{
	//	name:       name,
	//	characters: make(map[rune]*glCharacter, len(abc)),
	//}
	//
	//for _, r := range abc {
	//	if C.FT_Load_Char(face, r, C.FT_LOAD_RENDER) != 0 {
	//		return nil, errors.New("failed to load character")
	//	}
	//
	//	var tex uint32
	//	gl.GenTextures(1, &tex)
	//	gl.BindTexture(gl.TEXTURE_2D, tex)
	//	gl.TexImage2D(
	//		gl.TEXTURE_2D,
	//		0,
	//		gl.RED,
	//		face.glyph.bitmap.width,
	//		face.glyph.bitmap.rows,
	//		0,
	//		gl.RED,
	//		gl.UNSIGNED_BYTE,
	//		face.glyph.bitmap.buffer,
	//	)
	//
	//	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	//	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	//	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	//	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	//
	//	c := glCharacter{
	//		textureId: tex,
	//		char:      r,
	//		size:      common.NewIntVector3(face.glyph.bitmap.width, face.glyph.bitmap.rows, 0),
	//		bearing:   common.NewIntVector3(face.glyph.bitmap_left, face.glyph.bitmap_top, 0),
	//		advance:   face.glyph.advance.x,
	//	}
	//
	//	f.characters[r] = &c
	//}
	//
	//C.FT_Done_Face(face)
	//
	//return &f, nil
}
