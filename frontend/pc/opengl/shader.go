// +build windows linux darwin
// +build !android
// +build !ios

package opengl

import (
	_ "embed"
)

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

//go:embed shaders/ShapeFragShader.glsl
var ShapeFragShaderStr string

//go:embed shaders/ShapeVertexShader.glsl
var ShapeVertexShaderStr string

//go:embed shaders/TriangleFragShader.glsl
var TriangleFragShaderStr string

//go:embed shaders/TextFragShader.glsl
var TextFragShaderStr string

//go:embed shaders/TextVertexShader.glsl
var TextVertexShaderStr string

//go:embed shaders/CommonVertexShader.glsl
var CommonVertexShaderStr string

//go:embed shaders/CommonFragmentShader.glsl
var CommonFragmentShaderStr string

//go:embed shaders/PolygonVertex.glsl
var PolygonVertexShaderStr string

//go:embed shaders/PolygonFragment.glsl
var PolygonFragmentShaderStr string