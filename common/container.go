package common

type Container interface {
	SetValue(value interface{})
	GetValue() interface{}
}
