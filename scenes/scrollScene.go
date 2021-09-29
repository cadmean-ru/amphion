package scenes

import (
	. "github.com/cadmean-ru/amphion/common/a"
	. "github.com/cadmean-ru/amphion/engine"
	. "github.com/cadmean-ru/amphion/engine/builtin"
)

func ScrollScene(_ *AmphionEngine) *SceneObject {
	scene := NewSceneObject("scoll scene")
	scene.AddComponent(NewScroll())

	grid := NewSceneObject("grid")
	grid.Transform.SetSize(MatchParent, WrapContent)
	gridLayout := NewGridLayout()
	gridLayout.AddColumn(FillParent)
	grid.AddComponent(gridLayout)
	scene.AddChild(grid)

	scroll := NewSceneObject("scroll")
	scroll.Transform.SetPosition(100, 100)
	scroll.Transform.SetSize(200, 300)
	scroll.AddComponent(NewScroll())
	scroll.AddComponent(NewBoundaryView())
	grid.AddChild(scroll)

	text := NewSceneObject("text")
	text.Transform.SetSize(300, WrapContent)
	textView := NewTextView("fdgsdkj329487ieokdfjlkdsfkldsfkl;jdsfl;jkdl;fjdkg jdflgkdsfjsadjkl jl;jkldsjklgkjldaskljgsg")
	textView.FontSize = 50
	text.AddComponent(textView)
	scroll.AddChild(text)

	rect1 := NewSceneObject("rect1")
	rect1.Transform.SetSize(10, 200)
	rect1View := NewShapeView(ShapeRectangle)
	rect1View.FillColor = Blue()
	rect1.AddComponent(rect1View)
	grid.AddChild(rect1)

	rect2 := NewSceneObject("rect2")
	rect2.Transform.SetSize(10, 300)
	rect2View := NewShapeView(ShapeRectangle)
	rect2View.FillColor = Green()
	rect2.AddComponent(rect2View)
	grid.AddChild(rect2)

	BindEventHandler(EventKeyUp, func(event AmphionEvent) bool {
		textView.SetText(textView.Text + event.StringData())
		return true
	})

	return scene
}