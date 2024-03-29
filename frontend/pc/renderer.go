// +build windows linux darwin
// +build !android
// +build !ios

package pc

import (
	"fmt"
	"github.com/cadmean-ru/amphion/common/a"
	"github.com/cadmean-ru/amphion/rendering"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

// OpenGLRenderer implements the RendererDelegate interface for pc
type OpenGLRenderer struct {
	window     *glfw.Window
	wSize      a.IntVector3
	projection [16]float32
	front      *Frontend
	renderers  []*glPrimitiveRenderer
	clipArea   *rendering.ClipArea2D
}

func (r *OpenGLRenderer) OnPrepare() {
	var err error

	r.window.MakeContextCurrent()

	if err = gl.Init(); err != nil {
		panic(err)
	}

	fmt.Println(gl.GoStr(gl.GetString(gl.VERSION)))
	fmt.Println(gl.GoStr(gl.GetString(gl.VENDOR)))
	fmt.Println(gl.GoStr(gl.GetString(gl.RENDERER)))

	//gl.Viewport(0, 0, int32(r.wSize.X), int32(r.wSize.Y))

	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	gl.PixelStorei(gl.UNPACK_ALIGNMENT, 1)

	gl.ClearColor(1, 1, 1, 1)
	//
	//gl.Clear(gl.COLOR_BUFFER_BIT)

	//r.calculateProjection()

	//gl.Viewport(0, 0, 500, 500)

	textRenderer := &TextRenderer{glPrimitiveRenderer: &glPrimitiveRenderer{}}
	rectangleRenderer := &RectangleRenderer{glPrimitiveRenderer: &glPrimitiveRenderer{}}
	ellipseRenderer := &EllipseRenderer{glPrimitiveRenderer: &glPrimitiveRenderer{}}
	triangleRenderer := &TriangleRenderer{glPrimitiveRenderer: &glPrimitiveRenderer{}}
	lineRenderer := &LineRenderer{glPrimitiveRenderer: &glPrimitiveRenderer{}}
	imageRenderer := &ImageRenderer{glPrimitiveRenderer: &glPrimitiveRenderer{}}

	r.front.renderer.RegisterPrimitiveRendererDelegate(rendering.PrimitiveText, textRenderer)
	r.front.renderer.RegisterPrimitiveRendererDelegate(rendering.PrimitiveRectangle, rectangleRenderer)
	r.front.renderer.RegisterPrimitiveRendererDelegate(rendering.PrimitiveEllipse, ellipseRenderer)
	r.front.renderer.RegisterPrimitiveRendererDelegate(rendering.PrimitiveTriangle, triangleRenderer)
	r.front.renderer.RegisterPrimitiveRendererDelegate(rendering.PrimitiveLine, lineRenderer)
	r.front.renderer.RegisterPrimitiveRendererDelegate(rendering.PrimitiveImage, imageRenderer)

	r.renderers = []*glPrimitiveRenderer{
		textRenderer.glPrimitiveRenderer,
		rectangleRenderer.glPrimitiveRenderer,
		ellipseRenderer.glPrimitiveRenderer,
		triangleRenderer.glPrimitiveRenderer,
		lineRenderer.glPrimitiveRenderer,
		imageRenderer.glPrimitiveRenderer,
	}

	r.clipArea = rendering.NewClipArea2DEmpty()
}

func (r *OpenGLRenderer) OnPerformRenderingStart() {
	//fmt.Println("Start Rendering")
	//fmt.Println(engine.GetInstance().GetCurrentScene().Transform.Size)
	//fmt.Println(r.wSize)

	//gl.Viewport(0, 0, int32(r.wSize.X), int32(r.wSize.Y))
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
}

func (r *OpenGLRenderer) OnPerformRenderingEnd() {
	r.window.SwapBuffers()

	//fmt.Println("End rendering")
}

func (r *OpenGLRenderer) OnClear() {
	fmt.Println("OpenGL renderer clear")

	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
}

func (r *OpenGLRenderer) OnStop() {
	fmt.Println("OpenGL renderer stop")
}

func (r *OpenGLRenderer) handleWindowResize(w, h int) {
	//r.wSize = a.NewIntVector3(w, h, 0)
	//fmt.Printf("OpenGL renderer: handle resize: %d, %d\n", w, h)
	gl.Viewport(0, 0, int32(w), int32(h))
	//r.calculateProjection()
}

//func (r *OpenGLRenderer) calculateProjection() {
//	xs := float32(r.wSize.X)
//	ys := float32(r.wSize.Y)
//	c1 := 2 / xs
//	c2 := 2 / ys
//
//	r.projection = [16]float32 {
//		c1, 0,  0,  -1,
//		0,  c2, 0,  -1,
//		0,  0,  1,  0,
//		0,  0,  0,  1,
//	}
//
//	fmt.Println(r.projection)
//
//	for _, renderer := range r.renderers {
//		r.setProjectionUniform(renderer.id)
//	}
//}
//
//func (r *OpenGLRenderer) setProjectionUniform(id uint32) {
//	gl.UseProgram(id)
//	loc := gl.GetUniformLocation(id, gl.Str(zeroTerminatedString("uProjection")))
//	if loc >= 0 {
//		gl.UniformMatrix4fv(loc, 1, false, &r.projection[0])
//	}
//	gl.UseProgram(0)
//}
