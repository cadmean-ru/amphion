package frontend

import (
	"github.com/cadmean-ru/amphion/common"
	"github.com/cadmean-ru/amphion/common/a"
)

type ResourceManagerImpl struct {
	resources map[a.ResId]string
	idgen     *common.IdGenerator
}

func (r *ResourceManagerImpl) RegisterResource(path string) {
	r.resources[a.ResId(r.idgen.NextId())] = path
}

func (r *ResourceManagerImpl) IdOf(path string) a.ResId {
	for id, p := range r.resources {
		if p == path {
			return id
		}
	}

	return -1
}

func (r *ResourceManagerImpl) PathOf(id a.ResId) string {
	return r.resources[id]
}

func (r *ResourceManagerImpl) FullPathOf(id a.ResId) string {
	return "res/" + r.resources[id]
}

func NewResourceManagerImpl() *ResourceManagerImpl {
	return &ResourceManagerImpl{
		resources: map[a.ResId]string{},
		idgen:     common.NewIdGenerator(),
	}
}