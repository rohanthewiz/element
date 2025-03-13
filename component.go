package element

// Component defines an interface for rendering a group of HTML elements
//
//	It takes an element.Builder instance  but returns nothing  since  it writes directly to the builder
/*type Component func(b *Builder, components ...Component)
 */

// Component is some group of elements that can be rendered in a tree of elements
type Component interface {
	// Render writes to the builder,  but returns anything so it can be rendered in a tree of elements
	Render(b *Builder) (x any)
}

// RenderComponents is a convenience function for rendering a list of components in a tree of elements
func RenderComponents(b *Builder, comps ...Component) (x any) {
	for _, comp := range comps {
		comp.Render(b)
	}
	return
}
