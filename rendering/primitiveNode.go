package rendering

type PrimitiveNode struct {
	id        int
	layer     int
	primitive Primitive
	redraw    bool
	state     interface{}
}

func newPrimitiveContainer(id int) *PrimitiveNode {
	return &PrimitiveNode{
		id: id,
	}
}
