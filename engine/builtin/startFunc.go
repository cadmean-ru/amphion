package builtin

import "github.com/cadmean-ru/amphion/engine"

type StartFunc struct {
	engine.ComponentImpl
	action func()
}

func (s *StartFunc) OnStart() {
	if s.action != nil {
		s.action()
	}
}

func NewStartFunc(action func()) *StartFunc {
	return &StartFunc{
		action: action,
	}
}
