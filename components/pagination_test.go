package components

import (
	"strings"
	"testing"

	"github.com/rohanthewiz/element"
)

func TestPagination_Render(t *testing.T) {
	tests := []struct {
		name     string
		pg       Pagination
		contains []string
	}{
		{
			name: "basic pagination",
			pg: Pagination{
				CurrentPage: 3,
				TotalPages:  10,
				BaseURL:     "/page/%d",
			},
			contains: []string{
				"<nav",
				`class="pagination"`,
				`aria-label="Page navigation"`,
				`class="pagination-list"`,
				`href="/page/2"`, // Prev
				"‹ Prev",
				`href="/page/4"`, // Next
				"Next ›",
				`class="pagination-item active"`,
				`class="pagination-current"`,
				">3<", // Current page
			},
		},
		{
			name: "first page - no prev",
			pg: Pagination{
				CurrentPage: 1,
				TotalPages:  5,
				BaseURL:     "/page/%d",
			},
			contains: []string{
				"Next ›",
			},
		},
		{
			name: "last page - no next",
			pg: Pagination{
				CurrentPage: 5,
				TotalPages:  5,
				BaseURL:     "/page/%d",
			},
			contains: []string{
				"‹ Prev",
			},
		},
		{
			name: "with first and last links",
			pg: Pagination{
				CurrentPage: 5,
				TotalPages:  10,
				BaseURL:     "/page/%d",
				ShowFirst:   true,
				ShowLast:    true,
			},
			contains: []string{
				"« First",
				`href="/page/1"`,
				"Last »",
				`href="/page/10"`,
			},
		},
		{
			name: "page numbers shown",
			pg: Pagination{
				CurrentPage: 5,
				TotalPages:  10,
				BaseURL:     "/page/%d",
			},
			contains: []string{
				`href="/page/3"`, // CurrentPage - 2
				`href="/page/4"`, // CurrentPage - 1
				// Page 5 is current (no href)
				`href="/page/6"`, // CurrentPage + 1
				`href="/page/7"`, // CurrentPage + 2
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := element.NewBuilder()
			tt.pg.Render(b)
			got := b.String()

			for _, want := range tt.contains {
				if !strings.Contains(got, want) {
					t.Errorf("Pagination.Render() missing %q\ngot: %s", want, got)
				}
			}
		})
	}
}

func TestPagination_Render_SinglePage(t *testing.T) {
	b := element.NewBuilder()
	pg := Pagination{
		CurrentPage: 1,
		TotalPages:  1,
		BaseURL:     "/page/%d",
	}
	pg.Render(b)
	got := b.String()

	if got != "" {
		t.Errorf("Pagination with single page should render nothing, got: %s", got)
	}
}

func TestPagination_Render_ZeroPages(t *testing.T) {
	b := element.NewBuilder()
	pg := Pagination{
		CurrentPage: 1,
		TotalPages:  0,
		BaseURL:     "/page/%d",
	}
	pg.Render(b)
	got := b.String()

	if got != "" {
		t.Errorf("Pagination with zero pages should render nothing, got: %s", got)
	}
}

func TestPagination_Render_FirstPage_NoFirstLink(t *testing.T) {
	b := element.NewBuilder()
	pg := Pagination{
		CurrentPage: 1,
		TotalPages:  5,
		BaseURL:     "/page/%d",
		ShowFirst:   true,
	}
	pg.Render(b)
	got := b.String()

	if strings.Contains(got, "« First") {
		t.Errorf("First page should not show 'First' link, got: %s", got)
	}
}

func TestPagination_Render_LastPage_NoLastLink(t *testing.T) {
	b := element.NewBuilder()
	pg := Pagination{
		CurrentPage: 5,
		TotalPages:  5,
		BaseURL:     "/page/%d",
		ShowLast:    true,
	}
	pg.Render(b)
	got := b.String()

	if strings.Contains(got, "Last »") {
		t.Errorf("Last page should not show 'Last' link, got: %s", got)
	}
}

func TestPagination_Render_FirstPage_NoPrev(t *testing.T) {
	b := element.NewBuilder()
	pg := Pagination{
		CurrentPage: 1,
		TotalPages:  5,
		BaseURL:     "/page/%d",
	}
	pg.Render(b)
	got := b.String()

	if strings.Contains(got, "‹ Prev") {
		t.Errorf("First page should not show 'Prev' link, got: %s", got)
	}
}

func TestPagination_Render_LastPage_NoNext(t *testing.T) {
	b := element.NewBuilder()
	pg := Pagination{
		CurrentPage: 5,
		TotalPages:  5,
		BaseURL:     "/page/%d",
	}
	pg.Render(b)
	got := b.String()

	if strings.Contains(got, "Next ›") {
		t.Errorf("Last page should not show 'Next' link, got: %s", got)
	}
}

func TestPagination_Render_PageRangeAtStart(t *testing.T) {
	b := element.NewBuilder()
	pg := Pagination{
		CurrentPage: 1,
		TotalPages:  10,
		BaseURL:     "/page/%d",
	}
	pg.Render(b)
	got := b.String()

	// Should show pages 1, 2, 3 (current-2 to current+2, but min is 1)
	if !strings.Contains(got, `href="/page/2"`) {
		t.Errorf("Should show page 2, got: %s", got)
	}
	if !strings.Contains(got, `href="/page/3"`) {
		t.Errorf("Should show page 3, got: %s", got)
	}
}

func TestPagination_Render_PageRangeAtEnd(t *testing.T) {
	b := element.NewBuilder()
	pg := Pagination{
		CurrentPage: 10,
		TotalPages:  10,
		BaseURL:     "/page/%d",
	}
	pg.Render(b)
	got := b.String()

	// Should show pages 8, 9, 10 (current-2 to current+2, but max is 10)
	if !strings.Contains(got, `href="/page/8"`) {
		t.Errorf("Should show page 8, got: %s", got)
	}
	if !strings.Contains(got, `href="/page/9"`) {
		t.Errorf("Should show page 9, got: %s", got)
	}
}
