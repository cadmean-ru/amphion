package atest

import (
	"github.com/cadmean-ru/amphion/common/dispatch"
	"github.com/cadmean-ru/amphion/frontend"
)

var instance TestingFrontend

type TestingFrontend interface {
	frontend.Frontend
	SimulateCallback(message *dispatch.Message)
}
