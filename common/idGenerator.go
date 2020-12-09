package common

type IdGenerator struct {
	nextId int64
}

func (g *IdGenerator) NextId() int64 {
	id := g.nextId
	g.nextId++
	return id
}

func NewIdGenerator() *IdGenerator {
	return &IdGenerator{
		nextId: 0,
	}
}
