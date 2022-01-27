package common

import (
	"fmt"
	"github.com/cadmean-ru/amphion/common/a"
)

type EdgeInsets struct {
	Left, Top, Right, Bottom float32
}

func (e *EdgeInsets) ToMap() a.SiMap {
	return a.SiMap{
		"left": e.Left,
		"top": e.Top,
		"right": e.Right,
		"bottom": e.Bottom,
	}
}

func (e *EdgeInsets) FromMap(state a.SiMap) {
	e.Left = state.GetFloat32("left")
	e.Top = state.GetFloat32("top")
	e.Right = state.GetFloat32("right")
	e.Bottom = state.GetFloat32("bottom")
}

func (e *EdgeInsets) String() string {
	return fmt.Sprintf("left: %f, top: %f, right: %f, bottom: %f", e.Left, e.Top, e.Right, e.Bottom)
}

func NewEdgeInsets(left, top, right, bottom float32) *EdgeInsets {
	return &EdgeInsets{
		Left:   left,
		Top:    top,
		Right:  right,
		Bottom: bottom,
	}
}
