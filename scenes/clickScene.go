package scenes

import (
	. "github.com/cadmean-ru/amphion/common/a"
	"github.com/cadmean-ru/amphion/common/shape"
	. "github.com/cadmean-ru/amphion/engine"
	. "github.com/cadmean-ru/amphion/engine/builtin"
)

func ClickScene(_ *AmphionEngine) *SceneObject {
	scene := NewSceneObject("click scene")
	scene.AddComponent(NewAbsoluteLayout())
	scene.AddComponent(NewStartFunc(func() {
		LogDebug("Start func")
	}))
	scene.AddComponent(NewUpdateFunc(func(ctx UpdateContext) {
		LogDebug("Update func: %f", ctx.DeltaTime)
	}))

	container := NewSceneObject("container")
	container.Transform.SetSize(100, 100)
	container.Transform.CenterInParent()
	container.AddComponent(NewBoundaryView())
	container.AddComponent(NewClipArea(shape.Ellipse))
	scene.AddChild(container)

	text := NewSceneObject("text")
	textView := NewTextView("")
	text.AddComponent(textView)
	scene.AddChild(text)

	rect1 := NewSceneObject("rect1")
	rect1.Transform.SetSize(50, 50)
	rect1.Transform.SetPosition(-10, -10)
	rect1Shape := NewShapeView(ShapeRectangle)
	rect1Shape.FillColor = Blue()
	rect1.AddComponent(rect1Shape)
	rect1.AddComponent(NewRectBoundary())
	rect1.AddComponent(NewEventListener(EventMouseDown, func(event AmphionEvent) bool {
		LogDebug("Click")
		textView.Text += "1"
		textView.Redraw()
		RequestRendering()
		return true
	}))
	container.AddChild(rect1)

	circle1 := NewSceneObject("circle1")
	circle1.Transform.SetSize(30, 30)
	circle1.Transform.SetPosition(75, 75)
	circle1Shape := NewShapeView(ShapeEllipse)
	circle1Shape.FillColor = Red()
	circle1Shape.StrokeWeight = 3
	circle1.AddComponent(circle1Shape)
	circle1.AddComponent(NewRectBoundary())
	circle1.AddComponent(NewEventListener(EventMouseDown, func(event AmphionEvent) bool {
		LogDebug("Click")
		textView.Text += "2"
		textView.Redraw()
		RequestRendering()
		return true
	}))
	circle1.AddComponent(NewEventListener(EventMouseIn, func(event AmphionEvent) bool {
		circle1Shape.StrokeColor = Blue()
		circle1Shape.Redraw()
		RequestRendering()
		return true
	}))
	circle1.AddComponent(NewEventListener(EventMouseOut, func(event AmphionEvent) bool {
		circle1Shape.StrokeColor = Black()
		circle1Shape.Redraw()
		RequestRendering()
		return true
	}))
	container.AddChild(circle1)

	return scene
}
