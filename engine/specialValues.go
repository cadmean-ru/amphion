package engine

import "github.com/cadmean-ru/amphion/common/a"

var specialYamlValuesMap = a.SiMap{
	"$MatchParent":           a.MatchParent,
	"$WrapContent":           a.WrapContent,
	"$CenterInParent":        a.CenterInParent,
	"$FillParent":            a.FillParent,
	"$TextAlignLeft":         a.TextAlignLeft,
	"$TextAlignRight":        a.TextAlignRight,
	"$TextAlignCenter":       a.TextAlignCenter,
	"$TextAlignTop":          a.TextAlignTop,
	"$TextAlignBottom":       a.TextAlignBottom,
	"$BuiltinShapePoint":     1,
	"$BuiltinShapeLine":      2,
	"$BuiltinShapeRectangle": 3,
	"$BuiltinShapeEllipse":   4,
	"$BuiltinShapeTriangle":  5,
	"$EventMouseDown":        EventMouseDown,
	"$EventMouseUp":          EventMouseUp,
	"$EventMouseIn":          EventMouseIn,
	"$EventMouseOut":         EventMouseOut,
	"$EventMouseMove":  EventMouseMove,
	"$EventTouchDown":  EventTouchDown,
	"$EventTouchUp":    EventTouchUp,
	"$EventTouchMove":  EventTouchMove,
	"$EventAppHide":    EventAppHide,
	"$EventAppShow":    EventAppShow,
	"$EventKeyDown":    EventKeyDown,
	"$EventKeyUp":      EventKeyUp,
	"$EventTextInput":  EventTextInput,
	"$EventFocusGain":  EventFocusGain,
	"$EventFocusLoose": EventFocusLose,
}

//IsSpecialValueString checks if the provided values is a string and it is a special values string, e.g. "$MatchParent".
func IsSpecialValueString(i interface{}) bool {
	if str, ok := i.(string); ok {
		return specialYamlValuesMap.ContainsKey(str)
	}
	return false
}

//GetSpecialValueFromString returns the value corresponding to the given special string.
func GetSpecialValueFromString(i interface{}) interface{} {
	if str, ok := i.(string); ok {
		return specialYamlValuesMap[str]
	}
	return 0
}

//IsSpecialTransformVector3 checks if the given vector contains special transform values.
func IsSpecialTransformVector3(pos a.Vector3) bool {
	return IsSpecialTransformValue(pos.X) || IsSpecialTransformValue(pos.Y) || IsSpecialTransformValue(pos.Z)
}

//IsSpecialTransformValue checks if the given float32 value is special(MatchParent, WrapContent or CenterInParent).
func IsSpecialTransformValue(x float32) bool {
	return x == a.CenterInParent || x == a.MatchParent || x == a.WrapContent
}