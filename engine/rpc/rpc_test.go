package rpc

import (
	"fmt"
	"github.com/cadmean-ru/amphion/common"
	"github.com/cadmean-ru/amphion/engine"
	"testing"
)



func TestF(t *testing.T) {
	var p = common.PlatformFromString("web")
	var amphion = engine.Initialize(p)

	amphion.Start()

	Initialize("http://testrpc.cadmean.ru")

	F("weatherForecast.get").Then(func(res interface{}) {
		fmt.Printf("%+v\n", res)
		amphion.Stop()
	}).Err(func(err Error) {
		fmt.Printf("%+v\n", err)
		amphion.Stop()
	}).Args(1, 2).Call()

	amphion.WaitForStop()
}
