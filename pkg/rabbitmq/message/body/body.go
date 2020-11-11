package body

// Body represents the message body interface
type Body interface {
	String() string
	Raw() []byte
}

// body represents a body
type body struct {
	raw []byte
}

// New creates a new body
func New(raw []byte) Body {
	return &body{
		raw: raw,
	}
}

// Raw returns the raw body
func (b *body) Raw() []byte {
	return b.raw
}

// String returns the stringified body
func (b *body) String() string {
	return string(b.Raw())
}
