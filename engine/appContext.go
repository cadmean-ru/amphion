package engine

import (
	"github.com/cadmean-ru/amphion/common/a"
	"github.com/cadmean-ru/amphion/frontend"
)

type ScreenOrientation int

const (
	OrientationUnknown ScreenOrientation = iota
	OrientationStraight
	OrientationReverse
	OrientationLeft
	OrientationRight
)

// AppContext holds app specific data.
type AppContext struct {
	navigationArgs a.SiMap
	orientation ScreenOrientation
}

func (a AppContext) Orientation() ScreenOrientation {
	return a.orientation
}

func makeAppContext(app *frontend.App) *AppContext {
	return &AppContext {

	}
}