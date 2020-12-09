package builtin

import (
	"github.com/cadmean-ru/amphion/engine"
	"github.com/cadmean-ru/amphion/rendering"
)

func transformToRenderingTransform(t engine.Transform) rendering.Transform {
	rt := rendering.NewTransform()

	rt.Position = t.GetGlobalTopLeftPosition().Round()
	rt.Size = t.Size.Round()

	return rt
}
