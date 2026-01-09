package element

import (
	"bytes"
)

// Builder has everything we need for building and rendering HTML
type Builder struct {
	Ele  elementFunc   // function for building an element
	Text textFunc      // function for rendering literal text
	s    *bytes.Buffer // accumulates our output
}

// NewBuilder returns a builder for creating our HTML
func NewBuilder() (b *Builder) {
	b = &Builder{}
	b.s = bytes.NewBuffer(make([]byte, 0, 256))

	b.Ele = func(el string, attrPairs ...string) Element {
		return New(b.s, el, attrPairs...)
	}
	b.Text = func(attrPairs ...string) (x any) {
		return Text(b.s, attrPairs...)
	}
	return
}

// B is a convenience function for NewBuilder
func B() (b *Builder) {
	return NewBuilder()
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

// Bytes returns the accumulated buffer as bytes
func (b *Builder) Bytes() []byte {
	return b.s.Bytes()
}

// Reset clears the internal bytes.Buffer
func (b *Builder) Reset() {
	b.s.Reset()
}

// Pretty returns the HTML with pretty formatting (indentation and line breaks)
// This is useful for debugging or when you need human-readable HTML output
func (b *Builder) Pretty() string {
	return PrettyHTML(b.s.String())
}
