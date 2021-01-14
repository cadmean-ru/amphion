package engine

// Experimental
// Holds an entry of clipboard
type ClipboardData struct {
	mime  string
	kind  string
	str   string
	bytes []byte
}

func (c *ClipboardData) GetMime() string {
	return c.mime
}

func (c *ClipboardData) GetKind() string {
	return c.kind
}

func (c *ClipboardData) GetString() string {
	if c.kind == "string" {
		return c.str
	}
	return ""
}

func (c *ClipboardData) GetData() []byte {
	if c.kind == "file" {
		return c.bytes
	}
	return nil
}

func NewClipboardData(kind, mime string, str string, bytes []byte) *ClipboardData {
	return &ClipboardData{
		mime:  mime,
		kind:  kind,
		str:   str,
		bytes: bytes,
	}
}