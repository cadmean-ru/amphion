package commonFrontend

import "github.com/cadmean-ru/amphion/rendering"

type Frontend interface {
	Init()
	Start()
	Stop()
	Reset()
	SetCallback(handler CallbackHandler)
	GetInputManager() InputManager
	GetRenderer() rendering.Renderer
	GetContext() Context
	CommencePanic(reason, msg string)
}

type Callback struct {
	Code int
	Data string
}

type CallbackHandler func()

type InputManager interface {

}