// +build windows linux

package frontend

import (
	"github.com/cadmean-ru/amphion/frontend/commonFrontend"
	"github.com/cadmean-ru/amphion/frontend/pc"
)

func GetFrontend() commonFrontend.Frontend {
	return pc.NewFrontend()
}