// Package frontend provides interface for platform-specific code.
package frontend

import (
	"github.com/cadmean-ru/amphion/common"
	"github.com/cadmean-ru/amphion/common/a"
	"github.com/cadmean-ru/amphion/common/dispatch"
	"github.com/cadmean-ru/amphion/rendering"
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

	//GetApp should return the app data from app.yaml file.
	GetApp() *App

	//GetLaunchArgs should return arguments.
	GetLaunchArgs() a.SiMap
}

const (
	CallbackContextChange = -100-iota
	CallbackMouseDown
	CallbackKeyDown
	CallbackMouseUp
	CallbackAppHide
	CallbackAppShow
	CallbackMouseMove
	CallbackMouseScroll
	CallbackTouchDown
	CallbackTouchUp
	CallbackTouchMove
	CallbackReady
	CallbackKeyUp
	CallbackTextInput
	CallbackOrientationChange
	CallbackStop
)

//ResourceManager defines interface, that should be implemented in the frontend to provide functionality of
//working with files in res folder.
type ResourceManager interface {
	RegisterResource(path string)
	IdOf(path string) a.ResId
	PathOf(id a.ResId) string
	FullPathOf(id a.ResId) string
	ReadFile(id a.ResId) ([]byte, error)
}
