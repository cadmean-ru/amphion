// Package frontend provides interface for platform-specific code.
package frontend

import (
	"github.com/cadmean-ru/amphion/common"
	"github.com/cadmean-ru/amphion/common/a"
	"github.com/cadmean-ru/amphion/rendering"
)

type Frontend interface {
	Init()
	Run()
	//Deprecated
	Reset()
	SetCallback(handler CallbackHandler)
	//Deprecated
	GetInputManager() InputManager
	GetRenderer() *rendering.ARenderer
	GetContext() Context
	//Deprecated
	GetPlatform() common.Platform
	CommencePanic(reason, msg string)
	ReceiveMessage(message Message)
	GetResourceManager() ResourceManager
	GetApp() *App
	GetLaunchArgs() a.SiMap
}

const (
	CallbackContextChange = -100
	CallbackMouseDown     = -101
	CallbackKeyDown       = -102
	CallbackMouseUp       = -103
	CallbackAppHide       = -104
	CallbackAppShow       = -105
	CallbackMouseMove     = -106
	CallbackMouseScroll   = -107
	CallbackTouchDown     = -108
	CallbackTouchUp       = -109
	CallbackTouchMove     = -110
	CallbackReady         = -111
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

// Deprecated
type InputManager interface {
	GetMousePosition() a.IntVector2
}

type ResourceManager interface {
	RegisterResource(path string)
	IdOf(path string) a.ResId
	PathOf(id a.ResId) string
	FullPathOf(id a.ResId) string
	ReadFile(id a.ResId) ([]byte, error)
}