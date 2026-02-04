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
