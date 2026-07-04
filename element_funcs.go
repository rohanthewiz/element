package element

// ELEMENT CONVENIENCE FUNCTIONS

// Vars returns a builder plus it's convenience methods for creating elements and text.
// We are deprecating this. Just use the builder for all the things, so prefer NewBuilder() or B()
func Vars() (b *Builder, e elementFunc, t textFunc) {
	b = NewBuilder()
	return b, b.Ele, b.Text
}

// V is a short form of Vars which returns a builder and it's convenience methods
// We are deprecating this. Just use the builder for all the things, so prefer NewBuilder() or B()
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

// ForEach2 is like ForEach, but also provides the item's index.
// Example:
//
//	items := []string{"item1", "item2"}
//	ForEach2(items, func(item string, i int) {
//		b.P().F("%d. %s", i+1, item)
//	})
func ForEach2[T any](items []T, each func(item T, index int)) (x any) {
	for i, itm := range items {
		each(itm, i)
	}
	return
}
