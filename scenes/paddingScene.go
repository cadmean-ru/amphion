package scenes

import (
	"github.com/cadmean-ru/amphion/common"
	"github.com/cadmean-ru/amphion/common/a"
	"github.com/cadmean-ru/amphion/engine"
	"github.com/cadmean-ru/amphion/engine/builtin"
)

func PaddingScene(_ *engine.AmphionEngine) *engine.SceneObject {
	scene := engine.NewSceneObject("padding scene")
	scene.AddComponent(builtin.NewAbsoluteLayout())

	paddingObj := engine.NewSceneObject("paddin")
	paddingObj.AddComponent(builtin.NewPadding(common.NewEdgeInsets(10, 10, 10, 10)))
	paddingBg := builtin.NewShapeView(builtin.ShapeRectangle)
	paddingBg.FillColor = a.Pink()
	paddingObj.AddComponent(paddingBg)

	box := engine.NewSceneObject("big box")
	box.Transform.SetSize(50, 50)
	boxRect := builtin.NewShapeView(builtin.ShapeRectangle)
	boxRect.FillColor = a.Red()
	box.AddComponent(boxRect)

	box2 := engine.NewSceneObject("big box")
	box2.Transform.SetSize(39, 69)
	boxRect2 := builtin.NewShapeView(builtin.ShapeRectangle)
	boxRect2.FillColor = a.Blue()
	box2.AddComponent(boxRect2)

	paddingObj.AddChild(box)
	paddingObj.AddChild(box2)
	scene.AddChild(paddingObj)

	return scene
}
