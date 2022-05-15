package scenes

import (
	. "github.com/cadmean-ru/amphion/common"
	"github.com/cadmean-ru/amphion/common/a"
	. "github.com/cadmean-ru/amphion/engine"
	. "github.com/cadmean-ru/amphion/engine/builtin"
)

func InputScene(_ *AmphionEngine) *SceneObject {
	scene := NewSceneObject("input scene")
	scene.AddComponent(NewPadding(NewEdgeInsets(10, 10, 10, 10)))

	input := NewSceneObject("input")
	input.Transform.SetSize(500, a.WrapContent)
	inputView := NewTextInput()
	inputView.FontSize = 30
	inputView.SingleLine = true
	inputView.SetOnChangeListener(func(newValue string) {
		LogDebug("Text change: %s", newValue)
	})
	input.AddComponent(inputView)
	scene.AddChild(input)

	return scene
}
