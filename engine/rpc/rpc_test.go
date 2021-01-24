package rpc

import (
	"fmt"
	"github.com/cadmean-ru/amphion/engine"
	"github.com/cadmean-ru/amphion/frontend/pc"
	"github.com/stretchr/testify/assert"
	"testing"
)


func TestF(t *testing.T) {
	var p = pc.NewFrontend()
	var amphion = engine.Initialize(p)

	amphion.Start()

	Initialize("http://testrpc.cadmean.ru")

	var expected = 69.0
	var actual float64
	var actualError error

	F("sum").Then(func(res interface{}) {
		fmt.Printf("%+v\n", res)
		actual = res.(float64)
		amphion.Stop()
	}).Err(func(err error) {
		fmt.Printf("%+v\n", err)
		actualError = err
		amphion.Stop()
	}).Call(2, 67)

	amphion.WaitForStop()

	assert.Nil(t, actualError)
	assert.Equal(t, expected, actual)
}
