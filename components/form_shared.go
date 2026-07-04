package components

import "github.com/rohanthewiz/element"

// -----------------------------------------------------------------------------
// Shared form helpers
// -----------------------------------------------------------------------------

// SelectOption is a single choice for SelectField and RadioGroup.
type SelectOption struct {
	Value    string // Submitted value
	Label    string // Display text (defaults to Value)
	Disabled bool   // Whether this option can be chosen
}

// optionLabel returns the display text for an option
func (o SelectOption) optionLabel() string {
	if o.Label != "" {
		return o.Label
	}
	return o.Value
}

// fieldClass builds the wrapper class for a form field
func fieldClass(errMsg, extra string) string {
	c := "form-field"
	if errMsg != "" {
		c += " has-error"
	}
	if extra != "" {
		c += " " + extra
	}
	return c
}

// describedByAttrs returns aria attributes pointing screen readers
// at the field's error or help text
func describedByAttrs(id, helpText, errMsg string) []string {
	if errMsg != "" {
		return []string{"aria-invalid", "true", "aria-describedby", id + "-error"}
	}
	if helpText != "" {
		return []string{"aria-describedby", id + "-help"}
	}
	return nil
}

// renderFieldMessages writes the error message (or the help text when there
// is no error) below a form control
func renderFieldMessages(b *element.Builder, id, helpText, errMsg string) {
	if errMsg != "" {
		b.SpanClass("form-error", "id", id+"-error").T(errMsg)
	} else if helpText != "" {
		b.SmallClass("form-help", "id", id+"-help").T(helpText)
	}
}
