package a

type Color struct {
	R, G, B, A byte
}

func NewColor(r, g, b, a byte) Color {
	return Color{
		R: r,
		G: g,
		B: b,
		A: a,
	}
}

func (c Color) ToMap() SiMap {
	return map[string]interface{}{
		"r": c.R,
		"g": c.G,
		"b": c.B,
		"a": c.A,
	}
}

func (c *Color) FromMap(siMap SiMap) {
	c.R = requireByte(siMap["r"])
	c.G = requireByte(siMap["g"])
	c.B = requireByte(siMap["b"])
	c.A = requireByte(siMap["a"])
}

func (c Color) EncodeToByteArray() []byte {
	arr := make([]byte, 4)
	arr[0] = c.R
	arr[1] = c.G
	arr[2] = c.B
	arr[3] = c.A
	return arr
}

func (c Color) Normalize() Vector4 {
	x := float32(c.R) / 255
	y := float32(c.G) / 255
	z := float32(c.B) / 255
	w := float32(c.A) / 255
	return NewVector4(x, y, z, w)
}

func BlackColor() Color {
	return NewColor(0, 0, 0, 255)
}

func WhiteColor() Color {
	return NewColor(255, 255, 255, 255)
}

func RedColor() Color {
	return NewColor(255, 0, 0, 255)
}

func GreenColor() Color {
	return NewColor(0, 255, 0, 255)
}

func BlueColor() Color {
	return NewColor(0, 0, 255, 255)
}

func TransparentColor() Color {
	return NewColor(0, 0, 0, 0)
}

func PinkColor() Color {
	return NewColor(255,192,203, 255)
}

func requireByte(num interface{}) byte {
	switch num.(type) {
	case byte:
		return num.(byte)
	case int:
		return byte(num.(int))
	case int32:
		return byte(num.(int32))
	case int64:
		return byte(num.(int64))
	case float32:
		return byte(num.(float32))
	case float64:
		return byte(num.(float64))
	default:
		return 0
	}
}