package rendering

import "github.com/cadmean-ru/amphion/common"

const transformBytesSize = 48

type Transform struct {
	Position common.IntVector3
	Pivot    common.IntVector3
	Rotation common.IntVector3
	Size     common.IntVector3
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
	_ = common.CopyByteArray(t.Position.EncodeToByteArray(), bytes, 0,  12)
	_ = common.CopyByteArray(t.Pivot.EncodeToByteArray(),    bytes, 12, 12)
	_ = common.CopyByteArray(t.Rotation.EncodeToByteArray(), bytes, 24, 12)
	_ = common.CopyByteArray(t.Size.EncodeToByteArray(),     bytes, 36, 12)
	return bytes
}

func NewTransform() Transform {
	return Transform{
		Position:    common.ZeroIntVector(),
		Pivot:       common.ZeroIntVector(),
		Rotation:    common.ZeroIntVector(),
		Size:        common.OneIntVector(),
	}
}
