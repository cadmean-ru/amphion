package builtin

import (
	"github.com/cadmean-ru/amphion/common"
	"github.com/cadmean-ru/amphion/common/a"
	"github.com/cadmean-ru/amphion/common/atext"
	"github.com/cadmean-ru/amphion/common/dispatch"
	"github.com/cadmean-ru/amphion/engine"
	"github.com/cadmean-ru/amphion/frontend"
	"github.com/cadmean-ru/amphion/rendering"
	"time"
)

//region Selection

type Selection struct {
	Start, End int
	Direction  byte
}

func (r *Selection) Move(by int) {
	r.Start += by
	r.End += by
	r.Direction = 0
}

func (r *Selection) SetPosition(pos int) {
	r.Start = pos
	r.End = pos
	r.Direction = 0
}

func (r *Selection) MoveEnd(by int) {
	if r.End+by < r.Start {
		r.MoveStart(by)
		return
	}

	r.End += by
	r.Direction = 2
}

func (r *Selection) MoveStart(by int) {
	if r.Start+by > r.End {
		r.MoveEnd(by)
		return
	}

	r.Start += by
	r.Direction = 1
}

func (r *Selection) Length() int {
	return r.End - r.Start
}

func (r *Selection) IndexInSelection(index int) bool {
	return index >= r.Start && index <= r.End
}

//endregion

type TextInput struct {
	engine.ViewImpl
	CursorColor        a.Color             `state:"cursorColor"`
	SelectionColor     a.Color             `state:"selectionColor"`
	InitialText        string              `state:"initialText"`
	OnChangeListener   engine.EventHandler `state:"onChangeListener"`
	FontSize           byte                `state:"fontSize"`
	HTextAlign         a.TextAlign         `state:"hTextAlign"`
	VTextAlign         a.TextAlign         `state:"vTextAlign"`
	SingleLine         bool                `state:"singleLine"`
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
	backspacePressed   bool
	backspaceTime      time.Time
	currentCursorColor a.Color
	blinkTime          time.Time
	blinkFlag          bool
	firstFrame         bool
}

func (s *TextInput) OnInit(ctx engine.InitContext) {
	s.ViewImpl.OnInit(ctx)

	s.currentText = make([]rune, 0, 50)
	s.currentSelection = &Selection{}

	s.aFont, _ = atext.ParseFont(atext.DefaultFontData)
	if s.FontSize == 0 {
		s.FontSize = 14
	}
	s.aFace = s.aFont.NewFace(int(s.FontSize))
	s.layoutText()

	if !s.SceneObject.HasComponent("Rect") {
		s.SceneObject.AddComponent(NewRectBoundary())
	}

	s.currentCursorColor = s.CursorColor

	s.currentText = []rune(s.InitialText)

	s.firstFrame = true
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

//region Update

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
		s.currentSelection.Direction = 2
	} else if index < s.currentSelection.Start {
		s.currentSelection.Start = index
		s.currentSelection.Direction = 1
	}
	s.ShouldRedraw = true
	engine.RequestRendering()
}

func (s *TextInput) handleKeyHold() {
	if s.backspacePressed {
		if time.Since(s.backspaceTime) > time.Millisecond*100 {
			s.backspaceTime = time.Now()
			s.handleBackspace()
			s.ShouldRedraw = true
			engine.RequestRendering()
		}
	}
}

func (s *TextInput) handleCursorBlink() {
	if time.Since(s.blinkTime) < 500*time.Millisecond {
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

//endregion

func (s *TextInput) OnLateUpdate(_ engine.UpdateContext) {
	if !s.firstFrame && (!s.SceneObject.IsFocused() || !s.ShouldDraw() && s.prevTransform.ActualEquals(s.SceneObject.Transform)) {
		return
	}

	s.firstFrame = false
	s.prevTransform = s.SceneObject.Transform

	s.layoutText()
	s.selectionOffset = s.getSelectionOffset()
}

func (s *TextInput) OnDraw(ctx engine.DrawingContext) {
	rt := s.SceneObject.Transform.ToRenderingTransform()

	tp := rendering.NewTextPrimitive(string(s.currentText), s.aText)
	tp.Transform = rt
	tp.Appearance = rendering.Appearance{
		FillColor: a.Black(),
	}
	ctx.GetRenderingNode().SetPrimitive(s.PrimitiveId, tp)

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
				if c.GetLine().GetIndex() != line && i != s.currentSelection.End-1 {
					endChar = prevChar
				}

				startPos := a.NewVector3(float32(startChar.GetPosition().X), float32(startChar.GetLine().GetPosition().Y), 0)
				endPos := a.NewVector3(float32(endChar.GetPosition().X+endChar.GetSize().X), float32(endChar.GetLine().GetY()+s.aFace.GetSize()), 0)
				boundary := common.NewRectFromMinMaxPositions(startPos, endPos)

				//engine.LogDebug("%v %v", startPos, endPos)

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

func (s *TextInput) MeasureContents() a.Vector3 {
	if s.aText == nil {
		return a.ZeroVector()
	}
	return s.aText.GetSize().ToFloat3()
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
		s.handleEnter()
	case engine.KeyLeftArrow:
		s.handleArrow(-1)
	case engine.KeyRightArrow:
		s.handleArrow(1)
	case "a":
		s.handleSelectAll()
	case "c":
		s.handleCopy()
	case "v":
		s.handlePaste()
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
	if s.currentSelection.Start == 0 && s.currentSelection.Length() == 0 {
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

	s.handleTextChange()
}

func (s *TextInput) handleDelete() {
	if s.currentSelection.End >= len(s.currentText) && s.currentSelection.Length() == 0 {
		return
	}

	if s.currentSelection.Length() > 0 {
		s.deleteSelection()
	} else {
		s.currentText = append(s.currentText[:s.currentSelection.End], s.currentText[s.currentSelection.End+1:]...)
	}

	s.handleTextChange()
}

func (s *TextInput) handleTextInput(event engine.Event) bool {
	s.handleRune([]rune(event.StringData()))
	s.ShouldRedraw = true
	engine.RequestRendering()

	return true
}

func (s *TextInput) handleEnter() {
	if s.SingleLine {
		return
	}

	s.handleRune([]rune("\n"))
}

func (s *TextInput) handleRune(appendingText []rune) {
	s.replaceSelection(appendingText)
	s.handleTextChange()
}

func (s *TextInput) replaceSelection(newText []rune) {
	if s.currentSelection.Start+len(newText) > len(s.currentText) {
		s.currentText = append(s.currentText, make([]rune, s.currentSelection.Start+len(newText)-len(s.currentText))...)
	}
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
		if s.currentSelection.Direction == 1 {
			s.currentSelection.MoveStart(dir)
		} else {
			s.currentSelection.MoveEnd(dir)
		}
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

func (s *TextInput) handleCopy() {
	if !engine.GetInputManager().IsMainCombinationKeyPressed() || s.currentSelection.Length() == 0 {
		return
	}

	text := s.GetSelectedText()

	cm := engine.GetFeaturesManager().GetFeature(engine.FeatureClipboardManager).(*engine.ClipboardManager)
	entry := engine.NewClipboardEntry(engine.ClipboardEntryString, []byte(text))
	cm.Write(entry)
}

func (s *TextInput) handlePaste() {
	if !engine.GetInputManager().IsMainCombinationKeyPressed() {
		return
	}

	cm := engine.GetFeaturesManager().GetFeature(engine.FeatureClipboardManager).(*engine.ClipboardManager)
	entry := cm.Read(engine.ClipboardEntryString)
	if entry == nil {
		return
	}

	s.handleRune([]rune(string(entry.Data())))
}

func (s *TextInput) handleTextChange() {
	s.InitialText = string(s.currentText)

	if s.OnChangeListener == nil {
		return
	}

	s.OnChangeListener(engine.NewAmphionEvent(s, 0, s.InitialText))
}

func (s *TextInput) globalPositionToSelectionIndex(pos a.Vector3) int {
	for i, c := range s.aText.GetAllChars() {
		charRect := common.NewRectFromPositionAndSize(c.GetPosition().ToFloat3(), c.GetSize().ToFloat3())
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
		return s.SceneObject.Transform.GlobalTopLeftPosition()
	}

	var y = s.SceneObject.Transform.GlobalTopLeftPosition().Y
	i := 0

	for l := 0; l < s.aText.GetLinesCount(); l++ {
		line := s.aText.GetLineAt(l)

		for c := 0; c < line.GetCharsCount(); c++ {
			if s.currentSelection.Direction < 2 && i == s.currentSelection.Start-1 || s.currentSelection.Direction == 2 && i == s.currentSelection.End-1 {
				char := line.GetCharAt(c)
				return a.NewVector3(float32(char.GetX()+char.GetSize().X), y, 0)
			}

			i++
		}

		y += float32(line.GetSize().Y)
	}

	return s.SceneObject.Transform.GlobalTopLeftPosition()
}

func (s *TextInput) layoutText() {
	bounds := s.SceneObject.Transform.GlobalRect()
	wantedSize := s.SceneObject.Transform.WantedSize()
	if wantedSize.X == a.WrapContent {
		bounds.X.Max = atext.Unbounded
	}
	if wantedSize.Y == a.WrapContent {
		bounds.Y.Max = atext.Unbounded
	}
	s.aText = atext.LayoutRunes(s.aFace, s.currentText, bounds, atext.LayoutOptions{
		HTextAlign: s.HTextAlign,
		VTextAlign: s.VTextAlign,
		SingleLine: s.SingleLine,
	})
}

func (s *TextInput) GetSelectedText() string {
	if s.currentSelection.Length() == 0 {
		return ""
	}

	selected := s.currentText[s.currentSelection.Start:s.currentSelection.End]
	return string(selected)
}

func (s *TextInput) SetOnChangeListener(listener func(newValue string)) {
	s.OnChangeListener = func(event engine.Event) bool {
		listener(event.StringData())
		return true
	}
}

func (s *TextInput) GetText() string {
	return string(s.currentText)
}

func (s *TextInput) SetText(newText string) {
	s.currentText = []rune(newText)
	s.handleTextChange()
	s.ShouldRedraw = true
	engine.RequestRendering()
}

func (s *TextInput) SetFontSize(fontSize byte) {
	s.FontSize = fontSize
	s.aFace = s.aFont.NewFace(int(s.FontSize))
	s.ShouldRedraw = true
	engine.RequestRendering()
}

func (s *TextInput) SetHTextAlign(align a.TextAlign) {
	s.HTextAlign = align
	s.ShouldRedraw = true
	engine.RequestRendering()
}

func (s *TextInput) SetVTextAlign(align a.TextAlign) {
	s.VTextAlign = align
	s.ShouldRedraw = true
	engine.RequestRendering()
}

func (s *TextInput) SetSingleLine(singleLine bool) {
	s.SingleLine = singleLine
	s.ShouldRedraw = true
	engine.RequestRendering()
}

func NewTextInput() *TextInput {
	return &TextInput{
		CursorColor:    a.NewColor(0, 0, 255, 255),
		SelectionColor: a.NewColor(0, 0, 255, 50),
	}
}
