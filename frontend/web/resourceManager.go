// +build js

package web

import (
	"github.com/cadmean-ru/amphion/common"
	"io/ioutil"
	"net/http"
)

type ResourceManager struct {
	resources map[int]string
	idgen     *common.IdGenerator
}

func (r *ResourceManager) RegisterResource(path string) {
	r.resources[r.idgen.NextId()] = path
}

func (r *ResourceManager) IdOf(path string) int {
	for id, p := range r.resources {
		if p == path {
			return id
		}
	}

	return -1
}

func (r *ResourceManager) PathOf(id int) string {
	return r.resources[id]
}

func (r *ResourceManager) ReadFile(id int) ([]byte, error) {
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
		resources: make(map[int]string),
		idgen:     common.NewIdGenerator(),
	}
}