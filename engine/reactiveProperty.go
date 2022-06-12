package engine

type ReactiveProperty[T any] struct {
	value  T
	object *SceneObject
}

func (r *ReactiveProperty[T]) Get() T {
	return r.value
}

func (r *ReactiveProperty[T]) Set(value T) {
	r.value = value
	r.object.Redraw()
	RequestRendering()
}

func (r *ReactiveProperty[T]) GetObject() *SceneObject {
	return r.object
}

func NewReactiveProperty[T any](object *SceneObject, value T) *ReactiveProperty[T] {
	return &ReactiveProperty[T]{
		value:  value,
		object: object,
	}
}