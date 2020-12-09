package common

type Equatable interface {
	Equals(other interface{}) bool
}
