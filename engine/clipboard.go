package engine

import (
	"golang.design/x/clipboard"
)

type ClipboardEntryKind int

const (
	ClipboardEntryString ClipboardEntryKind = iota
	ClipBoardEntryImage
)

// ClipboardEntry holds an entry of the clipboard.
type ClipboardEntry struct {
	kind  ClipboardEntryKind
	bytes []byte
}

func (c *ClipboardEntry) Kind() ClipboardEntryKind {
	return c.kind
}

func (c *ClipboardEntry) String() string {
	return string(c.bytes)
}

func (c *ClipboardEntry) Data() []byte {
	return c.bytes
}

//NewClipboardEntry creates a new instance of ClipboardEntry with the given kind and data.
func NewClipboardEntry(kind ClipboardEntryKind, bytes []byte) *ClipboardEntry {
	return &ClipboardEntry{
		kind:  kind,
		bytes: bytes,
	}
}

//ClipboardManager implements working with the system clipboard.
//It supports strings and images.
type ClipboardManager struct {

}

//Write writes the given ClipboardEntry to the system clipboard.
func (m *ClipboardManager) Write(entry *ClipboardEntry) {
	var format = clipboard.FmtText
	if entry.kind == ClipBoardEntryImage {
		format = clipboard.FmtImage
	}
	clipboard.Write(format, entry.bytes)
}

//Read reads data of the given kind from the clipboard.
//If no data of the given kind is present returns nil.
func (m *ClipboardManager) Read(kind ClipboardEntryKind) *ClipboardEntry {
	var format = clipboard.FmtText
	if kind == ClipBoardEntryImage {
		format = clipboard.FmtImage
	}
	bytes := clipboard.Read(format)
	if bytes == nil {
		return nil
	}
	return NewClipboardEntry(kind, bytes)
}

func newClipBoardManager() *ClipboardManager {
	return &ClipboardManager{
	}
}