// +build windows linux darwin

package pc

import "github.com/go-gl/gl/all-core/gl"

type glContainer struct {
	primitive OpenGLPrimitive
	vbo       uint32
	vao       uint32
	ebo       uint32
}

func (c *glContainer) free() {
	if c.vbo != 0 {
		gl.DeleteBuffers(1, &c.vbo)
	}

	if c.ebo != 0 {
		gl.DeleteBuffers(1, &c.ebo)
	}

	if c.vao != 0 {
		gl.DeleteVertexArrays(1, &c.vao)
	}
}

func newGlContainer() *glContainer {
	return &glContainer{}
}
