//+build js

package web

import "github.com/cadmean-ru/amphion/rendering"

type P5RendererDelegate struct {
	p5        *p5
	aRenderer *rendering.ARenderer
}

func (r *P5RendererDelegate) OnPrepare() {
	r.p5.prepare()
	r.p5.onDraw = r.drawP5

	r.aRenderer.RegisterPrimitiveRendererDelegate(rendering.PrimitivePoint, newP5PrimitiveRendererDelegate(r.p5, drawPoint))
	r.aRenderer.RegisterPrimitiveRendererDelegate(rendering.PrimitiveLine, newP5PrimitiveRendererDelegate(r.p5, drawLine))
	r.aRenderer.RegisterPrimitiveRendererDelegate(rendering.PrimitiveRectangle, newP5PrimitiveRendererDelegate(r.p5, drawRectangle))
	r.aRenderer.RegisterPrimitiveRendererDelegate(rendering.PrimitiveEllipse, newP5PrimitiveRendererDelegate(r.p5, drawEllipse))
	r.aRenderer.RegisterPrimitiveRendererDelegate(rendering.PrimitiveTriangle, newP5PrimitiveRendererDelegate(r.p5, drawTriangle))
	r.aRenderer.RegisterPrimitiveRendererDelegate(rendering.PrimitiveText, newP5PrimitiveRendererDelegate(r.p5, drawText))
	r.aRenderer.RegisterPrimitiveRendererDelegate(rendering.PrimitiveImage, newP5PrimitiveRendererDelegate(r.p5, drawImage))
}

func (r *P5RendererDelegate) OnPerformRenderingStart() {
	r.p5.redraw()
}

func (r *P5RendererDelegate) OnPerformRenderingEnd() {

}

func (r *P5RendererDelegate) OnClear() {
	r.p5.clear()
}

func (r *P5RendererDelegate) OnStop() {

}

func (r *P5RendererDelegate) drawP5(p5 *p5) {
	p5.clear()
	p5.rectModeCorner()

	r.aRenderer.GetRenderingPerformer()()
}

func newP5RendererDelegate() *P5RendererDelegate {
	return &P5RendererDelegate{
		p5:       &p5{},
	}
}
