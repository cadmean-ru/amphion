package engine

import (
	"github.com/cadmean-ru/amphion/frontend/pc"
)


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

