package cli

import (
	"fmt"
	"github.com/cadmean-ru/amphion/rendering"
)

type PrimitiveRendererDelegateWrap struct {
	r PrimitiveRendererDelegate
}

func (p *PrimitiveRendererDelegateWrap) OnStart() {
	p.r.OnStart()
}

func (p *PrimitiveRendererDelegateWrap) OnSetPrimitive(ctx *rendering.PrimitiveRenderingContext) {
	p.r.OnSetPrimitive(newCliPrimitiveRenderingContext(ctx))
}

func (p *PrimitiveRendererDelegateWrap) OnRender(ctx *rendering.PrimitiveRenderingContext) {
	fmt.Printf("Go rendering wrap %d\n", ctx.PrimitiveId)
	p.r.OnRender(newCliPrimitiveRenderingContext(ctx))
}

func (p *PrimitiveRendererDelegateWrap) OnRemovePrimitive(ctx *rendering.PrimitiveRenderingContext) {
	p.r.OnRemovePrimitive(newCliPrimitiveRenderingContext(ctx))
}

func (p *PrimitiveRendererDelegateWrap) OnStop() {
	p.r.OnStop()
}

func NewPrimitiveRendererDelegateWrap(delegate PrimitiveRendererDelegate) *PrimitiveRendererDelegateWrap {
	return &PrimitiveRendererDelegateWrap{r: delegate}
}