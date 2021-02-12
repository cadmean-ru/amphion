package ios

import (
	"github.com/cadmean-ru/amphion/common/a"
	"github.com/cadmean-ru/amphion/frontend"
	"github.com/cadmean-ru/amphion/frontend/cli"
)

type ResourceManager struct {
	*frontend.ResourceManagerImpl
	rm cli.ResourceManagerCLI
}

func (r *ResourceManager) ReadFile(id a.ResId) ([]byte, error) {
	return r.rm.ReadFile(r.FullPathOf(id))
}
