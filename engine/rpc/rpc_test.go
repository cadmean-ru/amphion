package rpc

import (
	"fmt"
	"github.com/cadmean-ru/amphion/atest"
	"github.com/cadmean-ru/amphion/engine"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFunc(t *testing.T) {
	atest.RunEngineTest(t, func(e *engine.AmphionEngine) {
		Initialize("http://testrpc.cadmean.ru")

		var expected = 69.0
		var actual float64
		var actualError error

		Func("sum").Then(func(res interface{}) {
			fmt.Printf("%+v\n", res)
			actual = res.(float64)
		}).Err(func(err error) {
			fmt.Printf("%+v\n", err)
			actualError = err
		}).Finally(func() {
			e.Stop()
		}).Call(2, 67)

		atest.WaitForStop()

		assert.Nil(t, actualError)
		assert.Equal(t, expected, actual)
	})
}
