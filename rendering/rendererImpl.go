package rendering

import (
	"fmt"
	"github.com/cadmean-ru/amphion/common"
	"sort"
)

type ARenderer struct {
	primitives         map[int]*PrimitiveContainer
	idgen              *common.IdGenerator
	delegate           RendererDelegate
	primitiveDelegates map[byte]PrimitiveRendererDelegate
	layers             []*Layer
}

func (r *ARenderer) Prepare() {
	r.layers[0] = newLayer()

	r.delegate.OnPrepare()
	for _, delegate := range r.primitiveDelegates {
		delegate.OnStart()
	}
}

func (r *ARenderer) AddPrimitive() int {
	id := r.idgen.NextId()
	r.primitives[id] = &PrimitiveContainer{}
	r.layers[0].addPrimitiveId(id)
	return id
}

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

	if delegate, ok := r.primitiveDelegates[primitive.GetType()]; ok {
		ctx := r.makePrimitiveRenderingContext(r.primitives[id])
		delegate.OnSetPrimitive(ctx)
		r.primitives[id].state = ctx.State
	}
}

func (r *ARenderer) RemovePrimitive(id int) {
	if p, ok := r.primitives[id]; ok {
		if d, ok := r.primitiveDelegates[p.primitive.GetType()]; ok {
			d.OnRemovePrimitive(r.makePrimitiveRenderingContext(p))
		}
		delete(r.primitives, id)
	}
}

func (r *ARenderer) PerformRendering() {
	r.delegate.OnPerformRenderingStart()

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

	r.delegate.OnPerformRenderingEnd()
}

func (r *ARenderer) Clear() {
	r.primitives = make(map[int]*PrimitiveContainer)
	r.idgen = common.NewIdGenerator()
	r.delegate.OnClear()
}

func (r *ARenderer) Stop() {
	for _, delegate := range r.primitiveDelegates {
		delegate.OnStop()
	}
	r.delegate.OnStop()
}

func (r *ARenderer) makePrimitiveRenderingContext(container *PrimitiveContainer) *PrimitiveRenderingContext {
	return &PrimitiveRenderingContext{
		Renderer:      r,
		Primitive:     container.primitive,
		PrimitiveKind: container.primitive.GetType(),
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

func NewARenderer(delegate RendererDelegate) *ARenderer {
	return &ARenderer{
		delegate:           delegate,
		idgen:              common.NewIdGenerator(),
		primitives:         make(map[int]*PrimitiveContainer),
		primitiveDelegates: make(map[byte]PrimitiveRendererDelegate),
		layers:             make([]*Layer, 1),
	}
}
