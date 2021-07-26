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

func IsValidKindValue(value Kind) bool {
	return value >= 0 && value <= Triangle
}