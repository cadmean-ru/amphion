package rendering

type PrimitiveRenderingContext struct {
	Renderer      *ARenderer
	Primitive     IPrimitive
	PrimitiveKind byte
	State         interface{}
	Redraw        bool
}

type PrimitiveRendererDelegate interface {
	OnStart()
	OnSetPrimitive(ctx *PrimitiveRenderingContext)
	OnRender(ctx *PrimitiveRenderingContext)
	OnRemovePrimitive(ctx *PrimitiveRenderingContext)
	OnStop()
}
