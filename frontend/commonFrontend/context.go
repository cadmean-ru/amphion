package commonFrontend

import "github.com/cadmean-ru/amphion/common"

type Context struct {
	DeviceInfo common.DeviceInfo
	ScreenInfo common.ScreenInfo
	Domain     string
	Host       string
	Port       string
}