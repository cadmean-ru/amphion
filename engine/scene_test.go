package engine

import (
	"fmt"
	"github.com/cadmean-ru/amphion/common/a"
	"github.com/cadmean-ru/amphion/common/require"
	"github.com/stretchr/testify/assert"
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

func (m *testComponent) GetInstanceState() a.SiMap {
	return map[string]interface{}{
		"num": m.num,
	}
}

func (m *testComponent) SetInstanceState(state a.SiMap) {
	m.num = require.Int(state["num"])
}

func (m *testComponent) GetName() string {
	return NameOfComponent(m)
}


func createTestScene() *SceneObject {
	scene := NewSceneObject("Test scene")

	rect := NewSceneObject("rect")
	rect.Transform.Size = a.NewVector3(100, 100, 100)
	rect.Transform.Position = a.NewVector3(100, 100, -2)
	test := &testComponent{}
	test.num = 69
	test2 := &testStatefulWithTags{}
	test2.Bruh = "Nice"
	test2.Bruh2 = 2
	test2.Color = a.GreenColor()
	test2.Arr = []int {42, 69}
	test2.Hand = func(_ AmphionEvent) bool {
		fmt.Println("Handle breh")
		return false
	}
	rect.AddComponent(test)
	rect.AddComponent(test2)

	circle := NewSceneObject("circle")
	circle.Transform.Size = a.NewVector3(50, 50, 0)
	circle.Transform.Position = a.NewVector3(10, 10 , 1)

	rect.AddChild(circle)
	scene.AddChild(rect)

	return scene
}

func TestSceneObject_EncodeToYaml(t *testing.T) {
	//scene := createTestScene()
	//data, err := scene.EncodeToYaml()
	//if err != nil {
	//	t.Error(err)
	//}
	//fmt.Println(string(data))
}

func TestSceneObject_DecodeFromYaml(t *testing.T) {
//	var srcYaml = `
//children:
//- children:
//  - children: []
//    components: []
//    id: 2
//    name: circle
//    transform:
//      pivot:
//        x: 0
//        "y": 0
//        z: 0
//      position:
//        x: 10
//        "y": 10
//        z: 1
//      rotation:
//        x: 0
//        "y": 0
//        z: 0
//      size:
//        x: 50
//        "y": 50
//        z: 0
//  components:
//  - name: github.com/cadmean-ru/amphion/engine.testComponent
//    state:
//      num: 69
//  - name: github.com/cadmean-ru/amphion/engine.testStatefulWithTags
//    state:
//      Bruh: Nice
//      arr:
//      - 42
//      - 69
//      breh: 2
//      color:
//        a: 255
//        b: 0
//        g: 255
//        r: 0
//  id: 1
//  name: rect
//  transform:
//    pivot:
//      x: 0
//      "y": 0
//      z: 0
//    position:
//      x: 100
//      "y": 100
//      z: -2
//    rotation:
//      x: 0
//      "y": 0
//      z: 0
//    size:
//      x: 100
//      "y": 100
//      z: 100
//components: []
//id: 0
//name: Test scene
//transform:
//  pivot:
//    x: 0
//    "y": 0
//    z: 0
//  position:
//    x: 0
//    "y": 0
//    z: 0
//  rotation:
//    x: 0
//    "y": 0
//    z: 0
//  size:
//    x: 1
//    "y": 1
//    z: 1
//`
//
//	atest.RunEngineTest(t, func(e *AmphionEngine) {
//		cm := e.GetComponentsManager()
//		cm.RegisterComponentType(&testComponent{})
//		cm.RegisterComponentType(&testStatefulWithTags{})
//
//		data := []byte(srcYaml)
//		scene := &SceneObject{}
//
//		err := scene.DecodeFromYaml(data)
//		if err != nil {
//			t.Error(err)
//		}
//
//		rect := scene.GetChildByName("rect")
//		circle := rect.GetChildByName("circle")
//		comp := rect.GetComponentByName((&testComponent{}).GetName())
//		comp2 := rect.GetComponentByName((&testStatefulWithTags{}).GetName())
//
//		fmt.Printf("%+v\n", scene)
//		fmt.Printf("%+v\n", rect)
//		fmt.Printf("%+v\n", comp)
//		fmt.Printf("%+v\n", comp2)
//		fmt.Printf("%+v\n", circle)
//	})
}

func TestSceneObject_GetComponentByName(t *testing.T) {
	ass := assert.New(t)

	o := NewSceneObjectForTesting("Testing object")
	component := &testComponent{}
	o.AddComponent(component)

	res := o.GetComponentByName("testComponent")
	ass.NotNil(res)

	res = o.GetComponentByName(".+Component")
	ass.NotNil(res)

	res = o.GetComponentByName("github.com/cadmean-ru/amphion/engine.testComponent")
	ass.NotNil(res)
}

func TestSceneObject_FindObjectByName(t *testing.T) {
	o1 := NewSceneObjectForTesting("Testing object 1")
	o2 := NewSceneObjectForTesting("Testing object 2")
	o1.AddChild(o2)
	o3 := NewSceneObjectForTesting("Testing object 3")
	o1.AddChild(o3)
	o4 := NewSceneObjectForTesting("Testing object 4")
	o2.AddChild(o4)
	o5 := NewSceneObjectForTesting("Testing object 5")
	o2.AddChild(o5)

	ass := assert.New(t)

	target := o1
	found := o1.FindObjectByName("Testing object 1")
	ass.Equal(target, found)

	target = o3
	found = o1.FindObjectByName("Testing object 3")
	ass.Equal(target, found)

	target = o5
	found = o1.FindObjectByName("Testing object 5")
	ass.Equal(target, found)
}

func TestSceneObject_FindComponentByName(t *testing.T) {
	o1 := NewSceneObjectForTesting("Testing object 1")
	o2 := NewSceneObjectForTesting("Testing object 2")
	o1.AddChild(o2)
	o3 := NewSceneObjectForTesting("Testing object 3")
	o1.AddChild(o3)
	o4 := NewSceneObjectForTesting("Testing object 4")
	o2.AddChild(o4)

	target := &testComponent{}
	o5 := NewSceneObjectForTesting("Testing object 5", target)
	o2.AddChild(o5)

	found := o1.FindComponentByName("testComponent")

	assert.Equal(t, target, found)
}