package components

import (
	"strconv"

	"github.com/rohanthewiz/element"
)

// -----------------------------------------------------------------------------
// TextArea Field Component
// -----------------------------------------------------------------------------

// TextAreaField renders a labeled <textarea>, following the same
// conventions as FormField (form-field wrapper, error and help text).
type TextAreaField struct {
	Label       string   // Field label
	Name        string   // Textarea name attribute
	ID          string   // Textarea id attribute (defaults to Name)
	Placeholder string   // Placeholder text
	Value       string   // Current content
	Rows        int      // Visible rows (default 4)
	Required    bool     // Whether the field is required
	Disabled    bool     // Whether the field is disabled
	HelpText    string   // Optional help text below the field
	Error       string   // Error message to display
	Class       string   // Additional CSS classes for the wrapper div
	Attrs       []string // Extra attribute pairs for the <textarea>
}

// Render implements the element.Component interface.
func (t TextAreaField) Render(b *element.Builder) (x any) {
	id := t.ID
	if id == "" {
		id = t.Name
	}
	rows := t.Rows
	if rows <= 0 {
		rows = 4
	}

	b.DivClass(fieldClass(t.Error, t.Class)).R(
		b.LabelClass("form-label", "for", id).R(
			b.T(t.Label),
			b.Wrap(func() {
				if t.Required {
					b.SpanClass("required").T(" *")
				}
			}),
		),
		b.Wrap(func() {
			attrs := []string{
				"id", id, "name", t.Name,
				"class", "form-input form-textarea",
				"rows", strconv.Itoa(rows),
			}
			if t.Placeholder != "" {
				attrs = append(attrs, "placeholder", t.Placeholder)
			}
			if t.Required {
				attrs = append(attrs, "required", "required")
			}
			if t.Disabled {
				attrs = append(attrs, "disabled", "disabled")
			}
			attrs = append(attrs, describedByAttrs(id, t.HelpText, t.Error)...)
			attrs = append(attrs, t.Attrs...)

			b.TextArea(attrs...).T(t.Value)
		}),
		b.Wrap(func() {
			renderFieldMessages(b, id, t.HelpText, t.Error)
		}),
	)
	return
}
