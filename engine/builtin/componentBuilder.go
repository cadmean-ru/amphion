package builtin

import "github.com/cadmean-ru/amphion/engine"

type BuilderComponent struct {
	name     string
	eng      *engine.AmphionEngine
	onStart  func()
	onStop   func()
	onUpdate func(ctx engine.UpdateContext)
}

func (c *BuilderComponent) OnInit(ctx engine.InitContext) {
	c.eng = ctx.GetEngine()
}

func (c *BuilderComponent) OnStart() {
	if c.onStart != nil {
		c.onStart()
	}
}

func (c *BuilderComponent) OnUpdate(ctx engine.UpdateContext) {
	if c.onUpdate != nil {
		c.onUpdate(ctx)
	}
}

func (c *BuilderComponent) OnStop() {
	if c.onStop != nil {
		c.onStop()
	}
}

func (c *BuilderComponent) GetName() string {
	return c.name
}

type ComponentBuilder struct {
	comp *BuilderComponent
}

func (b *ComponentBuilder) OnStart(onStart func()) *ComponentBuilder {
	b.comp.onStart = onStart
	return b
}

func (b *ComponentBuilder) OnUpdate(onUpdate func(ctx engine.UpdateContext)) *ComponentBuilder {
	b.comp.onUpdate = onUpdate
	return b
}

func (b *ComponentBuilder) OnStop(onStop func()) *ComponentBuilder {
	b.comp.onStop = onStop
	return b
}

func (b *ComponentBuilder) Build() *BuilderComponent {
	return b.comp
}

func NewComponentBuilder(name string) *ComponentBuilder {
	return &ComponentBuilder{
		comp: &BuilderComponent{
			name: name,
		},
	}
}