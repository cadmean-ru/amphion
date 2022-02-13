// Package rpc is a wrapper around github.com/cadmean-ru/goRPCKit library adapted for Amphion tasks system.
package rpc

import (
	"fmt"
	"github.com/cadmean-ru/amphion/engine"
	"github.com/cadmean-ru/go-rpc/rpc"
)

var instance *rpc.Client

// Initialize creates a new instance of rpc client with specified url.
func Initialize(url string) {
	if instance != nil {
		return
	}

	instance = rpc.NewClient(url, rpc.DefaultConfiguration())
}

type FunctionCallBuilder struct {
	fName     string
	onSuccess func(res interface{})
	onError   func(err error)
	finally   func()
}

//Then specifies callback to be called when the RPC call finishes successfully.
func (f *FunctionCallBuilder) Then(onSuccess func(res interface{})) *FunctionCallBuilder {
	f.onSuccess = onSuccess
	return f
}

//Err specifies callback to be called when the RPC call finishes with an error.
func (f *FunctionCallBuilder) Err(onError func(err error)) *FunctionCallBuilder {
	f.onError = onError
	return f
}

//Finally specifies callback to be called when the RPC call is completed regardless of the result.
func (f *FunctionCallBuilder) Finally(finally func()) *FunctionCallBuilder {
	f.finally = finally
	return f
}

//Call creates and runs task, calling the RPC function with given arguments.
func (f *FunctionCallBuilder) Call(args ...interface{}) {
	engine.RunTask(engine.NewTaskBuilder().Run(func() (interface{}, error) {
		return instance.Func(f.fName).Call(args...)
	}).Then(func(res interface{}) {
		output := res.(*rpc.FunctionOutput)
		if output.Error == 0 && f.onSuccess != nil {
			f.onSuccess(output.Result)
		} else if f.onError != nil {
			f.onError(rpc.NewError(output.Error, fmt.Sprintf("function exited with code %d", output.Error)))
		}
	}).Err(func(err error) {
		if f.onError != nil {
			f.onError(err)
		}
	}).Finally(f.finally).Build())
}

//Func creates a new call builder to call an RPC function with the specified name.
func Func(fName string) *FunctionCallBuilder {
	return &FunctionCallBuilder{
		fName: fName,
	}
}
