package components

import (
	"strings"
	"testing"

	"github.com/rohanthewiz/element"
)

func TestBreadcrumb_Render(t *testing.T) {
	tests := []struct {
		name     string
		bc       Breadcrumb
		contains []string
	}{
		{
			name: "basic breadcrumb",
			bc: Breadcrumb{
				Items: []BreadcrumbItem{
					{Label: "Home", Href: "/"},
					{Label: "Products", Href: "/products"},
					{Label: "Current Page"},
				},
			},
			contains: []string{
				"<nav",
				`class="breadcrumb"`,
				`aria-label="breadcrumb"`,
				`class="breadcrumb-list"`,
				`class="breadcrumb-item"`,
				`href="/"`,
				"Home",
				`href="/products"`,
				"Products",
				`class="breadcrumb-current"`,
				`aria-current="page"`,
				"Current Page",
			},
		},
		{
			name: "breadcrumb with default separator",
			bc: Breadcrumb{
				Items: []BreadcrumbItem{
					{Label: "Home", Href: "/"},
					{Label: "Page"},
				},
			},
			contains: []string{
				`class="breadcrumb-separator"`,
				" / ",
			},
		},
		{
			name: "breadcrumb with custom separator",
			bc: Breadcrumb{
				Separator: ">",
				Items: []BreadcrumbItem{
					{Label: "Home", Href: "/"},
					{Label: "Page"},
				},
			},
			contains: []string{
				" > ",
			},
		},
		{
			name: "single item breadcrumb",
			bc: Breadcrumb{
				Items: []BreadcrumbItem{
					{Label: "Home"},
				},
			},
			contains: []string{
				`class="breadcrumb-current"`,
				"Home",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := element.NewBuilder()
			tt.bc.Render(b)
			got := b.String()

			for _, want := range tt.contains {
				if !strings.Contains(got, want) {
					t.Errorf("Breadcrumb.Render() missing %q\ngot: %s", want, got)
				}
			}
		})
	}
}

func TestBreadcrumb_Render_LastItemNoLink(t *testing.T) {
	b := element.NewBuilder()
	bc := Breadcrumb{
		Items: []BreadcrumbItem{
			{Label: "Home", Href: "/"},
			{Label: "Current", Href: "/current"}, // Has href but is last
		},
	}
	bc.Render(b)
	got := b.String()

	// The last item should be rendered as current (span), not as a link
	// Check that "Current" is in a span with breadcrumb-current class
	if !strings.Contains(got, `<span class="breadcrumb-current"`) {
		t.Errorf("Last breadcrumb item should be rendered as current span, got: %s", got)
	}
}

func TestBreadcrumb_Render_EmptyHref(t *testing.T) {
	b := element.NewBuilder()
	bc := Breadcrumb{
		Items: []BreadcrumbItem{
			{Label: "No Link"}, // Empty href
			{Label: "Current"},
		},
	}
	bc.Render(b)
	got := b.String()

	// Item with empty href should also render as current (span)
	if strings.Count(got, "breadcrumb-current") != 2 {
		t.Errorf("Items with empty href should render as current, got: %s", got)
	}
}

func TestBreadcrumb_Render_SeparatorAriaHidden(t *testing.T) {
	b := element.NewBuilder()
	bc := Breadcrumb{
		Items: []BreadcrumbItem{
			{Label: "Home", Href: "/"},
			{Label: "Page"},
		},
	}
	bc.Render(b)
	got := b.String()

	if !strings.Contains(got, `aria-hidden="true"`) {
		t.Errorf("Separator should have aria-hidden=true, got: %s", got)
	}
}
