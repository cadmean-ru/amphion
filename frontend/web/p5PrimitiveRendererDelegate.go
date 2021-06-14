//+build js

package web

import (
	"github.com/cadmean-ru/amphion/rendering"
)

type p5drawDelegate func(p5 *p5, primitive rendering.Primitive)

type P5PrimitiveRendererDelegate struct {
	drawFunc p5drawDelegate
	p5       *p5
}

func (r *P5PrimitiveRendererDelegate) OnStart() {

}

func (r *P5PrimitiveRendererDelegate) OnSetPrimitive(ctx *rendering.PrimitiveRenderingContext) {

}

func (r *P5PrimitiveRendererDelegate) OnRender(ctx *rendering.PrimitiveRenderingContext) {
	r.drawFunc(r.p5, ctx.Primitive)
}

func (r *P5PrimitiveRendererDelegate) OnRemovePrimitive(ctx *rendering.PrimitiveRenderingContext) {

}

func (r *P5PrimitiveRendererDelegate) OnStop() {

}

func newP5PrimitiveRendererDelegate(p5 *p5, delegate p5drawDelegate) *P5PrimitiveRendererDelegate {
	return &P5PrimitiveRendererDelegate{
		p5:       p5,
		drawFunc: delegate,
	}
}
