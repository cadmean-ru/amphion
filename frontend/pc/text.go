package pc

import (
	"github.com/cadmean-ru/amphion/common/a"
	"github.com/cadmean-ru/amphion/common/atext"
	"github.com/cadmean-ru/amphion/engine"
	"github.com/cadmean-ru/amphion/rendering"
	"github.com/go-gl/gl/v4.1-core/gl"
)

type TextRenderer struct {
	*glPrimitiveRenderer
	fonts map[string]*atext.Font
	charsTextures  map[int]map[rune]uint32
}

func (r *TextRenderer) OnStart() {
	r.fonts = make(map[string]*atext.Font)
	r.charsTextures = make(map[int]map[rune]uint32)
	r.program = createAndLinkProgramOrPanic(
		createAndCompileShaderOrPanic(TextVertexShaderStr, gl.VERTEX_SHADER),
		createAndCompileShaderOrPanic(TextFragShaderStr, gl.FRAGMENT_SHADER),
	)
}

func (r *TextRenderer) OnRender(ctx *rendering.PrimitiveRenderingContext) {
	tp := ctx.Primitive.(*rendering.TextPrimitive)

	fontName := tp.TextAppearance.Font
	if fontName == "" {
		fontName = "Arial"
	}

	var font *atext.Font
	var ok bool
	var err error
	if font, ok = r.fonts[fontName]; !ok {
		font, err = atext.LoadFont(fontName)
		if err != nil {
			panic(err)
		}
		r.fonts[fontName] = font
	}

	face := font.NewFace(int(tp.TextAppearance.FontSize))

	color := tp.Appearance.FillColor.Normalize()

	runes := []rune(tp.Text)

	atext.LayoutRunes(face, runes, tp.Transform.GetRect(), atext.LayoutOptions{
		VTextAlign: tp.VTextAlign,
		HTextAlign: tp.HTextAlign,
	}).ForEachChar(func(c *atext.Char) {
		if !c.IsVisible() {
			return
		}

		pos2d := c.GetPosition()
		wSize := engine.GetScreenSize3()
		pos := a.NewIntVector3(pos2d.X, pos2d.Y, 0)
		npos := pos.Ndc(wSize)

		w := c.GetGlyph().GetSize().X
		h := c.GetGlyph().GetSize().Y

		brpos := pos.Add(a.NewIntVector3(w, h, 0))
		nbrpos := brpos.Ndc(wSize)

		var textId uint32
		if m1, ok := r.charsTextures[int(tp.TextAppearance.FontSize)]; ok {
			if tx, ok := m1[c.GetRune()]; ok {
				textId = tx
			} else {
				textId = genGlyphTex(c.GetGlyph())
				m1[c.GetRune()] = textId
			}
		} else {
			textId = genGlyphTex(c.GetGlyph())
			r.charsTextures[int(tp.TextAppearance.FontSize)] = make(map[rune]uint32)
			r.charsTextures[int(tp.TextAppearance.FontSize)][c.GetRune()] = textId
		}

		drawTex(ctx, npos, nbrpos, textId, r.program, true, func() {
			gl.Uniform3f(gl.GetUniformLocation(r.program, gl.Str("textColor\x00")), color.X, color.Y, color.Z)
		})
	})

	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.BindVertexArray(0)
	gl.BindTexture(gl.TEXTURE_2D, 0)
}

func genGlyphTex(g *atext.Glyph) uint32 {
	var textId uint32
	gl.GenTextures(1, &textId)
	gl.BindTexture(gl.TEXTURE_2D, textId)
	gl.TexImage2D(
		gl.TEXTURE_2D,
		0,
		gl.RED,
		int32(g.GetSize().X),
		int32(g.GetSize().Y),
		0,
		gl.RED,
		gl.UNSIGNED_BYTE,
		gl.Ptr(g.GetPixels()),
	)

	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)

	return textId
}