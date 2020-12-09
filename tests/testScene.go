package tests

import (
	"fmt"
	"github.com/cadmean-ru/amphion/common"
	"github.com/cadmean-ru/amphion/engine"
	"github.com/cadmean-ru/amphion/engine/builtin"
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

func createTestScene() *engine.SceneObject {
	scene := engine.NewSceneObject("Test scene")

	rect := engine.NewSceneObject("rect")
	rect.Transform.Size = common.NewVector3(100, 100, 100)
	rect.Transform.Position = common.NewVector3(100, 100, -2)
	rect.AddComponent(builtin.NewShapeView(rendering.PrimitiveRectangle))
	rect.AddComponent(&mover{})

	circle := engine.NewSceneObject("circle")
	circle.Transform.Size = common.NewVector3(50, 50, 0)
	circle.Transform.Position = common.NewVector3(10, 10 , 1)
	circleRenderer := builtin.NewShapeView(rendering.PrimitiveEllipse)
	circleRenderer.Appearance.StrokeWeight = 0
	circleRenderer.Appearance.FillColor = common.GreenColor()
	circle.AddComponent(circleRenderer)
	//circle.AddComponent(builtin.NewBoundaryView())

	rect.AddChild(circle)
	scene.AddChild(rect)

	text := engine.NewSceneObject("centeredText")
	text.Transform.Position = common.NewVector3(engine.CenterInParent, engine.CenterInParent, engine.CenterInParent)
	text.Transform.Pivot = common.NewVector3(0.5, 0.5, 0.5)
	text.Transform.Size = common.NewVector3(200, 50, 0)
	textComponent := builtin.NewTextView("Center")
	textComponent.TextAppearance.FontSize = 30
	textComponent.Appearance.FillColor = common.BlackColor()
	text.AddComponent(textComponent)
	text.AddComponent(builtin.NewBoundaryView())

	for i := 0; i < 10; i++ {
		text1 := engine.NewSceneObject(fmt.Sprintf("text%d", i))
		text1.Transform.Position = common.NewVector3(10, float64(i) * 50, 0)
		text1.Transform.Size = common.NewVector3(200, 50, 0)
		textComponent1 := builtin.NewTextView(fmt.Sprintf("Bruh %d", i))
		textComponent1.TextAppearance.FontSize = 30
		textComponent1.Appearance.FillColor = common.BlackColor()
		text1.AddComponent(textComponent1)
		text1.AddComponent(builtin.NewBoundaryView())
		scene.AddChild(text1)
	}

	scene.AddChild(text)

	return scene
}

var testEngineInstance *engine.AmphionEngine

func startEngineTest() *engine.AmphionEngine {
	p := common.PlatformFromString("test")
	e := engine.Initialize(p)
	testEngineInstance = e
	e.Start()
	return e
}
