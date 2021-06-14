// +build windows linux darwin
// +build !android
// +build !ios

package pc

import (
	"fmt"
	"github.com/cadmean-ru/amphion/rendering"
	"github.com/go-gl/gl/v4.1-core/gl"
)

type glPrimitiveRenderer struct {
	program uint32
}

func (r *glPrimitiveRenderer) OnStart() {

}

func (r *glPrimitiveRenderer) OnSetPrimitive(ctx *rendering.PrimitiveRenderingContext) {
	//fmt.Printf("glPrimitiveRenderer OnSetPrimitive id: %d\n", ctx.PrimitiveId)
	if ctx.State == nil {
		ctx.State = &glPrimitiveState{}
	}
}

func (r *glPrimitiveRenderer) OnRender(ctx *rendering.PrimitiveRenderingContext) {
	if ctx.State == nil {
		panic(fmt.Sprintf("OnRender called before OnSetPrimitive was called for id: %d", ctx.PrimitiveId))
	}
}

func (r *glPrimitiveRenderer) OnRemovePrimitive(ctx *rendering.PrimitiveRenderingContext) {
	//fmt.Println("glPrimitiveRenderer OnRemovePrimitive")
	state := ctx.State.(*glPrimitiveState)
	state.free()
}

func (r *glPrimitiveRenderer) OnStop() {
	//fmt.Println("glPrimitiveRenderer OnStop")
	gl.DeleteProgram(r.program)
}

