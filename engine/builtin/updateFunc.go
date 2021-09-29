package builtin

import "github.com/cadmean-ru/amphion/engine"

type UpdateFunc struct {
	engine.ComponentImpl
	action func(ctx engine.UpdateContext)
}

func (u *UpdateFunc) OnUpdate(ctx engine.UpdateContext) {
	if u.action != nil {
		u.action(ctx)
	}
}

func NewUpdateFunc(action func(ctx engine.UpdateContext)) *UpdateFunc {
	return &UpdateFunc{
		action: action,
	}
}
