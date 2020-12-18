package rendering

type Renderer interface {
	Prepare()
	AddPrimitive() int64
	SetPrimitive(id int64, primitive interface{}, shouldRerender bool)
	RemovePrimitive(id int64)
	PerformRendering()
	Clear()
	Stop()
}