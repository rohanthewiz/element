package components

import (
	"strconv"

	"github.com/rohanthewiz/element"
)

// -----------------------------------------------------------------------------
// Radio Group Component
// -----------------------------------------------------------------------------

// RadioGroup renders a set of radio buttons inside a <fieldset> with a legend.
type RadioGroup struct {
	Legend   string         // Group label rendered as the fieldset legend
	Name     string         // Shared name attribute for all radios
	Options  []SelectOption // The choices
	Selected string         // Value of the checked option
	Inline   bool           // Lay options out horizontally (adds radio-group-inline)
	Required bool           // Whether a choice is required
	Disabled bool           // Disable the whole group
	HelpText string         // Optional help text below the group
	Error    string         // Error message to display
	Class    string         // Additional CSS classes for the fieldset
}

// Render implements the element.Component interface.
func (r RadioGroup) Render(b *element.Builder) (x any) {
	groupClass := fieldClass(r.Error, "radio-group")
	if r.Inline {
		groupClass += " radio-group-inline"
	}
	if r.Class != "" {
		groupClass += " " + r.Class
	}

	fsAttrs := []string{"class", groupClass}
	if r.Disabled {
		fsAttrs = append(fsAttrs, "disabled", "disabled")
	}

	b.FieldSet(fsAttrs...).R(
		b.Wrap(func() {
			if r.Legend != "" {
				b.LegendClass("form-label").R(
					b.T(r.Legend),
					b.Wrap(func() {
						if r.Required {
							b.SpanClass("required").T(" *")
						}
					}),
				)
			}

			for i, opt := range r.Options {
				id := r.Name + "-" + strconv.Itoa(i)
				attrs := []string{
					"type", "radio", "id", id,
					"name", r.Name, "value", opt.Value,
					"class", "form-check-input",
				}
				if opt.Value == r.Selected {
					attrs = append(attrs, "checked", "checked")
				}
				if opt.Disabled {
					attrs = append(attrs, "disabled", "disabled")
				}
				if r.Required {
					attrs = append(attrs, "required", "required")
				}

				b.DivClass("form-check radio-option").R(
					b.Wrap(func() { b.Input(attrs...) }),
					b.LabelClass("form-check-label", "for", id).T(opt.optionLabel()),
				)
			}

			renderFieldMessages(b, r.Name, r.HelpText, r.Error)
		}),
	)
	return
}
