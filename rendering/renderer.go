package rendering

type Renderer interface {
	Prepare()
	AddPrimitive() int
	SetPrimitive(id int, primitive IPrimitive, shouldRerender bool)
	RemovePrimitive(id int)
	PerformRendering()
	Clear()
	Stop()
}