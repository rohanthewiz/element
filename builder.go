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
	b.Ele = func(el string, p ...string) Element {
		return New(b.s, el, p...)
	}
	b.Text = func(p ...string) struct{} {
		return Text(b.s, p...)
	}
	return
}

// WriteString writes directly to the string builder
func (b *Builder) WriteString(s string) (err error) {
	_, err = b.s.WriteString(s)
	return
}

// String returns the accumulated string output
func (b *Builder) String() string {
	return b.s.String()
}

// elementFunc build an element
type elementFunc func(el string, p ...string) Element

// textFunc renders literal text
type textFunc func(p ...string) struct{}
