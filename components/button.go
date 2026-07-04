package components

import "github.com/rohanthewiz/element"

// -----------------------------------------------------------------------------
// Button Component
// -----------------------------------------------------------------------------

// Button renders a styled <button>, or an <a> styled as a button when Href
// is set.
type Button struct {
	Text     string   // Button text
	Type     string   // "button", "submit", or "reset" (default "button"; ignored when Href is set)
	Variant  string   // "primary" (default), "secondary", "danger", "ghost", ... — rendered as btn-<variant>
	Size     string   // Optional "sm" or "lg" — rendered as btn-<size>
	Href     string   // If set, renders a link styled as a button
	Disabled bool     // Whether the button is disabled
	Class    string   // Additional CSS classes
	Attrs    []string // Extra attribute pairs
}

// Render implements the element.Component interface.
func (bt Button) Render(b *element.Builder) (x any) {
	variant := bt.Variant
	if variant == "" {
		variant = "primary"
	}
	btnClass := "btn btn-" + variant
	if bt.Size != "" {
		btnClass += " btn-" + bt.Size
	}
	if bt.Class != "" {
		btnClass += " " + bt.Class
	}

	if bt.Href != "" {
		attrs := []string{"href", bt.Href, "role", "button"}
		if bt.Disabled {
			attrs = append(attrs, "aria-disabled", "true")
			btnClass += " btn-disabled"
		}
		attrs = append(attrs, bt.Attrs...)
		b.AClass(btnClass, attrs...).T(bt.Text)
		return
	}

	btnType := bt.Type
	if btnType == "" {
		btnType = "button"
	}
	attrs := []string{"type", btnType}
	if bt.Disabled {
		attrs = append(attrs, "disabled", "disabled")
	}
	attrs = append(attrs, bt.Attrs...)
	b.ButtonClass(btnClass, attrs...).T(bt.Text)
	return
}
