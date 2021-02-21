package rendering

import (
	"fmt"
	"github.com/cadmean-ru/amphion/common"
	"sort"
)

type RendererImpl struct {
	primitives         map[int]*PrimitiveContainer
	idgen              *common.IdGenerator
	delegate           RendererDelegate
	primitiveDelegates map[byte]PrimitiveRendererDelegate
}

func (r *RendererImpl) Prepare() {
	r.delegate.OnPrepare()
	for _, delegate := range r.primitiveDelegates {
		delegate.OnStart()
	}
}

func (r *RendererImpl) AddPrimitive() int {
	id := r.idgen.NextId()
	r.primitives[id] = &PrimitiveContainer{}
	return id
}

func (r *RendererImpl) SetPrimitive(id int, primitive IPrimitive, shouldRerender bool) {
	if !shouldRerender {
		return
	}

	if _, ok := r.primitives[id]; ok {
		if r.primitives[id].primitive != nil && primitive.GetType() != r.primitives[id].primitive.GetType() {
			fmt.Printf("Cannot change primitive type for id: %d\n", id)
			return
		}
		r.primitives[id].primitive = primitive
		r.primitives[id].redraw = true
	} else {
		fmt.Printf("Warning! Primitive with id %d was not found.\n", id)
	}

	if delegate, ok := r.primitiveDelegates[primitive.GetType()]; ok {
		ctx := r.makePrimitiveRenderingContext(r.primitives[id])
		delegate.OnSetPrimitive(ctx)
		r.primitives[id].state = ctx.State
	}
}

func (r *RendererImpl) RemovePrimitive(id int) {
	if p, ok := r.primitives[id]; ok {
		if d, ok := r.primitiveDelegates[p.primitive.GetType()]; ok {
			d.OnRemovePrimitive(r.makePrimitiveRenderingContext(p))
		}
		delete(r.primitives, id)
	}
}

func (r *RendererImpl) PerformRendering() {
	r.delegate.OnPerformRenderingStart()

	count := 0
	for _, p := range r.primitives {
		if p.primitive == nil {
			continue
		}

		count++
	}

	list := make([]*PrimitiveContainer, count)

	var i = 0
	for _, p := range r.primitives {
		if p.primitive == nil {
			continue
		}
		list[i] = p
		i++
	}

	sort.Slice(list, func(i, j int) bool {
		z1 := list[i].primitive.GetTransform().Position.Z
		z2 := list[j].primitive.GetTransform().Position.Z
		return z1 < z2
	})

	for _, p := range list {
		if d, ok := r.primitiveDelegates[p.primitive.GetType()]; ok {
			ctx := r.makePrimitiveRenderingContext(p)
			d.OnRender(ctx)
			p.state = ctx.State
		}

		p.redraw = false
	}

	r.delegate.OnPerformRenderingEnd()
}

func (r *RendererImpl) Clear() {
	r.primitives = make(map[int]*PrimitiveContainer)
	r.idgen = common.NewIdGenerator()
	r.delegate.OnClear()
}

func (r *RendererImpl) Stop() {
	for _, delegate := range r.primitiveDelegates {
		delegate.OnStop()
	}
	r.delegate.OnStop()
}

func (r *RendererImpl) makePrimitiveRenderingContext(container *PrimitiveContainer) *PrimitiveRenderingContext {
	return &PrimitiveRenderingContext{
		Renderer:      r,
		Primitive:     container.primitive,
		PrimitiveKind: container.primitive.GetType(),
		State:         container.state,
		Redraw:        container.redraw,
	}
}

func (r *RendererImpl) RegisterPrimitiveRendererDelegate(primitiveKind byte, delegate PrimitiveRendererDelegate) {
	if delegate == nil {
		return
	}
	r.primitiveDelegates[primitiveKind] = delegate
}

func NewRendererImpl(delegate RendererDelegate) *RendererImpl {
	return &RendererImpl{
		delegate:           delegate,
		idgen:              common.NewIdGenerator(),
		primitives:         make(map[int]*PrimitiveContainer),
		primitiveDelegates: make(map[byte]PrimitiveRendererDelegate),
	}
}
