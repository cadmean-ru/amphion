//+build ios

package iosCli

import (
	"fmt"
	"github.com/cadmean-ru/amphion/engine"
	"github.com/cadmean-ru/amphion/frontend"
	"github.com/cadmean-ru/amphion/frontend/cli"
	"github.com/cadmean-ru/amphion/frontend/ios"
	"github.com/cadmean-ru/amphion/scenes"
)

var front frontend.Frontend

func AmphionInit(f cli.FrontendDelegate, rm cli.ResourceManagerDelegate, rd cli.RendererDelegate) {
	fmt.Println("WHY1")
	front = ios.NewFrontend(f, rm, rd)
	front.Init()

	fmt.Println("WHY2")
	e := engine.Initialize(front)
	fmt.Println("WHY3")

	go func() {
		fmt.Println("Here")
		e.Start()
		fmt.Println("Here started")

		scenes.RegisterComponents(e)
		scenes.RegisterResources(e)

		fmt.Println("WHY5")

		//e.LoadApp()
		e.ShowScene(scenes.GifScene(e))

		e.WaitForStop()
	}()

	fmt.Println("WHY4")
	go front.Run()
}

func RegisterPrimitiveRendererDelegate(primitiveKind int, delegate cli.PrimitiveRendererDelegate) {
	front.GetRenderer().RegisterPrimitiveRendererDelegate(byte(primitiveKind), cli.NewPrimitiveRendererDelegateWrap(delegate))
}