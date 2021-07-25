package rendering

type Traversable interface {
	RenderTraverse(func(node *Node), func(node *Node))
}
