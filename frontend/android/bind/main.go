//+build android

package bind

import (
	"github.com/cadmean-ru/amphion/engine"
	"github.com/cadmean-ru/amphion/engine/builtin"
	"github.com/cadmean-ru/amphion/frontend"
	"github.com/cadmean-ru/amphion/frontend/cli"
	"github.com/cadmean-ru/amphion/frontend/ios"
)

var front frontend.Frontend

func AmphionInitAndroid(f cli.FrontendDelegate, rm cli.ResourceManagerDelegate, rd cli.RendererDelegate) {
	front = ios.NewFrontend(f, rm, rd)
	front.Init()

	e := engine.Initialize(front)

	registerComponents(e)
	registerResources(e)

	e.Start()
	e.LoadApp()

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
	cm.RegisterComponentType(&builtin.InputField{})
	cm.RegisterComponentType(&builtin.MouseMover{})
	cm.RegisterComponentType(&builtin.BuilderComponent{})
}