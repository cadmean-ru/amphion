package scenes

import (
	"github.com/cadmean-ru/amphion/common/a"
	"github.com/cadmean-ru/amphion/engine"
	"github.com/cadmean-ru/amphion/engine/builtin"
)

//region 1
type TestPrefabSceneController struct {
	engine.ComponentImpl
}

func (c *TestPrefabSceneController) OnInit(ctx engine.InitContext) {
	c.ComponentImpl.OnInit(ctx)
	engine.LogDebug("Init")

	prefab, err := engine.LoadPrefab(a.ResId(Res_prefabs_test))
	if err != nil {
		return
	}

	c.SceneObject.AddChild(prefab)
}

func (c *TestPrefabSceneController) OnStart() {
	engine.LogDebug("Start")
}

func (c TestPrefabSceneController) GetName() string {
	return engine.NameOfComponent(c)
}

//endregion 1

//region 2

type PrefabController struct {
	engine.ComponentImpl
}

func (l *PrefabController) OnInit(ctx engine.InitContext) {
	l.ComponentImpl.OnInit(ctx)
	engine.LogDebug("Init")
}

func (l *PrefabController) OnStart() {
	engine.LogDebug("Start")
	l1 := l.SceneObject.GetComponentByName("github.com/cadmean-ru/amphion/engine/builtin.OnClickListener").(*builtin.OnClickListener)
	l1.OnClick = func(event engine.Event) bool {
		engine.LogDebug("Breh")
		return true
	}
}

func (l *PrefabController) GetName() string {
	return engine.NameOfComponent(l)
}

//endregion 2

func prefabScene(e *engine.AmphionEngine) *engine.SceneObject {
	scene := engine.NewSceneObject("test prefab")

	scene.AddComponent(&TestPrefabSceneController{})

	return scene
}
