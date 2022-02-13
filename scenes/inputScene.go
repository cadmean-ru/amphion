package scenes

import (
	. "github.com/cadmean-ru/amphion/common"
	. "github.com/cadmean-ru/amphion/engine"
	. "github.com/cadmean-ru/amphion/engine/builtin"
)

func InputScene(_ *AmphionEngine) *SceneObject {
	scene := NewSceneObject("input scene")
	scene.AddComponent(NewPadding(NewEdgeInsets(10, 10, 10, 10)))

	input := NewSceneObject("input")
	input.Transform.SetSize(500, 300)
	inputView := NewTextInput()
	inputView.SetOnChangeListener(func(newValue string) {
		LogDebug("Text change: %s", newValue)
	})
	input.AddComponent(inputView)
	scene.AddChild(input)

	return scene
}