package engine

import (
	"github.com/cadmean-ru/amphion/common/a"
	"github.com/stretchr/testify/assert"
	"testing"
)

type SetComponentStateTestComponent struct {
	*ComponentImpl
	StateInt          int           `state:"stateInt"`
	StateResId        a.ResId       `state:"stateResId"`
	StateTextAlign    a.TextAlign   `state:"stateTextAlign"`
	StateString       string        `state:"stateString"`
	StateUnmappable   a.Vector3     `state:"stateUnmappable"`
	StateUnstringable a.Color       `state:"stateUnstringable"`
	StatePtr          *a.Vector3    `state:"statePtr"`
	StateSlice        []a.Vector3   `state:"stateSlice"`
	StateEventHandler EventHandler  `state:"stateEventHandler"`
	StateSpecialFloat float32       `state:"stateSpecialFloat"`
	StateSpecialAlign a.TextAlign   `state:"stateSpecialAlign"`
	StateSpecialShape byte          `state:"stateSpecialShape"`
}

func TestComponentsManager_SetComponentState(t *testing.T) {
	var comp = &SetComponentStateTestComponent{}
	var cm = newComponentsManager()
	cm.RegisterEventHandler(eh)

	var state = a.SiMap{
		"stateInt":       2,
		"stateResId":     12,
		"stateTextAlign": 3,
		"stateString":    "test",
		"stateUnmappable": a.SiMap{
			"x": 1,
			"y": 2,
			"z": 3,
		},
		"stateUnstringable": "#ff00ff",
		"statePtr": a.SiMap{
			"x": 10,
			"y": 11,
			"z": 12,
		},
		"stateSlice": []a.SiMap {
			{
				"x": 4,
				"y": 5,
				"z": 6,
			},
			{
				"x": 7,
				"y": 8,
				"z": 9,
			},
		},
		"stateEventHandler": "github.com/cadmean-ru/amphion/engine.eh",
		"stateSpecialFloat": "$CenterInParent",
		"stateSpecialAlign": "$TextAlignTop",
		"stateSpecialShape": "$BuiltinShapeRectangle",
	}
	cm.SetComponentState(comp, state)

	ass := assert.New(t)

	ass.Equal(2, comp.StateInt)
	ass.Equal(a.ResId(12), comp.StateResId)
	ass.Equal(a.TextAlign(3), comp.StateTextAlign)
	ass.Equal("test", comp.StateString)
	ass.Equal(a.NewVector3(1, 2, 3), comp.StateUnmappable)
	ass.Equal(a.NewColor("#ff00ff"), comp.StateUnstringable)
	ass.Equal(a.NewVector3(10, 11, 12), *comp.StatePtr)
	ass.Equal([]a.Vector3 {
		a.NewVector3(4, 5, 6),
		a.NewVector3(7, 8, 9),
	}, comp.StateSlice)
	ass.Equal(getFunctionName(eh), getFunctionName(comp.StateEventHandler))
	ass.Equal(float32(a.CenterInParent), comp.StateSpecialFloat)
	ass.Equal(a.TextAlignTop, comp.StateTextAlign)
	ass.Equal(byte(3), comp.StateSpecialShape)
}
