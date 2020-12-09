package engine

type Object interface {
	ToString() string
}

type NamedObject interface {
	GetName() string
}

type AmphionObject interface {
	Object
	NamedObject
	GetId()   int64
}