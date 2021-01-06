package common

type IdGenerator struct {
	nextId int
}

func (g *IdGenerator) NextId() int {
	id := g.nextId
	g.nextId++
	return id
}

func NewIdGenerator() *IdGenerator {
	return &IdGenerator{
		nextId: 0,
	}
}
