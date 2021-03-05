package engine

import (
	"fmt"
	"github.com/cadmean-ru/amphion/common/a"
	"github.com/cadmean-ru/amphion/common/require"
	"gopkg.in/yaml.v2"
	"regexp"
)

// An object in the scene. The scene itself is also a SceneObject.
type SceneObject struct {
	id                  int
	name                string
	children            []*SceneObject
	components          []*ComponentContainer
	messageListeners    []MessageListenerComponent
	updatingComponents  []*ComponentContainer
	renderingComponents []*ComponentContainer
	boundaryComponents  []*ComponentContainer
	layout              Layout
	Transform           Transform
	parent              *SceneObject
	enabled             bool
	initialized         bool
	started             bool
	inCurrentScene      bool
}

//GetName returns the name of the scene object.
func (o *SceneObject) GetName() string {
	return o.name
}

// GetId returns the id of the object in the scene.
// The id is not guaranteed to remain the same over time and serves for internal engine purposes.
// For identification use object's name.
func (o *SceneObject) GetId() int {
	return o.id
}

// Returns the string representation of the scene object.
func (o *SceneObject) ToString() string {
	return o.name
}

// For GoLand debugging.
func (o *SceneObject) DebugString() string {
	return o.ToString()
}

// Returns the parent object of this scene object. Returns nil if no parent object.
func (o *SceneObject) GetParent() *SceneObject {
	return o.parent
}

func (o *SceneObject) appendChild(object *SceneObject) {
	object.parent = o
	object.Transform.SceneObject = object
	object.Transform.parent = &o.Transform
	object.inCurrentScene = o.inCurrentScene
	o.children = append(o.children, object)
}

// Adds a child object to this scene object.
func (o *SceneObject) AddChild(object *SceneObject) {
	o.appendChild(object)

	if o.inCurrentScene {
		if !object.initialized {
			instance.updateRoutine.initSceneObject(object)
		}
		instance.rebuildMessageTree()
		instance.RequestRendering()
	}
}

// Removes a child from this scene object.
func (o *SceneObject) RemoveChild(object *SceneObject) {
	index := -1
	for i, c := range o.children {
		if c == object {
			index = i
		}
	}
	if index != -1 {
		object.SetEnabled(false)

		o.children[index] = o.children[len(o.children)-1]
		o.children = o.children[:len(o.children)-1]

		if o.inCurrentScene {
			instance.rebuildMessageTree()
			instance.RequestRendering()
		}
	}
}

// Returns the list of children of this object.
func (o *SceneObject) GetChildren() []*SceneObject {
	c := make([]*SceneObject, len(o.children))
	copy(c, o.children)
	return c
}

// Finds scene object in the list of children of the current object.
// Returns nil if no object with the specified name was not found.
func (o *SceneObject) GetChildByName(name string) *SceneObject {
	for _, c := range o.children {
		if c != nil && c.name == name {
			return c
		}
	}
	return nil
}

// Adds a component to this scene object.
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
	if _, ok := component.(Layout); ok {
		o.layout = component.(Layout)
	}

	if o.inCurrentScene {
		instance.updateRoutine.initSceneObject(o)
		instance.RequestRendering()
	}
}

// Searches for component with the specified name throughout the components attached to this object.
// The parameter n can be a regex string.
// If there are multiple components with the same name returns the first.
// Returns nil if no component with the name n was found.
func (o *SceneObject) GetComponentByName(n string) Component {
	for _, c := range o.components {
		comp := c.GetComponent()
		if matched, err := regexp.MatchString(n, comp.GetName()); matched && err == nil {
			return comp
		}
	}

	return nil
}

// Searches for all components with the specified name throughout the components attached to this object.
// The parameter n can be a regex string.
// Returns empty slice if no components with the name n was found.
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

// Returns a slice of all components attached to the object.
func (o *SceneObject) GetComponents() []Component {
	arr := make([]Component, len(o.components))
	for i, c := range o.components {
		arr[i] = c.component
	}
	return arr
}

// Set the enabled state of this object as specified.
func (o *SceneObject) SetEnabled(enabled bool) {
	if o.enabled == enabled {
		return
	}

	o.enabled = enabled

	for _, c := range o.components {
		c.SetEnabled(enabled)
	}

	if o.inCurrentScene {
		if enabled {
			instance.updateRoutine.startSceneObject(o)
		} else {
			instance.updateRoutine.stopSceneObject(o)
		}
	}

	for _, so := range o.children {
		so.SetEnabled(enabled)
	}
}

// Returns if this object is currently enabled or not.
func (o *SceneObject) IsEnabled() bool {
	return o.enabled
}

// Set the position of this object equal to the specified vector, requesting rendering.
func (o *SceneObject) SetPosition(v a.Vector3) {
	o.Transform.Position = v
	o.Redraw()
}

// Set the position of this object equal to a new vector with specified coordinates, requesting rendering.
func (o *SceneObject) SetPositionXyz(x, y, z float32) {
	o.SetPosition(a.NewVector3(x, y, z))
}

func (o *SceneObject) SetPositionXy(x, y float32) {
	o.SetPosition(a.NewVector3(x, y, o.Transform.Position.Z))
}

// Set the size of this object equal to the specified vector, requesting rendering.
func (o *SceneObject) SetSize(v a.Vector3) {
	o.Transform.Size = v
	o.Redraw()
}

// Set the size of this object equal to a new vector with specified coordinates, requesting rendering.
func (o *SceneObject) SetSizeXyz(x, y, z float32) {
	o.SetSize(a.NewVector3(x, y, z))
}

// Set the size of this object equal to a new vector with specified coordinates, requesting rendering.
func (o *SceneObject) SetSizeXy(x, y float32) {
	o.SetSize(a.NewVector3(x, y, o.Transform.Size.Z))
}

// Forces all views of this object to redraw and requests rendering.
func (o *SceneObject) Redraw() {
	if !o.inCurrentScene {
		return
	}

	for _, view := range o.renderingComponents {
		view.component.(ViewComponent).ForceRedraw()
	}
	instance.RequestRendering()
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
	for _, c := range o.components {
		if c.initialized {
			continue
		}
		instance.currentComponent = c.component
		c.component.OnInit(ctx)
		c.initialized = true
	}

	o.initialized = true
	instance.currentComponent = nil
}

func (o *SceneObject) start() {
	for _, c := range o.components {
		if !c.enabled || !c.initialized || c.started {
			continue
		}
		instance.currentComponent = c.component
		c.component.OnStart()
		c.started = true
	}

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

	if o.layout != nil {
		instance.currentComponent = o.layout
		o.layout.LayoutChildren()
	}

	instance.currentComponent = nil
}

func (o *SceneObject) draw(ctx DrawingContext) {
	for _, c := range o.renderingComponents {
		if !c.enabled || !c.initialized {
			continue
		}

		instance.currentComponent = c.component
		c.component.(ViewComponent).OnDraw(ctx)
	}
	instance.currentComponent = nil
}

func (o *SceneObject) stop() {
	for _, c := range o.components {
		if c.enabled || !c.started {
			continue
		}
		instance.currentComponent = c.component
		c.component.OnStop()
		c.started = false
	}
	o.started = false
	instance.currentComponent = nil
}

func (o *SceneObject) setInCurrentScene(b bool) {
	o.inCurrentScene = b
	for _, c := range o.children {
		c.setInCurrentScene(b)
	}
}

//IsRendering checks if the scene object has any view components.
func (o *SceneObject) IsRendering() bool {
	return len(o.renderingComponents) > 0
}

//HasBoundary checks if the scene object has any boundary components.
func (o *SceneObject) HasBoundary() bool {
	return len(o.boundaryComponents) > 0
}

func (o *SceneObject) IsPointInsideBoundaries(point a.Vector3) bool {
	for _, b := range o.boundaryComponents {
		if b.component.(BoundaryComponent).IsPointInside(point) {
			return true
		}
	}

	return false
}

func (o *SceneObject) IsPointInsideBoundaries2D(point a.Vector3) bool {
	for _, b := range o.boundaryComponents {
		if b.component.(BoundaryComponent).IsPointInside2D(point) {
			return true
		}
	}

	return false
}

//IsFocused checks if the scene object is currently focused.
func (o *SceneObject) IsFocused() bool {
	return instance.sceneContext.focusedObject == o
}

//IsHovered checks if the cursor is over the scene object.
func (o *SceneObject) IsHovered() bool {
	return instance.sceneContext.hoveredObject == o
}

func (o *SceneObject) IsVisibleInScene() bool {
	sceneRect := instance.GetCurrentScene().Transform.GetGlobalRect()
	rect := o.Transform.GetGlobalRect()
	return rect.X.Max >= sceneRect.X.Min && rect.X.Min <= sceneRect.X.Max && rect.Y.Max >= sceneRect.Y.Min && rect.Y.Min <= sceneRect.Y.Max
}

// ForEachObject traverses the scene object tree, calling the action function for each of the objects.
// The action is also called for the object on which the method was called.
// The method skips uninitialized or disabled objects.
func (o *SceneObject) ForEachObject(action func(object *SceneObject)) {
	if !o.enabled {
		return
	}

	action(o)

	for _, c := range o.children {
		c.ForEachObject(action)
	}
}

func (o *SceneObject) ForEachChild(action func(object *SceneObject)) {
	if !o.enabled {
		return
	}

	for _, c := range o.children {
		action(c)
	}
}

// ForEachComponent iterates over each component attached to the object, calling the action function for each of them.
// The method skips uninitialized or disabled components.
func (o *SceneObject) ForEachComponent(action func(component Component)) {
	for _, c := range o.components {
		if !c.enabled || !c.initialized {
			continue
		}
		action(c.component)
	}
}

func (o *SceneObject) ToMap() a.SiMap {
	mChildren := make([]map[string]interface{}, len(o.children))
	for i, c := range o.children {
		mChildren[i] = c.ToMap()
	}

	mComponents := make([]map[string]interface{}, len(o.components))
	for i, c := range o.components {
		var state map[string]interface{}

		if IsStatefulComponent(c.component) {
			state = instance.GetComponentsManager().GetComponentState(c.component)
			fmt.Println(state)
		}

		cMap := map[string]interface{} {
			"name":  c.GetComponent().GetName(),
			"state": state,
		}

		mComponents[i] = cMap
	}

	return map[string]interface{}{
		"name": o.name,
		"id": o.id,
		"children": mChildren,
		"components": mComponents,
		"transform": o.Transform.ToMap(),
	}
}

func (o *SceneObject) FromMap(siMap a.SiMap) {
	o.name = siMap["name"].(string)
	o.id = require.Int(siMap["id"])
	o.Transform = NewTransformFromMap(a.RequireSiMap(siMap["transform"]))

	// Decode components
	iComponents := siMap["components"].([]interface{})
	o.components = make([]*ComponentContainer, 0, len(iComponents))
	o.renderingComponents = make([]*ComponentContainer, 0, 1)
	o.updatingComponents = make([]*ComponentContainer, 0, 1)
	o.boundaryComponents = make([]*ComponentContainer, 0, 1)
	o.enabled = true
	for _, c := range iComponents {
		cMap := a.RequireSiMap(c)
		cName := cMap["name"].(string)
		cState := a.RequireSiMap(cMap["state"])
		component := instance.GetComponentsManager().MakeComponent(cName)
		if component == nil {
			continue
		}
		if IsStatefulComponent(component) {
			instance.GetComponentsManager().SetComponentState(component, cState)
		}
		o.AddComponent(component)
	}

	// Decode children
	iChildren := siMap["children"].([]interface{})
	o.children = make([]*SceneObject, 0, len(iChildren))
	for _, c := range iChildren {
		obj := &SceneObject{}
		obj.FromMap(a.RequireSiMap(c))
		o.appendChild(obj)
	}
}

func (o *SceneObject) EncodeToYaml() ([]byte, error) {
	return yaml.Marshal(o.ToMap())
}

func (o *SceneObject) DecodeFromYaml(data []byte) error {
	oMap := make(a.SiMap)
	err := yaml.Unmarshal(data, &oMap)
	if err != nil {
		return err
	}

	o.FromMap(oMap)

	return nil
}

// NewSceneObject creates a new instance of scene object.
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
