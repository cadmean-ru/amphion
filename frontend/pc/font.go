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
	size      a.IntVector2
	bearing   a.IntVector2
	advance   int
	ascent    int
	descent   int
}

type glFont struct {
	name       string
	characters map[rune]*glCharacter
	ttf        *truetype.Font
	scale      int
	capHeight  int
	lineHeight int
	xHeight    int
	ascent     int
	descent    int
}

func (f *glFont) kern(a, b rune) int {
	i1 := f.ttf.Index(a)
	i2 := f.ttf.Index(b)
	return f.ttf.Kern(fixed.Int26_6(f.scale<<6), i1, i2).Round()
}

func loadFont(name string, scale int) (*glFont, error) {
	fontPath, err := findfont.Find(fmt.Sprintf("%s.ttf", name))
	if err != nil {
		return nil, err
	}

	fontData, err := ioutil.ReadFile(fontPath)

	f, err := truetype.Parse(fontData)
	if err != nil {
		return nil, err
	}

	face := truetype.NewFace(f, &truetype.Options{
		Size:              float64(scale),
		DPI:               72,
		Hinting:           font.HintingFull,
	})

	metrics := face.Metrics()

	glFont := glFont{
		name:       name,
		characters: make(map[rune]*glCharacter),
		ttf:        f,
		scale:      scale,
		capHeight:  int(metrics.CapHeight) >> 6,
		lineHeight: int(metrics.Height) >> 6,
		xHeight:    int(metrics.XHeight) >> 6,
		ascent:     int(metrics.Ascent) >> 6,
		descent:    int(metrics.Descent) >> 6,
	}

	for _, r := range "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890\uE00A!@#$%^&**()-=_+{}[]|\\//?<>,.±§~`" {
		bounds, adv, ok := face.GlyphBounds(r)

		if !ok {
			return nil, errors.New("failed to get glyph")
		}

		width := int(bounds.Max.X - bounds.Min.X) >> 6
		height := int(bounds.Max.Y - bounds.Min.Y) >> 6

		//if glyph has no dimensions set to a max value
		if width == 0 || height == 0 {
			bounds = f.Bounds(fixed.Int26_6(scale))
			width = int((bounds.Max.X - bounds.Min.X) >> 6)
			height = int((bounds.Max.Y - bounds.Min.Y) >> 6)

			//above can sometimes yield 0 for font smaller than 48pt, 1 is minimum
			if width == 0 || height == 0 {
				width = 1
				height = 1
			}
		}

		ascent := int(-bounds.Min.Y) >> 6
		descent := int(bounds.Max.Y) >> 6

		c := glCharacter{
			textureId: 0,
			char:      r,
			size:      a.NewIntVector2(width, height),
			bearing:   a.NewIntVector2(int(bounds.Min.X) >> 6, descent),
			advance:   int(adv) >> 6,
			ascent:    ascent,
			descent:   descent,
		}

		dr, mask, maskp, _, ok := face.Glyph(fixed.Point26_6{}, r)
		img := image.NewGray(dr)
		draw.Draw(img, dr, mask, maskp, draw.Src)

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
			gl.Ptr(img.Pix),
		)

		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)

		glFont.characters[r] = &c
	}

	return &glFont, nil
}
