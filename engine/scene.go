package engine

import (
	"encoding/json"
	"fmt"
	"github.com/cadmean-ru/amphion/common/a"
	"github.com/cadmean-ru/amphion/common/require"
	"gopkg.in/yaml.v2"
	"reflect"
	"regexp"
	"strings"
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
	object.setInCurrentScene(o.inCurrentScene)
	o.children = append(o.children, object)
}

// Adds a child object to this scene object.
func (o *SceneObject) AddChild(object *SceneObject) {
	o.appendChild(object)

	if o.inCurrentScene {
		if !object.initialized {
			instance.updateRoutine.initSceneObject(object)
			object.Traverse(func(child *SceneObject) bool {
				instance.updateRoutine.initSceneObject(child)
				return true
			}, true)
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

//GetChildrenCount returns the number of children of this object.
func (o *SceneObject) GetChildrenCount() int {
	return len(o.children)
}

// Finds scene object in the list of children of the current object.
// Returns nil if no object with the specified name was not found.
func (o *SceneObject) GetChildByName(name string) *SceneObject {
	for _, c := range o.children {
		if c != nil && c.name == name {
			return c
		}
	}

	panic(fmt.Sprintf("child scene object with name %s was not found (is it dirty?)", name))
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

// GetComponentByName searches for component with the specified name throughout the components attached to this object.
// The parameter n can be on of the following:
//
// - full component name, e.g. github.com/cadmean-ru/amphion/engine/builtin.TextView,
//
// - short component name, e.g. TextView,
//
// - regex string.
//
// If there are multiple components with the same name returns the first.
// Returns nil if no component with the name n was found or it has not been initialized or is disabled.
func (o *SceneObject) GetComponentByName(n string, includeDirty ...bool) Component {
	dirty := getDirty(includeDirty...)

	for _, c := range o.components {
		if o.componentMatcher(c, n, dirty) {
			return c.GetComponent()
		}
	}

	panic(fmt.Sprintf("component with name %s was not found (is it dirty?)", n))
}

// GetComponentsByName searches for all components with the specified name throughout the components attached to this object.
// The parameter n can be on of the following:
//
// - full component name, e.g. github.com/cadmean-ru/amphion/engine/builtin.TextView,
//
// - short component name, e.g. TextView,
//
// - regex string.
//
// Returns empty slice if no components with the name n was found.
//Components that have not been initialized or are disabled wont be included.
func (o *SceneObject) GetComponentsByName(n string, includeDirty ...bool) []Component {
	arr := make([]Component, 0, 1)

	dirty := getDirty(includeDirty...)

	for _, c := range o.components {
		if o.componentMatcher(c, n, dirty) {
			arr = append(arr, c.GetComponent())
		}
	}

	return arr
}

func (o *SceneObject) componentNameMatcher(comp Component, n string) bool {
	if n == comp.GetName() {
		return true
	}

	shortName := strings.Split(comp.GetName(), ".")[1] //The name after .
	if n == shortName {
		return true
	}

	if matched, err := regexp.MatchString(n, comp.GetName()); matched && err == nil {
		return true
	}

	return false
}

func (o *SceneObject) componentMatcher(container *ComponentContainer, name string, dirty bool) bool {
	return (dirty || !container.IsDirty()) && o.componentNameMatcher(container.GetComponent(), name)
}

// Returns a slice of all components attached to the object.
func (o *SceneObject) GetComponents() []Component {
	arr := make([]Component, 0, len(o.components))
	for _, c := range o.components {
		if c.IsDirty() {
			continue
		}

		arr = append(arr, c.component)
	}
	return arr
}

func (o *SceneObject) GetComponentContainers() []*ComponentContainer {
	arr := make([]*ComponentContainer, len(o.components))
	arr = append(arr, o.components...)
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
	if o.IsDirty() || !o.inCurrentScene {
		return
	}

	for _, view := range o.renderingComponents {
		if view.IsDirty() {
			continue
		}

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
	if o.initialized {
		return
	}

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

// Traverse traverses the scene object tree (pre-order), calling the action function for each of the objects.
// If action returns false interrupts the process.
// By default the method skips dirty objects.
// To also include dirty objects pass true as the second argument.
func (o *SceneObject) Traverse(action func(object *SceneObject) bool, includeDirty ...bool) {
	dirty := getDirty(includeDirty...)

	if !dirty && o.IsDirty() {
		return
	}

	if !action(o) {
		return
	}

	for _, c := range o.children {
		c.Traverse(action, dirty)
	}
}

// ForEachObject traverses the scene object tree, calling the action function for each of the objects.
// The action is also called for the object on which the method was called.
// The method skips dirty objects.
func (o *SceneObject) ForEachObject(action func(object *SceneObject)) {
	o.Traverse(func(object *SceneObject) bool {
		action(object)
		return true
	})
}

//ForEachChild cycles through all direct children of the scene object, calling the specified action for each of them.
//The method skips dirty objects.
func (o *SceneObject) ForEachChild(action func(object *SceneObject)) {
	if o.IsDirty() {
		return
	}

	for _, c := range o.children {
		action(c)
	}
}

//FindObjectByName searches for an object with the specified name through all the scene object tree.
//Returns the first suitable object.
//Returns nil if no object with the name was found.
//By default the search does not include dirty objects.
//To also include dirty objects pass true as the second argument.
func (o *SceneObject) FindObjectByName(name string, includeDirty ...bool) *SceneObject {
	var found *SceneObject
	dirty := getDirty(includeDirty...)

	o.Traverse(func(object *SceneObject) bool {
		if object.name == name {
			found = object
			return false
		}

		return true
	}, dirty)

	if found != nil {
		return found
	}

	panic(fmt.Sprintf("scene object with name %s was not found (is it dirty?)", name))
}

//FindComponentByName searches for a component with the specified name through all the scene object tree.
//Returns the first suitable component.
//Returns nil if no component with the name was found.
//By default the search does not include dirty components.
//To also include dirty components pass true as the second argument.
func (o *SceneObject) FindComponentByName(name string, includeDirty ...bool) Component {
	var found Component
	dirty := getDirty(includeDirty...)

	o.Traverse(func(object *SceneObject) bool {
		for _, c := range object.components {
			if o.componentMatcher(c, name, dirty) {
				found = c.GetComponent()
				return false
			}
		}

		return true
	}, dirty)

	if found != nil {
		return found
	}

	panic(fmt.Sprintf("component with name %s was not found (is it dirty?)", name))
}

// ForEachComponent iterates over each component attached to the object, calling the action function for each of them.
// The method skips dirty components.
func (o *SceneObject) ForEachComponent(action func(component Component)) {
	for _, c := range o.components {
		if c.IsDirty() {
			continue
		}
		action(c.component)
	}
}

//IsDirty checks if the scene object is dirty.
//An object is considered dirty if it has not been initialized or if it is disabled.
//It is not safe to work with a dirty object.
func (o *SceneObject) IsDirty() bool {
	return !o.initialized || !o.enabled
}

//DeepCopy creates a new scene object with all components and children of the receiver.
//The new object is dirty. It is enabled if the receiver is enabled.
func (o *SceneObject) DeepCopy(copyName string) *SceneObject {
	var newObject *SceneObject
	if instance == nil {
		newObject = NewSceneObjectForTesting(copyName)
	} else {
		newObject = NewSceneObject(copyName)
	}

	for _, comp := range o.components {
		compType := reflect.TypeOf(comp.component)
		var compCopy reflect.Value
		if compType.Kind() == reflect.Ptr {
			compCopy = reflect.New(compType.Elem())
		} else {
			compCopy = reflect.Indirect(reflect.New(compType))
		}

		newObject.AddComponent(compCopy.Interface().(Component))
	}

	newObject.Transform = o.Transform
	newObject.enabled = o.enabled

	for _, child := range o.children {
		childCopy := child.DeepCopy(child.name)
		newObject.AddChild(childCopy)
	}

	return newObject
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

func (o *SceneObject) DumpToMap() a.SiMap {
	mChildren := make([]map[string]interface{}, len(o.children))
	for i, c := range o.children {
		mChildren[i] = c.DumpToMap()
	}

	mComponents := make([]map[string]interface{}, len(o.components))
	for i, c := range o.components {
		var state map[string]interface{}

		if IsStatefulComponent(c.component) {
			state = instance.GetComponentsManager().GetComponentState(c.component)
		}

		cMap := map[string]interface{} {
			"name":  c.GetComponent().GetName(),
			"state": state,
			"initialized": c.initialized,
			"started": c.started,
			"enabled": c.enabled,
		}

		mComponents[i] = cMap
	}

	return map[string]interface{}{
		"name": o.name,
		"id": o.id,
		"initialized": o.initialized,
		"started": o.started,
		"enabled": o.enabled,
		"children": mChildren,
		"components": mComponents,
		"transform": o.Transform.ToMap(),
		"renderingTransform": o.Transform.ToRenderingTransform().ToMap(),
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

func (o *SceneObject) EncodeToJson() ([]byte, error) {
	return json.Marshal(o.ToMap())
}

func (o *SceneObject) DumpToJson() ([]byte, error) {
	return json.Marshal(o.DumpToMap())
}

// NewSceneObject creates a new instance of scene object.
// Can be used only with engine running.
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
	obj.Transform = NewTransform2D(obj)
	return obj
}

//NewSceneObjectForTesting creates a new instance of scene object without engine running.
//Can be used for testing purposes.
//Do not use in actual Amphion apps!
func NewSceneObjectForTesting(name string, components ...Component) *SceneObject {
	obj := &SceneObject{
		id:                  0,
		name:                name,
		children:            make([]*SceneObject, 0, 10),
		components:          make([]*ComponentContainer, 0, 10),
		messageListeners:    make([]MessageListenerComponent, 0),
		renderingComponents: make([]*ComponentContainer, 0, 1),
		updatingComponents:  make([]*ComponentContainer, 0, 1),
		boundaryComponents:  make([]*ComponentContainer, 0, 1),
		enabled:             true,
		initialized:         true,
	}
	obj.Transform = NewTransform2D(obj)

	for _, c := range components {
		obj.AddComponent(c)
	}

	ctx := newInitContext(&AmphionEngine{}, obj)
	for _, component := range obj.components {
		component.GetComponent().OnInit(ctx)
		component.initialized = true
	}

	return obj
}

func getDirty(dirties ...bool) bool {
	if len(dirties) == 1 {
		return dirties[0]
	}

	return false
}