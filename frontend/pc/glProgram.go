package pc

import (
	"fmt"
	"github.com/cadmean-ru/amphion/common/a"
	"github.com/cadmean-ru/amphion/engine"
	"github.com/cadmean-ru/amphion/rendering"
	"github.com/go-gl/gl/v4.1-core/gl"
	"strings"
)

type GlProgram struct {
	vertexShaderSource   string
	fragmentShaderSource string
	compiled             bool
	vertexShader         uint32
	fragmentShader       uint32
	id                   uint32
	deleted              bool
	tag                  string
	uniforms             map[string]int32
}

func (p *GlProgram) CompileAndLink() {
	p.vertexShader = p.compileShader(p.vertexShaderSource, gl.VERTEX_SHADER)
	p.fragmentShader = p.compileShader(p.fragmentShaderSource, gl.FRAGMENT_SHADER)

	p.id = p.linkProgram(p.vertexShader, p.fragmentShader)

	p.vertexShaderSource = ""
	p.fragmentShaderSource = ""

	p.compiled = true
}

func (p *GlProgram) compileShader(source string, shaderType uint32) uint32 {
	source = p.preprocessShaderSource(source, shaderType)
	shaderSource, free := gl.Strs(source)

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

		fmt.Println(infoLog)
		panic(fmt.Sprintf("failed to compile shader for program '%s'", p.tag))
	}

	return shader
}

func (p *GlProgram) preprocessShaderSource(shaderSource string, shaderType uint32) string {
	switch shaderType {
	case gl.VERTEX_SHADER:
		shaderSource = CommonVertexShaderStr + shaderSource
	case gl.FRAGMENT_SHADER:
		shaderSource = CommonFragmentShaderStr + shaderSource
	}

	if !strings.HasSuffix(shaderSource, "\x00") {
		shaderSource = zeroTerminatedString(shaderSource)
	}

	return shaderSource
}

func (p *GlProgram) linkProgram(vShader, fShader uint32) uint32 {
	shaderProgram := gl.CreateProgram()

	gl.AttachShader(shaderProgram, vShader)
	gl.AttachShader(shaderProgram, fShader)

	gl.LinkProgram(shaderProgram)

	var success int32

	gl.GetProgramiv(shaderProgram, gl.LINK_STATUS, &success)
	if success == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(shaderProgram, gl.INFO_LOG_LENGTH, &logLength)

		infoLog := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(shaderProgram, 512, nil, gl.Str(infoLog))

		fmt.Println(infoLog)
		panic(fmt.Sprintf("failed to link program '%s'", p.tag))
	}

	gl.DeleteShader(vShader)
	gl.DeleteShader(fShader)

	return shaderProgram
}

func (p *GlProgram) Delete() {
	gl.DeleteProgram(p.id)
	p.deleted = true
}

func (p *GlProgram) IsUsable() bool {
	return p.compiled && !p.deleted
}

func (p *GlProgram) Use() {
	if !p.IsUsable() {
		panic(fmt.Sprintf("program '%s' is not usable", p.tag))
	}

	gl.UseProgram(p.id)
}

func (p *GlProgram) SetClipArea2DUniforms(area *rendering.ClipArea2D) {
	clipRectUniform := p.GetUniformLocation("uClipArea2dRect")
	clipShapeUniform := p.GetUniformLocation("uClipArea2dShape")


	wSize := engine.GetScreenSize3()
	tl := a.NewIntVector3(int(area.Rect.X.Min), int(area.Rect.Y.Min), 0).Ndc(wSize)
	br := a.NewIntVector3(int(area.Rect.X.Max), int(area.Rect.Y.Max), 0).Ndc(wSize)

	gl.Uniform4f(clipRectUniform, tl.X, tl.Y, br.X, br.Y)
	gl.Uniform1i(clipShapeUniform, int32(area.Shape))
}

func (p *GlProgram) GetUniformLocation(uName string) int32 {
	if loc, ok := p.uniforms[uName]; ok {
		return loc
	}

	loc := gl.GetUniformLocation(p.id, zeroTerminatedGlString(uName))

	if loc == -1 {
		panic(fmt.Sprintf("'%s' uniform cannot be found in program '%s'", uName, p.tag))
	}

	p.uniforms[uName] = loc

	return loc
}

func NewGlProgram(vertexShader, fragmentShader, tag string) *GlProgram {
	return &GlProgram{
		vertexShaderSource:   vertexShader,
		fragmentShaderSource: fragmentShader,
		tag:                  tag,
		uniforms:             map[string]int32{},
	}
}

func zeroTerminatedString(s string) string {
	if strings.HasSuffix(s, "\x00") {
		return s
	}

	return s + "\x00"
}

func zeroTerminatedGlString(s string) *uint8 {
	return gl.Str(zeroTerminatedString(s))
}
