package frontend

import (
	"github.com/cadmean-ru/amphion/common"
	"github.com/cadmean-ru/amphion/common/a"
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
	GetResourceManager() ResourceManager
	GetApp() *App
	SetWindowTitle(title string)
}

const (
	CallbackContextChange = -100
	CallbackMouseDown     = -101
	CallbackKeyDown       = -102
	CallbackMouseUp       = -103
	CallbackAppHide       = -104
	CallbackAppShow       = -105
	CallbackMouseMove     = -106
)

type Callback struct {
	Code int
	Data string
}

func NewCallback(code int, data string) Callback {
	return Callback{
		Code: code,
		Data: data,
	}
}

type CallbackHandler func(callback Callback)

type InputManager interface {
	GetMousePosition() a.IntVector2
}

type ResourceManager interface {
	RegisterResource(path string)
	IdOf(path string) int
	PathOf(id int) string
	FullPathOf(id int) string
	ReadFile(id int) ([]byte, error)
}