package engine

import (
	"github.com/cadmean-ru/amphion/common"
	"github.com/cadmean-ru/amphion/common/a"
	"github.com/cadmean-ru/amphion/common/require"
	"github.com/cadmean-ru/amphion/rendering"
)

// Transform describes how a SceneObject object is positioned on the screen.
type Transform struct {
	position       a.Vector3
	pivot          a.Vector3
	rotation       a.Vector3
	size           a.Vector3
	actualSize     a.Vector3
	measuredSize   a.Vector3
	actualPosition a.Vector3
	//inLayout     bool

	sceneObject *SceneObject
	parent      *Transform
}

func (t *Transform) WantedPosition() a.Vector3 {
	return t.position
}

func (t *Transform) SetPosition(position ...interface{}) {
	t.position = getVector3FromInterfaceValues(t.position, position...)
	t.actualPosition = t.vectorWithoutSpecialValues(t.position)
}

func (t *Transform) SetPositionCentered() {
	t.SetPosition(a.NewVector3(a.CenterInParent, a.CenterInParent, a.CenterInParent))
}

func (t *Transform) Translate(translation ...interface{}) {
	t.SetPosition(t.position.Add(getVector3FromInterfaceValues(a.ZeroVector(), translation...)))
}

func (t *Transform) Pivot() a.Vector3 {
	return t.pivot
}

func (t *Transform) SetPivot(pivot ...interface{}) {
	t.pivot = getVector3FromInterfaceValues(t.pivot, pivot...)
}

func (t *Transform) SetPivotCentered() {
	t.SetPivot(a.NewVector3(0.5, 0.5, 0.5))
}

func (t *Transform) Rotation() a.Vector3 {
	return t.rotation
}

func (t *Transform) SetRotation(rotation a.Vector3) {
	t.rotation = rotation
}

func (t *Transform) WantedSize() a.Vector3 {
	return t.size
}

func (t *Transform) SetSize(size ...interface{}) {
	t.size = getVector3FromInterfaceValues(t.size, size...)
	t.actualSize = t.vectorWithoutSpecialValues(t.size)
}

func (t *Transform) SetSizeWrapContent() {
	t.SetSize(a.NewVector3(a.WrapContent, a.WrapContent, a.WrapContent))
}

func (t *Transform) SetSizeMatchParent() {
	t.SetSize(a.NewVector3(a.MatchParent, a.MatchParent, a.MatchParent))
}

// ActualSize calculates the actual size of the Transform replacing the special values.
func (t *Transform) ActualSize() a.Vector3 {
	return t.actualSize
}

// LocalPosition calculates the actual local position related to this transform's parent.
func (t *Transform) LocalPosition() a.Vector3 {
	return t.actualPosition
}

// GlobalPosition calculates the actual global position in the scene.
func (t *Transform) GlobalPosition() a.Vector3 {
	if t.parent == nil {
		return t.actualPosition
	}

	return t.parent.GlobalPosition().Sub(t.parent.actualSize.Multiply(t.parent.pivot)).Add(t.actualPosition)
}

func (t *Transform) CenterInParent() {
	t.SetPivotCentered()
	t.SetPositionCentered()
}

// TopLeftPosition calculates local position of the top left point of the bounding box.
func (t *Transform) TopLeftPosition() a.Vector3 {
	return t.actualPosition.Sub(t.actualSize.Multiply(t.pivot))
}

// GlobalTopLeftPosition calculates global position of the top left point of the bounding box.
func (t *Transform) GlobalTopLeftPosition() a.Vector3 {
	return t.GlobalPosition().Sub(t.actualSize.Multiply(t.pivot))
}

// Rect calculates local rect of this transform.
func (t *Transform) Rect() *common.RectBoundary {
	return t.calculateRect(a.ZeroVector(), t.actualSize)
}

// GlobalRect calculates global rect of this transform.
func (t *Transform) GlobalRect() *common.RectBoundary {
	return t.calculateRect(t.GlobalTopLeftPosition(), t.actualSize)
}

// GetParent returns the parent transform of the current transform.
func (t Transform) GetParent() *Transform {
	return t.parent
}

// ToRenderingTransform calculates a transform that is ready to be rendered on screen with all absolute values in pixels calculated.
func (t *Transform) ToRenderingTransform() rendering.Transform {
	rt := rendering.NewTransform()

	rt.Position = t.GlobalTopLeftPosition().Round()
	rt.Size = t.actualSize.Round()

	return rt
}

func (t *Transform) Equals(other Transform) bool {
	return t.position.Equals(other.position) && t.size.Equals(other.size) && t.pivot.Equals(other.pivot) && t.rotation.Equals(other.rotation)
}

// Calculates the rect boundary of the transform given the top left point's position.
func (t *Transform) calculateRect(tlp, size a.Vector3) *common.RectBoundary {
	return common.NewRectBoundaryFromPositionAndSize(tlp, size)
}

func (t *Transform) positionForMeasurement() a.Vector3 {
	return t.vectorWithoutSpecialValues(t.position)
}

func (t *Transform) vectorWithoutSpecialValues(v a.Vector3) a.Vector3 {
	if IsSpecialTransformValue(v.X) {
		v.X = 0
	}
	if IsSpecialTransformValue(v.X) {
		v.X = 0
	}
	if IsSpecialTransformValue(v.X) {
		v.X = 0
	}
	return v
}

func (t *Transform) ToMap() a.SiMap {
	return map[string]interface{}{
		"position": t.position.ToMap(),
		"pivot":    t.pivot.ToMap(),
		"rotation": t.rotation.ToMap(),
		"size":     t.size.ToMap(),
	}
}

func (t *Transform) DumpToMap() a.SiMap {
	return map[string]interface{}{
		"position":       t.position.ToMap(),
		"pivot":          t.pivot.ToMap(),
		"rotation":       t.rotation.ToMap(),
		"size":           t.size.ToMap(),
		"measuredSize":   t.measuredSize.ToMap(),
		"actualSize":     t.actualSize.ToMap(),
		"actualPosition": t.actualPosition.ToMap(),
	}
}

func (t *Transform) FromMap(siMap a.SiMap) {
	t.position = a.NewVector3FromMap(t.decodeSpecialValuesInVector(a.RequireSiMap(siMap["position"])))
	t.pivot = a.NewVector3FromMap(a.RequireSiMap(siMap["pivot"]))
	t.rotation = a.NewVector3FromMap(a.RequireSiMap(siMap["rotation"]))
	t.size = a.NewVector3FromMap(t.decodeSpecialValuesInVector(a.RequireSiMap(siMap["size"])))
}

func NewTransformFromMap(siMap a.SiMap) Transform {
	var t Transform
	t.FromMap(siMap)
	return t
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

// NewTransform2D creates a new transform with default values.
func NewTransform2D(object *SceneObject) Transform {
	return Transform{
		position:    a.ZeroVector(),
		pivot:       a.ZeroVector(),
		rotation:    a.ZeroVector(),
		size:        a.NewVector3(a.WrapContent, a.WrapContent, a.WrapContent),
		sceneObject: object,
	}
}

func getVector3FromInterfaceValues(defaultVector a.Vector3, values ...interface{}) a.Vector3 {
	switch len(values) {
	case 1:
		switch values[0].(type) {
		case a.Vector3:
			return values[0].(a.Vector3)
		}
	case 2:
		return a.NewVector3(
			require.Float32(values[0]),
			require.Float32(values[1]),
			defaultVector.Z,
		)
	case 3:
		return a.NewVector3(
			require.Float32(values[0]),
			require.Float32(values[1]),
			require.Float32(values[2]),
		)
	}

	return defaultVector
}
