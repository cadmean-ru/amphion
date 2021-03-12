package cli

type Frontend interface {
	Init()
	Run()
	Reset()
	GetAppData() []byte
	CommencePanic(reason, msg string)
	GetContext() *Context
	SetCallback(handler *CallbackHandler)
}