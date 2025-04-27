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

// ForEach is a renderable function for performing an operation on a list of generic items.
// Note: breaking change - no need to pass the builder as this is meant for a render tree
// where an instance of builder is already available
// Example:
//
//	items := []string{"item1", "item2"}
//	ForEach(items, func(item string) {
//		b.P().T(item)
//	})
func ForEach[T any](items []T, each func(item T)) (x any) {
	for _, itm := range items {
		each(itm)
	}
	return
}
