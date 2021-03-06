package engine

import (
	"fmt"
	"github.com/cadmean-ru/amphion/common/a"
	"github.com/stretchr/testify/assert"
	"testing"
)

type SetComponentStateTestComponent struct {
	*ComponentImpl
	StateInt          int           `state:"stateInt"`
	StateBool         bool          `state:"stateBool"`
	StateResId        a.ResId       `state:"stateResId"`
	StateTextAlign    a.TextAlign   `state:"stateTextAlign"`
	StateString       string        `state:"stateString"`
	StateUnmappable   a.Vector3     `state:"stateUnmappable"`
	StateUnstringable a.Color       `state:"stateUnstringable"`
	StatePtr          *a.Vector3    `state:"statePtr"`
	StateSlice        []a.Vector3   `state:"stateSlice"`
	StateEventHandler EventHandler  `state:"stateEventHandler"`
	StateSpecialFloat float32        `state:"stateSpecialFloat"`
	StateSpecialAlign a.TextAlign   `state:"stateSpecialAlign"`
	StateSpecialShape byte          `state:"stateSpecialShape"`
}

func (s *SetComponentStateTestComponent) GetName() string {
	return NameOfComponent(s)
}

func TestComponentsManager_SetComponentState(t *testing.T) {
	var comp = &SetComponentStateTestComponent{}
	var cm = newComponentsManager()
	cm.RegisterEventHandler(eh)

	var state = a.SiMap{
		"stateInt":       2,
		"stateBool":      true,
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
	ass.Equal(true, comp.StateBool)
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

func TestComponentsManager_GetComponentState(t *testing.T) {
	var comp = &SetComponentStateTestComponent{
		StateBool: true,
		StateUnstringable: a.NewColor("#696969"),
		StateUnmappable:   a.NewVector3(69, 69, 69),
	}
	var cm = newComponentsManager()

	state := cm.GetComponentState(comp)

	ass := assert.New(t)

	ass.Equal("#696969", state.GetString("stateUnstringable"))
	ass.Equal(float32(69), state["stateUnmappable"].(a.SiMap)["x"])
	ass.Equal(true, state.GetBool("stateBool"))
}

func TestGetFunctionName(t *testing.T) {
	fmt.Println(getFunctionName(eh))
	c := testComponent{}
	fmt.Println(getFunctionName(c.OnInit))
}

func TestComponentsManager_MakeComponent(t *testing.T) {
	cm := newComponentsManager()
	cm.RegisterComponentType(&SetComponentStateTestComponent{})

	ass := assert.New(t)

	name := "github.com/cadmean-ru/amphion/engine.SetComponentStateTestComponent"
	actual := cm.MakeComponent(name)
	ass.NotNil(actual)

	name = "SetComponentStateTestComponent"
	actual = cm.MakeComponent(name)
	ass.NotNil(actual)
}