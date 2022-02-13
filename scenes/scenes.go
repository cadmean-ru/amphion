package scenes

import (
	"fmt"
	"github.com/cadmean-ru/amphion/common"
	"github.com/cadmean-ru/amphion/common/a"
	"github.com/cadmean-ru/amphion/common/shape"
	"github.com/cadmean-ru/amphion/engine"
	"github.com/cadmean-ru/amphion/engine/builtin"
	"math"
	"math/rand"
)

func Scene1(e *engine.AmphionEngine) *engine.SceneObject {
	scene := engine.NewSceneObject("Test scene")
	scene.AddComponent(builtin.NewAbsoluteLayout())

	rect := engine.NewSceneObject("rect")
	rect.Transform.SetSize(100, 100)
	shape := builtin.NewShapeView(builtin.ShapeRectangle)
	shape.FillColor = a.Pink()
	rect.AddComponent(shape)
	//rect.AddComponent(&Mover{})

	circle := engine.NewSceneObject("circle")
	circle.Transform.SetSize(50, 50)
	circle.Transform.SetPosition(10, 10, 1)
	circleRenderer := builtin.NewShapeView(builtin.ShapeEllipse)
	circleRenderer.StrokeWeight = 0
	circleRenderer.FillColor = a.Green()
	circle.AddComponent(circleRenderer)
	circle.AddComponent(builtin.NewCircleBoundary())
	circle.AddComponent(builtin.NewOnClickListener(HandleCircleClick(e)))

	rect.AddChild(circle)

	scene.AddChild(rect)

	text := engine.NewSceneObject("Close text")
	text.Transform.SetPosition(a.CenterInParent, a.CenterInParent)
	text.Transform.SetPivotCentered()
	text.Transform.SetSize(100, 50)
	textComponent := builtin.NewTextView("Hello Amphion! 2")
	textComponent.FontSize = 30
	textComponent.TextColor = a.Black()
	text.AddComponent(textComponent)
	text.AddComponent(builtin.NewRectBoundary())
	//text.AddComponent(builtin.NewBoundaryView())
	//text.AddComponent(builtin.NewOnClickListener(func(event engine.Event) bool {
	//	e.GetLogger().Info(nil, "close")
	//	e.Stop()
	//	return false
	//}))
	//
	//for i := 0; i < 10; i++ {
	//	text1 := engine.NewSceneObject(fmt.Sprintf("text%d", i))
	//	text1.Transform.position = common.NewVector3(10, float64(i) * 50, 0)
	//	text1.Transform.size = common.NewVector3(200, 50, 0)
	//	textComponent1 := builtin.NewTextView(fmt.Sprintf("Bruh %d", i))
	//	textComponent1.TextAppearance.FontSize = 30
	//	textComponent1.Appearance.FillColor = common.Black()
	//	text1.AddComponent(textComponent1)
	//	text1.AddComponent(builtin.NewBoundaryView())
	//	text1.AddComponent(builtin.NewOnClickListener(func(event engine.Event) bool {
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
	//point.Transform.position = common.NewVector3(500, 10, 0)
	//scene.AddChild(point)
	//
	line := engine.NewSceneObject("line")
	line.Transform.SetPosition(400, 400)
	line.Transform.SetSize(100, 10)
	lineView := builtin.NewShapeView(builtin.ShapeLine)
	lineView.StrokeColor = a.NewColor(0x2c, 0x68, 0xa8, 0xff)
	lineView.StrokeWeight = 5
	line.AddComponent(lineView)
	scene.AddChild(line)
	//
	triangle := engine.NewSceneObject("triangle")
	triangle.Transform.SetPosition(100, 100)
	triangle.Transform.SetSize(100, 300)
	triangleView := builtin.NewShapeView(builtin.ShapeTriangle)
	triangleView.FillColor = a.Blue()
	triangle.AddComponent(triangleView)
	triangle.AddComponent(builtin.NewTriangleBoundary())
	scene.AddChild(triangle)

	image := engine.NewSceneObject("image")
	image.Transform.SetPosition(200, 200)
	image.Transform.SetSize(500, 100)
	imageView := builtin.NewImageView(Res_images_babyyoda)
	image.AddComponent(imageView)
	scene.AddChild(image)

	image2 := engine.NewSceneObject("image2")
	image2.Transform.SetPosition(200, 300, -1)
	//image2.Transform.SetSize(100, 300)
	imageView2 := builtin.NewImageView("https://c.tenor.com/UZJd1pjj4NMAAAAM/surprised-pikachu.gif")
	//imageView2 := builtin.NewImageView(Res_images_giphy)
	image2.AddComponent(imageView2)
	scene.AddChild(image2)

	return scene
}

func RegisterResources(e *engine.AmphionEngine) {
	rm := e.GetResourceManager()
	rm.RegisterResource("2006.ttf")
	rm.RegisterResource("images/baby-yoda.jpg")
	rm.RegisterResource("images/babyyoda2.png")
	rm.RegisterResource("images/cyberpunk.jpg")
	rm.RegisterResource("images/gun.png")
	rm.RegisterResource("test.yaml")
	rm.RegisterResource("scenes/main.scene")
	rm.RegisterResource("scenes/second.scene")
	rm.RegisterResource("scenes/text.scene")
	rm.RegisterResource("prefabs/test.yaml")
	rm.RegisterResource("images/giphy.gif")
}

func RegisterComponents(e *engine.AmphionEngine) {
	cm := e.GetComponentsManager()
	cm.RegisterComponentType(&Mover{})
	cm.RegisterComponentType(&builtin.ShapeView{})
	cm.RegisterComponentType(&builtin.EllipseBoundary{})
	cm.RegisterComponentType(&builtin.OnClickListener{})
	cm.RegisterComponentType(&builtin.TextView{})
	cm.RegisterComponentType(&builtin.RectBoundary{})
	cm.RegisterComponentType(&builtin.TriangleBoundary{})
	cm.RegisterComponentType(&builtin.BezierView{})
	cm.RegisterComponentType(&builtin.DropdownView{})
	cm.RegisterComponentType(&builtin.ImageView{})
	cm.RegisterComponentType(&builtin.MouseMover{})
	cm.RegisterComponentType(&builtin.BuilderComponent{})
	cm.RegisterComponentType(&builtin.AbsoluteLayout{})
	cm.RegisterComponentType(&builtin.GridLayout{})
	cm.RegisterComponentType(&TestController{})
	cm.RegisterComponentType(&PrefabController{})

	cm.RegisterEventHandler(HandleCircleClick(e))
	cm.RegisterEventHandler(navigateOnClick)
	cm.RegisterEventHandler(navigateOnClick2)
}

func Scene2(e *engine.AmphionEngine) *engine.SceneObject {
	//var counter = 0
	scene2 := engine.NewSceneObject("scene 2")
	textScene2 := engine.NewSceneObject("text")
	textScene2.Transform.SetPosition(a.CenterInParent, a.CenterInParent, 0)
	textScene2.Transform.SetPivotCentered()
	textScene2.Transform.SetSize(800, 200)
	textScene2Renderer := builtin.NewTextView("This is scene 2")
	textScene2Renderer.FontSize = 100
	textScene2Renderer.TextColor = a.Black()
	textScene2.AddComponent(textScene2Renderer)
	textScene2.AddComponent(builtin.NewRectBoundary())
	//textScene2.AddComponent(builtin.NewOnClickListener(func(event engine.Event) bool {
	//	fmt.Println("Click")
	//	textScene2Renderer.SetText(strconv.Itoa(counter))
	//	textScene2Renderer.Redraw()
	//	e.RequestRendering()
	//	e.RequestUpdate()
	//	counter++
	//	return false
	//}))
	scene2.AddChild(textScene2)
	//scene2.AddComponent(&TestController{})

	input := engine.NewSceneObject("input")
	input.Transform.SetSize(500, 500)
	scene2.AddChild(input)

	dropdown := engine.NewSceneObject("dropdown")
	dropdown.Transform.SetPosition(10, 10, 2)
	dropdown.Transform.SetSize(100, 35, 0)
	dropdownView := builtin.NewDropdownView([]string{"opt1", "opt2", "opt3"})
	dropdown.AddComponent(dropdownView)
	dropdown.AddComponent(builtin.NewRectBoundary())
	dropdown.AddComponent(builtin.NewOnClickListener(func(event engine.Event) bool {
		dropdownView.HandleClick()
		return true
	}))

	dropdown1 := engine.NewSceneObject("dropdown1")
	dropdown1.Transform.SetPosition(10, 50, 1)
	dropdown1.Transform.SetSize(450, 35)
	dropdownView1 := builtin.NewDropdownView([]string{"bruh1", "bruh2", "bruh3"})
	dropdown1.AddComponent(dropdownView1)
	dropdown1.AddComponent(builtin.NewRectBoundary())
	dropdown1.AddComponent(builtin.NewOnClickListener(func(event engine.Event) bool {
		dropdownView1.HandleClick()
		return true
	}))

	box := engine.NewSceneObject("Moving box")
	box.Transform.SetPosition(10, 100, 10)
	box.Transform.SetSize(500, 500)
	boxBg := builtin.NewShapeView(builtin.ShapeRectangle)
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
	curve.Transform.SetPosition(10, 500, 5)
	curve.Transform.SetSize(100, 100)
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
	engine.LogDebug("Init")
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
	engine.LogDebug("Start")
}

func (c *TestController) GetName() string {
	return engine.NameOfComponent(c)
}

//endregion

func HandleCircleClick(e *engine.AmphionEngine) engine.EventHandler {
	return func(event engine.Event) bool {
		var mousePos = event.Data.(a.IntVector3)
		e.GetLogger().Info(nil, fmt.Sprintf("BIG CLICK ON CIRCLE. Mouse pos: %d %d", mousePos.X, mousePos.Y))
		return false
	}
}

func MakeRect(name string, x, y, sx, sy float32, fill a.Color) *engine.SceneObject {
	rect := engine.NewSceneObject(name)
	rect.Transform.SetPosition(x, y)
	rect.Transform.SetSize(sx, sy)

	view := builtin.NewShapeView(builtin.ShapeRectangle)
	view.FillColor = fill
	view.CornerRadius = 10

	rect.AddComponent(view)

	return rect
}

func GridScene(e *engine.AmphionEngine) *engine.SceneObject {
	scene := engine.NewSceneObject("grid scene")

	grid := builtin.NewGridLayout()
	grid.AddColumn(100)
	grid.AddColumn(a.FillParent)
	grid.AddColumn(100)
	grid.AddRow(a.WrapContent)
	//grid.ChildrenSizeFit = true
	grid.RowPadding = 10
	grid.ColumnPadding = 20
	scene.AddComponent(grid)

	//var counter int

	addBtn := MakeRect("add button", 0, 0, 100, 100, a.Green())
	addBtnText := engine.NewSceneObject("add text")
	addBtnText.Transform.SetPosition(0, 0, 1)
	addBtnText.Transform.SetSize(a.MatchParent, a.MatchParent)
	addBtnTextView := builtin.NewTextView("add")
	addBtnTextView.HTextAlign = a.TextAlignCenter
	addBtnTextView.VTextAlign = a.TextAlignCenter
	addBtnText.AddComponent(addBtnTextView)
	addBtn.AddChild(addBtnText)
	addBtn.AddComponent(builtin.NewRectBoundary())
	addBtn.AddComponent(builtin.NewOnClickListener(func(event engine.Event) bool {
		color := a.NewColor(byte(rand.Intn(256)), byte(rand.Intn(256)), byte(rand.Intn(256)), 255)
		rect := MakeRect(fmt.Sprintf("Rect"), 0, 0, 100, float32(rand.Intn(300)), color)
		rect.AddComponent(builtin.NewRectBoundary())
		scene.AddChild(rect)

		//counter++
		//
		//inputObj := engine.NewSceneObject(fmt.Sprintf("input %d", counter))
		//inputObj.Transform.size.Y = float32(rand.Intn(300))
		//input := builtin.NewNativeInputView("", fmt.Sprintf("Enter some text %d", counter))
		//input.SetOnTextChangeListener(func(text string) {
		//	engine.LogDebug(fmt.Sprintf("Hello from %s! Text: %s\n", inputObj.GetName(), text))
		//})
		//inputObj.AddComponent(input)
		//scene.AddChild(inputObj)

		return false
	}))
	addBtn.AddComponent(builtin.NewEventListener(engine.EventMouseIn, func(event engine.Event) bool {
		engine.LogInfo("Mouse in")
		shape := addBtn.GetComponentByName(".+ShapeView").(*builtin.ShapeView)
		shape.FillColor = a.Pink()
		shape.Redraw()
		engine.RequestRendering()
		return false
	}))
	addBtn.AddComponent(builtin.NewEventListener(engine.EventMouseOut, func(event engine.Event) bool {
		engine.LogInfo("Mouse out")
		shape := addBtn.GetComponentByName(".+ShapeView").(*builtin.ShapeView)
		shape.FillColor = a.Green()
		shape.Redraw()
		engine.RequestRendering()
		return false
	}))

	rmvButton := MakeRect("remove button", 0, 0, 100, 100, a.Red())
	rmvButtonText := engine.NewSceneObject("remove text")
	rmvButtonText.Transform.SetPosition(10, 10, 1)
	rmvButtonText.Transform.SetSize(a.MatchParent, a.MatchParent, 0)
	rmvButtonTextView := builtin.NewTextView("remove")
	rmvButtonText.AddComponent(rmvButtonTextView)
	//rmvButtonText.AddComponent(builtin.NewBoundaryView())
	rmvButton.AddChild(rmvButtonText)
	rmvButton.AddComponent(builtin.NewRectBoundary())
	rmvButton.AddComponent(builtin.NewOnClickListener(func(event engine.Event) bool {
		children := scene.GetChildren()
		last := children[len(children)-1]
		scene.RemoveChild(last)
		return false
	}))

	scene.AddChild(addBtn)
	scene.AddChild(rmvButton)

	e.BindEventHandler(engine.EventKeyDown, func(event engine.Event) bool {
		engine.LogDebug(fmt.Sprintf("%+v\n", event.Data))
		return false
	})

	e.BindEventHandler(engine.EventMouseDown, func(event engine.Event) bool {
		data := event.Data.(engine.MouseEventData)
		if data.SceneObject == nil {
			engine.LogDebug("Clicked on nothing")
		} else {
			engine.LogDebug("Clicked on %+v", event.Data.(engine.MouseEventData).SceneObject.GetName())
		}

		return true
	})

	var offset a.Vector3
	engine.BindEventHandler(engine.EventMouseScroll, func(event engine.Event) bool {
		o := event.Data.(a.Vector2)
		dOffset := a.NewVector3(o.X, o.Y, 0)

		engine.LogDebug("Scroll: %f %f", dOffset.X, dOffset.Y)
		//e.GetCurrentScene().ForEachObject(func(object *engine.sceneObject) {
		//	object.Transform.position = object.Transform.position.Add(dOffset)
		//})

		s := e.GetCurrentScene()
		ss := s.Transform.ActualSize()
		visibleArea := common.NewRect(-offset.X, -offset.X+ss.X, -offset.Y, -offset.Y+ss.Y, -999, 999)
		realArea := common.NewRect(0, 0, 0, 0, -999, 999)
		s.ForEachObject(func(object *engine.SceneObject) {
			rect := object.Transform.GlobalRect()
			if rect.X.Min < realArea.X.Min {
				realArea.X.Min = rect.X.Min
			}
			if rect.X.Max > realArea.X.Max {
				realArea.X.Max = rect.X.Max
			}
			if rect.Y.Min < realArea.Y.Min {
				realArea.Y.Min = rect.Y.Min
			}
			if rect.Y.Max > realArea.Y.Max {
				realArea.Y.Max = rect.Y.Max
			}
		})

		var scrollingDown, scrollingUp = dOffset.Y < 0, dOffset.Y > 0
		var canScrollDown, canScrollUp bool
		var minOutY float32 = -1

		s.ForEachObject(func(object *engine.SceneObject) {
			rect := object.Transform.GlobalRect()
			if !visibleArea.IsRectInside(rect) {
				//if dOffset.Y > 0 { // if scrolling up
				//	if rect.Y.Min < visibleArea.Y.Min {
				//		// can scroll
				//		//finalDY = common.ClampFloat32(dOffset.Y, -visibleArea.Y.Min + rect.Y.Min, 0)
				//		dOutY := visibleArea.Y.Min - rect.Y.Min
				//		if minOutY1 == -1 || dOutY < minOutY1 {
				//			minOutY1 = dOutY
				//		}
				//	}
				//} else if dOffset.Y < 0 { // scrolling down
				//	if rect.Y.Max > visibleArea.Y.Max {
				//		// can scroll
				//		//finalDY = common.ClampFloat32(dOffset.Y, -visibleArea.Y.Min + rect.Y.Min, 0)
				//		dOutY := rect.Y.Max - visibleArea.Y.Max
				//		if minOutY2 == -1 || dOutY < minOutY2 {
				//			minOutY2 = dOutY
				//		}
				//	}
				//}

				if scrollingDown && !canScrollDown { // down
					canScrollDown = rect.Y.Min > visibleArea.Y.Max || rect.Y.Max > visibleArea.Y.Max
					if canScrollDown {
						m := float32(math.Min(math.Abs(float64(rect.Y.Min-visibleArea.Y.Max)), math.Abs(float64(rect.Y.Max-visibleArea.Y.Max))))
						if minOutY == -1 || m < minOutY {
							minOutY = m
						}
					}
				} else if scrollingUp && !canScrollUp { // up
					canScrollUp = rect.Y.Min < visibleArea.Y.Min || rect.Y.Max < visibleArea.Y.Min
					if canScrollUp {
						m := float32(math.Min(math.Abs(float64(visibleArea.Y.Min-rect.Y.Min)), math.Abs(float64(visibleArea.Y.Min-rect.Y.Max))))
						if minOutY == -1 || m < minOutY {
							minOutY = m
						}
					}
				}
			}
		})

		if scrollingDown {
			if !canScrollDown {
				dOffset.Y = 0
			} else {
				dOffset.Y = float32(math.Min(float64(dOffset.Y), float64(-minOutY)))
			}
		}
		if scrollingUp {
			if !canScrollUp {
				dOffset.Y = 0
			} else {
				dOffset.Y = float32(math.Min(float64(dOffset.Y), float64(minOutY)))
			}
		}

		s.Transform.Translate(dOffset)
		offset = offset.Add(dOffset)
		e.ForceAllViewsRedraw()
		e.RequestRendering()
		return true
	})

	return scene
}

func navigateOnClick(_ engine.Event) bool {
	engine.LogDebug("Big click")

	err := engine.Navigate("second", nil)
	if err != nil {
		engine.LogError("Error navigating: %e", err)
	}
	return true
}

func navigateOnClick2(_ engine.Event) bool {
	engine.LogDebug("Big click")

	err := engine.Navigate("main", nil)
	if err != nil {
		engine.LogError("Error navigating: %e", err)
	}
	return true
}

func textScene(e *engine.AmphionEngine) *engine.SceneObject {
	scene := engine.NewSceneObject("text scene")

	text := engine.NewSceneObject("text")
	text.Transform.SetSize(200, 200)
	textView := builtin.NewTextView("Hello\nnext line\naaaaaaaaaaaaaa aaaaaaaaaaaaaa aaaaaaaaa aaaaaaaaaaaa\nывлоарывлаолывтёоылдомыволадыаааа дылаоыа\n\"!@#$%^&*()_+-={}[]🤢🌮")
	textView.FontSize = 16
	textView.HTextAlign = a.TextAlignCenter
	textView.VTextAlign = a.TextAlignCenter
	text.AddComponent(textView)
	text.AddComponent(builtin.NewRectBoundary())
	text.AddComponent(builtin.NewMouseMover())
	scene.AddChild(text)

	grogu := engine.NewSceneObject("grogu")
	grogu.Transform.SetSize(69, 69)
	grogu.Transform.SetPosition(0, 0, 69)
	grogu.AddComponent(builtin.NewImageView(Res_images_gun))
	grogu.AddComponent(builtin.NewRectBoundary())
	grogu.AddComponent(builtin.NewMouseMover())
	scene.AddChild(grogu)

	return scene
}

func TreeScene(e *engine.AmphionEngine) *engine.SceneObject {
	scene := engine.NewSceneObject("Tree")

	rect := engine.NewSceneObject("rect")
	rectView := builtin.NewShapeView(builtin.ShapeRectangle)
	rectView.FillColor = a.NewColor(0, 69, 0)
	rectView.CornerRadius = 100
	rect.AddComponent(rectView)
	rect.Transform.SetSize(120, 120)
	rect.AddComponent(builtin.NewMouseMover())
	rect.AddComponent(builtin.NewRectBoundary())

	rectInside := engine.NewSceneObject("rect inside")
	rectInsideView := builtin.NewShapeView(builtin.ShapeRectangle)
	rectInsideView.FillColor = a.NewColor(69, 0, 0)
	rectInsideView.CornerRadius = 5
	rectInside.AddComponent(rectInsideView)
	rectInside.Transform.SetSize(50, 50)
	rectInside.Transform.SetPosition(a.CenterInParent, a.CenterInParent)
	rectInside.Transform.SetPivotCentered()

	scene.AddChild(rect)
	rect.AddChild(rectInside)

	return scene
}

func ClipScene(e *engine.AmphionEngine) *engine.SceneObject {
	scene := engine.NewSceneObject("clip scene")

	image := engine.NewSceneObject("image")
	image.Transform.SetPosition(a.CenterInParent, a.CenterInParent)
	image.Transform.SetPivotCentered()
	image.Transform.SetSize(400, 400)
	imageView := builtin.NewImageView(Res_images_babyyoda)
	image.AddComponent(imageView)

	frame := engine.NewSceneObject("baby frame")
	frame.Transform.SetSize(a.MatchParent, a.MatchParent)
	circle := builtin.NewShapeView(builtin.ShapeEllipse)
	circle.FillColor = a.Transparent()
	circle.StrokeWeight = 10
	circle.StrokeColor = a.Pink()
	frame.AddComponent(circle)
	image.AddChild(frame)

	image.AddComponent(builtin.NewClipArea(shape.Circle))
	scene.AddChild(image)

	rect := engine.NewSceneObject("rect")
	rect.Transform.SetSize(50, 50)
	rect.Transform.SetPivotCentered()
	rect.Transform.SetPosition(a.CenterInParent, a.CenterInParent)
	textView := builtin.NewTextView("BRUH BRUH BRUH BRUH BRUH BRUH BRUH BRUH BRUH BRUH BRUH BRUH BRUH BRUH BRUH BRUH BRUH BRUH BRUH BRUH BRUH BRUH ")
	textView.FontSize = 10
	textView.TextColor = a.Red()
	rect.AddComponent(textView)
	rect.AddComponent(builtin.NewClipArea(shape.Circle))
	image.AddChild(rect)

	return scene
}
