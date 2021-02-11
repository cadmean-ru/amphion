package engine

import "github.com/cadmean-ru/amphion/common/a"

// This struct holds data related to a scene.
type SceneContext struct {
	Args          a.SiMap
	focusedObject *SceneObject
	hoveredObject *SceneObject
}

func makeSceneContext() *SceneContext {
	var args a.SiMap
	if instance.appContext != nil {
		args = instance.appContext.navigationArgs
	}

	return &SceneContext{
		Args: args,
	}
}
