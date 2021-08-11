// +build windows linux darwin
// +build !android
// +build !ios

package pc

import (
	"fmt"
	"github.com/cadmean-ru/amphion/rendering"
)

type glPrimitiveRenderer struct {
	program  *GlProgram
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

	r.program.Activate()

	if ctx.ClipArea2D != nil {
		r.program.SetClipArea2DUniforms(ctx.ClipArea2D)
	}
}

func (r *glPrimitiveRenderer) OnRemovePrimitive(ctx *rendering.PrimitiveRenderingContext) {
	//fmt.Println("glPrimitiveRenderer OnRemovePrimitive")
	state := ctx.State.(*glPrimitiveState)
	state.free()
}

func (r *glPrimitiveRenderer) OnStop() {
	//fmt.Println("glPrimitiveRenderer OnStop")
	r.program.Delete()
}
