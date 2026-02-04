package components

import "github.com/rohanthewiz/element"

// -----------------------------------------------------------------------------
// Nav Component
// -----------------------------------------------------------------------------

// NavItem represents a single navigation link.
type NavItem struct {
	Label  string // Display text
	Href   string // Link URL
	Active bool   // Whether this is the current page
}

// Nav renders a navigation bar with links.
type Nav struct {
	Items []NavItem // Navigation items
	Brand string    // Optional brand/logo text
	Class string    // Additional CSS classes
}

// Render implements the element.Component interface.
func (n Nav) Render(b *element.Builder) (x any) {
	navClass := "nav"
	if n.Class != "" {
		navClass += " " + n.Class
	}

	b.NavClass(navClass).R(
		// Brand
		func() (x any) {
			if n.Brand != "" {
				b.AClass("nav-brand", "href", "#").T(n.Brand)
			}
			return
		}(),
		// Nav items
		b.UlClass("nav-list").R(
			func() (x any) {
				for _, item := range n.Items {
					itemClass := "nav-item"
					if item.Active {
						itemClass += " active"
					}
					b.LiClass(itemClass).R(
						b.AClass("nav-link", "href", item.Href).T(item.Label),
					)
				}
				return
			}(),
		),
	)
	return
}
