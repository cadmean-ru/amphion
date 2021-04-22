package rendering

import (
	"fmt"
	"github.com/cadmean-ru/amphion/common"
	"github.com/cadmean-ru/amphion/common/dispatch"
	"sort"
)

type ManagementMode byte

const (
	EngineManaged ManagementMode = iota
	FrontendManaged
)

type ARenderer struct {
	primitives         map[int]*PrimitiveContainer
	idgen              *common.IdGenerator
	delegate           RendererDelegate
	primitiveDelegates map[byte]PrimitiveRendererDelegate
	layers             []*Layer
	management         ManagementMode
	renderDispatcher   dispatch.WorkDispatcher
	preparedCallback   func()
	prepared           bool
}

// PC - main, android - rendering thread
func (r *ARenderer) Prepare() {
	r.layers[0] = newLayer()

	r.renderDispatcher.Execute(dispatch.NewWorkItemFunc(func() {
		r.delegate.OnPrepare()

		for _, delegate := range r.primitiveDelegates {
			delegate.OnStart()
		}

		r.prepared = true
		if r.preparedCallback != nil {
			r.preparedCallback()
		}
	}))
}

// Any thread, update goroutine
func (r *ARenderer) AddPrimitive() int {
	id := r.idgen.NextId()
	r.primitives[id] = &PrimitiveContainer{}
	r.layers[0].addPrimitiveId(id)
	return id
}

// Any thread, update goroutine
func (r *ARenderer) SetPrimitive(id int, primitive IPrimitive, shouldRerender bool) {
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


	r.renderDispatcher.Execute(dispatch.NewWorkItemFunc(func() {
		if delegate, ok := r.primitiveDelegates[primitive.GetType()]; ok {
			ctx := r.makePrimitiveRenderingContext(r.primitives[id])
			delegate.OnSetPrimitive(ctx)
			r.primitives[id].state = ctx.State
		}
	}))
}

// Any thread, update goroutine
func (r *ARenderer) RemovePrimitive(id int) {
	if p, ok := r.primitives[id]; ok {
		if d, ok := r.primitiveDelegates[p.primitive.GetType()]; ok {
			r.renderDispatcher.Execute(dispatch.NewWorkItemFunc(func() {
				d.OnRemovePrimitive(r.makePrimitiveRenderingContext(p))
			}))
		}
		for _, l := range r.layers {
			l.removePrimitiveId(id)
		}
		delete(r.primitives, id)
	}
}

// Main or rendering thread
func (r *ARenderer) PerformRendering() {
	r.renderDispatcher.Execute(dispatch.NewWorkItemFunc(r.delegate.OnPerformRenderingStart))

	if r.management != EngineManaged {
		return
	}

	r.renderDispatcher.Execute(dispatch.NewWorkItemFunc(r.renderingPerformer))

	r.renderDispatcher.Execute(dispatch.NewWorkItemFunc(r.delegate.OnPerformRenderingEnd))
}

// Update routine
func (r *ARenderer) Clear() {
	r.primitives = make(map[int]*PrimitiveContainer)
	r.idgen = common.NewIdGenerator()
	r.delegate.OnClear()
}

func (r *ARenderer) Stop() {
	for _, delegate := range r.primitiveDelegates {
		r.renderDispatcher.Execute(dispatch.NewWorkItemFunc(delegate.OnStop))
	}
	r.renderDispatcher.Execute(dispatch.NewWorkItemFunc(r.delegate.OnStop))
}

func (r *ARenderer) makePrimitiveRenderingContext(container *PrimitiveContainer) *PrimitiveRenderingContext {
	return &PrimitiveRenderingContext{
		Renderer:      r,
		Primitive:     container.primitive,
		PrimitiveKind: container.primitive.GetType(),
		PrimitiveId:   container.id,
		State:         container.state,
		Redraw:        container.redraw,
	}
}

func (r *ARenderer) RegisterPrimitiveRendererDelegate(primitiveKind byte, delegate PrimitiveRendererDelegate) {
	if delegate == nil {
		return
	}
	r.primitiveDelegates[primitiveKind] = delegate
}

func (r *ARenderer) GetLayer(index int) *Layer {
	if index < 0 || index > len(r.layers) {
		return nil
	}
	return r.layers[index]
}

func (r *ARenderer) AddLayer() {
	r.layers = append(r.layers, newLayer())
}

func (r *ARenderer) RemoveLayer(index int) {
	if len(r.layers) <= 1 {
		return
	}

	var i2 int
	newLayers := make([]*Layer, len(r.layers) - 1)
	for i, l := range r.layers {
		if i == index {
			continue
		}

		newLayers[i2] = l
		i2++
	}
}

func (r *ARenderer) SetPrimitiveLayer(primitiveId, layerIndex int) {
	newLayer := r.GetLayer(layerIndex)
	if newLayer == nil {
		return
	}

	var oldLayer *Layer
	for _, l := range r.layers {
		for _, p := range l.primitives {
			if p == primitiveId {
				oldLayer = l
				break
			}
		}
		if oldLayer != nil {
			break
		}
	}

	if oldLayer == nil {
		return
	}

	oldLayer.removePrimitiveId(primitiveId)
	newLayer.addPrimitiveId(primitiveId)
}

func (r *ARenderer) SetManagementMode(mode ManagementMode) {
	r.management = mode
}

func (r *ARenderer) GetRenderingPerformer() func() {
	if r.management == EngineManaged {
		return nil
	}

	return r.renderingPerformer
}

func (r *ARenderer) SetPreparedCallback(f func()) {
	if r.prepared {
		f()
		return
	}
	r.preparedCallback = f
}

func (r *ARenderer) renderingPerformer() {
	for _, layer := range r.layers {
		count := 0
		for _, p := range layer.primitives {
			if r.primitives[p].primitive == nil {
				continue
			}

			count++
		}

		list := make([]*PrimitiveContainer, count)

		var i = 0
		for _, p := range layer.primitives {
			if r.primitives[p].primitive == nil {
				continue
			}

			list[i] = r.primitives[p]
			i++
		}

		sort.Slice(list, func(i, j int) bool {
			z1 := list[i].primitive.GetTransform().Position.Z
			z2 := list[j].primitive.GetTransform().Position.Z
			return z1 < z2
		})

		for _, p := range list {
			if d, ok := r.primitiveDelegates[p.primitive.GetType()]; ok {
				if p.redraw {
					tr := p.primitive.GetTransform()
					tr.Position = tr.Position.Add(layer.translation)
					p.primitive.SetTransform(tr)
				}

				ctx := r.makePrimitiveRenderingContext(p)
				d.OnRender(ctx)
				p.state = ctx.State
			}

			p.redraw = false
		}
	}
}

func NewARenderer(delegate RendererDelegate, renderDispatcher dispatch.WorkDispatcher) *ARenderer {
	return &ARenderer{
		delegate:           delegate,
		renderDispatcher:   renderDispatcher,
		idgen:              common.NewIdGenerator(),
		primitives:         make(map[int]*PrimitiveContainer),
		primitiveDelegates: make(map[byte]PrimitiveRendererDelegate),
		layers:             make([]*Layer, 1),
	}
}
