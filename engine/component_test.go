package engine

import (
	"fmt"
	"github.com/cadmean-ru/amphion/common/a"
	"github.com/stretchr/testify/assert"
	"testing"
)

type testStatefulWithTags struct {
	ComponentImpl
	Bruh  string       "state"
	Bruh2 float64      `state:"breh"`
	Color a.Color      `state:"color"`
	Arr   []int        `state:"arr"`
	Hand  EventHandler `state:"hand"`
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
	assertions := assert.New(t)

	comp := &testStatefulWithTags{}
	comp.Bruh = "abc"
	comp.Bruh2 = 6.9
	comp.Color = a.Pink()
	comp.Arr = []int{1, 3, 4}
	comp.Hand = handleTest

	state := (&ComponentsManager{}).GetComponentState(comp)

	fmt.Println(state)

	assertions.Contains(state, "Bruh", `State should contain key "Bruh"`)
	assertions.Equal("abc", state["Bruh"], `The value of "Bruh" should match`)

	assertions.Contains(state, "breh", `State should contain key "breh"`)
	assertions.Equal(6.9, state["breh"], `The value of "breh" should match`)
}

func TestSetComponentState(t *testing.T) {
	assertions := assert.New(t)

	state := a.SiMap{
		"Bruh":  "abc",
		"breh":  6.9,
		"color": a.SiMap{"r": 255, "g": 192, "b": 203, "a": 255},
		"arr":   []interface{}{1, 3, 4},
		"hand":  "github.com/cadmean-ru/amphion/engine.handleTest",
	}

	comp := &testStatefulWithTags{}

	cm := newComponentsManager()
	cm.RegisterEventHandler(handleTest)

	cm.SetComponentState(comp, state)

	fmt.Printf("%+v\n", comp)

	comp.Hand(AmphionEvent{})

	assertions.Equal("abc", comp.Bruh)
	assertions.Equal(6.9, comp.Bruh2)
}

func handleTest(_ AmphionEvent) bool {
	fmt.Println("Handled")
	return false
}