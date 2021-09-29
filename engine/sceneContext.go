package engine

import "github.com/cadmean-ru/amphion/common/a"

// SceneContext struct holds data related to a SceneObject instance.
type SceneContext struct {
	scene             *SceneObject
	Args              a.SiMap
	focusedObject     *SceneObject
	hoveredObject     *SceneObject
	messageDispatcher *MessageDispatcher
	//layouter      *LayoutImpl
}

func makeSceneContext(scene *SceneObject) *SceneContext {
	var args a.SiMap
	if instance.appContext != nil {
		args = instance.appContext.navigationArgs
	}

	return &SceneContext{
		Args:              args,
		scene:             scene,
		messageDispatcher: newMessageDispatcherForScene(scene),
	}
}
