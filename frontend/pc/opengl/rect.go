//go:build (windows || linux || darwin) && !android && !ios
// +build windows linux darwin
// +build !android
// +build !ios

package opengl

import (
	"github.com/cadmean-ru/amphion/rendering"
	"github.com/go-gl/gl/v3.3-core/gl"
)

type RectangleRenderer struct {
	*glPrimitiveRenderer
}

func (r *RectangleRenderer) OnStart() {
	r.program = NewGlProgram(ShapeVertexShaderStr, ShapeFragShaderStr, "rect")
	r.program.CompileAndLink()
}

func (r *RectangleRenderer) OnRender(ctx *rendering.PrimitiveRenderingContext) {
	r.glPrimitiveRenderer.OnRender(ctx)

	gp := ctx.Primitive.(*rendering.GeometryPrimitive)
	state := ctx.State.(*glPrimitiveState)

	state.gen()

	if ctx.Redraw {
		gl.BindVertexArray(state.vao)

		tlPos := gp.Transform.Position.ToFloat()                        // normalized top left
		brPos := gp.Transform.Position.Add(gp.Transform.Size).ToFloat() // normalized bottom right

		color := gp.Appearance.FillColor
		r1 := float32(color.R)
		g1 := float32(color.G)
		b1 := float32(color.B)
		a1 := float32(color.A)

		strokeColor := gp.Appearance.StrokeColor
		r2 := float32(strokeColor.R)
		g2 := float32(strokeColor.G)
		b2 := float32(strokeColor.B)
		a2 := float32(strokeColor.A)

		stroke := float32(gp.Appearance.StrokeWeight)
		corner := float32(gp.Appearance.CornerRadius)

		vertices := []float32{
			tlPos.X, tlPos.Y, 0, tlPos.X, tlPos.Y, 0, brPos.X, brPos.Y, 0, r1, g1, b1, a1, stroke, r2, g2, b2, a2, corner,
			tlPos.X, brPos.Y, 0, tlPos.X, tlPos.Y, 0, brPos.X, brPos.Y, 0, r1, g1, b1, a1, stroke, r2, g2, b2, a2, corner,
			brPos.X, brPos.Y, 0, tlPos.X, tlPos.Y, 0, brPos.X, brPos.Y, 0, r1, g1, b1, a1, stroke, r2, g2, b2, a2, corner,
			brPos.X, tlPos.Y, 0, tlPos.X, tlPos.Y, 0, brPos.X, brPos.Y, 0, r1, g1, b1, a1, stroke, r2, g2, b2, a2, corner,
		}

		indices := []uint32{
			0, 1, 2,
			0, 3, 2,
		}

		const stride int32 = 76

		gl.BindBuffer(gl.ARRAY_BUFFER, state.vbo)
		gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)

		gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, state.ebo)
		gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(indices)*4, gl.Ptr(indices), gl.STATIC_DRAW)

		gl.VertexAttribPointer(0, 3, gl.FLOAT, false, stride, nil)
		gl.EnableVertexAttribArray(0)

		gl.VertexAttribPointer(1, 3, gl.FLOAT, false, stride, gl.PtrOffset(12))
		gl.EnableVertexAttribArray(1)

		gl.VertexAttribPointer(2, 3, gl.FLOAT, false, stride, gl.PtrOffset(24))
		gl.EnableVertexAttribArray(2)

		gl.VertexAttribPointer(3, 4, gl.FLOAT, false, stride, gl.PtrOffset(36))
		gl.EnableVertexAttribArray(3)

		gl.VertexAttribPointer(4, 1, gl.FLOAT, false, stride, gl.PtrOffset(52))
		gl.EnableVertexAttribArray(4)

		gl.VertexAttribPointer(5, 4, gl.FLOAT, false, stride, gl.PtrOffset(56))
		gl.EnableVertexAttribArray(5)

		gl.VertexAttribPointer(6, 1, gl.FLOAT, false, stride, gl.PtrOffset(72))
		gl.EnableVertexAttribArray(6)

		gl.BindBuffer(gl.ARRAY_BUFFER, 0)
		gl.BindVertexArray(0)
	}

	gl.BindVertexArray(state.vao)

	gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, nil)
}
