package scenes

import (
	"github.com/cadmean-ru/amphion/engine"
	"github.com/cadmean-ru/amphion/engine/builtin"
)

func TextScene(_ *engine.AmphionEngine) *engine.SceneObject {
	scene := engine.NewSceneObject("Text scene")

	//textObject := engine.NewSceneObject("Text")
	//text := builtin.NewTextView("Hello")
	//textObject.AddComponent(text)
	//scene.AddChild(textObject)

	y := `
children:
  - children: []
    components:
      - name: TextView
        state:
          text: Hello world
          fontSize: 16
          textColor: #696969
          hTextAlign: $TextAlignCenter
    id: 1
    name: text
    transform:
      pivot:
        x: 0
        "y": 0
        z: 0
      position:
        x: 0
        "y": 0
        z: 0
      rotation:
        x: 0
        "y": 0
        z: 0
      size:
        x: $MatchParent
        "y": $WrapContent
        z: 0
components:
  - name: AbsoluteLayout
    state: {}
id: 0
name: Text test
transform:
  pivot:
    x: 0
    "y": 0
    z: 0
  position:
    x: 0
    "y": 0
    z: 0
  rotation:
    x: 0
    "y": 0
    z: 0
  size:
    x: 1
    "y": 1
    z: 1
`
	scene.DecodeFromYaml([]byte(y))

	scene.AddComponent(builtin.NewComponentBuilder("debug").OnStart(func() {
		engine.LogDebug("Bruh %+v", scene.GetChildAt(0))
	}).Build())

	return scene
}
