package rendering

import (
	"fmt"
	"github.com/cadmean-ru/amphion/common/a"
)

type primitiveNodeMap map[int]*PrimitiveNode

type NodeHost interface {
	a.NamedObject
	Traversable
}

type Node struct {
	Traversable
	host       a.NamedObject
	renderer   *ARenderer
	primitives primitiveNodeMap
	clipArea   *ClipArea2D
}

func (n *Node) GetPrimitives() []*PrimitiveNode {
	primitivesCopy := make([]*PrimitiveNode, len(n.primitives))
	i := 0
	for _, v := range n.primitives {
		primitivesCopy[i] = v
		i++
	}
	return primitivesCopy
}

func (n *Node) GetPrimitive(id int) *PrimitiveNode {
	for _, n := range n.primitives {
		if n.id == id {
			return n
		}
	}

	return nil
}

func (n *Node) AddPrimitive() int {
	newId := n.renderer.idgen.NextId()
	n.primitives[newId] = newPrimitiveContainer(newId)
	return newId
}

func (n *Node) SetPrimitive(id int, data Primitive) {
	pNode, found := n.primitives[id]
	if !found {
		fmt.Printf("Warning! Primitive with id %d was not found.\n", id)
		return
	}
	if pNode.primitive != nil && pNode.primitive.GetType() != data.GetType() {
		fmt.Printf("Cannot change primitive type for id: %d\n", id)
		return
	}

	n.renderer.setPrimitiveOnFrontend(pNode, data)
}

func (n *Node) GetPrimitivesInLayer(layerIndex int) []*PrimitiveNode {
	primitivesList := make([]*PrimitiveNode, 0, len(n.primitives))
	for _, p := range n.primitives {
		if p.layer != layerIndex  || p.primitive == nil {
			continue
		}

		primitivesList = append(primitivesList, p)
	}
	return primitivesList
}

func (n *Node) SetPrimitiveLayer(id, layerIndex int) {
	if layerIndex >= len(n.renderer.layers) {
		fmt.Printf("Warning! No layer index %d\n", layerIndex)
		return
	}

	p := n.GetPrimitive(id)
	if p == nil {
		return
	}

	p.layer = layerIndex
}

func (n *Node) RemovePrimitive(id int) {
	pNode, found := n.primitives[id]
	if !found {
		return
	}

	delete(n.primitives, id)
	n.renderer.removePrimitiveOnFrontend(pNode)
}

func (n *Node) SetClipArea2D(area *ClipArea2D) {
	n.clipArea = area
}

func (n *Node) RemoveClipArea2D() {
	n.clipArea = nil
}