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
	orientation    ScreenOrientation
	delegate       AppDelegate
}

func (c *AppContext) Orientation() ScreenOrientation {
	return c.orientation
}

func (c *AppContext) onAppLoaded() {
	if c.delegate == nil {
		return
	}

	c.delegate.OnAppLoaded()
}

func (c *AppContext) onAppStopping() {
	if c.delegate == nil {
		return
	}

	c.delegate.OnAppStopping()
}

func makeAppContext(app *frontend.App) *AppContext {
	return &AppContext{

	}
}
