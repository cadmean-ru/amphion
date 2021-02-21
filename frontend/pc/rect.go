package pc

import (
	"github.com/cadmean-ru/amphion/engine"
	"github.com/cadmean-ru/amphion/rendering"
	"github.com/go-gl/gl/all-core/gl"
)

type RectangleRenderer struct {
	*glPrimitiveRenderer
}

func (r *RectangleRenderer) OnStart() {
	r.program = createAndLinkProgramOrPanic(
		createAndCompileShaderOrPanic(ShapeVertexShaderStr, gl.VERTEX_SHADER),
		createAndCompileShaderOrPanic(ShapeFragShaderStr, gl.FRAGMENT_SHADER),
	)
}

func (r *RectangleRenderer) OnRender(ctx *rendering.PrimitiveRenderingContext) {
	gp := ctx.Primitive.(*rendering.GeometryPrimitive)
	state := ctx.State.(*glPrimitiveState)

	state.gen()

	gl.UseProgram(r.program)

	if ctx.Redraw {
		gl.BindVertexArray(state.vao)

		wSize := engine.GetScreenSize3()
		ntlPos := gp.Transform.Position.Ndc(wSize)                        // normalized top left
		nbrPos := gp.Transform.Position.Add(gp.Transform.Size).Ndc(wSize) // normalized bottom right

		color := gp.Appearance.FillColor
		r1 := float32(color.R)
		g1 := float32(color.G)
		b1 := float32(color.B)
		a1 := float32(color.A)

		vertices := []float32{
			ntlPos.X, ntlPos.Y, 0, r1, g1, b1, a1,
			ntlPos.X, nbrPos.Y, 0, r1, g1, b1, a1,
			nbrPos.X, nbrPos.Y, 0, r1, g1, b1, a1,
			nbrPos.X, ntlPos.Y, 0, r1, g1, b1, a1,
		}

		indices := []uint32{
			0, 1, 2,
			0, 3, 2,
		}

		const stride int32 = 28

		gl.BindBuffer(gl.ARRAY_BUFFER, state.vbo)
		gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)

		gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, state.ebo)
		gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(indices)*4, gl.Ptr(indices), gl.STATIC_DRAW)

		gl.VertexAttribPointer(0, 3, gl.FLOAT, false, stride, nil)
		gl.EnableVertexAttribArray(0)

		gl.VertexAttribPointer(1, 4, gl.FLOAT, false, stride, gl.PtrOffset(12))
		gl.EnableVertexAttribArray(1)

		gl.BindBuffer(gl.ARRAY_BUFFER, 0)
		gl.BindVertexArray(0)
	}

	gl.BindVertexArray(state.vao)

	gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, nil)
}
