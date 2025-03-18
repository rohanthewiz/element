package element

// elementFunc build an element
type elementFunc func(el string, attrPairs ...string) Element

// textFunc renders literal text
type textFunc func(attrPairs ...string) struct{}

// CONVENIENCE METHODS

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

// CONVENIENCE FUNCTIONS

// Vars returns a builder plus it's convenience methods for creating elements and text
func Vars() (b *Builder, e elementFunc, t textFunc) {
	b = NewBuilder()
	return b, b.Ele, b.Text
}

// V is a short form of Vars which returns a builder and it's convenience methods
func V() (b *Builder, e elementFunc, t textFunc) {
	return Vars()
}
