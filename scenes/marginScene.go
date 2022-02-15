package scenes

import (
	"github.com/cadmean-ru/amphion/common"
	"github.com/cadmean-ru/amphion/common/a"
	"github.com/cadmean-ru/amphion/engine"
	"github.com/cadmean-ru/amphion/engine/builtin"
)

func MarginScene(_ *engine.AmphionEngine) *engine.SceneObject {
	scene := engine.NewSceneObject("marging scene")
	scene.AddComponent(builtin.NewAbsoluteLayout())

	box := engine.NewSceneObject("big box")
	box.Transform.SetSize(50, 50)
	boxRect := builtin.NewShapeView(builtin.ShapeRectangle)
	boxRect.FillColor = a.Red()
	box.AddComponent(boxRect)
	marging := builtin.NewMargin(common.NewEdgeInsets(10, 10, 10, 10))
	box.AddComponents(marging)
	//box.AddComponent(builtin.NewPadding(common.NewEdgeInsets(10, 10, 10, 10)))

	//smolBox := engine.NewSceneObject("smol box")
	//smolBox.Transform.SetSize(10, 10)
	//smolRect := builtin.NewShapeView(builtin.ShapeRectangle)
	//smolRect.FillColor = a.Green()
	//smolBox.AddComponent(smolRect)
	box.AddComponent(builtin.NewRectBoundary())
	box.AddComponent(builtin.NewEventListener(engine.EventMouseDown, func(event engine.Event) bool {
		marging.SetMargin(common.NewEdgeInsets(20, 1, 3, 4))
		return true
	}))
	//box.AddChild(smolBox)

	scene.AddChild(box)

	return scene
}
