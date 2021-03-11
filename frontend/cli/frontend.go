package cli

type CallbackHandler interface {
	HandleCallback(code int, data string)
}

type FrontendCLI interface {
	Init()
	Run()
	Reset()
	GetAppData() []byte
	CommencePanic(reason, msg string)
	GetContext() *Context
}