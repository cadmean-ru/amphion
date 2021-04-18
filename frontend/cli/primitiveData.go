package cli

type GeometryPrimitiveData struct {
	GeometryType int
	TlPositionN  *Vector3
	BrPositionN  *Vector3
	FillColorN   *Vector4
	StrokeColorN *Vector4
	StrokeWeight int
	CornerRadius int
}

type ImagePrimitiveData struct {
	TlPositionN *Vector3
	BrPositionN *Vector3
	ImageUrl    string
}

type TextPrimitiveData struct {
	Text       string
	TlPosition *Vector3
	Size       *Vector3
	TextColorN *Vector4
}
