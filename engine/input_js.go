// +build js

package engine

import (
	"github.com/cadmean-ru/amphion/common"
	"syscall/js"
)

var getMousePositionJs = js.Global().Get("getMousePosition")

func GetMousePosition() common.IntVector3 {
	mouseJs := getMousePositionJs.Invoke()
	return common.IntVector3{
		X: mouseJs.Get("x").Int(),
		Y: mouseJs.Get("y").Int(),
	}
}
