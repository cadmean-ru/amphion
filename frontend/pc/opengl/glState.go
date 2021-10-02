// +build windows linux darwin
// +build !android
// +build !ios

package opengl

import (
	"github.com/go-gl/gl/v3.3-core/gl"
)

type glPrimitiveState struct {
	vbo uint32
	ebo uint32
	vao uint32
	tex []uint32
}

func (s *glPrimitiveState) gen() {
	if s.vbo == 0 {
		temp := make([]uint32, 2)
		gl.GenBuffers(2, &temp[0])
		s.vbo = temp[0]
		s.ebo = temp[1]
		gl.GenVertexArrays(1, &s.vao)
	}
}

func (s *glPrimitiveState) free() {
	if s.vbo != 0 {
		gl.DeleteBuffers(1, &s.vbo)
	}

	if s.ebo != 0 {
		gl.DeleteBuffers(1, &s.ebo)
	}

	if s.vao != 0 {
		gl.DeleteVertexArrays(1, &s.vao)
	}

	for _, t := range s.tex {
		gl.DeleteTextures(1, &t)
	}
}