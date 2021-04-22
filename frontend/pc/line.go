// +build windows linux darwin
// +build !android
// +build !ios

package pc

import (
	"github.com/cadmean-ru/amphion/engine"
	"github.com/cadmean-ru/amphion/rendering"
	"github.com/go-gl/gl/v4.1-core/gl"
)

type LineRenderer struct {
	*glPrimitiveRenderer
}

func (r *LineRenderer) OnStart() {
	r.program = createAndLinkProgramOrPanic(
		createAndCompileShaderOrPanic(zeroTerminated(DefaultVertexShaderStr), gl.VERTEX_SHADER),
		createAndCompileShaderOrPanic(zeroTerminated(DefaultFragShaderStr), gl.FRAGMENT_SHADER),
	)
}

func (r *LineRenderer) OnRender(ctx *rendering.PrimitiveRenderingContext) {
	gp := ctx.Primitive.(*rendering.GeometryPrimitive)
	state := ctx.State.(*glPrimitiveState)

	state.gen()

	gl.UseProgram(r.program)

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