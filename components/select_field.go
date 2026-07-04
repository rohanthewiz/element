package components

import "github.com/rohanthewiz/element"

// -----------------------------------------------------------------------------
// Select Field Component
// -----------------------------------------------------------------------------

// SelectField renders a labeled <select> with options, following the same
// conventions as FormField (form-field wrapper, error and help text).
type SelectField struct {
	Label    string         // Field label
	Name     string         // Select name attribute
	ID       string         // Select id attribute (defaults to Name)
	Options  []SelectOption // The choices
	Selected string         // Value of the currently selected option
	Prompt   string         // Optional first, unselectable option (e.g. "Choose one…")
	Required bool           // Whether a selection is required
	Disabled bool           // Whether the select is disabled
	HelpText string         // Optional help text below the field
	Error    string         // Error message to display
	Class    string         // Additional CSS classes for the wrapper div
	Attrs    []string       // Extra attribute pairs for the <select>
}

// Render implements the element.Component interface.
func (s SelectField) Render(b *element.Builder) (x any) {
	id := s.ID
	if id == "" {
		id = s.Name
	}

	b.DivClass(fieldClass(s.Error, s.Class)).R(
		b.LabelClass("form-label", "for", id).R(
			b.T(s.Label),
			b.Wrap(func() {
				if s.Required {
					b.SpanClass("required").T(" *")
				}
			}),
		),
		b.Wrap(func() {
			attrs := []string{"id", id, "name", s.Name, "class", "form-input form-select"}
			if s.Required {
				attrs = append(attrs, "required", "required")
			}
			if s.Disabled {
				attrs = append(attrs, "disabled", "disabled")
			}
			attrs = append(attrs, describedByAttrs(id, s.HelpText, s.Error)...)
			attrs = append(attrs, s.Attrs...)

			b.Select(attrs...).R(
				b.Wrap(func() {
					if s.Prompt != "" {
						promptAttrs := []string{"value", "", "disabled", "disabled"}
						if s.Selected == "" {
							promptAttrs = append(promptAttrs, "selected", "selected")
						}
						b.Option(promptAttrs...).T(s.Prompt)
					}
					for _, opt := range s.Options {
						optAttrs := []string{"value", opt.Value}
						if opt.Value == s.Selected {
							optAttrs = append(optAttrs, "selected", "selected")
						}
						if opt.Disabled {
							optAttrs = append(optAttrs, "disabled", "disabled")
						}
						b.Option(optAttrs...).T(opt.optionLabel())
					}
				}),
			)
		}),
		b.Wrap(func() {
			renderFieldMessages(b, id, s.HelpText, s.Error)
		}),
	)
	return
}
