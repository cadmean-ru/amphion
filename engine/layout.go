package engine

type Layout interface {
	Component
	//LayoutResponder
	LayoutChildren()
}

type LayoutResponder interface {
	OnMeasure(ctx *LayoutContext)
	OnLayout(ctx *LayoutContext)
}

type LayoutContext struct {

}

type Layouter struct {

}