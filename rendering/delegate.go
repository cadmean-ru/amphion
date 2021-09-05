package rendering

type RendererDelegate interface {
	OnPrepare()
	OnCreatePrimitiveRenderingContext(ctx *PrimitiveRenderingContext)
	OnPerformRenderingStart()
	OnPerformRenderingEnd()
	OnClear()
	OnStop()
}
