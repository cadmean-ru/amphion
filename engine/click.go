package engine

import (
	"github.com/cadmean-ru/amphion/common/a"
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

	clickPos := mouseEventData.MousePosition
	candidates := make([]*SceneObject, 0, 1)

	engine.currentScene.ForEachObject(func(o *SceneObject) {
		if o.HasBoundary() && o.IsPointInsideBoundaries2D(a.NewVector3(float32(clickPos.X), float32(clickPos.Y), 0)) {
			candidates = append(candidates, o)
		}
	})

	if len(candidates) > 0 {
		sort.Slice(candidates, func(i, j int) bool {
			return candidates[i].Transform.GlobalPosition().Z > candidates[j].Transform.GlobalPosition().Z
		})
		o := candidates[0]
		mouseEventData.SceneObject = o

		if engine.sceneContext.focusedObject != nil {
			engine.messageDispatcher.DispatchDirectly(
				engine.sceneContext.focusedObject,
				NewMessage(
					engine,
					MessageBuiltinEvent,
					NewAmphionEvent(engine.sceneContext.focusedObject, EventFocusLoose, nil),
				),
			)
		}
		engine.messageDispatcher.DispatchDirectly(o, NewMessage(o, MessageBuiltinEvent, NewAmphionEvent(o, code, mouseEventData)))
		engine.sceneContext.focusedObject = o
		engine.messageDispatcher.DispatchDirectly(o, NewMessage(o, MessageBuiltinEvent, NewAmphionEvent(o, EventFocusGain, nil)))

		event := NewAmphionEvent(engine, code, mouseEventData)
		engine.updateRoutine.enqueueEventAndRequestUpdate(event)
	} else {
		if engine.sceneContext.focusedObject != nil {
			engine.messageDispatcher.DispatchDirectly(
				engine.sceneContext.focusedObject,
				NewMessage(
					engine,
					MessageBuiltinEvent,
					NewAmphionEvent(engine.sceneContext.focusedObject, EventFocusLoose, nil),
				),
			)
		}
		engine.sceneContext.focusedObject = nil
		event := NewAmphionEvent(engine, code, mouseEventData)
		engine.updateRoutine.enqueueEventAndRequestUpdate(event)
	}
}