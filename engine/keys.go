package engine

type KeyName string

const (
	KeyUnknown      KeyName = "Unknown"
	KeyLeftShift    KeyName = "LeftShift"
	KeyRightShift   KeyName = "RightShift"
	KeyLeftAlt      KeyName = "LeftAlt"
	KeyRightAlt     KeyName = "RightAlt"
	KeyLeftSuper    KeyName = "LeftSuper"
	KeyRightSuper   KeyName = "RightSuper"
	KeyLeftControl  KeyName = "LeftControl"
	KeyRightControl KeyName = "RightControl"
	KeySpace        KeyName = "Space"
	KeyEnter        KeyName = "Enter"
	KeyNumEnter     KeyName = "NumEnter"
	KeyBackspace    KeyName = "Backspace"
	KeyLeftArrow    KeyName = "LeftArrow"
	KeyRightArrow   KeyName = "RightArrow"
	KeyUpArrow      KeyName = "UpArrow"
	KeyDownArrow    KeyName = "DownArrow"
	KeyInsert       KeyName = "Insert"
	KeyDelete       KeyName = "Delete"
	KeyHome         KeyName = "Home"
	KeyEnd          KeyName = "End"
	KeyPageUp       KeyName = "PageUp"
	KeyPageDown     KeyName = "PageDown"
	KeyCapsLock     KeyName = "CapsLock"
	KeyTab          KeyName = "Tab"
	KeyEscape       KeyName = "Escape"
	KeyFn           KeyName = "Fn"
)

type KeyEventData struct {
	KeyName KeyName
}