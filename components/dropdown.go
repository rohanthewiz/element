package components

import "github.com/rohanthewiz/element"

// -----------------------------------------------------------------------------
// Dropdown Component
// -----------------------------------------------------------------------------

// DropdownItem is one entry in a Dropdown menu.
type DropdownItem struct {
	Label   string // Display text
	Href    string // Link target
	Divider bool   // Render a divider instead of a link
}

// Dropdown renders a disclosure menu using native <details>/<summary> —
// no JavaScript required. The browser toggles it open and closed;
// Esc/outside-click dismissal can be added with a few lines of JS if desired.
//
// Style hooks: .dropdown, .dropdown-toggle, .dropdown-menu, .dropdown-link,
// .dropdown-divider. Position the menu with
// .dropdown{position:relative} .dropdown-menu{position:absolute}.
type Dropdown struct {
	Label string // The toggle button text
	Items []DropdownItem
	Class string // Additional CSS classes for the wrapper
}

// Render implements the element.Component interface.
func (d Dropdown) Render(b *element.Builder) (x any) {
	wrapClass := "dropdown"
	if d.Class != "" {
		wrapClass += " " + d.Class
	}

	b.DetailsClass(wrapClass).R(
		b.SummaryClass("dropdown-toggle").T(d.Label),
		b.UlClass("dropdown-menu").R(
			element.ForEach(d.Items, func(item DropdownItem) {
				if item.Divider {
					b.LiClass("dropdown-divider", "role", "separator").R()
				} else {
					b.Li().R(
						b.AClass("dropdown-link", "href", item.Href).T(item.Label),
					)
				}
			}),
		),
	)
	return
}
