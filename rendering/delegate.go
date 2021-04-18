package rendering

type RendererDelegate interface {
	OnPrepare()
	OnPerformRenderingStart()
	OnPerformRenderingEnd()
	OnClear()
	OnStop()
}
