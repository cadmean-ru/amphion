package common

// Interface for making two values comparable.
type Equatable interface {
	Equals(other interface{}) bool
}
