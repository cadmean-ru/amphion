package engine

type Layout interface {
	Component
	LayoutChildren()
}
