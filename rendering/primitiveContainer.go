package rendering

type PrimitiveContainer struct {
	id int
	primitive IPrimitive
	redraw bool
	state interface{}
}
