package cli

import (
	"github.com/cadmean-ru/amphion/rendering"
)

type RendererDelegate interface {
	OnPrepare()
	//OnCreatePrimitiveRenderingContext(ctx *PrimitiveRenderingContext)
	OnPerformRenderingStart()
	OnPerformRenderingEnd()
	OnClear()
	OnStop()
}

type RendererDelegateWrap struct {
	d RendererDelegate
}

func (r *RendererDelegateWrap) OnPrepare() {
	r.d.OnPrepare()
}

func (r *RendererDelegateWrap) OnCreatePrimitiveRenderingContext(ctx *rendering.PrimitiveRenderingContext) {

}

func (r *RendererDelegateWrap) OnPerformRenderingStart() {
	r.d.OnPerformRenderingStart()
}

func (r *RendererDelegateWrap) OnPerformRenderingEnd() {
	r.d.OnPerformRenderingEnd()
}

func (r *RendererDelegateWrap) OnClear() {
	r.d.OnClear()
}

func (r *RendererDelegateWrap) OnStop() {
	r.d.OnStop()
}

func NewRendererDelegateWrap(delegate RendererDelegate) *RendererDelegateWrap {
	return &RendererDelegateWrap{d: delegate}
}