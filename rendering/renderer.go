package rendering

type Renderer interface {
	Prepare()
	AddPrimitive() int64
	SetPrimitive(id int64, primitive IPrimitive, shouldRerender bool)
	RemovePrimitive(id int64)
	PerformRendering()
	Clear()
	Stop()
}