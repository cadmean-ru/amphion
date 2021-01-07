package a

type IntVector2 struct {
	X, Y int
}

func (v IntVector2) Add(other IntVector2) IntVector2 {
	return IntVector2{
		X: v.X + other.X,
		Y: v.Y + other.Y,
	}
}

func (v IntVector2) Sub(other IntVector2) IntVector2 {
	return IntVector2{
		X: v.X - other.X,
		Y: v.Y - other.Y,
	}
}

func (v IntVector2) ToFloat() Vector2 {
	return Vector2{
		X: float32(v.X),
		Y: float32(v.Y),
	}
}

func (v IntVector2) ToFloat3() Vector3 {
	return Vector3{
		X: float32(v.X),
		Y: float32(v.Y),
	}
}

type Int32Vector2 struct {
	X, Y int32
}

type Vector2 struct {
	X, Y float32
}

func NewIntVector2(x, y int) IntVector2 {
	return IntVector2{
		X: x,
		Y: y,
	}
}