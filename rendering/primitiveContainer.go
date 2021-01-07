package rendering

const (
	primitiveStatusNew = 0
	primitiveStatusIdle = 1
	primitiveStatusRedraw = 2
	primitiveStatusDeleting = 3
)

// Deprecated: no longer needed
type primitiveContainer struct {
	id        int64
	status    int
	primitive *Primitive
}

func (c *primitiveContainer) GetId() int64 {
	return c.id
}

func newPrimitiveContainer(id int64) *primitiveContainer {
	return &primitiveContainer{
		id:        id,
		status:    primitiveStatusNew,
		primitive: nil,
	}
}
