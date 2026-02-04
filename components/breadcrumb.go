package components

import "github.com/rohanthewiz/element"

// -----------------------------------------------------------------------------
// Breadcrumb Component
// -----------------------------------------------------------------------------

// BreadcrumbItem represents a single breadcrumb link.
type BreadcrumbItem struct {
	Label string // Display text
	Href  string // Link URL (empty for current page)
}

// Breadcrumb renders a navigation breadcrumb trail.
type Breadcrumb struct {
	Items     []BreadcrumbItem // Breadcrumb items (last is current)
	Separator string           // Separator character (default: "/")
}

// Render implements the element.Component interface.
func (bc Breadcrumb) Render(b *element.Builder) (x any) {
	separator := bc.Separator
	if separator == "" {
		separator = "/"
	}

	b.NavClass("breadcrumb", "aria-label", "breadcrumb").R(
		b.OlClass("breadcrumb-list").R(
			func() (x any) {
				for i, item := range bc.Items {
					isLast := i == len(bc.Items)-1

					b.LiClass("breadcrumb-item").R(
						func() (x any) {
							if isLast || item.Href == "" {
								// Current page - no link
								b.SpanClass("breadcrumb-current", "aria-current", "page").T(item.Label)
							} else {
								b.A("href", item.Href).T(item.Label)
								b.SpanClass("breadcrumb-separator", "aria-hidden", "true").T(" " + separator + " ")
							}
							return
						}(),
					)
				}
				return
			}(),
		),
	)
	return
}
