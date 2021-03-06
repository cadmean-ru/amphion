//+build windows linux darwin

package engine

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNameOfComponent(t *testing.T) {
	expected := "github.com/cadmean-ru/amphion/engine.testComponent"
	actual := NameOfComponent(&testComponent{})

	fmt.Println(expected)
	fmt.Println(actual)

	assert.Equal(t, expected, actual, "The names should match")
}