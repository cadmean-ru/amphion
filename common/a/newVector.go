package a

import "fmt"

type Number interface {
	byte | int | int32 | int64 | float32 | float64
}

type Vector[T Number] struct {
	length int
	values []T
}

func (v *Vector[T]) SetValueAt(index int, value T) {
	v.values[index] = value
}

func (v *Vector[T]) GetValueAt(index int) T {
	return v.values[index]
}

func (v *Vector[T]) Set(values ...T) {
	if len(values) > v.length {
		panic(fmt.Sprintf("Vector is of length %d, but %d values were provided", v.length, len(values)))
	}

	for i := 0; i < len(values); i++ {
		v.values[i] = values[i]
	}
}

func (v *Vector[T]) Length() int {
	return v.length
}

func (v *Vector[T]) Add(x Vector[T]) *Vector[T] {
	if v.length != x.length {
		panic("Vector length mismatch")
	}
	result := NewVector[T](v.length)
	for i := 0; i < v.length; i++ {
		result.values[i] = v.values[i] + x.values[i]
	}
	return result
}

func (v *Vector[T]) Subtract(x Vector[T]) *Vector[T] {
	if v.length != x.length {
		panic("Vector length mismatch")
	}
	result := NewVector[T](v.length)
	for i := 0; i < v.length; i++ {
		result.values[i] = v.values[i] - x.values[i]
	}
	return result
}

func (v *Vector[T]) Multiply(x Vector[T]) *Vector[T] {
	if v.length != x.length {
		panic("Vector length mismatch")
	}
	result := NewVector[T](v.length)
	for i := 0; i < v.length; i++ {
		result.values[i] = v.values[i] * x.values[i]
	}
	return result
}

func (v *Vector[T]) MultiplyScalar(x T) *Vector[T] {
	result := NewVector[T](v.length)
	for i := 0; i < v.length; i++ {
		result.values[i] = v.values[i] * x
	}
	return result
}

func (v *Vector[T]) Equals(x Vector[T]) bool {
	if v.length != x.length {
		return false
	}

	for i := 0; i < v.length; i++ {
		if v.values[i] != x.values[i] {
			return false
		}
	}

	return true
}

func NewVector[T Number](length int) *Vector[T] {
	return &Vector[T]{
		length: length,
		values: make([]T, length),
	}
}
