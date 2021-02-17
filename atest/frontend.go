// +build linux windows darwin
// +build !android
// +build !ios

package atest

import (
	"fmt"
	"github.com/cadmean-ru/amphion/common"
	"github.com/cadmean-ru/amphion/frontend"
	"github.com/cadmean-ru/amphion/frontend/pc"
	"time"
)

var instance *Frontend

// The testing frontend.
type Frontend struct {
	*pc.Frontend
	handler frontend.CallbackHandler
	clicked bool
}

func (f *Frontend) Init() {
	f.Frontend.Init()
	instance = f
	fmt.Println("This is a testing frontend")
}

func (f *Frontend) GetPlatform() common.Platform {
	return common.PlatformFromString("testing")
}

func (f *Frontend) SetCallback(handler frontend.CallbackHandler) {
	f.handler = handler
	f.Frontend.SetCallback(handler)
}

func (f *Frontend) simulateCallback(code int, data string) {
	f.handler(frontend.NewCallback(code, data))
}

func (f *Frontend) simulateClick(x, y int) {
	if f.clicked {
		return
	}

	data := fmt.Sprintf("%d;%d", x, y)
	go func() {
		f.simulateCallback(frontend.CallbackMouseDown, data)
		f.clicked = true
		time.Sleep(100)
		f.simulateCallback(frontend.CallbackMouseUp, data)
		f.clicked = false
	}()
}

func newTestingFrontend() *Frontend {
	return &Frontend{
		Frontend: pc.NewFrontend(),
	}
}