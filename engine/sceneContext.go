package engine

import "github.com/cadmean-ru/amphion/common/a"

// This struct holds data related to a scene.
type SceneContext struct {
	Args a.SiMap
}

func makeSceneContext() *SceneContext {
	return &SceneContext{
		Args: instance.appContext.navigationArgs,
	}
}