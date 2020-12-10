package engine

import (
	"github.com/cadmean-ru/amphion/common"
	"gopkg.in/yaml.v2"
	"regexp"
)

// An object in the scene. The scene itself is also a SceneObject.
type SceneObject struct {
	id                  int64
	name                string
	children            []*SceneObject
	components          []*ComponentContainer
	messageListeners    []MessageListenerComponent
	updatingComponents  []*ComponentContainer
	renderingComponents []*ComponentContainer
	boundaryComponents  []*ComponentContainer
	Transform           Transform
	parent              *SceneObject
	enabled             bool
	initialized         bool
	started             bool
}

func (o *SceneObject) GetName() string {
	return o.name
}

func (o *SceneObject) GetId() int64 {
	return o.id
}

func (o *SceneObject) ToString() string {
	return o.name
}

func (o *SceneObject) GetParent() *SceneObject {
	return o.parent
}

func (o *SceneObject) AddChild(object *SceneObject) {
	object.parent = o
	object.Transform.parent = &o.Transform
	o.children = append(o.children, object)
	if !o.initialized {
		instance.updateRoutine.initSceneObject(o)
	}
	instance.rebuildMessageTree()
}

func (o *SceneObject) RemoveChild(object *SceneObject) {
	index := -1
	for i, c := range o.children {
		if c == object {
			index = i
		}
	}
	if index != -1 {
		o.children[index] = o.children[len(o.children)-1]
		o.children = o.children[:len(o.children)-1]

		instance.rebuildMessageTree()
	}
}

func (o *SceneObject) GetChildren() []*SceneObject {
	c := make([]*SceneObject, len(o.children))
	copy(c, o.children)
	return c
}

func (o *SceneObject) AddComponent(component Component) {
	container := NewComponentContainer(o, component)
	o.components = append(o.components, container)

	if _, ok := component.(UpdatingComponent); ok {
		o.updatingComponents = append(o.updatingComponents, container)
	}
	if _, ok := component.(ViewComponent); ok {
		o.renderingComponents = append(o.renderingComponents, container)
	}
	if _, ok := component.(MessageListenerComponent); ok {
		o.messageListeners = append(o.messageListeners, component.(MessageListenerComponent))
	}
	if _, ok := component.(BoundaryComponent); ok {
		o.boundaryComponents = append(o.boundaryComponents, container)
	}
}

func (o *SceneObject) GetComponentByName(n string) Component {
	for _, c := range o.components {
		comp := c.GetComponent()
		if matched, err := regexp.MatchString(n, comp.GetName()); matched && err == nil {
			return comp
		}
	}

	return nil
}

func (o *SceneObject) GetComponentsByName(n string) []Component {
	arr := make([]Component, 0, 1)

	for _, c := range o.components {
		comp := c.GetComponent()
		if matched, err := regexp.MatchString(n, comp.GetName()); matched && err == nil {
			arr = append(arr, comp)
		}
	}

	return arr
}

func (o *SceneObject) GetComponents() []Component {
	arr := make([]Component, len(o.components))
	for i, c := range o.components {
		arr[i] = c.component
	}
	return arr
}

func (o *SceneObject) SetEnabled(enabled bool) {
	if o.enabled == enabled {
		return
	}
	o.enabled = enabled
	if enabled {
		instance.updateRoutine.startSceneObject(o)
	} else {
		instance.updateRoutine.stopSceneObject(o)
	}
	for _, c := range o.components {
		c.SetEnabled(enabled)
	}
	for _, so := range o.children {
		so.SetEnabled(enabled)
	}
}

func (o *SceneObject) IsEnabled() bool {
	return o.enabled
}

func (o *SceneObject) OnMessage(message Message) bool {
	if !o.enabled {
		return true
	}

	continuePropagation := true
	for _, l := range o.messageListeners {
		if !l.OnMessage(message) {
			continuePropagation = false
		}
	}

	return continuePropagation
}

func (o *SceneObject) init(ctx InitContext) {
	if o.initialized {
		return
	}
	for _, c := range o.components {
		instance.currentComponent = c.component
		c.component.OnInit(ctx)
	}
	o.initialized = true
	instance.currentComponent = nil
}

func (o *SceneObject) start() {
	if o.started {
		return
	}
	o.ForEachComponent(func(c Component) {
		instance.currentComponent = c
		c.OnStart()
	})
	o.started = true
	instance.currentComponent = nil
}

func (o *SceneObject) update(ctx UpdateContext) {
	for _, c := range o.updatingComponents {
		if !c.enabled {
			continue
		}

		instance.currentComponent = c.component
		c.component.(UpdatingComponent).OnUpdate(ctx)
	}
	instance.currentComponent = nil
}

func (o *SceneObject) draw(ctx DrawingContext) {
	for _, c := range o.renderingComponents {
		if !c.enabled {
			continue
		}

		instance.currentComponent = c.component
		c.component.(ViewComponent).OnDraw(ctx)
	}
	instance.currentComponent = nil
}

func (o *SceneObject) stop() {
	if !o.started {
		return
	}
	for _, c := range o.components {
		instance.currentComponent = c.component
		c.component.OnStop()
	}
	o.started = false
	instance.currentComponent = nil
}

func (o *SceneObject) IsRendering() bool {
	return len(o.renderingComponents) > 0
}

func (o *SceneObject) HasBoundary() bool {
	return len(o.boundaryComponents) > 0
}

func (o *SceneObject) IsPointInsideBoundaries(point common.Vector3) bool {
	for _, b := range o.boundaryComponents {
		if b.component.(BoundaryComponent).IsPointInside(point) {
			return true
		}
	}

	return false
}

func (o *SceneObject) IsPointInsideBoundaries2D(point common.Vector3) bool {
	for _, b := range o.boundaryComponents {
		if b.component.(BoundaryComponent).IsPointInside2D(point) {
			return true
		}
	}

	return false
}

func (o *SceneObject) ForEachObject(action func(object *SceneObject)) {
	action(o)
	for _, c := range o.children {
		c.ForEachObject(action)
	}
}

func (o *SceneObject) ForEachComponent(action func(component Component)) {
	for _, c := range o.components {
		if !c.enabled {
			continue
		}
		action(c.component)
	}
}

func (o *SceneObject) ToMap() map[string]interface{} {
	mChildren := make([]map[string]interface{}, len(o.children))
	for i, c := range o.children {
		mChildren[i] = c.ToMap()
	}

	return map[string]interface{}{
		"name": o.name,
		"children": mChildren,
		"transform": o.Transform.ToMap(),
	}
}

func (o *SceneObject) EncodeToYaml() ([]byte, error) {
	return yaml.Marshal(o.ToMap())
}

func (o *SceneObject) DecodeFromYaml(data []byte) error {
	return yaml.Unmarshal(data, o)
}

func NewSceneObject(name string) *SceneObject {
	obj := &SceneObject{
		id:                  instance.idgen.NextId(),
		name:                name,
		children:            make([]*SceneObject, 0, 10),
		components:          make([]*ComponentContainer, 0, 10),
		messageListeners:    make([]MessageListenerComponent, 0),
		renderingComponents: make([]*ComponentContainer, 0, 1),
		updatingComponents:  make([]*ComponentContainer, 0, 1),
		boundaryComponents:  make([]*ComponentContainer, 0, 1),
		enabled:             true,
	}
	obj.Transform = NewTransform(obj)
	return obj
}