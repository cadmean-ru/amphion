// +build !js

package engine

import "github.com/cadmean-ru/amphion/common"

func GetMousePosition() common.IntVector3 {
	return common.ZeroVector().Round()
}
