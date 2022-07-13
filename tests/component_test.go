package tests

import (
	"fmt"
	"github.com/cadmean-ru/amphion/common/a"
	"github.com/cadmean-ru/amphion/engine"
	"github.com/stretchr/testify/assert"
	"testing"
)

type testStatefulWithTags struct {
	engine.ComponentImpl
	Bruh  string              "state"
	Bruh2 float64             `state:"breh"`
	Color a.Color             `state:"color"`
	Arr   []int               `state:"arr"`
	Hand  engine.EventHandler `state:"hand"`
}

func (t testStatefulWithTags) GetName() string {
	return engine.NameOfComponent(t)
}

func TestIsStatefulComponent(t *testing.T) {
	actual := engine.IsStatefulComponent(&testComponent{})

	assert.True(t, actual, "The testComponent should be stateful")

	actual = engine.IsStatefulComponent(&testStatefulWithTags{})

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

	state := (&engine.ComponentsManager{}).GetComponentState(comp)

	fmt.Println(state)

	assertions.Contains(state, "Bruh", `State should contain key "Bruh"`)
	assertions.Equal("abc", state["Bruh"], `The value of "Bruh" should match`)

	assertions.Contains(state, "breh", `State should contain key "breh"`)
	assertions.Equal(6.9, state["breh"], `The value of "breh" should match`)
}

func handleTest(_ engine.Event) bool {
	fmt.Println("Handled")
	return false
}
