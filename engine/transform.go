package engine

import (
	"github.com/cadmean-ru/amphion/common"
	"github.com/cadmean-ru/amphion/common/a"
)

type Transform struct {
	Position a.Vector3
	Pivot    a.Vector3
	Rotation a.Vector3
	Size     a.Vector3

	SceneObject *SceneObject
	parent      *Transform
}

func NewTransform(object *SceneObject) Transform {
	return Transform{
		Position:    a.ZeroVector(),
		Pivot:       a.ZeroVector(),
		Rotation:    a.ZeroVector(),
		Size:        a.OneVector(),
		SceneObject: object,
	}
}

func (t Transform) RenderingRepresentation() map[string]interface{} {
	return map[string]interface{} {
		"position": t.Position.ToMap(),
		"pivot":    t.Pivot.ToMap(),
		"rotation": t.Rotation.ToMap(),
		"size":     t.Size.ToMap(),
	}
}

func (t Transform) ToMap() a.SiMap {
	return t.RenderingRepresentation()
}

func (t *Transform) FromMap(siMap a.SiMap) {
	t.Position = a.NewVector3FromMap(a.RequireSiMap(siMap["position"]))
	t.Pivot = a.NewVector3FromMap(a.RequireSiMap(siMap["pivot"]))
	t.Rotation = a.NewVector3FromMap(a.RequireSiMap(siMap["rotation"]))
	t.Size = a.NewVector3FromMap(a.RequireSiMap(siMap["size"]))
}

func (t Transform) GetLocalPosition() a.Vector3 {
	var x, y, z float32
	if t.parent != nil && IsSpecialPosition(t.Position) {
		pb := t.parent.GetRect()
		if t.Position.X == CenterInParent {
			x = pb.X.GetLength() / 2
		} else {
			x = t.Position.X
		}
		if t.Position.Y == CenterInParent {
			y = pb.Y.GetLength() / 2
		} else {
			y = t.Position.Y
		}
		if t.Position.Z == CenterInParent {
			z = pb.Z.GetLength() / 2
		} else {
			z = t.Position.Z
		}
	} else {
		x = t.Position.X
		y = t.Position.Y
		z = t.Position.Z
	}

	return a.NewVector3(x, y, z)
}

func (t Transform) GetGlobalPosition() a.Vector3 {
	if t.parent == nil {
		return t.Position
	}

	return t.parent.GetGlobalPosition().Sub(t.parent.Size.Multiply(t.parent.Pivot)).Add(t.GetLocalPosition())
}

func (t Transform) GetParent() *Transform {
	return t.parent
}

func (t Transform) GetTopLeftPosition() a.Vector3 {
	return t.Position.Sub(t.Size.Multiply(t.Pivot))
}

func (t Transform) GetGlobalTopLeftPosition() a.Vector3 {
	return t.GetGlobalPosition().Sub(t.Size.Multiply(t.Pivot))
}

func (t Transform) GetRect() common.RectBoundary {
	return t.calculateRect(a.ZeroVector())
}

func (t Transform) GetGlobalRect() common.RectBoundary {
	return t.calculateRect(t.GetGlobalTopLeftPosition())
}

func (t Transform) calculateRect(tlp a.Vector3) common.RectBoundary {
	minX := tlp.X
	maxX := tlp.X + t.Size.X
	minY := tlp.Y
	maxY := tlp.Y + t.Size.Y
	minZ := tlp.Z
	maxZ := tlp.Z + t.Size.Z
	return common.NewRectBoundary(minX, maxX, minY, maxY, minZ, maxZ)
}

func IsSpecialPosition(pos a.Vector3) bool {
	return IsSpecialPositionValue(pos.X) || IsSpecialPositionValue(pos.Y) || IsSpecialPositionValue(pos.Z)
}

func IsSpecialPositionValue(x float32) bool {
	return x == CenterInParent
}

func NewTransformFromMap(siMap a.SiMap) Transform {
	var t Transform
	t.FromMap(siMap)
	return t
}

const (
	MatchParent    = -1
	WrapContent    = -2
	CenterInParent = -3
)
