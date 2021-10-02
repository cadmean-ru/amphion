//+build windows linux darwin
//+build !android
//+build !ios

package opengl

import (
	"github.com/cadmean-ru/amphion/rendering"
	"github.com/go-gl/gl/v3.3-core/gl"
)

type PolygonRenderer struct {
	*glPrimitiveRenderer
}

func (r *PolygonRenderer) OnStart() {
	r.program = NewGlProgram(PolygonVertexShaderStr, PolygonFragmentShaderStr, "polygon")
	r.program.CompileAndLink()
}

func (r *PolygonRenderer) OnRender(ctx *rendering.PrimitiveRenderingContext) {
	r.glPrimitiveRenderer.OnRender(ctx)

	pp := ctx.Primitive.(*rendering.PolygonPrimitive)
	state := ctx.State.(*glPrimitiveState)

	state.gen()

	if ctx.Redraw {
		gl.BindVertexArray(state.vao)

		color := pp.Appearance.FillColor.ToFloat4()

		const stride int32 = 28
		const floatsPerVertex = 7

		vertices := make([]float32, int32(len(pp.Vertices)) * floatsPerVertex)
		for i, v := range pp.Vertices {
			j := int32(i) * floatsPerVertex
			vertices[j] = v.X
			vertices[j+1] = v.Y
			vertices[j+2] = v.Z
			vertices[j+3] = color.X
			vertices[j+4] = color.Y
			vertices[j+5] = color.Z
			vertices[j+6] = color.W
		}

		indices := make([]uint32, len(pp.Indexes))
		for i, index := range pp.Indexes {
			indices[i] = uint32(index)
		}

		//vertices := []float32{
		//	0, 0, 0, 0, 0, 0, 255,
		//	500, 0, 0, 0, 0, 0, 255,
		//	69, 69, 0, 0, 0, 0, 255,
		//	420, 69, 0, 0, 0, 0, 255,
		//	//0.1, 0.1, 0, 0, 0, 0, 255,
		//	//1, 1, 0, 0, 0, 0, 255,
		//	//0, 0, 0, 0, 0, 0, 255,
		//}
		//
		//indices := []uint32{
		//	0, 1, 2,
		//	1, 3, 2,
		//}

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

	gl.DrawElements(gl.TRIANGLES, int32(len(pp.Indexes)), gl.UNSIGNED_INT, nil)
}