package rpc

import (
	"errors"
	"fmt"
)

type Error struct {
	error
	Code int
}

func newError(code int, msg string) Error {
	return Error{
		error: errors.New(fmt.Sprintf("rpc error: %d. %s", code, msg)),
		Code:  code,
	}
}

const (
	ErrorFunctionNotCallable     = -100
	ErrorFunctionNotFound        = -101
	ErrorIncompatibleRPCVersion  = -102
	ErrorInvalidArguments        = -200
	ErrorEncode                  = -300
	ErrorDecode                  = -301
	ErrorCouldNotSendCall        = -400
	ErrorNotSuccessfulStatusCode = -401
	ErrorServer                  = -500
	ErrorAuth                    = -600
)
