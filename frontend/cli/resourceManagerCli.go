package cli

import (
	"github.com/cadmean-ru/amphion/common/a"
	"github.com/cadmean-ru/amphion/frontend"
)

type ResourceManagerDelegate interface {
	ReadFile(path string) ([]byte, error)
}

type ResourceManagerImpl struct {
	*frontend.ResourceManagerImpl
	delegate ResourceManagerDelegate
}

func (r *ResourceManagerImpl) ReadFile(id a.ResId) ([]byte, error) {
	return r.delegate.ReadFile(r.FullPathOf(id))
}

func NewResourceManagerImpl(delegate ResourceManagerDelegate) *ResourceManagerImpl {
	return &ResourceManagerImpl{
		ResourceManagerImpl: frontend.NewResourceManagerImpl(),
		delegate:            delegate,
	}
}