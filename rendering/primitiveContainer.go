package rendering

import "github.com/cadmean-ru/amphion/common"

const (
	primitiveStatusNew = 0
	primitiveStatusIdle = 1
	primitiveStatusRedraw = 2
	primitiveStatusDeleting = 3
)

type primitiveContainer struct {
	id        int64
	status    int
	primitive common.ByteArrayEncodable
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
