//+build windows darwin linux

package main

import (
	"flag"
	"fmt"
	"github.com/cadmean-ru/amphion/common/a"
	"github.com/cadmean-ru/amphion/engine"
	"github.com/cadmean-ru/amphion/engine/builtin"
	"github.com/cadmean-ru/amphion/frontend/pc"
	"github.com/cadmean-ru/amphion/rendering"
	"github.com/cadmean-ru/amphion/utils"
	"log"
	"runtime"
	"time"
)



type TestController struct {
	eng *engine.AmphionEngine
	log *engine.Logger
}

func (c *TestController) OnInit(ctx engine.InitContext) {
	c.eng = ctx.GetEngine()
	c.log = ctx.GetLogger()
}

func (c *TestController) OnStart() {
	c.eng.GetResourceManager().GetFile("/test.yaml").ReadDataAsync(func(data []byte, err error) {
		if err != nil {
			c.log.Error(c, err.Error())
		}
		c.log.Info(c, string(data))
	})
}

func (c *TestController) OnStop() {

}

func (c *TestController) GetName() string {
	return "TestController"
}

func init() {
	runtime.LockOSThread()
}

func main() {

	var generator string
	var inputPath string
	var outputPath string
	var packageName string

	flag.StringVar(&generator, "generate", "", "Launches the specified code generator instead of starting engine. " +
		"Available generators are: shaders.")
	flag.StringVar(&inputPath, "i", "", "Define input path")
	flag.StringVar(&outputPath, "o", "", "Define output path")
	flag.StringVar(&packageName, "package", "", "Define package name")
	flag.Parse()

	switch generator {
	case "shaders":
		utils.GenerateShaders(inputPath, outputPath, packageName)
		return
	default:
		break
	}

	front := pc.NewFrontend()
	front.Init()

	e := engine.Initialize(front)

	cm := e.GetComponentsManager()
	cm.RegisterComponentType(&Mover{})
	cm.RegisterComponentType(&CyberpunkCountdown{})
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
	cm.RegisterComponentType(&builtin.InputView{})
	cm.RegisterComponentType(&builtin.MouseMover{})
	cm.RegisterComponentType(&builtin.BuilderComponent{})



	//if data, err := scene.EncodeToYaml(); err == nil {
	//	fmt.Println(string(data))
	//}

	//triangle.AddComponent(builtin.NewOnClickListener(func(event engine.AmphionEvent) bool {
	//	e.CloseScene(func() {
	//		_ = e.ShowScene(scene2(e))
	//	})
	//	return false
	//}))

	go func() {
		e.Start()

		if err := e.ShowScene(scene); err != nil {
			log.Println(err)
		}

		e.WaitForStop()
	}()

	front.Run()
}

type CyberpunkCountdown struct {
	releaseTime time.Time
	eng         *engine.AmphionEngine
	lastTick    time.Time
	textView    *builtin.TextView
}

func (c *CyberpunkCountdown) OnInit(ctx engine.InitContext) {
	c.eng = ctx.GetEngine()
	c.textView = ctx.GetSceneObject().GetComponentByName("TextView").(*builtin.TextView)
	c.releaseTime = time.Date(2020, 12, 10, 3, 0, 0, 0, time.Local)
}

func (c *CyberpunkCountdown) OnStart() {

}

func (c *CyberpunkCountdown) OnStop() {

}

func (c *CyberpunkCountdown) OnUpdate(_ engine.UpdateContext) {
	now := time.Now()
	if now.Sub(c.lastTick) >= time.Second {
		c.lastTick = now

		timeTillRelease := c.releaseTime.Sub(now)
		s := int(timeTillRelease.Seconds()) % 60
		m := int(timeTillRelease.Minutes()) % 60
		h := int(timeTillRelease.Hours()) % 24
		d := int(timeTillRelease.Hours() / 24)
		c.textView.SetText(fmt.Sprintf("%dd %dh %dm %ds", d, h, m, s))
		c.eng.RequestRendering()
	}
	c.eng.RequestUpdate()
}

func (c *CyberpunkCountdown) GetName() string {
	return "CyberpunkCountdown"
}

func createCyberpunkScene(e *engine.AmphionEngine) *engine.SceneObject {
	scene := engine.NewSceneObject("cyberpunk")
	sceneBg := builtin.NewShapeView(rendering.PrimitiveRectangle)
	sceneBg.StrokeWeight = 0
	sceneBg.FillColor = a.NewColor(0xfc, 0xee, 0x0a, 0xff)
	scene.AddComponent(sceneBg)
	sceneImage := builtin.NewImageView(2)
	scene.AddComponent(sceneImage)

	centered := engine.NewSceneObject("centered")
	centered.Transform.Position = a.NewVector3(engine.CenterInParent, engine.CenterInParent, 1)
	centered.Transform.Size = a.NewVector3(800, 110, 0)
	centered.Transform.Pivot = a.NewVector3(0.5, 0.5, 0.5)
	//centered.AddComponent(builtin.NewBoundaryView())
	centered.AddComponent(builtin.NewRectBoundary())
	centered.AddComponent(builtin.NewOnClickListener(func(event engine.AmphionEvent) bool {
		fmt.Println(centered.Transform.GetGlobalPosition().ToString())
		return false
	}))
	scene.AddChild(centered)

	title := engine.NewSceneObject("title")
	title.Transform.Position = a.NewVector3(engine.CenterInParent, 0, 0)
	title.Transform.Pivot = a.NewVector3(0.5, 0, 0.5)
	title.Transform.Size = a.NewVector3(680, 55, 0)
	titleView := builtin.NewTextView("Time till Cyberpunk release")
	titleView.TextColor = a.WhiteColor()
	titleView.FontSize = 52
	title.AddComponent(titleView)
	title.AddComponent(builtin.NewRectBoundary())
	title.AddComponent(builtin.NewOnClickListener(func(event engine.AmphionEvent) bool {
		fmt.Println(title.Transform.GetGlobalPosition().ToString())
		return false
	}))
	centered.AddChild(title)



	return scene
}

func scene2(e *engine.AmphionEngine) *engine.SceneObject {
	//var counter = 0
	scene2 := engine.NewSceneObject("scene 2")
	textScene2 := engine.NewSceneObject("text")
	textScene2.Transform.Position = a.NewVector3(engine.CenterInParent, engine.CenterInParent, 0)
	textScene2.Transform.Pivot = a.NewVector3(0.5, 0.5, 0.5)
	textScene2.Transform.Size = a.NewVector3(800, 200, 0)
	textScene2Renderer := builtin.NewTextView("This is scene 2")
	textScene2Renderer.FontSize = 100
	textScene2Renderer.TextColor = a.BlackColor()
	textScene2.AddComponent(textScene2Renderer)
	textScene2.AddComponent(builtin.NewRectBoundary())
	//textScene2.AddComponent(builtin.NewOnClickListener(func(event engine.AmphionEvent) bool {
	//	fmt.Println("Click")
	//	textScene2Renderer.SetText(strconv.Itoa(counter))
	//	textScene2Renderer.ForceRedraw()
	//	e.RequestRendering()
	//	e.RequestUpdate()
	//	counter++
	//	return false
	//}))
	scene2.AddChild(textScene2)
	//scene2.AddComponent(&TestController{})

	input := engine.NewSceneObject("input")
	input.Transform.Position = a.NewVector3(0, 0, 0)
	input.Transform.Size = a.NewVector3(500, 500 ,0)
	inputView := builtin.NewInputField()
	inputView.AllowMultiline = true
	input.AddComponent(inputView)
	scene2.AddChild(input)

	dropdown := engine.NewSceneObject("dropdown")
	dropdown.Transform.Position = a.NewVector3(10, 10, 2)
	dropdown.Transform.Size = a.NewVector3(100, 35, 0)
	dropdownView := builtin.NewDropdownView([]string {"opt1", "opt2", "opt3"})
	dropdown.AddComponent(dropdownView)
	dropdown.AddComponent(builtin.NewRectBoundary())
	dropdown.AddComponent(builtin.NewOnClickListener(func(event engine.AmphionEvent) bool {
		dropdownView.HandleClick()
		return true
	}))

	dropdown1 := engine.NewSceneObject("dropdown1")
	dropdown1.Transform.Position = a.NewVector3(10, 50, 1)
	dropdown1.Transform.Size = a.NewVector3(450, 35, 0)
	dropdownView1 := builtin.NewDropdownView([]string {"bruh1", "bruh2", "bruh3"})
	dropdown1.AddComponent(dropdownView1)
	dropdown1.AddComponent(builtin.NewRectBoundary())
	dropdown1.AddComponent(builtin.NewOnClickListener(func(event engine.AmphionEvent) bool {
		dropdownView1.HandleClick()
		return true
	}))


	box := engine.NewSceneObject("Moving box")
	box.Transform.Position = a.NewVector3(10, 100, 10)
	box.Transform.Size = a.NewVector3(500, 500, 0)
	boxBg := builtin.NewShapeView(rendering.PrimitiveRectangle)
	boxBg.StrokeWeight = 0
	boxBg.FillColor = a.NewColor(0xc4, 0xc4, 0xc4, 0xff)
	boxBg.CornerRadius = 10
	box.AddComponent(boxBg)
	box.AddComponent(builtin.NewRectBoundary())
	box.AddComponent(builtin.NewMouseMover())
	box.AddChild(dropdown)
	box.AddChild(dropdown1)
	scene2.AddChild(box)

	curve := engine.NewSceneObject("Curve")
	curve.Transform.Position = a.NewVector3(10, 500, 5)
	curve.Transform.Size = a.NewVector3(100, 100, 0)
	curve.AddComponent(builtin.NewBezierView(a.NewVector3(50, 0, 0), a.NewVector3(50, 100, 0)))
	scene2.AddChild(curve)

	return scene2
}