package main

import (
	"fmt"
	"github.com/cadmean-ru/amphion/common/a"
	"github.com/cadmean-ru/amphion/engine"
	"github.com/cadmean-ru/amphion/engine/builtin"
	"github.com/cadmean-ru/amphion/rendering"
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
	circle.AddComponent(builtin.NewOnClickListener(func(event engine.AmphionEvent) bool {
		var mousePos = event.Data.(a.IntVector3)
		e.GetLogger().Info(nil, fmt.Sprintf("BIG CLICK ON CIRCLE. Mouse pos: %d %d", mousePos.X, mousePos.Y))
		//e.CloseScene(func() {
		//	_ = e.ShowScene(createCyberpunkScene(e))
		//})
		return false
	}))

	rect.AddChild(circle)


	scene.AddChild(rect)

	text := engine.NewSceneObject("Close text")
	text.Transform.Position = a.NewVector3(engine.CenterInParent, engine.CenterInParent, engine.CenterInParent)
	text.Transform.Pivot = a.NewVector3(0.5, 0.5, 0.5)
	text.Transform.Size = a.NewVector3(300, 50, 0)
	textComponent := builtin.NewTextView("Hello Amphion!")
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
	imageView := builtin.NewImageView(3)
	image.AddComponent(imageView)
	scene.AddChild(image)

	image2 := engine.NewSceneObject("image2")
	image2.Transform.Position = a.NewVector3(200, 300, -1)
	image2.Transform.Size = a.NewVector3(100, 300, 0)
	imageView2 := builtin.NewImageView(2)
	image2.AddComponent(imageView2)
	scene.AddChild(image2)

	return scene
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
}