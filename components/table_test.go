package components

import (
	"strings"
	"testing"

	"github.com/rohanthewiz/element"
)

func TestTable_Render(t *testing.T) {
	tests := []struct {
		name     string
		table    Table
		contains []string
	}{
		{
			name: "basic table with headers and rows",
			table: Table{
				Headers: []string{"Name", "Age"},
				Rows: [][]any{
					{"Alice", 30},
					{"Bob", 25},
				},
			},
			contains: []string{
				"<table",
				"<thead>",
				"<th>Name</th>",
				"<th>Age</th>",
				"</thead>",
				"<tbody>",
				"<tr>",
				"<td>Alice</td>",
				"<td>30</td>",
				"<td>Bob</td>",
				"<td>25</td>",
				"</tbody>",
				"</table>",
			},
		},
		{
			name: "table without headers",
			table: Table{
				Rows: [][]any{
					{"Cell1", "Cell2"},
				},
			},
			contains: []string{
				"<table",
				"<tbody>",
				"<td>Cell1</td>",
				"<td>Cell2</td>",
				"</tbody>",
			},
		},
		{
			name: "striped table",
			table: Table{
				Headers: []string{"Col1"},
				Rows:    [][]any{{"Data"}},
				Striped: true,
			},
			contains: []string{
				`class="table table-striped"`,
			},
		},
		{
			name: "bordered table",
			table: Table{
				Headers:  []string{"Col1"},
				Rows:     [][]any{{"Data"}},
				Bordered: true,
			},
			contains: []string{
				`class="table table-bordered"`,
			},
		},
		{
			name: "custom class table",
			table: Table{
				Headers: []string{"Col1"},
				Rows:    [][]any{{"Data"}},
				Class:   "custom-table",
			},
			contains: []string{
				`class="custom-table"`,
			},
		},
		{
			name: "striped and bordered table",
			table: Table{
				Headers:  []string{"Col1"},
				Rows:     [][]any{{"Data"}},
				Striped:  true,
				Bordered: true,
			},
			contains: []string{
				"table-striped",
				"table-bordered",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := element.NewBuilder()
			tt.table.Render(b)
			got := b.String()

			for _, want := range tt.contains {
				if !strings.Contains(got, want) {
					t.Errorf("Table.Render() missing %q\ngot: %s", want, got)
				}
			}
		})
	}
}

func TestTable_Render_NoHeaders(t *testing.T) {
	b := element.NewBuilder()
	table := Table{
		Rows: [][]any{{"A", "B"}},
	}
	table.Render(b)
	got := b.String()

	if strings.Contains(got, "<thead>") {
		t.Errorf("Table without headers should not render thead, got: %s", got)
	}
}

type tdLink struct{ href, label string }

func (l tdLink) Render(b *element.Builder) (x any) {
	b.A("href", l.href).T(l.label)
	return
}

func TestTable_Render_Advanced(t *testing.T) {
	tests := []struct {
		name        string
		table       Table
		contains    []string
		notContains []string
	}{
		{
			name: "columns with alignment and classes",
			table: Table{
				Columns: []Column{
					{Header: "Item"},
					{Header: "Qty", Align: "right", Class: "qty-cell", HeaderClass: "qty-head"},
				},
				Rows: [][]any{{"Widget", 42}},
			},
			contains: []string{
				`<th class="qty-head" style="text-align:right">Qty</th>`,
				`<td class="qty-cell" style="text-align:right">42</td>`,
				"<td>Widget</td>",
			},
		},
		{
			name: "caption and footer rows",
			table: Table{
				Headers:    []string{"Product", "Price"},
				Rows:       [][]any{{"Widget", 9.99}},
				FooterRows: [][]any{{"Total", 9.99}},
				Caption:    "Q3 Sales",
			},
			contains: []string{
				"<caption>Q3 Sales</caption>",
				"<tfoot>",
				"<td>Total</td>",
				"</tfoot>",
			},
		},
		{
			name: "empty message when no rows",
			table: Table{
				Headers:      []string{"A", "B", "C"},
				EmptyMessage: "No records found",
			},
			contains: []string{
				`<td class="table-empty" colspan="3">No records found</td>`,
			},
		},
		{
			name: "component cells render in place",
			table: Table{
				Headers: []string{"Name", "Action"},
				Rows: [][]any{
					{"Alice", tdLink{href: "/users/1", label: "Edit"}},
				},
			},
			contains: []string{
				`<td><a href="/users/1">Edit</a></td>`,
			},
		},
		{
			name: "row class hook",
			table: Table{
				Headers: []string{"Status"},
				Rows:    [][]any{{"ok"}, {"failed"}},
				RowClass: func(i int, row []any) string {
					if row[0] == "failed" {
						return "row-error"
					}
					return ""
				},
			},
			contains: []string{
				`<tr class="row-error">`,
				"<tr>",
			},
		},
		{
			name: "responsive wrapper and hover",
			table: Table{
				Headers:    []string{"A"},
				Rows:       [][]any{{"x"}},
				Hover:      true,
				Responsive: true,
			},
			contains: []string{
				`<div class="table-responsive">`,
				"table-hover",
			},
		},
		{
			name: "extra attrs pass through",
			table: Table{
				Headers: []string{"A"},
				Rows:    [][]any{{"x"}},
				Attrs:   []string{"id", "results", "data-page", "2"},
			},
			contains: []string{
				`id="results"`,
				`data-page="2"`,
			},
		},
		{
			name: "scalar cell types render without quotes or padding",
			table: Table{
				Rows: [][]any{{int64(7), 3.14, true, nil}},
			},
			contains: []string{
				"<td>7</td>",
				"<td>3.14</td>",
				"<td>true</td>",
				"<td></td>",
			},
			notContains: []string{"<thead>"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := element.NewBuilder()
			tt.table.Render(b)
			got := b.String()

			for _, want := range tt.contains {
				if !strings.Contains(got, want) {
					t.Errorf("Table.Render() missing %q\ngot: %s", want, got)
				}
			}
			for _, bad := range tt.notContains {
				if strings.Contains(got, bad) {
					t.Errorf("Table.Render() should not contain %q\ngot: %s", bad, got)
				}
			}
		})
	}
}
