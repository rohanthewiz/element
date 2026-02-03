package form_field

import "github.com/rohanthewiz/element"

// -----------------------------------------------------------------------------
// Form Field Component
// -----------------------------------------------------------------------------

// FormField renders a form field with label, input, and optional help text.
type FormField struct {
	Label       string // Field label
	Name        string // Input name attribute
	Type        string // Input type (text, email, password, etc.)
	Placeholder string // Placeholder text
	Value       string // Current value
	Required    bool   // Whether field is required
	HelpText    string // Optional help text below field
	Error       string // Error message to display
}

// Render implements the element.Component interface.
func (f FormField) Render(b *element.Builder) (x any) {
	inputType := f.Type
	if inputType == "" {
		inputType = "text"
	}

	fieldClass := "form-field"
	if f.Error != "" {
		fieldClass += " has-error"
	}

	b.DivClass(fieldClass).R(
		// Label
		b.LabelClass("form-label", "for", f.Name).R(
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
				"id", f.Name,
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
			b.Input(attrs...)
			return
		}(),
		// Error message
		func() (x any) {
			if f.Error != "" {
				b.SpanClass("form-error").T(f.Error)
			}
			return
		}(),
		// Help text
		func() (x any) {
			if f.HelpText != "" && f.Error == "" {
				b.SmallClass("form-help").T(f.HelpText)
			}
			return
		}(),
	)
	return
}
