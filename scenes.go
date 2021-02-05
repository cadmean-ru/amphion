package main

import (
	"fmt"
	"github.com/cadmean-ru/amphion/common/a"
	"github.com/cadmean-ru/amphion/engine"
	"github.com/cadmean-ru/amphion/engine/builtin"
	"github.com/cadmean-ru/amphion/rendering"
	"math/rand"
)

func scene1(e *engine.AmphionEngine) *engine.SceneObject {
	scene := engine.NewSceneObject("Test scene")

	rect := engine.NewSceneObject("rect")
	rect.Transform.Size = a.NewVector3(100, 100, 100)
	rect.Transform.Position = a.NewVector3(100, 100, -2)
	shape := builtin.NewShapeView(rendering.PrimitiveRectangle)
	shape.FillColor = a.PinkColor()
	rect.AddComponent(shape)
	rect.AddComponent(&Mover{})

	circle := engine.NewSceneObject("circle")
	circle.Transform.Size = a.NewVector3(50, 50, 0)
	circle.Transform.Position = a.NewVector3(10, 10 , 1)
	circleRenderer := builtin.NewShapeView(rendering.PrimitiveEllipse)
	circleRenderer.StrokeWeight = 0
	circleRenderer.FillColor = a.GreenColor()
	circle.AddComponent(circleRenderer)
	circle.AddComponent(builtin.NewCircleBoundary())
	circle.AddComponent(builtin.NewOnClickListener(handleCircleClick(e)))

	rect.AddChild(circle)


	scene.AddChild(rect)

	text := engine.NewSceneObject("Close text")
	text.Transform.Position = a.NewVector3(engine.CenterInParent, engine.CenterInParent, engine.CenterInParent)
	text.Transform.Pivot = a.NewVector3(0.5, 0.5, 0.5)
	text.Transform.Size = a.NewVector3(300, 50, 0)
	textComponent := builtin.NewTextView("Hello Amphion! 2")
	textComponent.FontSize = 30
	textComponent.TextColor = a.BlackColor()
	text.AddComponent(textComponent)
	text.AddComponent(builtin.NewRectBoundary())
	text.AddComponent(builtin.NewBoundaryView())
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
	scene.AddChild(text)
	//
	//point := engine.NewSceneObject("point")
	//point.AddComponent(builtin.NewShapeView(rendering.PrimitivePoint))
	//point.Transform.Position = common.NewVector3(500, 10, 0)
	//scene.AddChild(point)
	//
	line := engine.NewSceneObject("line")
	line.Transform.Position = a.NewVector3(400, 400, 0)
	line.Transform.Size = a.NewVector3(100, 10, 0)
	lineView := builtin.NewShapeView(rendering.PrimitiveLine)
	lineView.StrokeColor = a.NewColor(0x2c, 0x68, 0xa8, 0xff)
	lineView.StrokeWeight = 5
	line.AddComponent(lineView)
	scene.AddChild(line)
	//
	triangle := engine.NewSceneObject("triangle")
	triangle.Transform.Position = a.NewVector3(100, 100, 0)
	triangle.Transform.Size = a.NewVector3(100, 300, 0)
	triangleView := builtin.NewShapeView(rendering.PrimitiveTriangle)
	triangleView.FillColor = a.BlueColor()
	triangle.AddComponent(triangleView)
	triangle.AddComponent(builtin.NewTriangleBoundary())
	scene.AddChild(triangle)

	image := engine.NewSceneObject("image")
	image.Transform.Position = a.NewVector3(200, 200, 0)
	image.Transform.Size = a.NewVector3(500, 100, 0)
	imageView := builtin.NewImageView(Res_images_babyyoda)
	image.AddComponent(imageView)
	scene.AddChild(image)

	image2 := engine.NewSceneObject("image2")
	image2.Transform.Position = a.NewVector3(200, 300, -1)
	image2.Transform.Size = a.NewVector3(100, 300, 0)
	imageView2 := builtin.NewImageView(Res_images_cyberpunk)
	image2.AddComponent(imageView2)
	scene.AddChild(image2)

	return scene
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
	cm.RegisterComponentType(&Mover{})
	//cm.RegisterComponentType(&CyberpunkCountdown{})
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

	cm.RegisterEventHandler(handleCircleClick(e))
	cm.RegisterEventHandler(navigateOnClick)
	cm.RegisterEventHandler(navigateOnClick2)
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

//region TestController

type TestController struct {
	engine.ComponentImpl
}

func (c *TestController) OnInit(ctx engine.InitContext) {
	c.ComponentImpl.OnInit(ctx)
}

func (c *TestController) OnStart() {
	c.Engine.RunTask(engine.NewTaskBuilder().Run(func() (interface{}, error) {
		return c.Engine.GetResourceManager().ReadFile(5)
	}).Then(func(res interface{}) {
		bytes := res.([]byte)
		str := string(bytes)
		c.Logger.Info(c, str)
	}).Err(func(err error) {
		c.Logger.Error(c, err.Error())
	}).Build())
}

func (c *TestController) GetName() string {
	return engine.NameOfComponent(c)
}

//endregion

func handleCircleClick(e *engine.AmphionEngine) engine.EventHandler {
	return func(event engine.AmphionEvent) bool {
		var mousePos = event.Data.(a.IntVector3)
		e.GetLogger().Info(nil, fmt.Sprintf("BIG CLICK ON CIRCLE. Mouse pos: %d %d", mousePos.X, mousePos.Y))
		return false
	}
}

func makeRect(name string, x, y, sx, sy float32, fill a.Color) *engine.SceneObject {
	rect := engine.NewSceneObject(name)
	rect.Transform.Position = a.NewVector3(x, y, 0)
	rect.Transform.Size = a.NewVector3(sx, sy, 0)

	view := builtin.NewShapeView(builtin.ShapeRectangle)
	view.FillColor = fill

	rect.AddComponent(view)

	return rect
}

func gridScene(e *engine.AmphionEngine) *engine.SceneObject {
	scene := engine.NewSceneObject("grid scene")

	grid := builtin.NewGridLayout()
	grid.Cols = 3
	grid.RowPadding = 10
	grid.ColPadding = 50
	scene.AddComponent(grid)

	//var counter int

	addBtn := makeRect("add button", 0, 0, 100, 100, a.GreenColor())
	addBtnText := engine.NewSceneObject("add text")
	addBtnText.Transform.Position = a.NewVector3(10, 10, 1)
	addBtnText.Transform.Size = a.NewVector3(100, 30, 0)
	addBtnTextView := builtin.NewTextView("addkj")
	addBtnText.AddComponent(addBtnTextView)
	addBtn.AddChild(addBtnText)
	addBtn.AddComponent(builtin.NewRectBoundary())
	addBtn.AddComponent(builtin.NewOnClickListener(func(event engine.AmphionEvent) bool {
		color := a.NewColor(byte(rand.Intn(256)), byte(rand.Intn(256)), byte(rand.Intn(256)), 255)
		rect := makeRect(fmt.Sprintf("Rect"), 0, 0, 100, engine.MatchParent, color)
		rect.AddComponent(builtin.NewRectBoundary())
		scene.AddChild(rect)

		//counter++
		//
		//inputObj := engine.NewSceneObject(fmt.Sprintf("input %d", counter))
		//inputObj.Transform.Size.Y = float32(rand.Intn(300))
		//input := builtin.NewNativeInputView("", fmt.Sprintf("Enter some text %d", counter))
		//input.SetOnTextChangeListener(func(text string) {
		//	engine.LogDebug(fmt.Sprintf("Hello from %s! Text: %s\n", inputObj.GetName(), text))
		//})
		//inputObj.AddComponent(input)
		//scene.AddChild(inputObj)

		return false
	}))
	addBtn.AddComponent(builtin.NewEventListener(engine.EventMouseIn, func(event engine.AmphionEvent) bool {
		engine.LogInfo("Mouse in")
		shape := addBtn.GetComponentByName(".+ShapeView").(*builtin.ShapeView)
		shape.FillColor = a.PinkColor()
		shape.ForceRedraw()
		engine.RequestRendering()
		return false
	}))
	addBtn.AddComponent(builtin.NewEventListener(engine.EventMouseOut, func(event engine.AmphionEvent) bool {
		engine.LogInfo("Mouse out")
		shape := addBtn.GetComponentByName(".+ShapeView").(*builtin.ShapeView)
		shape.FillColor = a.GreenColor()
		shape.ForceRedraw()
		engine.RequestRendering()
		return false
	}))

	rmvButton := makeRect("remove button", 0, 0, 100, 100, a.RedColor())
	rmvButtonText := engine.NewSceneObject("remove text")
	rmvButtonText.Transform.Position = a.NewVector3(10, 10, 1)
	rmvButtonText.Transform.Size = a.NewVector3(engine.MatchParent, engine.MatchParent, 0)
	rmvButtonTextView := builtin.NewTextView("remove")
	rmvButtonText.AddComponent(rmvButtonTextView)
	rmvButtonText.AddComponent(builtin.NewBoundaryView())
	rmvButton.AddChild(rmvButtonText)
	rmvButton.AddComponent(builtin.NewRectBoundary())
	rmvButton.AddComponent(builtin.NewOnClickListener(func(event engine.AmphionEvent) bool {
		children := scene.GetChildren()
		last := children[len(children)-1]
		scene.RemoveChild(last)
		return false
	}))

	scene.AddChild(addBtn)
	scene.AddChild(rmvButton)

	e.BindEventHandler(engine.EventKeyDown, func(event engine.AmphionEvent) bool {
		engine.LogDebug(fmt.Sprintf("%+v\n", event.Data))
		return false
	})

	e.BindEventHandler(engine.EventMouseDown, func(event engine.AmphionEvent) bool {
		data := event.Data.(engine.MouseEventData)
		if data.SceneObject == nil {
			engine.LogDebug("Clicked on nothing")
		} else {
			engine.LogDebug("Clicked on %+v", event.Data.(engine.MouseEventData).SceneObject.GetName())
		}

		return true
	})

	return scene
}

func navigateOnClick(_ engine.AmphionEvent) bool {
	engine.LogDebug("Big click")

	err := engine.Navigate("second", nil)
	if err != nil {
		engine.LogError("Error navigating: %e", err)
	}
	return true
}

func navigateOnClick2(_ engine.AmphionEvent) bool {
	engine.LogDebug("Big click")

	err := engine.Navigate("main", nil)
	if err != nil {
		engine.LogError("Error navigating: %e", err)
	}
	return true
}

//func dropScene(e *engine.AmphionEngine) *engine.SceneObject {
//	scene := engine.NewSceneObject("drop scene")
//
//	gridLayout := builtin.NewGridLayout()
//	gridLayout.RowPadding = 10
//
//	scene.AddComponent(gridLayout)
//
//	title := engine.NewSceneObject("title")
//	title.Transform.Size.Y = 30
//	title.AddComponent(builtin.NewTextView("Drop a file here:"))
//	scene.AddChild(title)
//
//	fileName := engine.NewSceneObject("file name")
//	fileName.Transform.Size.Y = 30
//	fileNameText := builtin.NewTextView("")
//	fileName.AddComponent(fileNameText)
//
//	fileContents := engine.NewSceneObject("file contents")
//	fileContents.Transform.Size.Y = 1000
//	fileContentsText := builtin.NewTextView("")
//	fileContents.AddComponent(fileContentsText)
//
//	preview := engine.NewSceneObject("preview")
//	preview.Transform.Size.Y = 100
//	previewImage := builtin.NewImageView(Res_images_babyyoda)
//	preview.AddComponent(previewImage)
//
//	dropZone := engine.NewSceneObject("drop zone")
//	dropZone.Transform.Size = a.NewVector3(0, 100, 0)
//	dropZone.Transform.Position = a.NewVector3(0, 0, 10)
//	dropZone.AddComponent(builtin.NewBoundaryView())
//	dropZone.AddComponent(builtin.NewFileDropZone(func(event engine.AmphionEvent) bool {
//		data := event.Data.(engine.InputFileData)
//
//		if strings.Contains(data.Mime, "image") {
//			fileNameText.SetText("")
//			fileContentsText.SetText("")
//			url := fmt.Sprintf("data:%s;base64,%s", data.Mime, base64.StdEncoding.EncodeToString(data.Data))
//			previewImage.SetImageUrl(url)
//		} else {
//			fileNameText.SetText(data.Name)
//			fileContentsText.SetText(string(data.Data))
//		}
//		return false
//	}))
//
//	scene.AddChild(dropZone)
//	scene.AddChild(preview)
//	scene.AddChild(fileName)
//	scene.AddChild(fileContents)
//
//	return scene
//}