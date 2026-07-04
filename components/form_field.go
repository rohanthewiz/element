package components

import "github.com/rohanthewiz/element"

// -----------------------------------------------------------------------------
// Form Field Component
// -----------------------------------------------------------------------------

// FormField renders a form field with label, input, and optional help text.
type FormField struct {
	Label       string   // Field label
	Name        string   // Input name attribute
	ID          string   // Input id attribute (defaults to Name)
	Type        string   // Input type (text, email, password, etc.)
	Placeholder string   // Placeholder text
	Value       string   // Current value
	Required    bool     // Whether field is required
	Disabled    bool     // Whether field is disabled
	HelpText    string   // Optional help text below field
	Error       string   // Error message to display
	Class       string   // Additional CSS classes for the wrapper div
	Attrs       []string // Extra attribute pairs for the <input> (e.g. "autocomplete", "off")
}

// Render implements the element.Component interface.
func (f FormField) Render(b *element.Builder) (x any) {
	inputType := f.Type
	if inputType == "" {
		inputType = "text"
	}
	id := f.ID
	if id == "" {
		id = f.Name
	}

	fieldClass := "form-field"
	if f.Error != "" {
		fieldClass += " has-error"
	}
	if f.Class != "" {
		fieldClass += " " + f.Class
	}

	b.DivClass(fieldClass).R(
		// Label
		b.LabelClass("form-label", "for", id).R(
			b.T(f.Label),
			func() (x any) {
				if f.Required {
					b.SpanClass("required").T(" *")
				}
				return
			}(),
		),
		// Input
		func() (x any) {
			attrs := []string{
				"type", inputType,
				"id", id,
				"name", f.Name,
				"class", "form-input",
			}
			if f.Placeholder != "" {
				attrs = append(attrs, "placeholder", f.Placeholder)
			}
			if f.Value != "" {
				attrs = append(attrs, "value", f.Value)
			}
			if f.Required {
				attrs = append(attrs, "required", "required")
			}
			if f.Disabled {
				attrs = append(attrs, "disabled", "disabled")
			}
			// Point screen readers at the error or help text
			if f.Error != "" {
				attrs = append(attrs, "aria-invalid", "true", "aria-describedby", id+"-error")
			} else if f.HelpText != "" {
				attrs = append(attrs, "aria-describedby", id+"-help")
			}
			attrs = append(attrs, f.Attrs...)
			b.Input(attrs...)
			return
		}(),
		// Error message
		func() (x any) {
			if f.Error != "" {
				b.SpanClass("form-error", "id", id+"-error").T(f.Error)
			}
			return
		}(),
		// Help text
		func() (x any) {
			if f.HelpText != "" && f.Error == "" {
				b.SmallClass("form-help", "id", id+"-help").T(f.HelpText)
			}
			return
		}(),
	)
	return
}
