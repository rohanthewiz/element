package components

import "github.com/rohanthewiz/element"

// -----------------------------------------------------------------------------
// Badge Component
// -----------------------------------------------------------------------------

// Badge renders a small status label (count, tag, state indicator).
type Badge struct {
	Text  string    // Badge text
	Type  AlertType // Color scheme (info, success, warning, error) — default "neutral"
	Pill  bool      // Fully rounded ends
	Class string    // Additional CSS classes
}

// Render implements the element.Component interface.
func (bd Badge) Render(b *element.Builder) (x any) {
	badgeType := string(bd.Type)
	if badgeType == "" {
		badgeType = "neutral"
	}
	badgeClass := "badge badge-" + badgeType
	if bd.Pill {
		badgeClass += " badge-pill"
	}
	if bd.Class != "" {
		badgeClass += " " + bd.Class
	}

	b.SpanClass(badgeClass).T(bd.Text)
	return
}
