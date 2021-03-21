package cli

type FrontendDelegate interface {
	Init()
	Run()
	Reset()
	GetAppData() []byte
	CommencePanic(reason, msg string)
	GetContext() *Context
	SetCallback(handler *CallbackHandler)
	ExecuteOnMainThread(delegate *ExecDelegate)
}