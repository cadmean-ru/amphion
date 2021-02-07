// +build js

package web

import (
	"github.com/cadmean-ru/amphion/common"
	"github.com/cadmean-ru/amphion/common/a"
	"github.com/cadmean-ru/amphion/engine"
	"io/ioutil"
	"net/http"
)

type ResourceManager struct {
	resources map[a.ResId]string
	idgen     *common.IdGenerator
}

func (r *ResourceManager) RegisterResource(path string) {
	r.resources[a.ResId(r.idgen.NextId())] = path
}

func (r *ResourceManager) IdOf(path string) a.ResId {
	for id, p := range r.resources {
		if p == path {
			return id
		}
	}

	return -1
}

func (r *ResourceManager) PathOf(id a.ResId) string {
	return r.resources[id]
}

func (r *ResourceManager) FullPathOf(id a.ResId) string {
	return "res/" + r.resources[id]
}

func (r *ResourceManager) ReadFile(id a.ResId) ([]byte, error) {
	resp, err := http.Get("http://" + engine.GetInstance().GetCurrentApp().PublicUrl + "/res/" + r.resources[id])
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func newResourceManager() *ResourceManager {
	return &ResourceManager{
		resources: make(map[a.ResId]string),
		idgen:     common.NewIdGenerator(),
	}
}