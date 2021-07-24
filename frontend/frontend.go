// Package frontend provides interface for platform-specific code.
package frontend

import (
	"github.com/cadmean-ru/amphion/common"
	"github.com/cadmean-ru/amphion/common/a"
	"github.com/cadmean-ru/amphion/common/dispatch"
	"github.com/cadmean-ru/amphion/rendering"
	"github.com/cadmean-ru/amphion/rendering/gpu"
)

//Frontend interface defines functions that should be implemented in an Amphion frontend natively.
//It defines the way engine and frontend will communicate.
type Frontend interface {
	//Init is called when the frontend is created.
	Init()

	//Run is like the main function for the frontend.
	Run()

	SetEngineDispatcher(disp dispatch.MessageDispatcher)

	//GetRenderer should return a fully configured ARenderer. Must be deterministic.
	GetRenderer() *rendering.ARenderer

	//GetContext should return the Context. It is called by the engine whenever the frontend sends CallbackContextChange.
	GetContext() Context

	//GetPlatform should return the current Platform.
	GetPlatform() common.Platform

	//CommencePanic should indicate in a native way, that an error has occurred in either the app or the engine.
	CommencePanic(reason, msg string)

	GetMessageDispatcher() dispatch.MessageDispatcher

	GetWorkDispatcher() dispatch.WorkDispatcher

	//GetResourceManager should return an implementation of ResourceManager interface.
	GetResourceManager() ResourceManager

	//GetApp should returns the app data from app.yaml file.
	GetApp() *App

	//GetLaunchArgs should return arguments.
	GetLaunchArgs() a.SiMap

	//GetGpu should return the gpu implementation for the current frontend.
	GetGpu() gpu.Gpu
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
	CallbackKeyUp         = -112
	CallbackTextInput     = -113
)

//Deprecated: use dispatch.Message instead
//Callback contains a code and string data, that is accepted by the engine through the CallbackHandler passed with SetCallback.
type Callback struct {
	Code int
	Data string
}

//Deprecated: use dispatch.Message instead
//NewCallback creates a new Callback instance with the specified code and string data.
func NewCallback(code int, data string) Callback {
	return Callback{
		Code: code,
		Data: data,
	}
}

//Deprecated: use dispatch.MessageHandler instead
//CallbackHandler is a function that is provided by the engine through SetCallback and
//called to pass data from frontend to the engine.
type CallbackHandler func(callback Callback)

//ResourceManager defines interface, that should be implemented in the frontend to provide functionality of
//working with files in res folder.
type ResourceManager interface {
	RegisterResource(path string)
	IdOf(path string) a.ResId
	PathOf(id a.ResId) string
	FullPathOf(id a.ResId) string
	ReadFile(id a.ResId) ([]byte, error)
}