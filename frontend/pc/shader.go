package pc

import (
	"github.com/go-gl/gl/all-core/gl"
	"log"
	"strings"
)

const DefaultVertexShaderText = "#version 330 core\nlayout (location = 0) in vec3 aPos;\n\nvoid main()\n{\n    gl_Position = vec4(aPos.x, aPos.y, aPos.z, 1.0);\n}\x00"

const DefaultFragShaderText = "#version 330 core\nout vec4 FragColor;\n\nvoid main()\n{\n    FragColor = vec4(1.0f, 0.5f, 0.2f, 1.0f);\n}\x00"

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

func createAndLinkDefaultProgramOrPanic() uint32 {
	vertexShader := createAndCompileShaderOrPanic(DefaultVertexShaderText, gl.VERTEX_SHADER)
	fragShader := createAndCompileShaderOrPanic(DefaultFragShaderText, gl.FRAGMENT_SHADER)

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
