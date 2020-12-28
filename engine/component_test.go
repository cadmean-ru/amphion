package engine

import (
	"fmt"
	"github.com/cadmean-ru/amphion/common"
	"github.com/stretchr/testify/assert"
	"testing"
)

type testStatefulWithTags struct {
	ComponentImpl
	Bruh  string       "state"
	Bruh2 float64      `state:"breh"`
	Color common.Color `state:"color"`
}

func (t testStatefulWithTags) GetName() string {
	return NameOfComponent(t)
}

func TestIsStatefulComponent(t *testing.T) {
	actual := IsStatefulComponent(&testComponent{})

	assert.True(t, actual, "The testComponent should be stateful")

	actual = IsStatefulComponent(&testStatefulWithTags{})

	assert.True(t, actual, "The testStatefulWithTags should also be stateful")
}

func TestGetComponentState(t *testing.T) {
	a := assert.New(t)

	comp := &testStatefulWithTags{}
	comp.Bruh = "abc"
	comp.Bruh2 = 6.9
	comp.Color = common.PinkColor()

	state := GetComponentState(comp)

	fmt.Println(state)

	a.Contains(state, "Bruh", `State should contain key "Bruh"`)
	a.Equal("abc", state["Bruh"], `The value of "Bruh" should match`)

	a.Contains(state, "breh", `State should contain key "breh"`)
	a.Equal(6.9, state["breh"], `The value of "breh" should match`)
}

func TestSetComponentState(t *testing.T) {
	a := assert.New(t)

	state := common.SiMap {
		"Bruh": "abc",
		"breh": 6.9,
		"color": common.SiMap{"r": 255, "g": 192, "b": 203, "a": 255},
	}

	comp := &testStatefulWithTags{}

	SetComponentState(comp, state)

	fmt.Printf("%+v\n", comp)

	a.Equal("abc", comp.Bruh)
	a.Equal(6.9, comp.Bruh2)
}