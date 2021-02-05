package engine

import (
	"errors"
	"github.com/cadmean-ru/amphion/common/a"
	"strings"
)

func Navigate(path string, args a.SiMap) (err error) {
	pathTokens := strings.Split(path, "/")

	if len(pathTokens) == 0 {
		err = errors.New("invalid path")
		return
	}

	if pathTokens[0] == "" {
		pathTokens = pathTokens[1:]

		if len(pathTokens) == 0 {
			err = errors.New("invalid path")
			return
		}
	}

	for _, p := range pathTokens {
		if p == "" {
			err = errors.New("invalid path")
			return
		}
	}

	scenePath := "scenes/" + strings.Join(pathTokens, "/") + ".scene"

	sceneId := instance.GetResourceManager().IdOf(scenePath)

	if sceneId == -1 {
		err = errors.New("scene not found")
		return
	}

	instance.appContext.navigationArgs = args
	instance.LoadScene(sceneId, true)

	return
}