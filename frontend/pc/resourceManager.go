// +build windows linux darwin
// +build !android
// +build !ios

package pc

import (
	"github.com/cadmean-ru/amphion/common"
	"github.com/cadmean-ru/amphion/common/a"
	"io/ioutil"
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
	return ioutil.ReadFile("res/" + r.resources[id])
}

func newResourceManager() *ResourceManager {
	return &ResourceManager{
		resources: make(map[a.ResId]string),
		idgen:     common.NewIdGenerator(),
	}
}