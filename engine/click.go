package engine

import "github.com/cadmean-ru/amphion/common/a"

// Represents the data of mouse related event.
type MouseEventData struct {
	MousePosition a.IntVector2
	SceneObject   *SceneObject
}
