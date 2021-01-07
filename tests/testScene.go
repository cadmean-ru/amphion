package tests

import (
	"fmt"
	"github.com/cadmean-ru/amphion/common/a"
	"github.com/cadmean-ru/amphion/engine"
	"github.com/cadmean-ru/amphion/engine/builtin"
	"github.com/cadmean-ru/amphion/frontend/pc"
	"github.com/cadmean-ru/amphion/rendering"
)

type mover struct {
	obj  *engine.SceneObject
	log  *engine.Logger
	prc  *builtin.ShapeView
	eng  *engine.AmphionEngine
	dir  bool
}

func (m *mover) OnInit(ctx engine.InitContext) {
	m.obj = ctx.GetSceneObject()
	m.log = ctx.GetLogger()
	m.prc = m.obj.GetComponentByName("ShapeView").(*builtin.ShapeView)
	m.eng = ctx.GetEngine()
}

func (m *mover) OnStart() {
	m.log.Info(m, "Start")
}

func (m *mover) OnUpdate(ctx engine.UpdateContext) {
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

func (m *mover) OnStop() {
	m.log.Info(m, "Stop")
}

func (m *mover) GetName() string {
	return "mover"
}

func CreateTestScene() *engine.SceneObject {
	scene := engine.NewSceneObject("Test scene")

	rect := engine.NewSceneObject("rect")
	rect.Transform.Size = a.NewVector3(100, 100, 100)
	rect.Transform.Position = a.NewVector3(100, 100, -2)
	rect.AddComponent(builtin.NewShapeView(rendering.PrimitiveRectangle))
	rect.AddComponent(&mover{})

	circle := engine.NewSceneObject("circle")
	circle.Transform.Size = a.NewVector3(50, 50, 0)
	circle.Transform.Position = a.NewVector3(10, 10 , 1)
	circleRenderer := builtin.NewShapeView(rendering.PrimitiveEllipse)
	circleRenderer.Appearance.StrokeWeight = 0
	circleRenderer.Appearance.FillColor = a.GreenColor()
	circle.AddComponent(circleRenderer)
	//circle.AddComponent(builtin.NewBoundaryView())

	rect.AddChild(circle)
	scene.AddChild(rect)

	text := engine.NewSceneObject("centeredText")
	text.Transform.Position = a.NewVector3(engine.CenterInParent, engine.CenterInParent, engine.CenterInParent)
	text.Transform.Pivot = a.NewVector3(0.5, 0.5, 0.5)
	text.Transform.Size = a.NewVector3(200, 50, 0)
	textComponent := builtin.NewTextView("Center")
	textComponent.TextAppearance.FontSize = 30
	textComponent.Appearance.FillColor = a.BlackColor()
	text.AddComponent(textComponent)
	text.AddComponent(builtin.NewBoundaryView())

	for i := 0; i < 10; i++ {
		text1 := engine.NewSceneObject(fmt.Sprintf("text%d", i))
		text1.Transform.Position = a.NewVector3(10, float32(i) * 50, 0)
		text1.Transform.Size = a.NewVector3(200, 50, 0)
		textComponent1 := builtin.NewTextView(fmt.Sprintf("Bruh %d", i))
		textComponent1.TextAppearance.FontSize = 30
		textComponent1.Appearance.FillColor = a.BlackColor()
		text1.AddComponent(textComponent1)
		text1.AddComponent(builtin.NewBoundaryView())
		scene.AddChild(text1)
	}

	scene.AddChild(text)

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
	textScene2Renderer.TextAppearance.FontSize = 100
	textScene2Renderer.Appearance.FillColor = a.BlackColor()
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
	inputView := builtin.NewInputView()
	inputView.TextAppearance.FontSize = 10
	inputView.Appearance.FillColor = a.BlackColor()
	input.AddComponent(inputView)
	scene2.AddChild(input)

	dropdown := engine.NewSceneObject("dropdown")
	dropdown.Transform.Position = a.NewVector3(10, 10, 1)
	dropdown.Transform.Size = a.NewVector3(450, 35, 0)
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
	dropdown1.AddComponent(builtin.NewBoundaryView())
	dropdown1.AddComponent(builtin.NewRectBoundary())
	dropdown1.AddComponent(builtin.NewTextView("dfdsf"))
	dropdown1.AddComponent(builtin.NewOnClickListener(func(event engine.AmphionEvent) bool {
		dropdownView1.HandleClick()
		return true
	}))


	box := engine.NewSceneObject("Moving box")
	box.Transform.Position = a.NewVector3(10, 100, 1)
	box.Transform.Size = a.NewVector3(500, 500, 0)
	boxBg := builtin.NewShapeView(rendering.PrimitiveRectangle)
	boxBg.Appearance.StrokeWeight = 0
	boxBg.Appearance.FillColor = a.NewColor(0xc4, 0xc4, 0xc4, 0xff)
	boxBg.Appearance.CornerRadius = 10
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

var testEngineInstance *engine.AmphionEngine

func startEngineTest() *engine.AmphionEngine {
	front := pc.NewFrontend()
	e := engine.Initialize(front)
	testEngineInstance = e
	e.Start()
	return e
}
