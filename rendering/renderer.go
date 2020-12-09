package rendering

type Renderer interface {
	Prepare()
	AddPrimitive() int64
	SetPrimitive(id int64, primitive PrimitiveBuilder, shouldRerender bool)
	RemovePrimitive(id int64)
	PerformRendering()
	Clear()
}