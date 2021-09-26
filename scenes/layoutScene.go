package scenes

import (
	. "github.com/cadmean-ru/amphion/common/a"
	. "github.com/cadmean-ru/amphion/engine"
	. "github.com/cadmean-ru/amphion/engine/builtin"
	"time"
)

type LayoutSceneController struct {
    ComponentImpl
}

func (s *LayoutSceneController) OnStart() {
    LogInfo("Started bruh")

	RunTask(NewTaskBuilder().Run(func() (interface{}, error) {
		time.Sleep(3 * time.Second)
		json, _ := GetCurrentScene().DumpToJson()
		LogDebug(string(json))
		return nil, nil
	}).Build())
}

func NewLayoutSceneController() *LayoutSceneController {
    return &LayoutSceneController{}
}

func LayoutScene(e *AmphionEngine) *SceneObject {
	scene := NewSceneObject("layout scene")
	scene.AddComponent(NewLayoutSceneController())
	scene.AddComponent(NewAbsoluteLayout())
	bg := NewShapeView(ShapeRectangle)
	bg.FillColor = Blue()
	scene.AddComponent(bg)

	rect := NewSceneObject("rect")
	rect.Transform.SetPosition(10, 10)
	rect.Transform.SetSize(WrapContent, WrapContent)
	rectView := NewShapeView(ShapeRectangle)
	rectView.FillColor = Pink()
	rect.AddComponent(rectView)
	scene.AddChild(rect)

	textContainer := NewSceneObject("text container")
	textContainer.Transform.CenterInParent()
	textContainerBg := NewShapeView(ShapeRectangle)
	textContainerBg.FillColor = White()
	textContainer.AddComponent(textContainerBg)
	scene.AddChild(textContainer)

	text := NewSceneObject("text")
	textView := NewTextView("hello")
	textView.FontSize = 69
	text.AddComponent(textView)
	textContainer.AddChild(text)

	bigContainer := NewSceneObject("big container")
	bigContainerBg := NewShapeView(ShapeRectangle)
	bigContainerBg.FillColor = Red()
	bigContainerBg.StrokeColor = White()
	bigContainerBg.StrokeWeight = 3
	bigContainer.AddComponent(bigContainerBg)
	circle1 := NewSceneObject("circle 1")
	circle1.Transform.SetSize(20, 20)
	circleView1 := NewShapeView(ShapeEllipse)
	circleView1.FillColor = NewColor("#123")
	circle1.AddComponent(circleView1)
	bigContainer.AddChild(circle1)
	circle2 := NewSceneObject("circle 2")
	circle2.Transform.SetSize(30, 30)
	circle2.Transform.SetPosition(20, 20)
	circleView2 := NewShapeView(ShapeEllipse)
	circleView2.FillColor = NewColor("#456")
	circle2.AddComponent(circleView2)
	bigContainer.AddChild(circle2)
	scene.AddChild(bigContainer)

	grid := NewSceneObject("grid")
	grid.Transform.SetPosition(100, 100)
	grid.Transform.SetSize(MatchParent, WrapContent)
	gridLayout := NewGridLayout()
	gridLayout.AddColumn(WrapContent)
	gridLayout.AddColumn(FillParent)
	gridLayout.RowPadding = 10
	gridLayout.ColumnPadding = 10
	grid.AddComponent(gridLayout)
	gridBg := NewShapeView(ShapeRectangle)
	gridBg.FillColor = Pink()
	grid.AddComponent(gridBg)

	for i := 0; i < 10; i++ {
		gridRect := NewSceneObject("grid rect")
		gridRect.Transform.SetSize(10, 10)
		gridRectView := NewShapeView(ShapeRectangle)
		gridRectView.FillColor = Green()
		gridRect.AddComponent(gridRectView)
		grid.AddChild(gridRect)
	}

	scene.AddChild(grid)

	return scene
}
