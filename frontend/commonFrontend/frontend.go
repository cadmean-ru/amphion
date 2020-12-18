package commonFrontend

import (
	"github.com/cadmean-ru/amphion/common"
	"github.com/cadmean-ru/amphion/rendering"
)

type Frontend interface {
	Init()
	Run()
	Reset()
	SetCallback(handler CallbackHandler)
	GetInputManager() InputManager
	GetRenderer() rendering.Renderer
	GetContext() Context
	GetPlatform() common.Platform
	CommencePanic(reason, msg string)
	ReceiveMessage(message Message)
}

type Callback struct {
	Code int
	Data string
}

type CallbackHandler func()

type InputManager interface {

}