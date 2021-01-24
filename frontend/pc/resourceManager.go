// +build windows linux darwin
// +build !android

package pc

import (
	"github.com/cadmean-ru/amphion/common"
	"io/ioutil"
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

func (r *ResourceManager) FullPathOf(id int) string {
	return "res/" + r.resources[id]
}

func (r *ResourceManager) ReadFile(id int) ([]byte, error) {
	return ioutil.ReadFile("res/" + r.resources[id])
}

func newResourceManager() *ResourceManager {
	return &ResourceManager{
		resources: make(map[int]string),
		idgen:     common.NewIdGenerator(),
	}
}