package a

// Interface for values that can be serialized as string.
type Stringable interface {
	ToString() string
}

// Interface for values that can be deserialized from string.
type Unstringable interface {
	FromString(src string)
}
