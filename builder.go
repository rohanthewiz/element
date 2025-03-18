package element

import (
	"strings"
)

// Builder has everything we need for building and rendering HTML
type Builder struct {
	Ele  elementFunc      // function for building an element
	Text textFunc         // function for rendering literal text
	s    *strings.Builder // accumulates our output
}

// NewBuilder returns a builder for creating our HTML
func NewBuilder() (b *Builder) {
	b = &Builder{}
	b.s = &strings.Builder{}

	b.Ele = func(el string, attrPairs ...string) Element {
		return New(b.s, el, attrPairs...)
	}
	b.Text = func(attrPairs ...string) struct{} {
		return Text(b.s, attrPairs...)
	}
	return
}

// Input / 0utput

// WriteString writes directly to the string builder
func (b *Builder) WriteString(s string) (err error) {
	_, err = b.s.WriteString(s)
	return
}

// WriteBytes writes bytes directly to the string builder
func (b *Builder) WriteBytes(byts []byte) (err error) {
	_, _ = b.s.Write(byts)
	return
}

// String returns the accumulated string output
func (b *Builder) String() string {
	return b.s.String()
}
