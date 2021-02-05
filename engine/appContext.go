package engine

import (
	"github.com/cadmean-ru/amphion/common/a"
	"github.com/cadmean-ru/amphion/frontend"
)

// App context holds app specific data.
type AppContext struct {
	navigationArgs a.SiMap
}

func makeAppContext(app *frontend.App) *AppContext {
	return &AppContext {

	}
}