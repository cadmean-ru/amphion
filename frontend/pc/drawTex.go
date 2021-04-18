package pc

import (
	"github.com/cadmean-ru/amphion/common/a"
	"github.com/cadmean-ru/amphion/rendering"
	"github.com/go-gl/gl/v4.1-core/gl"
)

func drawTex(ctx *rendering.PrimitiveRenderingContext, nPos, nbrPos a.Vector3, texId, progId uint32, alwaysRedraw bool, beforeDraw func()) {
	state := ctx.State.(*glPrimitiveState)
	state.gen()

	if ctx.Redraw || alwaysRedraw {
		gl.BindVertexArray(state.vao)

		vertices := []float32 {
			nPos.X,   nPos.Y,   0,	0, 0, // top left
			nPos.X,   nbrPos.Y, 0,	0, 1, // bottom left
			nbrPos.X, nbrPos.Y, 0,	1, 1, // top right
			nbrPos.X, nPos.Y,   0,	1, 0, // bottom right
		}

		indices := []uint32 {
			0, 1, 2,
			0, 3, 2,
		}

		gl.BindBuffer(gl.ARRAY_BUFFER, state.vbo)
		gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.DYNAMIC_DRAW)

		gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, state.ebo)
		gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(indices)*4, gl.Ptr(indices), gl.DYNAMIC_DRAW)

		gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 20, nil)
		gl.EnableVertexAttribArray(0)

		gl.VertexAttribPointer(1, 2, gl.FLOAT, false, 20, gl.PtrOffset(12))
		gl.EnableVertexAttribArray(1)

		gl.BindBuffer(gl.ARRAY_BUFFER, 0)
		gl.BindVertexArray(0)
	}

	gl.UseProgram(progId)
	gl.BindVertexArray(state.vao)
	gl.BindTexture(gl.TEXTURE_2D, texId)

	if beforeDraw != nil {
		beforeDraw()
	}

	gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, nil)

	gl.BindTexture(gl.TEXTURE_2D, 0)
	gl.BindVertexArray(0)
}
