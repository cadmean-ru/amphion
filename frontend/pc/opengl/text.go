//+build windows linux darwin
//+build !ios
//+build !android

package opengl

import (
	"github.com/cadmean-ru/amphion/common/a"
	"github.com/cadmean-ru/amphion/common/atext"
	"github.com/cadmean-ru/amphion/rendering"
	"github.com/go-gl/gl/v3.3-core/gl"
)

type TextRenderer struct {
	*glPrimitiveRenderer
	fonts map[string]*atext.Font
	charsTextures  map[int]map[rune]uint32
}

func (r *TextRenderer) OnStart() {
	r.fonts = make(map[string]*atext.Font)
	r.charsTextures = make(map[int]map[rune]uint32)

	r.program = NewGlProgram(TextVertexShaderStr, TextFragShaderStr, "text")
	r.program.CompileAndLink()
}

func (r *TextRenderer) OnRender(ctx *rendering.PrimitiveRenderingContext) {
	r.glPrimitiveRenderer.OnRender(ctx)

	tp := ctx.Primitive.(*rendering.TextPrimitive)

	color := tp.Appearance.FillColor.Normalize()

	tp.TextProvider.GetAText().ForEachChar(func(c *atext.Char) {
		if !c.IsVisible() {
			return
		}

		pos2d := c.GetPosition()
		pos := a.NewIntVector3(pos2d.X, pos2d.Y, 0)
		npos := pos.ToFloat()

		w := c.GetGlyph().GetSize().X
		h := c.GetGlyph().GetSize().Y

		brpos := pos.Add(a.NewIntVector3(w, h, 0))
		nbrpos := brpos.ToFloat()

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

		drawTex(ctx, npos, nbrpos, textId, r.program.id, true, func() {
			gl.Uniform4f(r.program.GetUniformLocation("uTextColor"), color.X, color.Y, color.Z, color.W)
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