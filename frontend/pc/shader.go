// +build windows linux darwin
// +build !android

package pc

import (
	"github.com/go-gl/gl/v3.3-core/gl"
	"log"
	"strings"
)

func createAndCompileShaderOrPanic(text string, shaderType uint32) uint32 {
	shaderSource, free := gl.Strs(text)

	shader := gl.CreateShader(shaderType)

	gl.ShaderSource(shader, 1, shaderSource, nil)
	free()
	gl.CompileShader(shader)

	var success int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &success)

	if success == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		infoLog := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, 512, nil, gl.Str(infoLog))

		log.Println(infoLog)
		panic("failed to compile shader")
	}

	return shader
}

func createAndLinkProgramOrPanic(vertexShader, fragShader uint32) uint32 {
	shaderProgram := gl.CreateProgram()
	gl.AttachShader(shaderProgram, vertexShader)
	gl.AttachShader(shaderProgram, fragShader)
	gl.LinkProgram(shaderProgram)

	var success int32

	gl.GetProgramiv(shaderProgram, gl.LINK_STATUS, &success);
	if success == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(shaderProgram, gl.INFO_LOG_LENGTH, &logLength)

		infoLog := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(shaderProgram, 512, nil, gl.Str(infoLog))

		log.Println(infoLog)
		panic("failed to link program")
	}

	gl.DeleteShader(vertexShader)
	gl.DeleteShader(fragShader)

	return shaderProgram
}

func createAndLinkDefaultProgramOrPanic() uint32 {
	vertexShader := createAndCompileShaderOrPanic(DefaultVertexShaderStr, gl.VERTEX_SHADER)
	fragShader := createAndCompileShaderOrPanic(DefaultFragShaderStr, gl.FRAGMENT_SHADER)

	return createAndLinkProgramOrPanic(vertexShader, fragShader)
}
