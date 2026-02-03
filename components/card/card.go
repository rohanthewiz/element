package card

import "github.com/rohanthewiz/element"

// -----------------------------------------------------------------------------
// Card Component
// -----------------------------------------------------------------------------

// Card renders a styled card container with optional header and footer.
type Card struct {
	Title         string            // Card title (rendered in header)
	Body          string            // Simple text body content
	BodyComponent element.Component // Alternative: render a component as body
	Footer        string            // Optional footer text
	Class         string            // Additional CSS classes
}

// Render implements the element.Component interface.
func (c Card) Render(b *element.Builder) (x any) {
	cardClass := "card"
	if c.Class != "" {
		cardClass += " " + c.Class
	}

	b.DivClass(cardClass).R(
		// Header with title
		func() (x any) {
			if c.Title != "" {
				b.DivClass("card-header").R(
					b.H4Class("card-title").T(c.Title),
				)
			}
			return
		}(),
		// Body
		b.DivClass("card-body").R(
			func() (x any) {
				if c.BodyComponent != nil {
					c.BodyComponent.Render(b)
				} else if c.Body != "" {
					b.P().T(c.Body)
				}
				return
			}(),
		),
		// Footer
		func() (x any) {
			if c.Footer != "" {
				b.DivClass("card-footer").R(
					b.T(c.Footer),
				)
			}
			return
		}(),
	)
	return
}
