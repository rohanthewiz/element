package components

import "github.com/rohanthewiz/element"

// -----------------------------------------------------------------------------
// Checkbox Field Component
// -----------------------------------------------------------------------------

// CheckboxField renders a single checkbox with its label on the right,
// following the same conventions as FormField.
type CheckboxField struct {
	Label    string   // Text shown beside the checkbox
	Name     string   // Input name attribute
	ID       string   // Input id attribute (defaults to Name)
	Value    string   // Submitted value when checked (default "on" per HTML)
	Checked  bool     // Whether the box starts checked
	Required bool     // Whether the box must be checked
	Disabled bool     // Whether the box is disabled
	HelpText string   // Optional help text below the field
	Error    string   // Error message to display
	Class    string   // Additional CSS classes for the wrapper div
	Attrs    []string // Extra attribute pairs for the <input>
}

// Render implements the element.Component interface.
func (c CheckboxField) Render(b *element.Builder) (x any) {
	id := c.ID
	if id == "" {
		id = c.Name
	}

	extra := "form-check"
	if c.Class != "" {
		extra += " " + c.Class
	}

	b.DivClass(fieldClass(c.Error, extra)).R(
		b.Wrap(func() {
			attrs := []string{"type", "checkbox", "id", id, "name", c.Name, "class", "form-check-input"}
			if c.Value != "" {
				attrs = append(attrs, "value", c.Value)
			}
			if c.Checked {
				attrs = append(attrs, "checked", "checked")
			}
			if c.Required {
				attrs = append(attrs, "required", "required")
			}
			if c.Disabled {
				attrs = append(attrs, "disabled", "disabled")
			}
			attrs = append(attrs, describedByAttrs(id, c.HelpText, c.Error)...)
			attrs = append(attrs, c.Attrs...)
			b.Input(attrs...)
		}),
		b.LabelClass("form-check-label", "for", id).R(
			b.T(c.Label),
			b.Wrap(func() {
				if c.Required {
					b.SpanClass("required").T(" *")
				}
			}),
		),
		b.Wrap(func() {
			renderFieldMessages(b, id, c.HelpText, c.Error)
		}),
	)
	return
}
