// This package is a wrapper around github.com/cadmean-ru/goRPCKit library adapted for Amphion tasks system.
package rpc

import (
	"github.com/cadmean-ru/amphion/engine"
	rpckit "github.com/cadmean-ru/goRPCKit/rpc"
)

var instance *rpckit.Client

// Initialize creates a new instance of rpc client with specified url.
func Initialize(url string) {
	if instance != nil {
		return
	}

	instance = rpckit.NewClient(url)
}

type FunctionCallBuilder struct {
	fName     string
	onSuccess func(res interface{})
	onError   func(err error)
}

// Specifies callback to be called when the RPC call finishes successfully.
func (f *FunctionCallBuilder) Then(onSuccess func(res interface{})) *FunctionCallBuilder {
	f.onSuccess = onSuccess
	return f
}

// Specifies callback to be called when the RPC call finishes with an error.
func (f *FunctionCallBuilder) Err(onError func(err error)) *FunctionCallBuilder {
	f.onError = onError
	return f
}

// Creates and runs task, calling the RPC function with given arguments.
func (f *FunctionCallBuilder) Call(args ...interface{}) {
	engine.RunTask(engine.NewTaskBuilder().Run(func() (interface{}, error) {
		return instance.F(f.fName).Call(args...)
	}).Then(func(res interface{}) {
		if f.onSuccess != nil {
			f.onSuccess(res)
		}
	}).Err(func(err error) {
		if f.onError != nil {
			f.onError(err)
		}
	}).Build())
}

// Creates a new call builder to call an RPC function with the specified name.
func F(fName string) *FunctionCallBuilder {
	return &FunctionCallBuilder{
		fName:     fName,
	}
}
