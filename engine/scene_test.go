package engine

import (
	"fmt"
	"github.com/cadmean-ru/amphion/common"
	"testing"
)

type testComponent struct {
	ComponentImpl
	num  int
}

func (m *testComponent) OnStart() {
	m.Logger.Info(m, "Start")
}

func (m *testComponent) OnStop() {
	m.Logger.Info(m, "Stop")
}

func (m *testComponent) GetInstanceState() common.SiMap {
	return map[string]interface{}{
		"num": m.num,
	}
}

func (m *testComponent) SetInstanceState(state common.SiMap) {
	m.num = common.RequireInt(state["num"])
}

func (m *testComponent) GetName() string {
	return m.NameOf(m)
}


func createTestScene() *SceneObject {
	scene := NewSceneObject("Test scene")

	rect := NewSceneObject("rect")
	rect.Transform.Size = common.NewVector3(100, 100, 100)
	rect.Transform.Position = common.NewVector3(100, 100, -2)
	test := &testComponent{}
	test.num = 69
	rect.AddComponent(test)

	circle := NewSceneObject("circle")
	circle.Transform.Size = common.NewVector3(50, 50, 0)
	circle.Transform.Position = common.NewVector3(10, 10 , 1)

	rect.AddChild(circle)
	scene.AddChild(rect)

	return scene
}

func TestSceneObject_EncodeToYaml(t *testing.T) {
	runEngineTest(func(e *AmphionEngine) {
		scene := createTestScene()
		data, err := scene.EncodeToYaml()
		if err != nil {
			t.Error(err)
		}
		fmt.Println(string(data))
	})
}

func TestSceneObject_DecodeFromYaml(t *testing.T) {
	var srcYaml = `
children:
- children:
  - children: []
    components: []
    id: 2
    name: circle
    transform:
      pivot:
        x: 0
        "y": 0
        z: 0
      position:
        x: 10
        "y": 10
        z: 1
      rotation:
        x: 0
        "y": 0
        z: 0
      size:
        x: 50
        "y": 50
        z: 0
  components:
  - name: github.com/cadmean-ru/amphion/engine.testComponent
    state:
      num: 69
  id: 1
  name: rect
  transform:
    pivot:
      x: 0
      "y": 0
      z: 0
    position:
      x: 100
      "y": 100
      z: -2
    rotation:
      x: 0
      "y": 0
      z: 0
    size:
      x: 100
      "y": 100
      z: 100
components: []
id: 0
name: Test scene
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

	runEngineTest(func(e *AmphionEngine) {
		cm := e.GetComponentsManager()
		cm.RegisterComponentType(&testComponent{})

		data := []byte(srcYaml)
		scene := &SceneObject{}

		err := scene.DecodeFromYaml(data)
		if err != nil {
			t.Error(err)
		}

		rect := scene.GetChildByName("rect")
		circle := rect.GetChildByName("circle")
		comp := rect.GetComponentByName((&testComponent{}).GetName())

		fmt.Printf("%+v\n", scene)
		fmt.Printf("%+v\n", rect)
		fmt.Printf("%+v\n", comp)
		fmt.Printf("%+v\n", circle)
	})
}