//+build ios

package iosCli

import (
	"fmt"
	"github.com/cadmean-ru/amphion/engine"
	"github.com/cadmean-ru/amphion/engine/builtin"
	"github.com/cadmean-ru/amphion/frontend"
	"github.com/cadmean-ru/amphion/frontend/cli"
	"github.com/cadmean-ru/amphion/frontend/ios"
)

var front frontend.Frontend

func AmphionInit(f cli.FrontendDelegate, rm cli.ResourceManagerDelegate, rd cli.RendererDelegate) {
	front = ios.NewFrontend(f, rm, rd)
	front.Init()

	e := engine.Initialize(front)

	go func() {
		fmt.Println("Here")
		e.Start()
		fmt.Println("Here started")

		registerComponents(e)
		registerResources(e)

		e.LoadApp()

		e.WaitForStop()
	}()

	go front.Run()
}

func RegisterPrimitiveRendererDelegate(primitiveKind int, delegate cli.PrimitiveRendererDelegate) {
	front.GetRenderer().RegisterPrimitiveRendererDelegate(byte(primitiveKind), cli.NewPrimitiveRendererDelegateWrap(delegate))
}

func registerResources(e *engine.AmphionEngine) {
	rm := e.GetResourceManager()
	rm.RegisterResource("2006.ttf")
	rm.RegisterResource("images/baby-yoda.jpg")
	rm.RegisterResource("images/babyyoda2.png")
	rm.RegisterResource("images/cyberpunk.jpg")
	rm.RegisterResource("images/gun.png")
	rm.RegisterResource("test.yaml")
	rm.RegisterResource("scenes/main.scene")
	rm.RegisterResource("scenes/second.scene")
}

func registerComponents(e *engine.AmphionEngine) {
	cm := e.GetComponentsManager()
	cm.RegisterComponentType(&builtin.ShapeView{})
	cm.RegisterComponentType(&builtin.CircleBoundary{})
	cm.RegisterComponentType(&builtin.OnClickListener{})
	cm.RegisterComponentType(&builtin.TextView{})
	cm.RegisterComponentType(&builtin.RectBoundary{})
	cm.RegisterComponentType(&builtin.TriangleBoundary{})
	cm.RegisterComponentType(&builtin.BezierView{})
	cm.RegisterComponentType(&builtin.DropdownView{})
	cm.RegisterComponentType(&builtin.ImageView{})
	cm.RegisterComponentType(&builtin.MouseMover{})
	cm.RegisterComponentType(&builtin.BuilderComponent{})
	cm.RegisterComponentType(&IosTestController{})
}

type IosTestController struct {
	engine.ComponentImpl
}

func (i *IosTestController) OnStart() {
	engine.LogDebug("Breh")

	//engine.BindEventHandler(engine.EventTouchDown, func(event engine.AmphionEvent) bool {
	//	data := event.Data.(engine.MouseEventData)
	//	engine.LogDebug("iOS T O U C H   down  %d %d", data.MousePosition.X, data.MousePosition.Y)
	//	return true
	//})

	//engine.BindEventHandler(engine.EventTouchMove, func(event engine.AmphionEvent) bool {
	//	data := event.Data.(engine.MouseEventData)
	//	//engine.LogDebug("iOS T O U C H   move  %d %d", data.MousePosition.X, data.MousePosition.Y)
	//	return true
	//})

	//engine.BindEventHandler(engine.EventTouchDown, func(event engine.AmphionEvent) bool {
	//	//data := event.Data.(engine.MouseEventData)
	//	//engine.LogDebug("iOS T O U C H   up %d %d", data.MousePosition.X, data.MousePosition.Y)
	//	return true
	//})

	i.SceneObject.AddComponent(builtin.NewOnClickListener(func(event engine.AmphionEvent) bool {
		engine.LogDebug("Ios CLICK")
		return true
	}))
}

func (i *IosTestController) GetName() string {
	return engine.NameOfComponent(i)
}