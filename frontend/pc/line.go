// +build windows linux darwin
// +build !android
// +build !ios

package pc

import (
	"github.com/cadmean-ru/amphion/engine"
	"github.com/cadmean-ru/amphion/rendering"
	"github.com/go-gl/gl/v3.3-core/gl"
)

type LineRenderer struct {
	*glPrimitiveRenderer
}

func (r *LineRenderer) OnStart() {
	r.program = NewGlProgram(DefaultVertexShaderStr, DefaultFragShaderStr, "line")
	r.program.CompileAndLink()
}

func (r *LineRenderer) OnRender(ctx *rendering.PrimitiveRenderingContext) {
	r.glPrimitiveRenderer.OnRender(ctx)

	gp := ctx.Primitive.(*rendering.GeometryPrimitive)
	state := ctx.State.(*glPrimitiveState)

	state.gen()

	if ctx.Redraw {
		gl.BindVertexArray(state.vao)

		wSize := engine.GetScreenSize3()
		nPos := gp.Transform.Position.Ndc(wSize)
		nSize := gp.Transform.Position.Add(gp.Transform.Size).Ndc(wSize)

		vertices := []float32 {
			nPos.X, nPos.Y, 0,
			nSize.X, nSize.Y, 0,
		}

		gl.BindBuffer(gl.ARRAY_BUFFER, state.vbo)
		gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)

		gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 12, nil)
		gl.EnableVertexAttribArray(0)

		gl.BindBuffer(gl.ARRAY_BUFFER, 0)
		gl.BindVertexArray(0)
	}

	//color := gp.Appearance.StrokeColor.Normalize()

	gl.BindVertexArray(state.vao)

	gl.DrawArrays(gl.LINE, 0, 2)
}