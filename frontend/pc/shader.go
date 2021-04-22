// +build windows linux darwin
// +build !android
// +build !ios

package pc

import (
	_ "embed"
	"github.com/go-gl/gl/v4.1-core/gl"
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

func createAndLinkProgramOrPanic(shaders ...uint32) uint32 {
	shaderProgram := gl.CreateProgram()

	for _, shader := range shaders {
		gl.AttachShader(shaderProgram, shader)
	}

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

	for _, shader := range shaders {
		gl.DeleteShader(shader)
	}

	return shaderProgram
}

func createAndLinkDefaultProgramOrPanic() uint32 {
	vertexShader := createAndCompileShaderOrPanic(zeroTerminated(DefaultVertexShaderStr), gl.VERTEX_SHADER)
	fragShader := createAndCompileShaderOrPanic(zeroTerminated(DefaultFragShaderStr), gl.FRAGMENT_SHADER)

	return createAndLinkProgramOrPanic(vertexShader, fragShader)
}

func zeroTerminated(s string) string {
	return s + "\x00"
}

//go:embed shaders/DefaultFragShader.glsl
var DefaultFragShaderStr string

//go:embed shaders/DefaultVertexShader.glsl
var DefaultVertexShaderStr string

//go:embed shaders/EllipseFragShader.glsl
var EllipseFragShaderStr string

//go:embed shaders/EllipseVertexShader.glsl
var EllipseVertexShaderStr string

//go:embed shaders/ImageFragShader.glsl
var ImageFragShaderStr string

//go:embed shaders/ImageVertexShader.glsl
var ImageVertexShaderStr string

//go:embed shaders/RectFragShader.glsl
var RectFragShaderStr string

//go:embed shaders/ShapeFragShader.glsl
var ShapeFragShaderStr string

//go:embed shaders/ShapeVertexShader.glsl
var ShapeVertexShaderStr string

//go:embed shaders/TextFragShader.glsl
var TextFragShaderStr string

//go:embed shaders/TextVertexShader.glsl
var TextVertexShaderStr string

//go:embed shaders/CommonVertexShader.glsl
var CommonVertexShaderStr string

//go:embed shaders/CommonFragmentShader.glsl
var CommonFragmentShaderStr string