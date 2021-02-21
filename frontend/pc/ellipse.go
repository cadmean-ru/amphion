package pc

import (
	"github.com/cadmean-ru/amphion/engine"
	"github.com/cadmean-ru/amphion/rendering"
	"github.com/go-gl/gl/all-core/gl"
)

type EllipseRenderer struct {
	*glPrimitiveRenderer
}

func (r *EllipseRenderer) OnStart() {
	r.program = createAndLinkProgramOrPanic(
		createAndCompileShaderOrPanic(EllipseVertexShaderStr, gl.VERTEX_SHADER),
		createAndCompileShaderOrPanic(EllipseFragShaderStr, gl.FRAGMENT_SHADER),
	)
}

func (r *EllipseRenderer) OnRender(ctx *rendering.PrimitiveRenderingContext) {
	gp := ctx.Primitive.(*rendering.GeometryPrimitive)
	state := ctx.State.(*glPrimitiveState)

	state.gen()

	gl.UseProgram(r.program)

	if ctx.Redraw {
		gl.BindVertexArray(state.vao)

		wSize := engine.GetScreenSize3()
		tlPos := gp.Transform.Position.Ndc(wSize)
		brPos := gp.Transform.Position.Add(gp.Transform.Size).Ndc(wSize)

		color := gp.Appearance.FillColor
		r1 := float32(color.R)
		g1 := float32(color.G)
		b1 := float32(color.B)
		a1 := float32(color.A)

		vertices := []float32 {
			tlPos.X, tlPos.Y, 0, r1, g1, b1, a1, tlPos.X, tlPos.Y, 0, brPos.X, brPos.Y, 0,
			tlPos.X, brPos.Y, 0, r1, g1, b1, a1, tlPos.X, tlPos.Y, 0, brPos.X, brPos.Y, 0,
			brPos.X, brPos.Y, 0, r1, g1, b1, a1, tlPos.X, tlPos.Y, 0, brPos.X, brPos.Y, 0,
			brPos.X, tlPos.Y, 0, r1, g1, b1, a1, tlPos.X, tlPos.Y, 0, brPos.X, brPos.Y, 0,
		}

		indices := []uint32 {
			0, 1, 2,
			0, 3, 2,
		}

		const stride int32 = 52

		gl.BindBuffer(gl.ARRAY_BUFFER, state.vbo)
		gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)

		gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, state.ebo)
		gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(indices)*4, gl.Ptr(indices), gl.STATIC_DRAW)

		gl.VertexAttribPointer(0, 3, gl.FLOAT, false, stride, nil)
		gl.EnableVertexAttribArray(0)

		gl.VertexAttribPointer(1, 4, gl.FLOAT, false, stride, gl.PtrOffset(12))
		gl.EnableVertexAttribArray(1)

		gl.VertexAttribPointer(2, 3, gl.FLOAT, false, stride, gl.PtrOffset(28))
		gl.EnableVertexAttribArray(2)

		gl.VertexAttribPointer(3, 3, gl.FLOAT, false, stride, gl.PtrOffset(40))
		gl.EnableVertexAttribArray(3)

		gl.BindBuffer(gl.ARRAY_BUFFER, 0)
		gl.BindVertexArray(0)
	}

	gl.BindVertexArray(state.vao)

	gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, nil)
}
