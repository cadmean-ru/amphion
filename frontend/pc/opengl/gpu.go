package opengl

import (
	"fmt"
	"github.com/cadmean-ru/amphion/common"
	"github.com/cadmean-ru/amphion/rendering/gpu"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

type Gpu struct {
	window *glfw.Window

}

func (g *Gpu) Init() {
	var err error

	g.window.MakeContextCurrent()

	if err = gl.Init(); err != nil {
		panic(err)
	}

	fmt.Println(gl.GoStr(gl.GetString(gl.VERSION)))
	fmt.Println(gl.GoStr(gl.GetString(gl.VENDOR)))
	fmt.Println(gl.GoStr(gl.GetString(gl.RENDERER)))

	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	gl.PixelStorei(gl.UNPACK_ALIGNMENT, 1)

	gl.ClearColor(1, 1, 1, 1)
}

func (g *Gpu) AllocateBuffer(size int) gpu.Buffer {
	return nil
}

func (g *Gpu) SetClippingArea(rect *common.RectBoundary) {

}

func (g *Gpu) ClearClippingArea() {

}

func NewOpenGlGpu() *Gpu {
	return &Gpu{

	}
}