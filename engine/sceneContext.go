package engine

import "github.com/cadmean-ru/amphion/common/a"

// This struct holds data related to a scene.
type SceneContext struct {
	Args          a.SiMap
	focusedObject *SceneObject
	hoveredObject *SceneObject
}

func makeSceneContext() *SceneContext {
	return &SceneContext{
		Args: instance.appContext.navigationArgs,
	}
}
