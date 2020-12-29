package a

import (
	"github.com/cadmean-ru/amphion/common"
	"github.com/cadmean-ru/amphion/common/require"
)

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

func (c Color) ToMap() common.SiMap {
	return map[string]interface{}{
		"r": c.R,
		"g": c.G,
		"b": c.B,
		"a": c.A,
	}
}

func (c *Color) FromMap(siMap common.SiMap) {
	c.R = require.Byte(siMap["r"])
	c.G = require.Byte(siMap["g"])
	c.B = require.Byte(siMap["b"])
	c.A = require.Byte(siMap["a"])
}

func (c Color) EncodeToByteArray() []byte {
	arr := make([]byte, 4)
	arr[0] = c.R
	arr[1] = c.G
	arr[2] = c.B
	arr[3] = c.A
	return arr
}

func (c Color) Normalize() common.Vector4 {
	x := float32(c.R) / 255
	y := float32(c.G) / 255
	z := float32(c.B) / 255
	w := float32(c.A) / 255
	return common.NewVector4(x, y, z, w)
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