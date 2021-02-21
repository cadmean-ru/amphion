package rendering

import (
	"github.com/cadmean-ru/amphion/common"
	"github.com/cadmean-ru/amphion/common/a"
)

type Transform struct {
	Position a.IntVector3
	Rotation a.IntVector3
	Size     a.IntVector3
}

func (t Transform) ToMap() map[string]interface{} {
	return map[string]interface{} {
		"position": t.Position.ToMap(),
		"rotation": t.Rotation.ToMap(),
		"size":     t.Size.ToMap(),
	}
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
		Rotation: a.ZeroIntVector(),
		Size:     a.OneIntVector(),
	}
}
