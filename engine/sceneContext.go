package engine

import "github.com/cadmean-ru/amphion/common/a"

// SceneContext struct holds data related to a SceneObject instance.
type SceneContext struct {
	Args          a.SiMap
	focusedObject *SceneObject
	hoveredObject *SceneObject
	//layouter      *LayoutImpl
}

func makeSceneContext() *SceneContext {
	var args a.SiMap
	if instance.appContext != nil {
		args = instance.appContext.navigationArgs
	}

	return &SceneContext{
		Args:     args,
		//layouter: newLayouter(),
	}
}
