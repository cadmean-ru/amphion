// +build windows linux darwin
// +build !android
// +build !ios

package opengl

import (
	"github.com/cadmean-ru/amphion/rendering"
	"github.com/go-gl/gl/v3.3-core/gl"

)

type ImageRenderer struct {
	*glPrimitiveRenderer
}

func (r *ImageRenderer) OnStart() {
	r.program = NewGlProgram(ImageVertexShaderStr, ImageFragShaderStr, "image")
	r.program.CompileAndLink()
}

func (r *ImageRenderer) OnRender(ctx *rendering.PrimitiveRenderingContext) {
	r.glPrimitiveRenderer.OnRender(ctx)

	ip := ctx.Primitive.(*rendering.ImagePrimitive)
	state := ctx.State.(*glPrimitiveState)

	if state.tex == nil {
		state.tex = make([]uint32, len(ip.Bitmaps))

		gl.GenTextures(int32(len(ip.Bitmaps)), &state.tex[0])

		for i, bitmap := range ip.Bitmaps {

			gl.BindTexture(gl.TEXTURE_2D, state.tex[i])

			gl.TexImage2D(
				gl.TEXTURE_2D,
				0,
				gl.SRGB_ALPHA,
				int32(bitmap.Width),
				int32(bitmap.Height),
				0,
				gl.RGBA,
				gl.UNSIGNED_BYTE,
				gl.Ptr(bitmap.Pixels),
			)

			gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
			gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)
			gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
			gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)

			bitmap.Dispose()
		}

		gl.BindTexture(gl.TEXTURE_2D, 0)
	}

	texId := state.tex[ip.Index]

	nPos := ip.Transform.Position.ToFloat()
	brPos := ip.Transform.Position.Add(ip.Transform.Size).ToFloat()

	drawTex(ctx, nPos, brPos, texId, r.program.id, false, nil)
}
