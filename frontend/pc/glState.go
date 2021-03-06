// +build windows linux darwin
// +build !android
// +build !ios

package pc

import (
	"github.com/go-gl/gl/v4.1-core/gl"
)

type glPrimitiveState struct {
	vbo uint32
	ebo uint32
	vao uint32
	tex uint32
}

func (c *glPrimitiveState) gen() {
	if c.vbo == 0 {
		temp := make([]uint32, 2)
		gl.GenBuffers(2, &temp[0])
		c.vbo = temp[0]
		c.ebo = temp[1]
		gl.GenVertexArrays(1, &c.vao)
	}
}

func (c *glPrimitiveState) free() {
	if c.vbo != 0 {
		gl.DeleteBuffers(1, &c.vbo)
	}

	if c.ebo != 0 {
		gl.DeleteBuffers(1, &c.ebo)
	}

	if c.vao != 0 {
		gl.DeleteVertexArrays(1, &c.vao)
	}

	if c.tex != 0 {
		gl.DeleteTextures(1, &c.tex)
	}
}