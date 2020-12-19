package main

import (
	"fmt"
	"github.com/cadmean-ru/amphion/common"
	"github.com/cadmean-ru/amphion/engine"
	"github.com/cadmean-ru/amphion/engine/builtin"
	"github.com/cadmean-ru/amphion/frontend/pc"
	"github.com/cadmean-ru/amphion/rendering"
	"log"
	"runtime"
	"time"
)

type Mover struct {
	obj  *engine.SceneObject
	log  *engine.Logger
	prc  *builtin.ShapeView
	eng  *engine.AmphionEngine
	dir  bool
}

func (m *Mover) OnInit(ctx engine.InitContext) {
	m.obj = ctx.GetSceneObject()
	m.log = ctx.GetLogger()
	m.prc = m.obj.GetComponentByName("ShapeView").(*builtin.ShapeView)
	m.eng = ctx.GetEngine()
}

func (m *Mover) OnStart() {
	m.log.Info(m, "Start")
}

func (m *Mover) OnUpdate(ctx engine.UpdateContext) {
	maxX := m.eng.GetCurrentScene().Transform.Size.X - m.obj.Transform.Size.X

	if m.obj.Transform.Position.X <= 0 {
		m.dir = true
	} else if m.obj.Transform.Position.X >= maxX {
		m.obj.Transform.Position.X = maxX
		m.dir = false
	}
	dX := 100 * ctx.DeltaTime
	if m.dir {
		m.obj.Transform.Position.X += dX
	} else {
		m.obj.Transform.Position.X -= dX
	}
	m.prc.ForceRedraw()
	m.eng.RequestRendering()
}

func (m *Mover) OnStop() {
	m.log.Info(m, "Stop")
}

func (m *Mover) GetName() string {
	return "Mover"
}

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

	front := pc.NewFrontend()
	front.Init()

	e := engine.Initialize(front)


	scene := engine.NewSceneObject("Test scene")

	rect := engine.NewSceneObject("rect")
	rect.Transform.Size = common.NewVector3(100, 100, 100)
	rect.Transform.Position = common.NewVector3(100, 100, -2)
	shape := builtin.NewShapeView(rendering.PrimitiveRectangle)
	shape.FillColor = common.PinkColor()
	rect.AddComponent(shape)
	rect.AddComponent(&Mover{})

	circle := engine.NewSceneObject("circle")
	circle.Transform.Size = common.NewVector3(50, 50, 0)
	circle.Transform.Position = common.NewVector3(10, 10 , 1)
	circleRenderer := builtin.NewShapeView(rendering.PrimitiveEllipse)
	circleRenderer.Appearance.StrokeWeight = 0
	circleRenderer.Appearance.FillColor = common.GreenColor()
	circle.AddComponent(circleRenderer)
	circle.AddComponent(builtin.NewCircleBoundary())
	circle.AddComponent(builtin.NewOnClickListener(func(event engine.AmphionEvent) bool {
		var mousePos = event.Data.(common.IntVector3)
		e.GetLogger().Info(nil, fmt.Sprintf("BIG CLICK ON CIRCLE. Mouse pos: %d %d", mousePos.X, mousePos.Y))
		//e.CloseScene(func() {
		//	_ = e.ShowScene(createCyberpunkScene(e))
		//})
		return false
	}))

	rect.AddChild(circle)


	scene.AddChild(rect)

	//text := engine.NewSceneObject("Close text")
	//text.Transform.Position = common.NewVector3(engine.CenterInParent, engine.CenterInParent, engine.CenterInParent)
	//text.Transform.Pivot = common.NewVector3(0.5, 0.5, 0.5)
	//text.Transform.Size = common.NewVector3(200, 50, 0)
	//textComponent := builtin.NewTextView("Close")
	//textComponent.TextAppearance.FontSize = 30
	//textComponent.Appearance.FillColor = common.BlackColor()
	//text.AddComponent(textComponent)
	//text.AddComponent(builtin.NewRectBoundary())
	////text.AddComponent(builtin.NewBoundaryView())
	//text.AddComponent(builtin.NewOnClickListener(func(event engine.AmphionEvent) bool {
	//	e.GetLogger().Info(nil, "close")
	//	e.Stop()
	//	return false
	//}))
	//
	//for i := 0; i < 10; i++ {
	//	text1 := engine.NewSceneObject(fmt.Sprintf("text%d", i))
	//	text1.Transform.Position = common.NewVector3(10, float64(i) * 50, 0)
	//	text1.Transform.Size = common.NewVector3(200, 50, 0)
	//	textComponent1 := builtin.NewTextView(fmt.Sprintf("Bruh %d", i))
	//	textComponent1.TextAppearance.FontSize = 30
	//	textComponent1.Appearance.FillColor = common.BlackColor()
	//	text1.AddComponent(textComponent1)
	//	text1.AddComponent(builtin.NewBoundaryView())
	//	text1.AddComponent(builtin.NewOnClickListener(func(event engine.AmphionEvent) bool {
	//		e.GetLogger().Info(nil, fmt.Sprintf("BIG CLICK ON %s", text1.GetName()))
	//		return false
	//	}))
	//	scene.AddChild(text1)
	//}
	//
	//scene.AddChild(text)
	//
	//point := engine.NewSceneObject("point")
	//point.AddComponent(builtin.NewShapeView(rendering.PrimitivePoint))
	//point.Transform.Position = common.NewVector3(500, 10, 0)
	//scene.AddChild(point)
	//
	line := engine.NewSceneObject("line")
	line.Transform.Position = common.NewVector3(400, 400, 0)
	line.Transform.Size = common.NewVector3(100, 10, 0)
	lineView := builtin.NewShapeView(rendering.PrimitiveLine)
	lineView.Appearance.StrokeColor = common.NewColor(0x2c, 0x68, 0xa8, 0xff)
	lineView.Appearance.StrokeWeight = 5
	line.AddComponent(lineView)
	scene.AddChild(line)
	//
	triangle := engine.NewSceneObject("triangle")
	triangle.Transform.Position = common.NewVector3(100, 100, 0)
	triangle.Transform.Size = common.NewVector3(100, 300, 0)
	triangleView := builtin.NewShapeView(rendering.PrimitiveTriangle)
	triangleView.Appearance.FillColor = common.BlueColor()
	triangle.AddComponent(triangleView)
	triangle.AddComponent(builtin.NewTriangleBoundary())
	scene.AddChild(triangle)

	image := engine.NewSceneObject("image")
	image.Transform.Position = common.NewVector3(200, 200, 0)
	image.Transform.Size = common.NewVector3(100, 100, 0)
	imageView := builtin.NewImageView(3)
	image.AddComponent(imageView)
	scene.AddChild(image)
	//
	//image2 := engine.NewSceneObject("image2")
	//image2.Transform.Position = common.NewVector3(500, 500, -1)
	//image2.Transform.Size = common.NewVector3(500, 200, 0)
	//imageView2 := builtin.NewImageView(0)
	//image2.AddComponent(imageView2)
	//scene.AddChild(image2)

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
	sceneBg.Appearance.StrokeWeight = 0
	sceneBg.Appearance.FillColor = common.NewColor(0xfc, 0xee, 0x0a, 0xff)
	scene.AddComponent(sceneBg)
	sceneImage := builtin.NewImageView(2)
	scene.AddComponent(sceneImage)

	centered := engine.NewSceneObject("centered")
	centered.Transform.Position = common.NewVector3(engine.CenterInParent, engine.CenterInParent, 1)
	centered.Transform.Size = common.NewVector3(800, 110, 0)
	centered.Transform.Pivot = common.NewVector3(0.5, 0.5, 0.5)
	//centered.AddComponent(builtin.NewBoundaryView())
	centered.AddComponent(builtin.NewRectBoundary())
	centered.AddComponent(builtin.NewOnClickListener(func(event engine.AmphionEvent) bool {
		fmt.Println(centered.Transform.GetGlobalPosition().ToString())
		return false
	}))
	scene.AddChild(centered)

	title := engine.NewSceneObject("title")
	title.Transform.Position = common.NewVector3(engine.CenterInParent, 0, 0)
	title.Transform.Pivot = common.NewVector3(0.5, 0, 0.5)
	title.Transform.Size = common.NewVector3(680, 55, 0)
	titleView := builtin.NewTextView("Time till Cyberpunk release")
	titleView.Appearance.FillColor = common.WhiteColor()
	titleView.TextAppearance.FontSize = 52
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
	textScene2.Transform.Position = common.NewVector3(engine.CenterInParent, engine.CenterInParent, 0)
	textScene2.Transform.Pivot = common.NewVector3(0.5, 0.5, 0.5)
	textScene2.Transform.Size = common.NewVector3(800, 200, 0)
	textScene2Renderer := builtin.NewTextView("This is scene 2")
	textScene2Renderer.TextAppearance.FontSize = 100
	textScene2Renderer.Appearance.FillColor = common.BlackColor()
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
	scene2.AddComponent(&TestController{})

	input := engine.NewSceneObject("input")
	input.Transform.Position = common.NewVector3(0, 0, 0)
	input.Transform.Size = common.NewVector3(500, 500 ,0)
	inputView := builtin.NewInputField()
	inputView.AllowMultiline = true
	input.AddComponent(inputView)
	scene2.AddChild(input)

	dropdown := engine.NewSceneObject("dropdown")
	dropdown.Transform.Position = common.NewVector3(10, 10, 2)
	dropdown.Transform.Size = common.NewVector3(100, 35, 0)
	dropdownView := builtin.NewDropdownView([]string {"opt1", "opt2", "opt3"})
	dropdown.AddComponent(dropdownView)
	dropdown.AddComponent(builtin.NewRectBoundary())
	dropdown.AddComponent(builtin.NewOnClickListener(func(event engine.AmphionEvent) bool {
		dropdownView.HandleClick()
		return true
	}))

	dropdown1 := engine.NewSceneObject("dropdown1")
	dropdown1.Transform.Position = common.NewVector3(10, 50, 1)
	dropdown1.Transform.Size = common.NewVector3(450, 35, 0)
	dropdownView1 := builtin.NewDropdownView([]string {"bruh1", "bruh2", "bruh3"})
	dropdown1.AddComponent(dropdownView1)
	dropdown1.AddComponent(builtin.NewRectBoundary())
	dropdown1.AddComponent(builtin.NewOnClickListener(func(event engine.AmphionEvent) bool {
		dropdownView1.HandleClick()
		return true
	}))


	box := engine.NewSceneObject("Moving box")
	box.Transform.Position = common.NewVector3(10, 100, 10)
	box.Transform.Size = common.NewVector3(500, 500, 0)
	boxBg := builtin.NewShapeView(rendering.PrimitiveRectangle)
	boxBg.Appearance.StrokeWeight = 0
	boxBg.Appearance.FillColor = common.NewColor(0xc4, 0xc4, 0xc4, 0xff)
	boxBg.Appearance.CornerRadius = 10
	box.AddComponent(boxBg)
	box.AddComponent(builtin.NewRectBoundary())
	box.AddComponent(builtin.NewMouseMover())
	box.AddChild(dropdown)
	box.AddChild(dropdown1)
	scene2.AddChild(box)

	curve := engine.NewSceneObject("Curve")
	curve.Transform.Position = common.NewVector3(10, 500, 5)
	curve.Transform.Size = common.NewVector3(100, 100, 0)
	curve.AddComponent(builtin.NewBezierView(common.NewVector3(50, 0, 0), common.NewVector3(50, 100, 0)))
	scene2.AddChild(curve)

	return scene2
}