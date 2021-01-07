// +build windows linux darwin
// +build !android

package pc

import (
	"errors"
	"fmt"
	"github.com/cadmean-ru/amphion/common/a"
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
	size      a.IntVector3
	bearing   a.IntVector3
	advance   a.IntVector3
}

type glFont struct {
	name       string
	characters map[rune]*glCharacter
	maxHeight  int
	ttf        *truetype.Font
}

func (f *glFont) kern(a, b rune) int {
	i1 := f.ttf.Index(a)
	i2 := f.ttf.Index(b)
	return f.ttf.Kern(fixed.Int26_6(f.maxHeight<<6), i1, i2).Round()
}

func loadFont(name string) (*glFont, error) {
	fontPath, err := findfont.Find(fmt.Sprintf("%s.ttf", name))
	if err != nil {
		return nil, err
	}

	fontData, err := ioutil.ReadFile(fontPath)

	f, err := truetype.Parse(fontData)
	if err != nil {
		return nil, err
	}

	fupe := fixed.Int26_6(18<<6)

	face := truetype.NewFace(f, &truetype.Options{
		Size:              18,
		DPI:               72,
		Hinting:           font.HintingNone,
	})

	glFont := glFont{
		name:       name,
		characters: make(map[rune]*glCharacter),
		ttf:        f,
	}

	maxH := 0

	for _, r := range "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890\uE00A" {
		face.Glyph(fixed.Point26_6{}, r)
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

		c := glCharacter{
			textureId: 0,
			char:      r,
			size:      a.NewIntVector3(dr.Size().X, dr.Size().Y, 0),
			bearing:   a.NewIntVector3(hm.LeftSideBearing.Round(), vm.TopSideBearing.Round(), 0),
			advance:   a.NewIntVector3(hm.AdvanceWidth.Round(), vm.AdvanceHeight.Round(), 0),
		}

		if dr.Size().Y > maxH {
			maxH = dr.Size().Y
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

	glFont.maxHeight = maxH

	return &glFont, nil
}
