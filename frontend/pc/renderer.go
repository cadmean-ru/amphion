// +build windows linux darwin
// +build !android

package pc

import (
	"fmt"
	"github.com/cadmean-ru/amphion/common/a"
	"github.com/cadmean-ru/amphion/engine"
	"github.com/cadmean-ru/amphion/rendering"
	"github.com/go-gl/gl/all-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

//go:generate ../../build/darwin/test --generate shaders -i ./shaders -o ./shaders.gen.go --package pc


type OpenGLRenderer struct {
	window         *glfw.Window
	wSize          a.IntVector3
	projection     [16]float32
	front          *Frontend
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

	gl.Viewport(0, 0, int32(r.wSize.X), int32(r.wSize.Y))

	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	gl.PixelStorei(gl.UNPACK_ALIGNMENT, 1)

	gl.ClearColor(1, 1, 1, 1)
	//
	//gl.Clear(gl.COLOR_BUFFER_BIT)

	//r.calculateProjection()

	r1 := engine.GetInstance().GetRenderer()
	r1.RegisterPrimitiveRendererDelegate(rendering.PrimitiveText, &TextRenderer{glPrimitiveRenderer: &glPrimitiveRenderer{}})
}

func (r *OpenGLRenderer) OnPerformRenderingStart() {
	fmt.Println("Rendering")
	fmt.Println(engine.GetInstance().GetCurrentScene().Transform.Size)

	//gl.ClearColor(1, 1, 1, 1)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
}

func (r *OpenGLRenderer) OnPerformRenderingEnd() {
	r.window.SwapBuffers()
	fmt.Println("End rendering")
}

func (r *OpenGLRenderer) OnClear() {

}

func (r *OpenGLRenderer) OnStop() {

}

func (r *OpenGLRenderer) handleWindowResize(w, h int) {
	//r.calculateProjection()
}



//func (r *OpenGLRenderer) calculateProjection() {
//	xs := float32(r.wSize.X)
//	ys := float32(r.wSize.Y)
//	c1 := 2 / xs
//	c2 := 2 / ys
//	r.projection = [16]float32 {
//		c1, 0,  0,  -1,
//		0,  c2, 0,  -1,
//		0,  0,  1,  0,
//		0,  0,  0,  1,
//	}
//
//	//fmt.Println(r.projection)
//
//	r.setProjectionUniform(r.shapeProgram)
//}
//
//func (r *OpenGLRenderer) setProjectionUniform(program uint32) {
//	gl.UseProgram(r.shapeProgram)
//	gl.UniformMatrix4fv(gl.GetUniformLocation(program, gl.Str("uProjection\x00")), 1, false, &r.projection[0])
//	gl.UseProgram(0)
//}