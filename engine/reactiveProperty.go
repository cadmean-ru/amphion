package engine

import "github.com/cadmean-ru/amphion/common"

type ReactiveProperty[T any] struct {
	value       T
	object      *SceneObject
	subscribers []func(T)
}

func (r *ReactiveProperty[T]) Value() T {
	return r.value
}

func (r *ReactiveProperty[T]) Set(value T) {
	if r.object == nil {
		panic("reactive property not attached to an object")
	}
	r.value = value
	r.object.Redraw()
	r.notify()
}

func (r *ReactiveProperty[T]) GetObject() *SceneObject {
	return r.object
}

func (r *ReactiveProperty[T]) Subscribe(delegate func(T)) {
	r.subscribers = append(r.subscribers, delegate)
}

func (r *ReactiveProperty[T]) Unsubscribe(delegate func(T)) {
	for i, d := range r.subscribers {
		if common.PtrEquals(d, delegate) {
			r.subscribers = append(r.subscribers[:i], r.subscribers[i+1:]...)
			return
		}
	}
}

func (r *ReactiveProperty[T]) notify() {
	for _, s := range r.subscribers {
		s(r.value)
	}
}

func NewReactiveProperty[T any](object *SceneObject, value T) *ReactiveProperty[T] {
	return &ReactiveProperty[T]{
		value:       value,
		object:      object,
		subscribers: make([]func(T), 0),
	}
}
