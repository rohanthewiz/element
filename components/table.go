package components

import (
	"fmt"
	"strconv"

	"github.com/rohanthewiz/element"
)

// -----------------------------------------------------------------------------
// Table Component
// -----------------------------------------------------------------------------

// Column describes a single table column when more control than a plain
// header string is needed. Use Table.Columns (instead of Table.Headers)
// to get per-column alignment and classes.
type Column struct {
	Header      string // Header text
	Align       string // "left", "center", or "right" — applied to header and body cells
	Class       string // Optional class applied to every body cell in this column
	HeaderClass string // Optional class applied to the header cell
}

// Table renders an HTML table with customizable headers and data rows.
//
// Cells in Rows and FooterRows may be:
//   - element.Component — rendered in place (links, buttons, badges, etc.)
//   - string, ints, floats, bool — written directly without reflection
//   - fmt.Stringer — its String() output
//   - anything else — formatted with %v
type Table struct {
	Headers      []string // Simple column headers (use Columns for more control)
	Columns      []Column // Full column definitions — takes precedence over Headers
	Rows         [][]any  // Table body data (2D array)
	FooterRows   [][]any  // Optional footer rows rendered in <tfoot>
	Caption      string   // Optional <caption> text
	EmptyMessage string   // Shown in a single row spanning all columns when Rows is empty
	Class        string   // Optional CSS class for the table (default "table")
	Striped      bool     // Whether to add striped row styling
	Bordered     bool     // Whether to add borders
	Hover        bool     // Whether to add row hover styling
	Responsive   bool     // Wrap the table in a horizontally scrollable container
	Attrs        []string // Extra attribute pairs for the <table> element
	// RowClass, if set, returns a class for each body row (empty string for none)
	RowClass func(rowIndex int, row []any) string
}

// Render implements the element.Component interface.
func (t Table) Render(b *element.Builder) (x any) {
	cols := t.Columns
	if len(cols) == 0 {
		for _, h := range t.Headers {
			cols = append(cols, Column{Header: h})
		}
	}

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
	if t.Hover {
		tableClass += " table-hover"
	}

	if t.Responsive {
		b.DivClass("table-responsive").R(
			t.renderTable(b, cols, tableClass),
		)
	} else {
		t.renderTable(b, cols, tableClass)
	}
	return
}

func (t Table) renderTable(b *element.Builder, cols []Column, tableClass string) (x any) {
	b.TableClass(tableClass, t.Attrs...).R(
		b.Wrap(func() {
			if t.Caption != "" {
				b.Caption().T(t.Caption)
			}

			// Header — rendered when any column has header text
			if hasHeader(cols) {
				b.THead().R(
					b.Tr().R(
						element.ForEach(cols, func(c Column) {
							b.Th(cellAttrs(c.HeaderClass, c.Align)...).T(c.Header)
						}),
					),
				)
			}

			// Body
			b.TBody().R(
				b.Wrap(func() {
					if len(t.Rows) == 0 && t.EmptyMessage != "" {
						span := max(len(cols), 1)
						b.Tr().R(
							b.TdClass("table-empty", "colspan", strconv.Itoa(span)).T(t.EmptyMessage),
						)
						return
					}

					for i, row := range t.Rows {
						rowCls := ""
						if t.RowClass != nil {
							rowCls = t.RowClass(i, row)
						}

						// Elements open as soon as they are created,
						// so choose the row's form before creating it
						if rowCls != "" {
							b.TrClass(rowCls).R(t.renderCells(b, cols, row))
						} else {
							b.Tr().R(t.renderCells(b, cols, row))
						}
					}
				}),
			)

			// Footer
			if len(t.FooterRows) > 0 {
				b.TFoot().R(
					b.Wrap(func() {
						for _, row := range t.FooterRows {
							b.Tr().R(t.renderCells(b, cols, row))
						}
					}),
				)
			}
		}),
	)
	return
}

// renderCells writes the <td> cells for one row, applying column class/alignment
func (t Table) renderCells(b *element.Builder, cols []Column, row []any) (x any) {
	for j, cell := range row {
		var attrs []string
		if j < len(cols) {
			attrs = cellAttrs(cols[j].Class, cols[j].Align)
		}
		b.Td(attrs...).R(
			renderCell(b, cell),
		)
	}
	return
}

// cellAttrs builds the attribute pairs for a th/td from a class and alignment
func cellAttrs(class, align string) (attrs []string) {
	if class != "" {
		attrs = append(attrs, "class", class)
	}
	if align != "" {
		attrs = append(attrs, "style", "text-align:"+align)
	}
	return
}

// hasHeader reports whether any column has header text
func hasHeader(cols []Column) bool {
	for _, c := range cols {
		if c.Header != "" {
			return true
		}
	}
	return false
}

// renderCell writes a single cell value to the builder.
// Components render themselves; common scalar types avoid fmt allocations.
func renderCell(b *element.Builder, v any) (x any) {
	switch c := v.(type) {
	case nil:
	case element.Component:
		c.Render(b)
	case string:
		b.T(c)
	case int:
		b.T(strconv.Itoa(c))
	case int64:
		b.T(strconv.FormatInt(c, 10))
	case uint64:
		b.T(strconv.FormatUint(c, 10))
	case float64:
		b.T(strconv.FormatFloat(c, 'f', -1, 64))
	case float32:
		b.T(strconv.FormatFloat(float64(c), 'f', -1, 32))
	case bool:
		b.T(strconv.FormatBool(c))
	case fmt.Stringer:
		b.T(c.String())
	default:
		b.F("%v", c)
	}
	return
}
