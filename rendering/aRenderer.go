package rendering

import (
	"github.com/cadmean-ru/amphion/common"
	"github.com/cadmean-ru/amphion/common/dispatch"
)

type ManagementMode byte

const (
	EngineManaged ManagementMode = iota
	FrontendManaged
)

type ARenderer struct {
	primitives         map[int]*PrimitiveNode
	idgen              *common.IdGenerator
	delegate           RendererDelegate
	primitiveDelegates map[byte]PrimitiveRendererDelegate
	layers             []*Layer
	management         ManagementMode
	renderDispatcher   dispatch.WorkDispatcher
	preparedCallback   func()
	prepared           bool
	root               *Node
	clipStack          *ClipStack2D
}

// PC - main, android - rendering thread
func (r *ARenderer) Prepare() {
	r.layers[0] = newLayer(0)

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

func (r *ARenderer) SetRoot(n *Node) {
	r.root = n
}

func (r *ARenderer) MakeNode(host NodeHost) *Node {
	return &Node{
		renderer:    r,
		primitives:  make(primitiveNodeMap, 1),
		Traversable: host,
		host:        host,
	}
}

// Any thread, update goroutine
//func (r *ARenderer) AddPrimitive() int {
//	id := r.idgen.NextId()
//	r.primitives[id] = &PrimitiveNode{}
//	r.layers[0].addPrimitiveId(id)
//	return id
//}

// Any thread, update goroutine
func (r *ARenderer) setPrimitiveOnFrontend(pNode *PrimitiveNode, data Primitive) {
	if delegate, ok := r.primitiveDelegates[data.GetType()]; ok {
		r.renderDispatcher.Execute(dispatch.NewWorkItemFunc(func() {
			pNode.primitive = data
			pNode.redraw = true

			ctx := r.makePrimitiveRenderingContext(pNode)
			delegate.OnSetPrimitive(ctx)
			pNode.state = ctx.State
		}))
	}
}

// Any thread, update goroutine
func (r *ARenderer) removePrimitiveOnFrontend(pNode *PrimitiveNode) {
	if d, ok := r.primitiveDelegates[pNode.primitive.GetType()]; ok {
		r.renderDispatcher.Execute(dispatch.NewWorkItemFunc(func() {
			d.OnRemovePrimitive(r.makePrimitiveRenderingContext(pNode))
		}))
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

func (r *ARenderer) Clear() {
	r.primitives = make(map[int]*PrimitiveNode)
	r.layers = make([]*Layer, 1)
	r.layers[0] = newLayer(0)
	r.idgen = common.NewIdGenerator()
	r.renderDispatcher.Execute(dispatch.NewWorkItemFunc(r.delegate.OnClear))
}

func (r *ARenderer) Stop() {
	for _, delegate := range r.primitiveDelegates {
		r.renderDispatcher.Execute(dispatch.NewWorkItemFunc(delegate.OnStop))
	}
	r.renderDispatcher.Execute(dispatch.NewWorkItemFunc(r.delegate.OnStop))
}

func (r *ARenderer) makePrimitiveRenderingContext(container *PrimitiveNode) *PrimitiveRenderingContext {
	return &PrimitiveRenderingContext{
		Renderer:      r,
		Primitive:     container.primitive,
		PrimitiveKind: container.primitive.GetType(),
		PrimitiveId:   container.id,
		State:         container.state,
		Redraw:        container.redraw,
		ClipArea2D:    r.clipStack.Peek(),
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
	r.layers = append(r.layers, newLayer(len(r.layers)))
}

func (r *ARenderer) RemoveLayer(index int) {
	if len(r.layers) <= 1 || index >= len(r.layers) {
		return
	}

	for i := index; i < len(r.layers)-1; i++ {
		r.layers[i] = r.layers[i+1]
		r.layers[i].index = i
	}

	r.layers = r.layers[:len(r.layers)-1]
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
	if r.root == nil {
		panic("No root node")
	}

	for _, layer := range r.layers {
		r.root.RenderTraverse(func(node *Node) {
			if node.clipArea != nil {
				r.clipStack.Push(node.clipArea)
			}

			for _, p := range node.GetPrimitivesInLayer(layer.index) {
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
		}, func(node *Node) {
			if node.clipArea != nil {
				r.clipStack.Pop()
			}
		})

		r.clipStack.Clear()
	}
}

func NewARenderer(delegate RendererDelegate, renderDispatcher dispatch.WorkDispatcher) *ARenderer {
	return &ARenderer{
		delegate:           delegate,
		renderDispatcher:   renderDispatcher,
		idgen:              common.NewIdGenerator(),
		primitives:         make(map[int]*PrimitiveNode),
		primitiveDelegates: make(map[byte]PrimitiveRendererDelegate),
		layers:             make([]*Layer, 1),
		clipStack:          NewClipStack2D(),
	}
}
