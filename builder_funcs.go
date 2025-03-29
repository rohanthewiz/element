package element

import "fmt"

// elementFunc build an element
type elementFunc func(el string, attrPairs ...string) Element

// textFunc renders literal text
type textFunc func(attrPairs ...string) struct{}

// BUILDER CONVENIENCE METHODS

// Funcs is a convenience method for returning
// the builder methods b.Ele and b.Text
// that are used for creating elements and literal text
func (b *Builder) Funcs() (ele elementFunc, text textFunc) {
	return b.Ele, b.Text
}

// Vars is a convenience method for returning
// the builder methods b.Ele and b.Text
// that are used for creating elements and literal text
func (b *Builder) Vars() (ele elementFunc, text textFunc) {
	return b.Ele, b.Text
}

// V is a convenience method for returning
// the builder methods b.Ele and b.Text
// that are used for creating elements and literal text
func (b *Builder) V() (ele elementFunc, text textFunc) {
	return b.Ele, b.Text
}

// F formats and writes a string to the builder.
// It's a convenience method for fmt.Sprintf that writes directly to the builder.
//
// Example:
//
//	b := NewBuilder()
//	b.F("Hello, %s!", "world") // Writes "Hello, world!" to the builder
//
// Returns an empty interface for method chaining.
func (b *Builder) F(formatString string, args ...any) (x any) {
	b.WriteString(fmt.Sprintf(formatString, args...))
	return
}

// ELEMENT CONVENIENCE FUNCTIONS

// Vars returns a builder plus it's convenience methods for creating elements and text
func Vars() (b *Builder, e elementFunc, t textFunc) {
	b = NewBuilder()
	return b, b.Ele, b.Text
}

// V is a short form of Vars which returns a builder and it's convenience methods
func V() (b *Builder, e elementFunc, t textFunc) {
	return Vars()
}
