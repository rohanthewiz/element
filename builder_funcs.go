package element

import "fmt"

// elementFunc build an element
type elementFunc func(el string, attrPairs ...string) Element

// textFunc renders literal text
type textFunc func(attrPairs ...string) any

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
func (b *Builder) F(format string, args ...any) (x any) {
	_ = b.WriteString(fmt.Sprintf(format, args...))
	return
}

// T renders a list of strings directly to the builder.
// It is the fastest way to render text
func (b *Builder) T(strs ...string) (x any) {
	for _, str := range strs {
		_ = b.WriteString(str)
	}
	return
}

// Wrap allows some Go code inside a render tree, so some logic can be performed in the process of rendering
// Example:
//
//		b, _, t := Vars()
//	 	isEvening := true
//
//		b.Div().R(
//			b.H2().T("Testing some things"),
//			b.Wrap(func() {
//				if isEvening {
//					b.H3().T("Good evening!")
//				} else {
//					b.H3().T("Good day!")
//				}
//			})
//		)
func (b *Builder) Wrap(fn func()) (x any) {
	fn()
	return
}

// RenderComps renders a list of components inside a render tree
func (b *Builder) RenderComps(bc ...Component) (x any) {
	for _, comp := range bc {
		comp.Render(b)
	}
	return
}
