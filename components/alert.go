package components

import (
	"fmt"

	"github.com/rohanthewiz/element"
)

// -----------------------------------------------------------------------------
// Alert Component
// -----------------------------------------------------------------------------

// AlertType defines the style/severity of an alert.
type AlertType string

const (
	AlertInfo    AlertType = "info"
	AlertSuccess AlertType = "success"
	AlertWarning AlertType = "warning"
	AlertError   AlertType = "error"
)

// Alert renders a styled notification message.
type Alert struct {
	Type        AlertType // Alert style (info, success, warning, error)
	Title       string    // Optional title
	Message     string    // Alert message
	Dismissible bool      // Whether to show a close button
}

// Render implements the element.Component interface.
func (a Alert) Render(b *element.Builder) (x any) {
	alertType := a.Type
	if alertType == "" {
		alertType = AlertInfo
	}
	alertClass := fmt.Sprintf("alert alert-%s", alertType)
	if a.Dismissible {
		alertClass += " alert-dismissible"
	}

	b.DivClass(alertClass, "role", "alert").R(
		// Optional title
		func() (x any) {
			if a.Title != "" {
				b.Strong().T(a.Title)
				b.T(" ")
			}
			return
		}(),
		// Message
		b.T(a.Message),
		// Dismiss button
		func() (x any) {
			if a.Dismissible {
				b.ButtonClass("alert-close", "type", "button", "aria-label", "Close").T("×")
			}
			return
		}(),
	)
	return
}
