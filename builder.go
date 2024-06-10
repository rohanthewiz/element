package element

import (
	"strings"
)

// Builder has everything we need for building and rendering HTML
type Builder struct {
	Ele         elementFunc      // function for building an element
	E           elementFunc      // shorter version of Ele
	EleNoRender elementFunc      // function for building an element without  automatically rendering an open tag
	ENR         elementFunc      // shorter version of EleNoReader
	Text        textFunc         // function for rendering literal text
	T           textFunc         // shorter version of Text
	sb          *strings.Builder // accumulates our output
}

// NewBuilder returns a builder for creating our HTML
func NewBuilder() (b *Builder) {
	b = &Builder{}
	b.sb = &strings.Builder{}

	b.Ele = func(el string, attrPairs ...string) Element {
		return New(b.sb, el, attrPairs...)
	}
	b.E = b.Ele

	b.EleNoRender = func(el string, attrPairs ...string) Element {
		return NewNoRender(b.sb, el, attrPairs...)
	}
	b.ENR = b.EleNoRender

	b.Text = func(attrPairs ...string) struct{} {
		return Text(b.sb, attrPairs...)
	}
	b.T = b.Text

	return
}

// NB is alias method for NewBuilder which returns a builder for composing our HTML
func NB() *Builder {
	return NewBuilder()
}

// WriteString writes directly to the string builder
func (b *Builder) WriteString(s string) (err error) {
	_, err = b.sb.WriteString(s)
	return
}

// WS is an alias method for WriteString which writes directly to the string builder
func (b *Builder) WS(s string) (err error) {
	_, err = b.sb.WriteString(s)
	return
}

// WriteBytes writes bytes directly to the string builder
func (b *Builder) WriteBytes(byts []byte) (err error) {
	_, _ = b.sb.Write(byts)
	return
}

// S is effectively an alias for Builder.String() which returns the accumulated string output
func (b *Builder) S() string {
	return b.sb.String()
}

// String returns the accumulated string output
func (b *Builder) String() string {
	return b.sb.String()
}

// elementFunc build an element
type elementFunc func(el string, attrPairs ...string) Element

// textFunc renders literal text
type textFunc func(attrPairs ...string) struct{}
