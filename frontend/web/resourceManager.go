// +build js

package web

import (
	"github.com/cadmean-ru/amphion/common"
	"github.com/cadmean-ru/amphion/common/a"
	"io/ioutil"
	"net/http"
)

type ResourceManager struct {
	resources map[a.Int]string
	idgen     *common.IdGenerator
}

func (r *ResourceManager) RegisterResource(path string) {
	r.resources[a.Int(r.idgen.NextId())] = path
}

func (r *ResourceManager) IdOf(path string) a.Int {
	for id, p := range r.resources {
		if p == path {
			return id
		}
	}

	return -1
}

func (r *ResourceManager) PathOf(id a.Int) string {
	return r.resources[id]
}

func (r *ResourceManager) ReadFile(id a.Int) ([]byte, error) {
	resp, err := http.Get("http://localhost:8080/res/" + r.resources[id])
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
		resources: make(map[a.Int]string),
		idgen:     common.NewIdGenerator(),
	}
}