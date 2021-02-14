package cli

type RendererCLI interface {
	Prepare()
	AddPrimitive() int
	SetGeometryPrimitive(id int, primitiveData *GeometryPrimitiveData)
	RemovePrimitive(id int)
	PerformRendering()
	Clear()
	Stop()
}

type Vector3 struct {
	X, Y, Z float32
}

type Vector4 struct {
	X, Y, Z, W float32
}

type GeometryPrimitiveData struct {
	GeometryType int
	TlPositionN  *Vector3
	BrPositionN  *Vector3
	FillColorN   *Vector4
	StrokeColorN *Vector4
	StrokeWeight int
	CornerRadius int
}