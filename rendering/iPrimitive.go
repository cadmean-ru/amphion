package rendering

type Primitive interface {
	GetType() byte
	GetTransform() Transform
	SetTransform(t Transform)
}
