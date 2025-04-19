package element

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

// ForEach is a convenience method for performing an operation on a list of items of any type
func ForEach[T any](b *Builder, items []T, each func(b *Builder, item T)) (x any) {
	for _, itm := range items {
		each(b, itm)
	}
	return
}
