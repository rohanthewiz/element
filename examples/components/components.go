// Package main demonstrates reusable Element components with no external dependencies.
// These components showcase the Component interface and can be used as templates
// for building your own component library.
package main

import (
	"fmt"

	"github.com/rohanthewiz/element"
)

// -----------------------------------------------------------------------------
// Table Component
// -----------------------------------------------------------------------------

// Table renders an HTML table with customizable headers and data rows.
// It accepts headers as strings and body data as a 2D slice of any type.
type Table struct {
	Headers []string   // Column headers
	Rows    [][]any    // Table body data (2D array)
	Class   string     // Optional CSS class for the table
	Striped bool       // Whether to add striped row styling
	Bordered bool      // Whether to add borders
}

// Render implements the element.Component interface.
func (t Table) Render(b *element.Builder) (x any) {
	// Build table classes
	tableClass := t.Class
	if tableClass == "" {
		tableClass = "table"
	}
	if t.Striped {
		tableClass += " table-striped"
	}
	if t.Bordered {
		tableClass += " table-bordered"
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
// Card Component
// -----------------------------------------------------------------------------

// Card renders a styled card container with optional header and footer.
type Card struct {
	Title   string           // Card title (rendered in header)
	Body    string           // Simple text body content
	BodyComponent element.Component // Alternative: render a component as body
	Footer  string           // Optional footer text
	Class   string           // Additional CSS classes
}

// Render implements the element.Component interface.
func (c Card) Render(b *element.Builder) (x any) {
	cardClass := "card"
	if c.Class != "" {
		cardClass += " " + c.Class
	}

	b.DivClass(cardClass).R(
		// Header with title
		func() (x any) {
			if c.Title != "" {
				b.DivClass("card-header").R(
					b.H4Class("card-title").T(c.Title),
				)
			}
			return
		}(),
		// Body
		b.DivClass("card-body").R(
			func() (x any) {
				if c.BodyComponent != nil {
					c.BodyComponent.Render(b)
				} else if c.Body != "" {
					b.P().T(c.Body)
				}
				return
			}(),
		),
		// Footer
		func() (x any) {
			if c.Footer != "" {
				b.DivClass("card-footer").R(
					b.T(c.Footer),
				)
			}
			return
		}(),
	)
	return
}

// -----------------------------------------------------------------------------
// Nav Component
// -----------------------------------------------------------------------------

// NavItem represents a single navigation link.
type NavItem struct {
	Label  string // Display text
	Href   string // Link URL
	Active bool   // Whether this is the current page
}

// Nav renders a navigation bar with links.
type Nav struct {
	Items []NavItem // Navigation items
	Brand string    // Optional brand/logo text
	Class string    // Additional CSS classes
}

// Render implements the element.Component interface.
func (n Nav) Render(b *element.Builder) (x any) {
	navClass := "nav"
	if n.Class != "" {
		navClass += " " + n.Class
	}

	b.NavClass(navClass).R(
		// Brand
		func() (x any) {
			if n.Brand != "" {
				b.AClass("nav-brand", "href", "#").T(n.Brand)
			}
			return
		}(),
		// Nav items
		b.UlClass("nav-list").R(
			func() (x any) {
				for _, item := range n.Items {
					itemClass := "nav-item"
					if item.Active {
						itemClass += " active"
					}
					b.LiClass(itemClass).R(
						b.AClass("nav-link", "href", item.Href).T(item.Label),
					)
				}
				return
			}(),
		),
	)
	return
}

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
	Type    AlertType // Alert style (info, success, warning, error)
	Title   string    // Optional title
	Message string    // Alert message
	Dismissible bool  // Whether to show a close button
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

// -----------------------------------------------------------------------------
// Breadcrumb Component
// -----------------------------------------------------------------------------

// BreadcrumbItem represents a single breadcrumb link.
type BreadcrumbItem struct {
	Label string // Display text
	Href  string // Link URL (empty for current page)
}

// Breadcrumb renders a navigation breadcrumb trail.
type Breadcrumb struct {
	Items     []BreadcrumbItem // Breadcrumb items (last is current)
	Separator string           // Separator character (default: "/")
}

// Render implements the element.Component interface.
func (bc Breadcrumb) Render(b *element.Builder) (x any) {
	separator := bc.Separator
	if separator == "" {
		separator = "/"
	}

	b.NavClass("breadcrumb", "aria-label", "breadcrumb").R(
		b.OlClass("breadcrumb-list").R(
			func() (x any) {
				for i, item := range bc.Items {
					isLast := i == len(bc.Items)-1

					b.LiClass("breadcrumb-item").R(
						func() (x any) {
							if isLast || item.Href == "" {
								// Current page - no link
								b.SpanClass("breadcrumb-current", "aria-current", "page").T(item.Label)
							} else {
								b.A("href", item.Href).T(item.Label)
								b.SpanClass("breadcrumb-separator", "aria-hidden", "true").T(" " + separator + " ")
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
// Definition List Component
// -----------------------------------------------------------------------------

// Definition represents a term and its definition.
type Definition struct {
	Term       string // The term being defined
	Definition string // The definition
}

// DefinitionList renders a definition list (<dl>) element.
type DefinitionList struct {
	Items []Definition // Term/definition pairs
	Class string       // Optional CSS class
}

// Render implements the element.Component interface.
func (dl DefinitionList) Render(b *element.Builder) (x any) {
	dlClass := "definition-list"
	if dl.Class != "" {
		dlClass += " " + dl.Class
	}

	b.DlClass(dlClass).R(
		func() (x any) {
			for _, item := range dl.Items {
				b.DtClass("definition-term").T(item.Term)
				b.DdClass("definition-desc").T(item.Definition)
			}
			return
		}(),
	)
	return
}

// -----------------------------------------------------------------------------
// Pagination Component
// -----------------------------------------------------------------------------

// Pagination renders page navigation controls.
type Pagination struct {
	CurrentPage int    // Current page number (1-based)
	TotalPages  int    // Total number of pages
	BaseURL     string // URL pattern with %d placeholder for page number
	ShowFirst   bool   // Show "First" link
	ShowLast    bool   // Show "Last" link
}

// Render implements the element.Component interface.
func (p Pagination) Render(b *element.Builder) (x any) {
	if p.TotalPages <= 1 {
		return // Don't render if only one page
	}

	b.NavClass("pagination", "aria-label", "Page navigation").R(
		b.UlClass("pagination-list").R(
			// First page link
			func() (x any) {
				if p.ShowFirst && p.CurrentPage > 1 {
					b.Li().R(
						b.AClass("pagination-link", "href", fmt.Sprintf(p.BaseURL, 1)).T("« First"),
					)
				}
				return
			}(),
			// Previous page link
			func() (x any) {
				if p.CurrentPage > 1 {
					b.Li().R(
						b.AClass("pagination-link", "href", fmt.Sprintf(p.BaseURL, p.CurrentPage-1)).T("‹ Prev"),
					)
				}
				return
			}(),
			// Page numbers (show current and neighbors)
			func() (x any) {
				start := max(1, p.CurrentPage-2)
				end := min(p.TotalPages, p.CurrentPage+2)

				for i := start; i <= end; i++ {
					liClass := "pagination-item"
					if i == p.CurrentPage {
						liClass += " active"
						b.LiClass(liClass).R(
							b.SpanClass("pagination-current").T(fmt.Sprintf("%d", i)),
						)
					} else {
						b.LiClass(liClass).R(
							b.AClass("pagination-link", "href", fmt.Sprintf(p.BaseURL, i)).T(fmt.Sprintf("%d", i)),
						)
					}
				}
				return
			}(),
			// Next page link
			func() (x any) {
				if p.CurrentPage < p.TotalPages {
					b.Li().R(
						b.AClass("pagination-link", "href", fmt.Sprintf(p.BaseURL, p.CurrentPage+1)).T("Next ›"),
					)
				}
				return
			}(),
			// Last page link
			func() (x any) {
				if p.ShowLast && p.CurrentPage < p.TotalPages {
					b.Li().R(
						b.AClass("pagination-link", "href", fmt.Sprintf(p.BaseURL, p.TotalPages)).T("Last »"),
					)
				}
				return
			}(),
		),
	)
	return
}

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
