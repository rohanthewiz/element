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

// Funcs is a convenience method for returning
// the builder functions b.Ele and b.Text
// that are used for creating elements and literal text
func (b *Builder) Funcs() (ele elementFunc, text textFunc) {
	return b.Ele, b.Text
}

// EleFuncs is another convenience method for returning
// the builder functions b.Ele and b.Text
// that are used for creating elements and literal text
func (b *Builder) EleFuncs() (ele elementFunc, text textFunc) {
	return b.Ele, b.Text
}

// Fn is another convenience method for returning
// the builder functions b.Ele and b.Text
// that are used for creating elements and literal text
func (b *Builder) Fn() (ele elementFunc, text textFunc) {
	return b.Ele, b.Text
}

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
