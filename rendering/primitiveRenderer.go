package rendering

type PrimitiveRenderingContext struct {
	Renderer      *ARenderer
	Primitive     Primitive
	PrimitiveKind byte
	PrimitiveId   int
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
