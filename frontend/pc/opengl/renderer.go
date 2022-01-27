//go:build (windows || linux || darwin) && !android && !ios
// +build windows linux darwin
// +build !android
// +build !ios

package opengl

import (
	"fmt"
	"github.com/cadmean-ru/amphion/common/a"
	"github.com/cadmean-ru/amphion/rendering"
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

//Renderer implements the RendererDelegate interface for pc
type Renderer struct {
	Window     *glfw.Window
	WSize      a.IntVector3
	ARenderer  *rendering.ARenderer
	projection a.Matrix4
	renderers  []*glPrimitiveRenderer
	clipArea   *rendering.ClipArea2D
}

func (r *Renderer) OnPrepare() {
	var err error

	r.Window.MakeContextCurrent()

	if err = gl.Init(); err != nil {
		panic(err)
	}

	fmt.Println(gl.GoStr(gl.GetString(gl.VERSION)))
	fmt.Println(gl.GoStr(gl.GetString(gl.VENDOR)))
	fmt.Println(gl.GoStr(gl.GetString(gl.RENDERER)))

	//gl.Viewport(0, 0, int32(r.WSize.X), int32(r.WSize.Y))

	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	gl.PixelStorei(gl.UNPACK_ALIGNMENT, 1)

	gl.ClearColor(1, 1, 1, 1)
	//
	//gl.Clear(gl.COLOR_BUFFER_BIT)

	//gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE)

	r.calculateProjection()

	//gl.Viewport(0, 0, 500, 500)

	textRenderer := &TextRenderer{glPrimitiveRenderer: &glPrimitiveRenderer{}}
	rectangleRenderer := &RectangleRenderer{glPrimitiveRenderer: &glPrimitiveRenderer{}}
	ellipseRenderer := &EllipseRenderer{glPrimitiveRenderer: &glPrimitiveRenderer{}}
	triangleRenderer := &TriangleRenderer{glPrimitiveRenderer: &glPrimitiveRenderer{}}
	lineRenderer := &LineRenderer{glPrimitiveRenderer: &glPrimitiveRenderer{}}
	imageRenderer := &ImageRenderer{glPrimitiveRenderer: &glPrimitiveRenderer{}}
	polygonRenderer := &PolygonRenderer{glPrimitiveRenderer: &glPrimitiveRenderer{}}

	r.ARenderer.RegisterPrimitiveRendererDelegate(rendering.PrimitiveText, textRenderer)
	r.ARenderer.RegisterPrimitiveRendererDelegate(rendering.PrimitiveRectangle, rectangleRenderer)
	r.ARenderer.RegisterPrimitiveRendererDelegate(rendering.PrimitiveEllipse, ellipseRenderer)
	r.ARenderer.RegisterPrimitiveRendererDelegate(rendering.PrimitiveTriangle, triangleRenderer)
	r.ARenderer.RegisterPrimitiveRendererDelegate(rendering.PrimitiveLine, lineRenderer)
	r.ARenderer.RegisterPrimitiveRendererDelegate(rendering.PrimitiveImage, imageRenderer)
	r.ARenderer.RegisterPrimitiveRendererDelegate(rendering.PrimitivePolygon, polygonRenderer)

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

func (r *Renderer) OnCreatePrimitiveRenderingContext(ctx *rendering.PrimitiveRenderingContext) {
	ctx.Projection = r.projection
}

func (r *Renderer) OnPerformRenderingStart() {
	//fmt.Println("Start Rendering")
	//fmt.Println(engine.GetInstance().GetCurrentScene().Transform.size)
	//fmt.Println(r.WSize)

	//gl.Viewport(0, 0, int32(r.WSize.X), int32(r.WSize.Y))
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
}

func (r *Renderer) OnPerformRenderingEnd() {
	r.Window.SwapBuffers()

	//fmt.Println("End rendering")
}

func (r *Renderer) OnClear() {
	fmt.Println("OpenGL renderer clear")

	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
}

func (r *Renderer) OnStop() {
	fmt.Println("OpenGL renderer stop")
}

func (r *Renderer) HandleWindowResize(w, h int) {
	//r.WSize = a.NewIntVector3(w, h, 0)
	//fmt.Printf("OpenGL renderer: handle resize: %d, %d\n", w, h)
	gl.Viewport(0, 0, int32(w), int32(h))
	r.calculateProjection()
}

func (r *Renderer) calculateProjection() {
	xs := float32(r.WSize.X)
	ys := float32(r.WSize.Y)
	zs := float32(2)
	c1 := 2 / xs
	c2 := -2 / ys
	c3 := -2 / zs

	r.projection = a.Matrix4{
		c1, 0, 0, -1,
		0, c2, 0, 1,
		0, 0, c3, 0,
		0, 0, 0, 1,
	}

	//fmt.Println(r.projection)
	//
	//fmt.Println(r.projection.MulVector(a.NewVector4(float32(r.WSize.X), float32(r.WSize.X), 0, 1)))
	//fmt.Println(r.projection.MulVector(a.NewVector4(250, 250, 0, 1)))
}
