package widgets

import (
	"github.com/cadmean-ru/amphion/common/a"
	"github.com/cadmean-ru/amphion/engine"
)

type Options struct {
	Position a.Vector3
	Pivot    a.Vector3
	Size     a.Vector3
	Rotation a.Vector3
}

type ContainerOptions struct {
	Options
	Children []*engine.SceneObject
}

type Widget[T engine.Component] struct {
	SceneObject *engine.SceneObject
	Component   T
}
