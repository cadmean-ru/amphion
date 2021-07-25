package shape

type Kind int

const (
	Empty Kind = iota
	Rectangle
	RoundedRectangle
	Circle
	Ellipse
	Triangle
)
