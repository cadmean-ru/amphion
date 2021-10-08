package engine

import (
	"github.com/cadmean-ru/amphion/common/a"
	"github.com/cadmean-ru/amphion/common/dispatch"
	"sort"
)

type MouseButton byte

const (
	MouseUnknown MouseButton = iota
	MouseLeft
	MouseRight
	MouseMiddle
	MouseBack
	MouseForward
)

// MouseEventData represents the data of mouse related event.
// Contains mouse coordinates, SceneObject, that was clicked and the mouse button.
type MouseEventData struct {
	SceneObject   *SceneObject
	MousePosition a.IntVector2
	MouseButton   MouseButton
}

func (engine *AmphionEngine) handleClickEvent(mouseEventData MouseEventData, code int) {
	if engine.currentScene == nil {
		return
	}

	clickPos := mouseEventData.MousePosition.ToFloat3()
	candidates := make([]*SceneObject, 0, 1)

	engine.currentScene.TraversePreOrder(func(object *SceneObject) bool {
		if !object.Transform.GlobalRect().IsPointInside2D(clickPos) {
			return false
		}

		if object.HasBoundary() {
			if object.IsPointInsideSolidBoundaries2D(clickPos) {
				candidates = append(candidates, object)
				return true
			}

			return object.IsPointInsideBoundaries(clickPos)
		}

		return true
	})

	if len(candidates) > 0 {
		sort.Slice(candidates, func(i, j int) bool {
			return candidates[i].Transform.GlobalPosition().Z > candidates[j].Transform.GlobalPosition().Z
		})
		o := candidates[0]
		mouseEventData.SceneObject = o

		if engine.sceneContext.focusedObject != nil {
			engine.sceneContext.messageDispatcher.DispatchDirectly(
				engine.sceneContext.focusedObject,
				dispatch.NewMessageFromWithAnyData(
					engine,
					MessageBuiltinEvent,
					NewAmphionEvent(engine.sceneContext.focusedObject, EventFocusLose, nil),
				),
			)
		}
		engine.sceneContext.messageDispatcher.DispatchDirectly(o, dispatch.NewMessageFromWithAnyData(o, MessageBuiltinEvent, NewAmphionEvent(o, code, mouseEventData)))
		engine.sceneContext.focusedObject = o
		engine.sceneContext.messageDispatcher.DispatchDirectly(o, dispatch.NewMessageFromWithAnyData(o, MessageBuiltinEvent, NewAmphionEvent(o, EventFocusGain, nil)))

		event := NewAmphionEvent(engine, code, mouseEventData)
		engine.updateRoutine.enqueueEventAndRequestUpdate(event)
	} else {
		if engine.sceneContext.focusedObject != nil {
			engine.sceneContext.messageDispatcher.DispatchDirectly(
				engine.sceneContext.focusedObject,
				dispatch.NewMessageFromWithAnyData(
					engine,
					MessageBuiltinEvent,
					NewAmphionEvent(engine.sceneContext.focusedObject, EventFocusLose, nil),
				),
			)
		}
		engine.sceneContext.focusedObject = nil
		event := NewAmphionEvent(engine, code, mouseEventData)
		engine.updateRoutine.enqueueEventAndRequestUpdate(event)
	}
}

func (engine *AmphionEngine) handleMouseMove(mousePos a.IntVector2) {
	if engine.currentScene == nil {
		return
	}

	var candidate *SceneObject
	mousePos3 := mousePos.ToFloat3()

	engine.currentScene.TraversePreOrder(func(object *SceneObject) bool {
		if !object.Transform.GlobalRect().IsPointInside2D(mousePos3) {
			return false
		}

		if object.HasBoundary() {
			if object.IsPointInsideBoundaries2D(mousePos3) {
				candidate = object
				return true
			}

			return false
		}

		return true
	})

	if candidate != nil {
		o := candidate

		if o == engine.sceneContext.hoveredObject {
			return
		}

		if engine.sceneContext.hoveredObject != nil {
			engine.sceneContext.messageDispatcher.DispatchDirectly(
				engine.sceneContext.hoveredObject,
				dispatch.NewMessageFromWithAnyData(
					engine.sceneContext.hoveredObject,
					MessageBuiltinEvent,
					NewAmphionEvent(engine.sceneContext.hoveredObject, EventMouseOut, nil),
				),
			)
		}

		engine.sceneContext.messageDispatcher.DispatchDirectly(
			o,
			dispatch.NewMessageFromWithAnyData(
				o,
				MessageBuiltinEvent,
				NewAmphionEvent(o, EventMouseIn, nil),
			),
		)

		engine.sceneContext.hoveredObject = o
	} else {
		if engine.sceneContext.hoveredObject != nil {
			engine.sceneContext.messageDispatcher.DispatchDirectly(
				engine.sceneContext.hoveredObject,
				dispatch.NewMessageFromWithAnyData(
					engine.sceneContext.hoveredObject,
					MessageBuiltinEvent,
					NewAmphionEvent(engine.sceneContext.hoveredObject, EventMouseOut, nil),
				),
			)

			engine.sceneContext.hoveredObject = nil
		}
	}
}