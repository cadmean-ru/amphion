package builtin

import (
	"github.com/cadmean-ru/amphion/common"
	"github.com/cadmean-ru/amphion/common/a"
	"github.com/cadmean-ru/amphion/common/atext"
	"github.com/cadmean-ru/amphion/common/dispatch"
	"github.com/cadmean-ru/amphion/engine"
	"github.com/cadmean-ru/amphion/frontend"
	"github.com/cadmean-ru/amphion/rendering"
	"math"
	"time"
)

type Selection struct {
	Start, End int
}

func (r *Selection) Move(by int) {
	r.Start += by
	r.End += by
}

func (r *Selection) SetPosition(pos int) {
	r.Start = pos
	r.End = pos
}

func (r *Selection) MoveEnd(by int) {
	if r.End+by < r.Start {
		r.MoveStart(by)
		return
	}

	r.End += by
}

func (r *Selection) MoveStart(by int) {
	if r.Start+by > r.End {
		r.MoveEnd(by)
		return
	}

	r.Start += by
}

func (r *Selection) Length() int {
	return r.End - r.Start
}

func (r *Selection) IndexInSelection(index int) bool {
	return index >= r.Start && index <= r.End
}

type TextInput struct {
	engine.ViewImpl
	CursorColor        a.Color `state:"cursorColor"`
	SelectionColor     a.Color `state:"selectionColor"`
	currentSelection   *Selection
	currentText        []rune
	aFont              *atext.Font
	aFace              *atext.Face
	aText              *atext.Text
	selectionPrimitive int
	cursorPrimitive    int
	prevTransform      engine.Transform
	selectionOffset    a.Vector3
	mousePressed       bool
	mouseDownPosition  a.IntVector2
	shouldDrag         bool
	padding            float32
	backspacePressed   bool
	backspaceTime      time.Time
	currentCursorColor a.Color
	blinkTime          time.Time
	blinkFlag          bool
}

func (s *TextInput) OnInit(ctx engine.InitContext) {
	s.ViewImpl.OnInit(ctx)

	s.currentText = make([]rune, 0, 50)
	s.currentSelection = &Selection{}

	s.aFont, _ = atext.ParseFont(atext.DefaultFontData)
	s.aFace = s.aFont.NewFace(20)
	s.layoutText()

	if !s.SceneObject.HasComponent("RectBoundary") {
		s.SceneObject.AddComponent(NewRectBoundary())
	}

	s.currentCursorColor = s.CursorColor
}

func (s *TextInput) OnStart() {
	engine.BindEventHandler(engine.EventKeyDown, s.handleKeyDown)
	engine.BindEventHandler(engine.EventKeyUp, s.handleKeyUp)
	engine.BindEventHandler(engine.EventMouseDown, s.handleClick)
	engine.BindEventHandler(engine.EventTextInput, s.handleTextInput)
	engine.BindEventHandler(engine.EventMouseUp, s.handleMouseUp)

	s.PrimitiveId = s.SceneObject.GetRenderingNode().AddPrimitive()
	s.selectionPrimitive = s.SceneObject.GetRenderingNode().AddPrimitive()
	s.cursorPrimitive = s.SceneObject.GetRenderingNode().AddPrimitive()
}

func (s *TextInput) OnMessage(msg *dispatch.Message) bool {
	if !engine.IsBuiltinEventMessage(msg) {
		return true
	}

	event := engine.EventFromMessage(msg)
	if event.Sender != s.SceneObject {
		return true
	}

	switch event.Code {
	case engine.EventMouseIn:
		s.handleMouseIn(event)
	case engine.EventMouseOut:
		s.handleMouseOut(event)
	case engine.EventFocusLose:
		s.handleFocusLose(event)
	}

	return false
}

func (s *TextInput) OnUpdate(_ engine.UpdateContext) {
	if !s.SceneObject.IsFocused() {
		return
	}

	s.handleMouseDrag()
	s.handleKeyHold()
	s.handleCursorBlink()

	engine.RequestUpdate()
}

func (s *TextInput) handleMouseDrag() {
	if !s.mousePressed {
		return
	}

	mousePos := engine.GetInputManager().GetCursorPosition()

	if !s.shouldDrag {
		if mousePos.Distance(s.mouseDownPosition) > 5 {
			s.shouldDrag = true
		} else {
			return
		}
	}

	index := s.globalPositionToSelectionIndex(mousePos.ToFloat3())
	if index == -1 {
		return
	}

	if index > s.currentSelection.End {
		s.currentSelection.End = index
	} else if index < s.currentSelection.Start {
		s.currentSelection.Start = index
	}
	s.ShouldRedraw = true
	engine.RequestRendering()
}

func (s *TextInput) handleKeyHold() {
	if s.backspacePressed {
		if time.Since(s.backspaceTime) > time.Second {
			s.backspaceTime = time.Now()
			s.handleBackspace()
			s.ShouldRedraw = true
			engine.RequestRendering()
		}
	}
}

func (s *TextInput) handleCursorBlink() {
	if time.Since(s.blinkTime) < 500 * time.Millisecond {
		return
	}

	s.blinkTime = time.Now()
	s.blinkFlag = !s.blinkFlag
	if s.blinkFlag {
		s.currentCursorColor = a.Transparent()
	} else {
		s.currentCursorColor = s.CursorColor
	}
	s.ShouldRedraw = true
	engine.RequestRendering()
}

func (s *TextInput) OnLateUpdate(_ engine.UpdateContext) {
	if !s.SceneObject.IsFocused() || !s.ShouldDraw() && s.prevTransform.ActualEquals(s.SceneObject.Transform) {
		return
	}

	s.prevTransform = s.SceneObject.Transform

	s.layoutText()
	s.selectionOffset = s.getSelectionOffset()
}

func (s *TextInput) OnDraw(ctx engine.DrawingContext) {
	tp := rendering.NewTextPrimitive(string(s.currentText), s.aText)
	tp.Transform = s.SceneObject.Transform.ToRenderingTransform()
	tp.Appearance = rendering.Appearance{
		FillColor: a.Black(),
	}
	ctx.GetRenderingNode().SetPrimitive(s.PrimitiveId, tp)

	//globalRect := s.SceneObject.Transform.GlobalRect()

	gp := rendering.NewGeometryPrimitive(rendering.PrimitiveRectangle)
	gp.Transform.Position = s.selectionOffset.Round()
	gp.Transform.Size = a.NewIntVector3(2, s.aFace.GetSize(), 0)
	gp.Appearance.FillColor = s.currentCursorColor
	gp.Appearance.StrokeWeight = 0
	ctx.GetRenderingNode().SetPrimitive(s.cursorPrimitive, gp)

	pp := rendering.NewPolygonPrimitive()
	pp.Transform = s.SceneObject.Transform.ToRenderingTransform()
	pp.Vertices = []a.Vector3{}
	pp.Indexes = []uint{}
	z := s.SceneObject.Transform.LocalPosition().Z

	if s.currentSelection.Length() == 0 {
		pp.Vertices = append(pp.Vertices,
			a.NewVector3(s.selectionOffset.X, s.selectionOffset.Y, z),
			a.NewVector3(s.selectionOffset.X, s.selectionOffset.Y+float32(s.aFace.GetSize()), z),
			a.NewVector3(s.selectionOffset.X+2, s.selectionOffset.Y, z),
			a.NewVector3(s.selectionOffset.X+2, s.selectionOffset.Y+float32(s.aFace.GetSize()), z),
		)
		pp.Indexes = append(pp.Indexes,
			0, 1, 2,
			1, 3, 2,
		)
		pp.Appearance.FillColor = a.Transparent()
	} else {
		startChar := s.aText.GetCharAt(s.currentSelection.Start)
		var rectsCount uint = 0

		line := startChar.GetLine().GetIndex()
		prevChar := startChar
		for i := s.currentSelection.Start; i < s.currentSelection.End; i++ {
			c := s.aText.GetCharAt(i)
			if i == s.currentSelection.End-1 || c.GetLine().GetIndex() != line {
				endChar := c
				if c.GetLine().GetIndex() != line &&  i != s.currentSelection.End-1 {
					endChar = prevChar
				}

				startPos := a.NewVector3(float32(startChar.GetPosition().X), float32(startChar.GetLine().GetPosition().Y), 0)
				endPos := a.NewVector3(float32(endChar.GetPosition().X + endChar.GetSize().X), float32(endChar.GetLine().GetY() + s.aFace.GetSize()), 0)
				boundary := common.NewRectBoundaryFromMinMaxPositions(startPos, endPos)

				engine.LogDebug("%v %v", startPos, endPos)

				pp.Vertices = append(pp.Vertices,
					a.NewVector3(boundary.X.Min, boundary.Y.Min, z),
					a.NewVector3(boundary.X.Min, boundary.Y.Max, z),
					a.NewVector3(boundary.X.Max, boundary.Y.Min, z),
					a.NewVector3(boundary.X.Max, boundary.Y.Max, z),
				)

				pp.Indexes = append(pp.Indexes,
					rectsCount*4,
					rectsCount*4+1,
					rectsCount*4+2,
					rectsCount*4+1,
					rectsCount*4+3,
					rectsCount*4+2,
				)

				rectsCount++
				line = c.GetLine().GetIndex()

				if i == s.currentSelection.End-1 {
					break
				} else {
					startChar = c
				}
			}
			prevChar = c
		}

		pp.Appearance.FillColor = s.SelectionColor
	}

	ctx.GetRenderingNode().SetPrimitive(s.selectionPrimitive, pp)

	s.ShouldRedraw = false
}

func (s *TextInput) OnStop() {
	engine.UnbindEventHandler(engine.EventKeyDown, s.handleKeyDown)
	engine.UnbindEventHandler(engine.EventKeyUp, s.handleKeyUp)
	engine.UnbindEventHandler(engine.EventMouseDown, s.handleClick)
	engine.UnbindEventHandler(engine.EventTextInput, s.handleTextInput)
	engine.UnbindEventHandler(engine.EventMouseUp, s.handleMouseUp)

	s.SceneObject.GetRenderingNode().RemovePrimitive(s.PrimitiveId)
	s.SceneObject.GetRenderingNode().RemovePrimitive(s.selectionPrimitive)
	s.SceneObject.GetRenderingNode().RemovePrimitive(s.cursorPrimitive)
}

func (s *TextInput) handleKeyDown(event engine.Event) bool {
	if !s.SceneObject.IsFocused() {
		return true
	}

	data := event.KeyEventData()
	//engine.LogDebug("%v", data.KeyName)

	switch data.KeyName {
	case engine.KeyBackspace:
		s.handleBackspace()
	case engine.KeyDelete:
		s.handleDelete()
	case engine.KeyEnter, engine.KeyNumEnter:
		s.handleRune([]rune("\n"))
	case engine.KeyLeftArrow:
		s.handleArrow(-1)
	case engine.KeyRightArrow:
		s.handleArrow(1)
	case "a":
		s.handleSelectAll()
	default:
	}

	s.ShouldRedraw = true
	engine.RequestRendering()

	return true
}

func (s *TextInput) handleKeyUp(event engine.Event) bool {
	keyEvent := event.KeyEventData()
	switch keyEvent.KeyName {
	case engine.KeyBackspace:
		s.backspacePressed = false
	}

	return true
}

func (s *TextInput) handleBackspace() {
	if s.currentSelection.Start == 0 {
		return
	}

	if s.currentSelection.Length() > 0 {
		s.deleteSelection()
	} else {
		s.currentText = append(s.currentText[:s.currentSelection.Start-1], s.currentText[s.currentSelection.Start:]...)
		s.currentSelection.Move(-1)
	}

	s.backspaceTime = time.Now()
	s.backspacePressed = true
}

func (s *TextInput) handleDelete() {
	if s.currentSelection.End >= len(s.currentText) {
		return
	}

	if s.currentSelection.Length() > 0 {
		s.deleteSelection()
	} else {
		s.currentText = append(s.currentText[:s.currentSelection.End], s.currentText[s.currentSelection.End+1:]...)
	}
}

func (s *TextInput) handleTextInput(event engine.Event) bool {
	s.handleRune([]rune(event.StringData()))
	s.ShouldRedraw = true
	engine.RequestRendering()

	return true
}

func (s *TextInput) handleRune(appendingText []rune) {
	if s.currentSelection.Length() > 0 {
		s.replaceSelection(appendingText)
	} else {
		if s.currentSelection.Start == 0 {
			s.currentText = append(appendingText, s.currentText...)
		} else {
			s.currentText = append(s.currentText[:s.currentSelection.Start], s.currentText[s.currentSelection.Start-1:]...)
			s.currentText[s.currentSelection.Start] = appendingText[0]
		}
		s.currentSelection.Move(1)
	}
}

func (s *TextInput) replaceSelection(newText []rune) {
	s.currentText = append(s.currentText[:s.currentSelection.Start+len(newText)], s.currentText[s.currentSelection.End:]...)
	for i := s.currentSelection.Start; i < s.currentSelection.Start+len(newText); i++ {
		s.currentText[i] = newText[i-s.currentSelection.Start]
	}
	s.currentSelection.SetPosition(s.currentSelection.Start + len(newText))
}

func (s *TextInput) deleteSelection() {
	s.currentText = append(s.currentText[:s.currentSelection.Start], s.currentText[s.currentSelection.End:]...)
	s.currentSelection.SetPosition(s.currentSelection.Start)
}

func (s *TextInput) handleArrow(dir int) {
	if s.currentSelection.Start <= 0 && dir <= -1 || s.currentSelection.End >= len(s.currentText) && dir >= 1 {
		return
	}

	if engine.GetInputManager().IsShiftPressed() {
		s.currentSelection.MoveEnd(dir)
	} else {
		s.currentSelection.Move(dir)
	}
}

func (s *TextInput) handleSelectAll() {
	if !engine.GetInputManager().IsMainCombinationKeyPressed() {
		return
	}

	s.currentSelection.Start = 0
	s.currentSelection.End = len(s.currentText)
}

func (s *TextInput) handleClick(event engine.Event) bool {
	mouseEvent := event.MouseEventData()
	pos := mouseEvent.MousePosition.ToFloat3()
	if mouseEvent.MouseButton != engine.MouseLeft || !s.SceneObject.Transform.GlobalRect().IsPointInside2D(pos) {
		return true
	}

	index := s.globalPositionToSelectionIndex(pos)
	if index != -1 {
		s.currentSelection.SetPosition(index)
		s.ShouldRedraw = true
		engine.RequestRendering()
	}

	s.mousePressed = true
	s.mouseDownPosition = mouseEvent.MousePosition

	return true
}

func (s *TextInput) handleMouseUp(event engine.Event) bool {
	mouseEvent := event.MouseEventData()
	if mouseEvent.MouseButton != engine.MouseLeft {
		return true
	}

	s.mousePressed = false
	s.shouldDrag = false

	return true
}

func (s *TextInput) handleMouseIn(_ engine.Event) bool {
	msg := dispatch.NewMessage(frontend.MessageSetStandardCursor)
	msg.IntData = engine.StandardCursorIBeam

	s.Engine.GetFrontend().GetMessageDispatcher().SendMessage(msg)

	return true
}

func (s *TextInput) handleMouseOut(_ engine.Event) bool {
	msg := dispatch.NewMessage(frontend.MessageSetStandardCursor)
	msg.IntData = engine.StandardCursorArrow

	s.Engine.GetFrontend().GetMessageDispatcher().SendMessage(msg)

	return true
}

func (s *TextInput) handleFocusLose(_ engine.Event) bool {
	s.currentCursorColor = s.CursorColor
	s.blinkFlag = false
	s.ShouldRedraw = true
	engine.RequestRendering()
	return true
}

func (s *TextInput) globalPositionToSelectionIndex(pos a.Vector3) int {
	for i, c := range s.aText.GetAllChars() {
		charRect := common.NewRectBoundaryFromPositionAndSize(c.GetPosition().ToFloat3(), c.GetSize().ToFloat3())
		if charRect.IsPointInside2D(pos) {
			if pos.X < charRect.X.Min+charRect.X.GetLength()/2 {
				return i
			} else {
				return i + 1
			}
		}
	}

	return -1
}

func (s *TextInput) getSelectionOffset() a.Vector3 {
	if s.aText.GetCharsCount() == 0 {
		return a.ZeroVector()
	}

	var y float32 = 0
	i := 0

	for l := 0; l < s.aText.GetLinesCount(); l++ {
		line := s.aText.GetLineAt(l)

		for c := 0; c < line.GetCharsCount(); c++ {
			if i == s.currentSelection.Start-1 {
				char := line.GetCharAt(c)
				return a.NewVector3(float32(char.GetX()+char.GetSize().X), y, 0)
			}

			i++
		}

		y += float32(line.GetSize().Y)
	}

	return a.ZeroVector()
}

func (s *TextInput) layoutText() {
	bounds := s.SceneObject.Transform.GlobalRect()
	bounds.Shrink(a.NewVector3(s.padding, s.padding, 0))
	wantedSize := s.SceneObject.Transform.WantedSize()
	if wantedSize.X == a.WrapContent {
		bounds.X.Max = atext.Unbounded
	}
	if wantedSize.Y == a.WrapContent {
		bounds.Y.Max = atext.Unbounded
	}
	s.aText = atext.LayoutRunes(s.aFace, s.currentText, bounds, atext.LayoutOptions{})
}

func (s *TextInput) calculatePadding() {
	s.padding = float32(math.Ceil(float64(s.aFace.GetSize()) * 0.15))
}

func NewTextInput() *TextInput {
	return &TextInput{
		CursorColor: a.NewColor(0, 0, 255, 255),
		SelectionColor: a.NewColor(0, 0, 255, 100),
	}
}
