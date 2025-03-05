package element

// Component defines a function type for rendering a group of HTML elements
//
//	It takes an element.Builder instance  but returns nothing  since  it writes directly to the builder
type Component func(b *Builder, components ...Component)
