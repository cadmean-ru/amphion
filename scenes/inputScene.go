package scenes

import (
	. "github.com/cadmean-ru/amphion/engine"
	. "github.com/cadmean-ru/amphion/engine/builtin"
)

func InputScene(_ *AmphionEngine) *SceneObject {
	scene := NewSceneObject("input scene")

	input := NewSceneObject("input")
	input.Transform.SetSize(300, 50)
	inputView := NewTextInput()
	input.AddComponent(inputView)
	scene.AddChild(input)

	return scene
}
