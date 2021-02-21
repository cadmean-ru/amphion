package rendering

import (
	"github.com/cadmean-ru/amphion/common"
	"github.com/cadmean-ru/amphion/common/a"
)

const transformBytesSize = 48

type Transform struct {
	Position a.IntVector3
	Pivot    a.IntVector3
	Rotation a.IntVector3
	Size     a.IntVector3
}

func (t Transform) ToMap() map[string]interface{} {
	return map[string]interface{} {
		"position": t.Position.ToMap(),
		"pivot":    t.Pivot.ToMap(),
		"rotation": t.Rotation.ToMap(),
		"size":     t.Size.ToMap(),
	}
}

func (t Transform) EncodeToByteArray() []byte {
	bytes := make([]byte, transformBytesSize)
	_ = a.CopyByteArray(t.Position.EncodeToByteArray(), bytes, 0,  12)
	_ = a.CopyByteArray(t.Pivot.EncodeToByteArray(),    bytes, 12, 12)
	_ = a.CopyByteArray(t.Rotation.EncodeToByteArray(), bytes, 24, 12)
	_ = a.CopyByteArray(t.Size.EncodeToByteArray(),     bytes, 36, 12)
	return bytes
}

func (t Transform) GetRect() common.RectBoundary {
	tlp := t.Position
	minX := tlp.X
	maxX := tlp.X + t.Size.X
	minY := tlp.Y
	maxY := tlp.Y + t.Size.Y
	minZ := tlp.Z
	maxZ := tlp.Z + t.Size.Z
	return common.NewRectBoundary(float32(minX), float32(maxX), float32(minY), float32(maxY), float32(minZ), float32(maxZ))
}

func NewTransform() Transform {
	return Transform{
		Position: a.ZeroIntVector(),
		Pivot:    a.ZeroIntVector(),
		Rotation: a.ZeroIntVector(),
		Size:     a.OneIntVector(),
	}
}
