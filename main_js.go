//go:build js
// +build js

package main

import (
	"github.com/cadmean-ru/amphion/engine"
	"github.com/cadmean-ru/amphion/frontend/web"
	"github.com/cadmean-ru/amphion/scenes"
)

func main() {
	front := web.NewFrontend()
	front.Init()

	e := engine.Initialize(front)

	scenes.RegisterComponents(e)
	scenes.RegisterResources(e)

	go func() {
		e.Start()
		_ = e.ShowScene(scenes.InputScene(e))
		//e.LoadApp()
		e.WaitForStop()
	}()

	front.Run()
}
