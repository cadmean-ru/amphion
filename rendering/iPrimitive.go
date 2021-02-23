package rendering

type IPrimitive interface {
	GetType() byte
	GetTransform() Transform
	SetTransform(t Transform)
}
