// +build windows linux darwin
// +build !android

package pc

import (
	"github.com/go-gl/gl/all-core/gl"
	"log"
	"strings"
)

const DefaultVertexShaderText = "#version 330 core\nlayout (location = 0) in vec3 aPos;\n\nout vec4 trashPos;\n\nvoid main()\n{\n    trashPos = vec4(aPos.x, aPos.y, aPos.z, 1.0);\n    gl_Position = trashPos;\n}\x00"

const DefaultFragShaderText = "#version 330 core\nout vec4 FragColor;\n\nuniform vec4 ourColor;\n\nvoid main()\n{\n    FragColor = ourColor;\n}\x00"

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
	vertexShader := createAndCompileShaderOrPanic(DefaultVertexShaderText, gl.VERTEX_SHADER)
	fragShader := createAndCompileShaderOrPanic(DefaultFragShaderText, gl.FRAGMENT_SHADER)

	return createAndLinkProgramOrPanic(vertexShader, fragShader)
}

const EllipseFragShaderText = "#version 330 core\nout vec4 FragColor;\n\nuniform vec4 ourColor;\n\nuniform vec3 tlPos;\nuniform vec3 brPos;\n\nin vec4 trashPos;\n\nvoid main()\n{\n    float a = (brPos.x - tlPos.x) / 2; //0.1\n    float b = (brPos.y - tlPos.y) / 2; //0.1\n    float x = trashPos.x;\n    float y = trashPos.y;\n    float xc = tlPos.x + a; // 0.3\n    float yc = tlPos.y + b; // -0.3\n    float x2 = x - xc;\n    float y2 = y - yc;\n    float c = x2 / a;\n    float d = y2 / b;\n    float res = c * c + d * d;\n    if (res <= 1)\n        FragColor = ourColor;\n    else\n        discard;\n}\x00"

const ImageVertexShader = "#version 330 core\nlayout (location = 0) in vec3 aPos;\nlayout (location = 1) in vec2 aTexCoord;\n\nout vec2 texCoord;\n\nvoid main()\n{\n    gl_Position = vec4(aPos, 1.0);\n    texCoord = aTexCoord;\n}\x00"
const ImageFragShader = "#version 330 core\n\nout vec4 FragColor;\n\nin vec2 texCoord;\n\nuniform sampler2D ourTexture;\n\nvoid main()\n{\n    FragColor = texture(ourTexture, texCoord);\n}\x00"
