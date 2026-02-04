package components

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
	Headers  []string // Column headers
	Rows     [][]any  // Table body data (2D array)
	Class    string   // Optional CSS class for the table
	Striped  bool     // Whether to add striped row styling
	Bordered bool     // Whether to add borders
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
