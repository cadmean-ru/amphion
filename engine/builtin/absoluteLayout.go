package builtin

import "github.com/cadmean-ru/amphion/engine"

type AbsoluteLayout struct {
	engine.LayoutImpl
}

func NewAbsoluteLayout() *AbsoluteLayout {
	return &AbsoluteLayout{}
}
