// +build !js,!pc,!linux

package rendering

import "github.com/cadmean-ru/amphion/common"

type EmptyRenderer struct {
	idgen *common.IdGenerator
}

func (r *EmptyRenderer) AddPrimitive() int64 { return r.idgen.NextId() }

func (r *EmptyRenderer) Prepare() {}

func (r *EmptyRenderer) SetPrimitive(_ int64, _ interface{}, _ bool) {}

func (r *EmptyRenderer) RemovePrimitive(_ int64) {}

func (r *EmptyRenderer) PerformRendering() {}

func (r *EmptyRenderer) Clear() {}

func (r *EmptyRenderer) Stop() {}

func NewRenderer() Renderer {
	return &EmptyRenderer{
		idgen: common.NewIdGenerator(),
	}
}