package pc

import (
	"github.com/cadmean-ru/amphion/rendering"
	"github.com/go-gl/gl/all-core/gl"
)

type glPrimitiveRenderer struct {
	program uint32
}

func (r *glPrimitiveRenderer) OnStart() {

}

func (r *glPrimitiveRenderer) OnSetPrimitive(ctx *rendering.PrimitiveRenderingContext) {
	if ctx.State == nil {
		ctx.State = &glPrimitiveState{}
	}
}

func (r *glPrimitiveRenderer) OnRender(ctx *rendering.PrimitiveRenderingContext) {

}

func (r *glPrimitiveRenderer) OnRemovePrimitive(ctx *rendering.PrimitiveRenderingContext) {
	state := ctx.State.(*glPrimitiveState)
	state.free()
}

func (r *glPrimitiveRenderer) OnStop() {
	gl.DeleteProgram(r.program)
}

