// +build js

package engine

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func loadScene(scene string) (*SceneObject, error) {
	url := fmt.Sprintf("http://%s:%s/scenes/%s", instance.globalContext.host, instance.globalContext.port, scene)
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	s := &SceneObject{}
	err = s.DecodeFromYaml(data)
	if err != nil {
		return nil, err
	}

	return s, nil
}
