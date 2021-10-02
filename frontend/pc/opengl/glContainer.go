// +build windows linux darwin
// +build !android
// +build !ios

package opengl

import (
	"github.com/cadmean-ru/amphion/rendering"
	"github.com/go-gl/gl/v3.3-core/gl"
)

type glContainer struct {
	primitive rendering.Primitive
	vbo       uint32
	vao       uint32
	ebo       uint32
	redraw    bool
	colorLoc  int32
	other     map[string]interface{}
}

func (c *glContainer) gen() {
	if c.vbo == 0 {
		temp := make([]uint32, 2)
		gl.GenBuffers(2, &temp[0])
		c.vbo = temp[0]
		c.ebo = temp[1]
		gl.GenVertexArrays(1, &c.vao)
	}
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
	return &glContainer{
		colorLoc: -1,
		other:    make(map[string]interface{}),
	}
}
