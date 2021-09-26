package engine

import (
	"errors"
	"github.com/cadmean-ru/amphion/common/a"
	"github.com/cadmean-ru/amphion/common/dispatch"
	"github.com/cadmean-ru/amphion/frontend"
	"strings"
)

// Opens the SceneObject corresponding to the specified path.
// Path "/" corresponds to the app's main SceneObject.
// Other paths correspond to the SceneObject path in the SceneObject folder.
// For example "res/scenes/test.SceneObject" corresponds to "/test", "res/scenes/hello/test.SceneObject" to "/hello/test".
// Args can be passed, that can be later read from the SceneObject.
func Navigate(path string, args a.SiMap) (err error) {
	if instance.currentApp == nil {
		err = errors.New("cannot navigate without loaded app")
		return
	}

	var scenePath string

	if path == "/" {
		scenePath = "scenes/" + instance.currentApp.MainScene + ".SceneObject"
	} else {
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

		scenePath = "scenes/" + strings.Join(pathTokens, "/") + ".SceneObject"
	}


	sceneId := instance.GetResourceManager().IdOf(scenePath)

	if sceneId == -1 {
		err = errors.New("SceneObject not found")
		return
	}

	instance.appContext.navigationArgs = args
	instance.LoadScene(sceneId, true)

	instance.front.GetMessageDispatcher().SendMessage(dispatch.NewMessageWithStringData(frontend.MessageNavigate, path))

	return
}