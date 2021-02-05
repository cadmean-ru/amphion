package a

import (
	"fmt"
	"strings"
)

// Color represents a color in RGBA format.
type Color struct {
	R, G, B, A byte
}

// Creates new Color struct using the specified arguments.
// This function can accept the following arguments:
// - a color hex string,
// - 1 byte value for grayscale color,
// - 3 byte values for RGB color,
// - 4 byte values for RGBA color.
// If the given arguments cannot be interpreted as a color returns black color.
func NewColor(params ...interface{}) Color {
	switch len(params) {
	case 1:
		if hex, ok := params[0].(string); ok {
			return ParseHexColor(hex)
		} else if b, ok := params[0].(byte); ok {
			return Color{
				R: b,
				G: b,
				B: b,
				A: 255,
			}
		} else {
			return BlackColor()
		}
	case 3:
		var r, g, b, a byte = requireByte(params[0]), requireByte(params[1]), requireByte(params[2]), 255
		return Color{
			R: r,
			G: g,
			B: b,
			A: a,
		}
	case 4:
		var r, g, b, a = requireByte(params[0]), requireByte(params[1]), requireByte(params[2]), requireByte(params[3])
		return Color{
			R: r,
			G: g,
			B: b,
			A: a,
		}
	default:
		return BlackColor()
	}
}

// Parses color hex string in format #rrggbbaa, #rrggbb, #rgba or #rgb to a Color struct.
func ParseHexColor(hex string) Color {
	c := BlackColor()

	if !strings.HasPrefix(hex, "#") {
		return c
	}

	switch len(hex) {
	case 9:
		_, _ = fmt.Sscanf(hex, "#%02x%02x%02x%02x", &c.R, &c.G, &c.B, &c.A)
	case 7:
		_, _ = fmt.Sscanf(hex, "#%02x%02x%02x", &c.R, &c.G, &c.B)
	case 5:
		_, _ = fmt.Sscanf(hex, "#%1x%1x%1x%1x", &c.R, &c.G, &c.B, &c.A)
	case 4:
		_, _ = fmt.Sscanf(hex, "#%1x%1x%1x", &c.R, &c.G, &c.B)
		c.R *= 17
		c.G *= 17
		c.B *= 17
	}

	return c
}

// Returns a color hex string in format #rrggbbaa or #rrggbb.
func (c *Color) GetHex() string {
	if c.A == 255 {
		return fmt.Sprintf("#%02x%02x%02x", c.R, c.G, c.B)
	} else {
		return fmt.Sprintf("#%02x%02x%02x%02x", c.R, c.G, c.B, c.A)
	}
}

func (c *Color) Equals(other interface{}) bool {
	if otherColor, ok := other.(Color); ok {
		return c.R == otherColor.R && c.G == otherColor.G && c.B == otherColor.B && c.A == otherColor.A
	} else if otherColorPtr, ok := other.(*Color); ok {
		return c.R == otherColorPtr.R && c.G == otherColorPtr.G && c.B == otherColorPtr.B && c.A == otherColorPtr.A
	} else {
		return false
	}
}

//region Stringable implementation

func (c *Color) ToString() string {
	return c.GetHex()
}

func (c *Color) FromString(src string) {
	c1 := ParseHexColor(src)
	c.R = c1.R
	c.G = c1.G
	c.B = c1.B
	c.A = c1.A
}

//endregion

//region Mappable implementation

func (c *Color) ToMap() SiMap {
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

//endregion

func (c *Color) EncodeToByteArray() []byte {
	arr := make([]byte, 4)
	arr[0] = c.R
	arr[1] = c.G
	arr[2] = c.B
	arr[3] = c.A
	return arr
}

// Returns Vector4 with normalized color values.
func (c *Color) Normalize() Vector4 {
	x := float32(c.R) / 255
	y := float32(c.G) / 255
	z := float32(c.B) / 255
	w := float32(c.A) / 255
	return NewVector4(x, y, z, w)
}

//region color presets

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

//endregion

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