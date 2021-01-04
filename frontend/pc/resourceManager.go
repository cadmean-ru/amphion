// +build windows linux darwin
// +build !android

package pc

import (
	"github.com/cadmean-ru/amphion/common"
	"github.com/cadmean-ru/amphion/common/a"
	"io/ioutil"
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
	return ioutil.ReadFile("res/" + r.resources[id])
}

func newResourceManager() *ResourceManager {
	return &ResourceManager{
		resources: make(map[a.Int]string),
		idgen:     common.NewIdGenerator(),
	}
}