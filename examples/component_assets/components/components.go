// Package components demonstrates Element components with embedded assets.
// This example shows how to implement the AssetProvider interface to manage
// component-scoped CSS and JavaScript.
package components

import (
	_ "embed"
	"fmt"

	"github.com/rohanthewiz/element"
)

// Embed component assets
//
//go:embed alert.css
var alertCSS []byte

//go:embed alert.js
var alertJS []byte

//go:embed table.css
var tableCSS []byte

// -----------------------------------------------------------------------------
// Alert Component with CSS and JavaScript
// -----------------------------------------------------------------------------

// AlertType defines the style/severity of an alert.
type AlertType string

const (
	AlertInfo    AlertType = "info"
	AlertSuccess AlertType = "success"
	AlertWarning AlertType = "warning"
	AlertError   AlertType = "error"
)

// Alert renders a styled notification message with optional dismiss functionality.
type Alert struct {
	Type        AlertType
	Title       string
	Message     string
	Dismissible bool
}

// AssetID implements AssetProvider - returns unique identifier for deduplication.
func (a Alert) AssetID() string {
	return "Alert"
}

// CSS implements AssetProvider - returns component CSS.
func (a Alert) CSS() []byte {
	return alertCSS
}

// JS implements AssetProvider - returns component JavaScript.
func (a Alert) JS() []byte {
	if a.Dismissible {
		return alertJS // Only include JS if dismissible functionality is needed
	}
	return nil
}

// Render implements the element.Component interface.
func (a Alert) Render(b *element.Builder) (x any) {
	// Register this component's assets with the builder
	b.RegisterAssets(a)

	alertType := a.Type
	if alertType == "" {
		alertType = AlertInfo
	}

	alertClass := fmt.Sprintf("ele-alert ele-alert-%s", alertType)
	if a.Dismissible {
		alertClass += " ele-alert-dismissible"
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
				b.ButtonClass("ele-alert-close", "type", "button", "aria-label", "Close").T("×")
			}
			return
		}(),
	)
	return
}

// -----------------------------------------------------------------------------
// Table Component with CSS only
// -----------------------------------------------------------------------------

// Table renders an HTML table with customizable styling.
type Table struct {
	Headers  []string
	Rows     [][]any
	Class    string
	Striped  bool
	Bordered bool
}

// AssetID implements AssetProvider.
func (t Table) AssetID() string {
	return "Table"
}

// CSS implements AssetProvider.
func (t Table) CSS() []byte {
	return tableCSS
}

// JS implements AssetProvider - Table doesn't need JavaScript.
func (t Table) JS() []byte {
	return nil
}

// Render implements the element.Component interface.
func (t Table) Render(b *element.Builder) (x any) {
	// Register this component's assets
	b.RegisterAssets(t)

	// Build table classes
	tableClass := t.Class
	if tableClass == "" {
		tableClass = "ele-table"
	} else {
		tableClass = "ele-table " + tableClass
	}
	if t.Striped {
		tableClass += " ele-table-striped"
	}
	if t.Bordered {
		tableClass += " ele-table-bordered"
	}

	b.TableClass(tableClass).R(
		// Render header if provided
		func() (x any) {
			if len(t.Headers) > 0 {
				b.THead().R(
					b.Tr().R(
						func() (x any) {
							for _, header := range t.Headers {
								b.Th().T(header)
							}
							return
						}(),
					),
				)
			}
			return
		}(),
		// Render body
		b.TBody().R(
			func() (x any) {
				for _, row := range t.Rows {
					b.Tr().R(
						func() (x any) {
							for _, cell := range row {
								b.Td().T(fmt.Sprintf("%v", cell))
							}
							return
						}(),
					)
				}
				return
			}(),
		),
	)
	return
}

// -----------------------------------------------------------------------------
// Helper: Component without assets (for comparison)
// -----------------------------------------------------------------------------

// SimpleCard is a basic card component that doesn't provide assets.
// This demonstrates that components can work without implementing AssetProvider.
type SimpleCard struct {
	Title string
	Body  string
}

// Render implements the element.Component interface.
// Note: This component does NOT implement AssetProvider.
func (c SimpleCard) Render(b *element.Builder) (x any) {
	// No asset registration needed - this component uses inline styles
	b.Div("style", "border: 1px solid #ccc; padding: 1rem; border-radius: 0.25rem; margin: 1rem 0;").R(
		func() (x any) {
			if c.Title != "" {
				b.H3("style", "margin-top: 0;").T(c.Title)
			}
			return
		}(),
		b.P().T(c.Body),
	)
	return
}
