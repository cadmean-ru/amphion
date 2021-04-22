package cli

import "github.com/cadmean-ru/amphion/common/dispatch"

type FrontendDelegate interface {
	Init()
	Run()
	Reset()
	GetAppData() []byte
	CommencePanic(reason, msg string)
	GetContext() *Context
	SetCallbackDispatcher(dispatcher dispatch.MessageDispatcher)
	GetMainThreadDispatcher() dispatch.WorkDispatcher
	GetRenderingThreadDispatcher() dispatch.WorkDispatcher
}