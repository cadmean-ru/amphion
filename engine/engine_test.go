//+build windows linux darwin

package engine

import (
	"fmt"
	"github.com/cadmean-ru/amphion/frontend/pc"
	"github.com/stretchr/testify/assert"
	"testing"
)

// Tests for functions of struct AmphionEngine and public package APIs

func runEngineTest(testDelegate func(e *AmphionEngine)) {
	front := pc.NewFrontend()
	front.Init()
	e := Initialize(front)

	//go func() {
		e.Start()

		testDelegate(e)

		e.Stop()
		e.WaitForStop()
	//}()

	//front.Run()
}

func TestNameOfComponent(t *testing.T) {
	expected := "github.com/cadmean-ru/amphion/engine.testComponent"
	actual := NameOfComponent(&testComponent{})

	fmt.Println(expected)
	fmt.Println(actual)

	assert.Equal(t, expected, actual, "The names should match")
}