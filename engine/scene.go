package engine

import (
	"encoding/json"
	"fmt"
	"github.com/cadmean-ru/amphion/common/a"
	"github.com/cadmean-ru/amphion/common/require"
	"github.com/cadmean-ru/amphion/rendering"
	"gopkg.in/yaml.v2"
	"reflect"
)

const MaxSceneObjects = 1000

// SceneObject is an object in the scene. The scene itself is also a SceneObject.
type SceneObject struct {
	id                  int
	name                string
	children            []*SceneObject
	components          []*ComponentContainer
	messageListeners    []*ComponentContainer
	updatingComponents  []*ComponentContainer
	renderingComponents []*ComponentContainer
	boundaryComponents  []*ComponentContainer
	layout              *ComponentContainer
	renderingNode       *rendering.Node
	Transform           Transform
	parent              *SceneObject
	enabled             bool
	initialized         bool
	started             bool
	inCurrentScene      bool
	willBeRemoved       bool
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

// ToString returns the string representation of the scene object.
func (o *SceneObject) ToString() string {
	return o.name
}

// DebugString return string for debugging.
func (o *SceneObject) DebugString() string {
	return o.ToString()
}

//GetParent returns the parent object of this scene object. Returns nil if no parent object.
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

//AddChild adds a child object to this scene object.
func (o *SceneObject) AddChild(object *SceneObject) {
	o.appendChild(object)

	if o.inCurrentScene {
		if !object.initialized {
			object.Traverse(func(child *SceneObject) bool {
				instance.updateRoutine.initSceneObject(child)
				return true
			}, true)
		}
		instance.rebuildMessageTree()
		instance.RequestRendering()
	}
}

//AddNewChild creates a new scene object with the given name and adds it as a child to the current object.
func (o *SceneObject) AddNewChild(name string) *SceneObject {
	obj := NewSceneObject(name)
	o.AddChild(obj)
	return obj
}

//RemoveChild removes the specified child from this scene object.
func (o *SceneObject) RemoveChild(object *SceneObject) {
	if !o.removeChildFromList(object) {
		return
	}

	object.SetEnabled(false)
	object.markForRemoval()

	if o.inCurrentScene {
		instance.rebuildMessageTree()
		instance.RequestRendering()
	}
}

func (o *SceneObject) removeChildFromList(object *SceneObject) bool {
	index := -1
	for i, c := range o.children {
		if c == object {
			index = i
		}
	}

	if index != -1 {
		o.children[index] = o.children[len(o.children)-1]
		o.children = o.children[:len(o.children)-1]
		return true
	}

	return false
}

//RemoveAllChildren removes all children from the receiver scene object.
func (o *SceneObject) RemoveAllChildren() {
	for _, o1 := range o.children {
		o.RemoveChild(o1)
	}
}

//GetChildren returns the list of children of this scene object.
//Modifying the returned list wont modify the actual list of children of this scene object.
func (o *SceneObject) GetChildren() []*SceneObject {
	c := make([]*SceneObject, len(o.children))
	copy(c, o.children)
	return c
}

//GetChildAt returns the child scene object at the given index.
//If there is no child with that index returns nil.
func (o *SceneObject) GetChildAt(index int) *SceneObject {
	if index >= 0 && index < len(o.children) {
		return o.children[index]
	} else {
		return nil
	}
}

//GetChildrenCount returns the number of children of this object.
func (o *SceneObject) GetChildrenCount() int {
	return len(o.children)
}

//GetChildByName finds scene object in the list of children of the current object.
//Returns nil if no object with the specified name was not found.
func (o *SceneObject) GetChildByName(name string) *SceneObject {
	for _, c := range o.children {
		if c != nil && c.name == name {
			return c
		}
	}

	panic(fmt.Sprintf("child scene object with name %s was not found (is it dirty?)", name))
}

//AddComponent adds a component to this scene object.
//Returns the given component.
func (o *SceneObject) AddComponent(component Component) Component {
	container := NewComponentContainer(o, component)
	o.components = append(o.components, container)

	if _, ok := component.(UpdatingComponent); ok {
		o.updatingComponents = append(o.updatingComponents, container)
	}
	if _, ok := component.(ViewComponent); ok {
		o.renderingComponents = append(o.renderingComponents, container)
	}
	if _, ok := component.(MessageListenerComponent); ok {
		o.messageListeners = append(o.messageListeners, container)
	}
	if _, ok := component.(BoundaryComponent); ok {
		o.boundaryComponents = append(o.boundaryComponents, container)
	}
	if _, ok := component.(Layout); ok {
		o.layout = container
	}

	if o.inCurrentScene {
		instance.updateRoutine.initSceneObject(o)
		instance.RequestRendering()
	}

	return component
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
//Components that have not been initialized or are disabled will not be included.
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

func (o *SceneObject) componentMatcher(container *ComponentContainer, name string, dirty bool) bool {
	return (dirty || !container.IsDirty()) && ComponentNameMatches(NameOfComponent(container.GetComponent()), name)
}

//GetComponents returns a slice of all components attached to the object.
//Modifying the returned list will not change the actual list of components of this scene object.
func (o *SceneObject) GetComponents(includeDirty ...bool) []Component {
	dirty := getDirty(includeDirty...)

	arr := make([]Component, 0, len(o.components))
	for _, c := range o.components {
		if !dirty && c.IsDirty() {
			continue
		}

		arr = append(arr, c.component)
	}
	return arr
}

//RemoveComponent removes the given component from the scene object.
func (o *SceneObject) RemoveComponent(comp Component) {
	var container *ComponentContainer
	for _, c := range o.components {
		if c.component == comp {
			container = c
			break
		}
	}

	o.removeComponentFromAllLists(container)
}

//RemoveComponentByName removes a component with the given name from the scene object&
//If there are more than one component with the same name only the first encountered component will be removed.
func (o *SceneObject) RemoveComponentByName(name string) {
	var container *ComponentContainer
	for _, c := range o.components {
		if o.componentMatcher(c, name, true) {
			container = c
			break
		}
	}

	o.removeComponentFromAllLists(container)
}

func (o *SceneObject) removeComponentFromAllLists(container *ComponentContainer) {
	if container == nil {
		return
	}

	component := container.component
	name := NameOfComponent(component)

	o.components = o.removeComponentFromList(o.components, name)

	if _, ok := component.(UpdatingComponent); ok {
		o.updatingComponents = o.removeComponentFromList(o.updatingComponents, name)
	}
	if _, ok := component.(ViewComponent); ok {
		o.renderingComponents = o.removeComponentFromList(o.renderingComponents, name)
	}
	if _, ok := component.(MessageListenerComponent); ok {
		o.messageListeners = o.removeComponentFromList(o.messageListeners, name)
	}
	if _, ok := component.(BoundaryComponent); ok {
		o.boundaryComponents = o.removeComponentFromList(o.boundaryComponents, name)
	}
	if _, ok := component.(Layout); ok {
		o.layout = nil
	}

	if o.inCurrentScene {
		instance.updateRoutine.stopComponent(container)
		instance.RequestRendering()
	}
}

func (o *SceneObject) removeComponentFromList(arr []*ComponentContainer, name string) []*ComponentContainer {
	index := -1
	for i, c := range o.components {
		if o.componentMatcher(c, name, true) {
			index = i
			break
		}
	}

	if index == -1 {
		return arr
	}

	arr[index] = arr[len(arr)-1]
	return arr[:len(arr)-1]
}

//RemoveAllComponents removes all components of the scene object.
func (o *SceneObject) RemoveAllComponents() {
	for _, c := range o.components {
		o.RemoveComponent(c.component)
	}
}

func (o *SceneObject) markForRemoval() {
	o.willBeRemoved = true

	for _, c := range o.children {
		c.markForRemoval()
	}
}

//GetComponentContainers returns the list of containers, that makes possible to enable/disable specific components.
func (o *SceneObject) GetComponentContainers() []*ComponentContainer {
	arr := make([]*ComponentContainer, 0, len(o.components))
	arr = append(arr, o.components...)
	return arr
}

//SetEnabled sets the enabled state of this object as specified.
func (o *SceneObject) SetEnabled(enabled bool) {
	if o.willBeRemoved || o.enabled == enabled {
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

//IsEnabled returns if this object is currently enabled or not.
func (o *SceneObject) IsEnabled() bool {
	return o.enabled
}

// SetPosition sets the position of this object equal to the specified vector, requesting rendering.
func (o *SceneObject) SetPosition(v a.Vector3) {
	o.Transform.Position = v
	o.Redraw()
}

// SetPositionXyz sets the position of this object equal to a new vector with specified coordinates, requesting rendering.
func (o *SceneObject) SetPositionXyz(x, y, z float32) {
	o.SetPosition(a.NewVector3(x, y, z))
}

//SetPositionXy sets the position of this object equal to a new vector with specified coordinates, requesting rendering.
func (o *SceneObject) SetPositionXy(x, y float32) {
	o.SetPosition(a.NewVector3(x, y, o.Transform.Position.Z))
}

// SetSize Set the size of this object equal to the specified vector, requesting rendering.
func (o *SceneObject) SetSize(v a.Vector3) {
	o.Transform.Size = v
	o.Redraw()
}

// SetSizeXyz Set the size of this object equal to a new vector with specified coordinates, requesting rendering.
func (o *SceneObject) SetSizeXyz(x, y, z float32) {
	o.SetSize(a.NewVector3(x, y, z))
}

// SetSizeXy Set the size of this object equal to a new vector with specified coordinates, requesting rendering.
func (o *SceneObject) SetSizeXy(x, y float32) {
	o.SetSize(a.NewVector3(x, y, o.Transform.Size.Z))
}

// Redraw forces all views of this object to redraw and requests rendering.
func (o *SceneObject) Redraw() {
	if o.IsDirty() || !o.inCurrentScene {
		return
	}

	for _, view := range o.renderingComponents {
		if view.IsDirty() {
			continue
		}

		viewComp := view.component.(ViewComponent)
		viewComp.Redraw()
	}
	instance.RequestRendering()
}

//func (o *SceneObject) Measure() a.IntVector3 {
//	rect := common.NewRectBoundary(0, 0, 0, 0, 0, 0)
//	for _, c := range o.children {
//
//	}
//}

func (o *SceneObject) OnMessage(message Message) bool {
	if !o.enabled {
		return true
	}

	continuePropagation := true
	for _, l := range o.messageListeners {
		if !l.component.(MessageListenerComponent).OnMessage(message) {
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
		if NameOfComponent(c.component) == "github.com/cadmean-ru/amphion/engine/builtin.ClipArea" {
			LogDebug("init %+v", c)
		}
	}

	o.initialized = true
	instance.currentComponent = nil
}

func (o *SceneObject) start() {
	for _, c := range o.components {
		if c.IsDirty() || c.started {
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
		if c.IsDirty() || !c.started {
			continue
		}

		instance.currentComponent = c.component
		c.component.(UpdatingComponent).OnUpdate(ctx)
	}

	if o.layout != nil && !o.layout.IsDirty() {
		instance.currentComponent = o.layout.component
		o.layout.component.(Layout).LayoutChildren()
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

func (o *SceneObject) draw(ctx DrawingContext) {
	for _, c := range o.renderingComponents {
		if c.IsDirty() || !c.started {
			continue
		}

		view := c.component.(ViewComponent)
		instance.currentComponent = c.component

		if !view.ShouldDraw() {
			continue
		}

		view.OnDraw(ctx)
	}
	instance.currentComponent = nil
}

func (o *SceneObject) RenderTraverse(action func(node *rendering.Node), afterChildrenAction func(node *rendering.Node)) {
	if o.IsDirty() {
		return
	}

	action(o.renderingNode)

	for _, c := range o.children {
		c.RenderTraverse(action, afterChildrenAction)
	}

	afterChildrenAction(o.renderingNode)
}

func (o *SceneObject) GetRenderingNode() *rendering.Node {
	return o.renderingNode
}

func (o *SceneObject) setInCurrentScene(b bool) {
	o.inCurrentScene = b
	for _, c := range o.children {
		c.setInCurrentScene(b)
	}
}

//HasViews checks if the scene object has any view components.
func (o *SceneObject) HasViews() bool {
	return len(o.renderingComponents) > 0
}

//HasBoundary checks if the scene object has any boundary components.
func (o *SceneObject) HasBoundary() bool {
	return len(o.boundaryComponents) > 0
}

func (o *SceneObject) IsPointInsideBoundaries(point a.Vector3) bool {
	for _, b := range o.boundaryComponents {
		if b.IsDirty() {
			continue
		}

		if b.component.(BoundaryComponent).IsPointInside(point) {
			return true
		}
	}

	return false
}

func (o *SceneObject) IsPointInsideBoundaries2D(point a.Vector3) bool {
	for _, b := range o.boundaryComponents {
		if b.IsDirty() {
			continue
		}

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
func (o *SceneObject) ForEachObject(action func(object *SceneObject), includeDirty ...bool) {
	dirty := getDirty(includeDirty...)

	o.Traverse(func(object *SceneObject) bool {
		action(object)
		return true
	}, dirty)
}

//ForEachChild cycles through all direct children of the scene object, calling the specified action for each of them.
//The method skips dirty objects.
func (o *SceneObject) ForEachChild(action func(object *SceneObject), includeDirty ...bool) {
	dirty := getDirty(includeDirty...)

	if !dirty && o.IsDirty() {
		return
	}

	for _, c := range o.children {
		if !dirty && c.IsDirty() {
			continue
		}

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
//An object is considered dirty if it has not been initialized or if it is disabled or
//it will be deleted in the next update.
//It is not safe to work with a dirty object.
func (o *SceneObject) IsDirty() bool {
	return !o.initialized || !o.enabled || o.willBeRemoved
}

//Copy creates a new scene object with all components and children of the receiver.
//The new object is dirty. It is enabled if the receiver is enabled.
func (o *SceneObject) Copy(copyName string) *SceneObject {
	var newObject *SceneObject
	if instance == nil {
		newObject = NewSceneObjectForTesting(copyName)
	} else {
		newObject = NewSceneObject(copyName)
	}

	for _, comp := range o.components {
		compType := reflect.TypeOf(comp.component)
		var compCopyValue reflect.Value
		if compType.Kind() == reflect.Ptr {
			compCopyValue = reflect.New(compType.Elem())
		} else {
			compCopyValue = reflect.Indirect(reflect.New(compType))
		}

		compCopy := compCopyValue.Interface().(Component)

		if instance != nil {
			cm := instance.GetComponentsManager()
			state := cm.GetComponentState(comp.component)
			cm.SetComponentState(compCopy, state)
		}

		newObject.AddComponent(compCopy)
	}

	newObject.Transform = o.Transform
	newObject.enabled = o.enabled

	for _, child := range o.children {
		childCopy := child.Copy(child.name)
		newObject.AddChild(childCopy)
	}

	return newObject
}

//RemoveFromScene removes the scene object from the current scene.
//After that the object is considered dirty.
func (o *SceneObject) RemoveFromScene() {
	if o.parent == nil || !o.inCurrentScene {
		return
	}

	o.GetParent().RemoveChild(o)
}

//SetParent changes the parent of the scene object to the specified object.
//Can be used to move a child from one object to another.
func (o *SceneObject) SetParent(newParent *SceneObject) {
	if o.parent != nil {
		o.parent.removeChildFromList(o)
	}

	newParent.AddChild(o)
	o.Redraw()
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
			"name":  NameOfComponent(c.GetComponent()),
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
			"name":  NameOfComponent(c.GetComponent()),
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
		"willBeRemoved": o.willBeRemoved,
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
	o.renderingNode = instance.renderer.MakeNode(o)
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
		messageListeners:    make([]*ComponentContainer, 0),
		renderingComponents: make([]*ComponentContainer, 0, 1),
		updatingComponents:  make([]*ComponentContainer, 0, 1),
		boundaryComponents:  make([]*ComponentContainer, 0, 1),
		enabled:             true,
	}
	obj.renderingNode = instance.renderer.MakeNode(obj)
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
		messageListeners:    make([]*ComponentContainer, 0),
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

	ctx := newInitContext(obj)
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