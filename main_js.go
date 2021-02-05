//+build js

package main

import (
	"github.com/cadmean-ru/amphion/engine"
	"github.com/cadmean-ru/amphion/frontend/web"
)

func main() {
	front := web.NewFrontend()
	front.Init()

	e := engine.Initialize(front)

	registerComponents(e)
	registerResources(e)

	go func() {
		e.Start()
		e.ShowScene(gridScene(e))
		//e.LoadApp()
		e.WaitForStop()
	}()

	front.Run()
}