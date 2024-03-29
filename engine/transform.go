package engine

import (
	"github.com/cadmean-ru/amphion/common"
	"github.com/cadmean-ru/amphion/common/a"
	"github.com/cadmean-ru/amphion/rendering"
)

// Transform describes how a scene object is positioned on the screen.
type Transform struct {
	Position a.Vector3
	Pivot    a.Vector3
	Rotation a.Vector3
	Size     a.Vector3

	SceneObject *SceneObject
	parent      *Transform
}

// NewTransform2D creates a new transform with default values.
func NewTransform2D(object *SceneObject) Transform {
	return Transform{
		Position:    a.ZeroVector(),
		Pivot:       a.ZeroVector(),
		Rotation:    a.ZeroVector(),
		Size:        a.OneVector(),
		SceneObject: object,
	}
}

func (t Transform) ToMap() a.SiMap {
	return map[string]interface{} {
		"position": t.Position.ToMap(),
		"pivot":    t.Pivot.ToMap(),
		"rotation": t.Rotation.ToMap(),
		"size":     t.Size.ToMap(),
	}
}

func (t *Transform) FromMap(siMap a.SiMap) {
	t.Position = a.NewVector3FromMap(t.decodeSpecialValuesInVector(a.RequireSiMap(siMap["position"])))
	t.Pivot = a.NewVector3FromMap(a.RequireSiMap(siMap["pivot"]))
	t.Rotation = a.NewVector3FromMap(a.RequireSiMap(siMap["rotation"]))
	t.Size = a.NewVector3FromMap(t.decodeSpecialValuesInVector(a.RequireSiMap(siMap["size"])))
}

func (t *Transform) decodeSpecialValuesInVector(siMap a.SiMap) a.SiMap {
	if IsSpecialValueString(siMap["x"]) {
		siMap["x"] = GetSpecialValueFromString(siMap["x"])
	}

	if IsSpecialValueString(siMap["y"]) {
		siMap["y"] = GetSpecialValueFromString(siMap["y"])
	}

	if IsSpecialValueString(siMap["z"]) {
		siMap["z"] = GetSpecialValueFromString(siMap["z"])
	}

	return siMap
}

// GetLocalPosition calculates the actual local position related to this transform's parent.
func (t Transform) GetLocalPosition() a.Vector3 {
	var x, y, z float32

	if t.parent != nil && IsSpecialTransformVector3(t.Position) {
		pb := t.parent.GetRect()

		if t.Position.X == a.CenterInParent {
			x = pb.X.GetLength() / 2
		} else {
			x = t.Position.X
		}

		if t.Position.Y == a.CenterInParent {
			y = pb.Y.GetLength() / 2
		} else {
			y = t.Position.Y
		}

		if t.Position.Z == a.CenterInParent {
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

// GetGlobalPosition calculates the actual global position in the scene.
func (t Transform) GetGlobalPosition() a.Vector3 {
	if t.parent == nil {
		return t.Position
	}

	return t.parent.GetGlobalPosition().Sub(t.parent.Size.Multiply(t.parent.Pivot)).Add(t.GetLocalPosition())
}

// GetParent returns the parent transform of the current transform.
func (t Transform) GetParent() *Transform {
	return t.parent
}

// GetTopLeftPosition calculates local position of the top left point of the bounding box.
func (t Transform) GetTopLeftPosition() a.Vector3 {
	return t.Position.Sub(t.Size.Multiply(t.Pivot))
}

// GetGlobalTopLeftPosition calculates global position of the top left point of the bounding box.
func (t Transform) GetGlobalTopLeftPosition() a.Vector3 {
	return t.GetGlobalPosition().Sub(t.Size.Multiply(t.Pivot))
}

// GetSize calculates the actual size of the Transform replacing the special values.
func (t Transform) GetSize() a.Vector3 {
	var x, y, z float32
	var parentSize a.Vector3
	if t.parent != nil {
		parentSize = t.parent.GetSize()
	}
	var tlp = t.GetTopLeftPosition()

	if IsSpecialTransformVector3(t.Size) {
		if t.Size.X == a.MatchParent {
			x = common.ClampFloat32(parentSize.X, 0, parentSize.X - tlp.X)
		} else {
			x = t.Size.X
		}

		if t.Size.Y == a.MatchParent {
			y = common.ClampFloat32(parentSize.Y, 0, parentSize.Y - tlp.Y)
		} else {
			y = t.Size.Y
		}

		if t.Size.Z == a.MatchParent {
			z = common.ClampFloat32(parentSize.Z, 0, parentSize.Z - tlp.Z)
		} else {
			z = t.Size.Z
		}
	} else {
		x = t.Size.X
		y = t.Size.Y
		z = t.Size.Z
	}

	return a.NewVector3(x, y, z)
}

// GetRect calculates local rect of this transform.
func (t Transform) GetRect() *common.RectBoundary {
	return t.calculateRect(a.ZeroVector())
}

// GetGlobalRect calculates global rect of this transform.
func (t Transform) GetGlobalRect() *common.RectBoundary {
	return t.calculateRect(t.GetGlobalTopLeftPosition())
}

// ToRenderingTransform calculates a transform that is ready to be rendered on screen with all absolute values in pixels calculated.
func (t *Transform) ToRenderingTransform() rendering.Transform {
	rt := rendering.NewTransform()

	rt.Position = t.GetGlobalTopLeftPosition().Round()
	rt.Size = t.GetSize().Round()

	return rt
}

func (t Transform) Equals(other Transform) bool {
	return t.Position.Equals(other.Position) && t.Size.Equals(other.Size) && t.Pivot.Equals(other.Pivot) && t.Rotation.Equals(other.Rotation)
}

// Calculates the rect boundary of the transform given the top left point's position.
func (t Transform) calculateRect(tlp a.Vector3) *common.RectBoundary {
	minX := tlp.X
	maxX := tlp.X + t.GetSize().X
	minY := tlp.Y
	maxY := tlp.Y + t.GetSize().Y
	minZ := tlp.Z
	maxZ := tlp.Z + t.GetSize().Z
	return common.NewRectBoundary(minX, maxX, minY, maxY, minZ, maxZ)
}

func NewTransformFromMap(siMap a.SiMap) Transform {
	var t Transform
	t.FromMap(siMap)
	return t
}