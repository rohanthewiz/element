package element

// Component defines a function type for rendering a group of HTML elements
//
//	It takes an element.Builder instance  but returns nothing  since  it writes directly to the builder
type Component func(b *Builder, components ...Component)

// RenderComponents is used to render a list of components
func RenderComponents(b *Builder, comps ...Component) (x any) {
	for _, comp := range comps {
		comp(b)
	}
	return
}
