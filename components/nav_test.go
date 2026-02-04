package components

import (
	"strings"
	"testing"

	"github.com/rohanthewiz/element"
)

func TestNav_Render(t *testing.T) {
	tests := []struct {
		name     string
		nav      Nav
		contains []string
	}{
		{
			name: "basic nav with items",
			nav: Nav{
				Items: []NavItem{
					{Label: "Home", Href: "/"},
					{Label: "About", Href: "/about"},
				},
			},
			contains: []string{
				"<nav",
				`class="nav"`,
				`class="nav-list"`,
				`class="nav-item"`,
				`class="nav-link"`,
				`href="/"`,
				"Home",
				`href="/about"`,
				"About",
			},
		},
		{
			name: "nav with brand",
			nav: Nav{
				Brand: "MyBrand",
				Items: []NavItem{
					{Label: "Home", Href: "/"},
				},
			},
			contains: []string{
				`class="nav-brand"`,
				"MyBrand",
			},
		},
		{
			name: "nav with active item",
			nav: Nav{
				Items: []NavItem{
					{Label: "Home", Href: "/", Active: true},
					{Label: "About", Href: "/about"},
				},
			},
			contains: []string{
				`class="nav-item active"`,
			},
		},
		{
			name: "nav with custom class",
			nav: Nav{
				Class: "main-nav",
				Items: []NavItem{
					{Label: "Home", Href: "/"},
				},
			},
			contains: []string{
				`class="nav main-nav"`,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := element.NewBuilder()
			tt.nav.Render(b)
			got := b.String()

			for _, want := range tt.contains {
				if !strings.Contains(got, want) {
					t.Errorf("Nav.Render() missing %q\ngot: %s", want, got)
				}
			}
		})
	}
}

func TestNav_Render_NoBrand(t *testing.T) {
	b := element.NewBuilder()
	nav := Nav{
		Items: []NavItem{{Label: "Home", Href: "/"}},
	}
	nav.Render(b)
	got := b.String()

	if strings.Contains(got, "nav-brand") {
		t.Errorf("Nav without brand should not render brand, got: %s", got)
	}
}

func TestNav_Render_EmptyItems(t *testing.T) {
	b := element.NewBuilder()
	nav := Nav{
		Brand: "Brand",
		Items: []NavItem{},
	}
	nav.Render(b)
	got := b.String()

	// Should still render the nav structure
	if !strings.Contains(got, "<nav") {
		t.Errorf("Nav.Render() should render nav element even with empty items, got: %s", got)
	}
	if !strings.Contains(got, "nav-list") {
		t.Errorf("Nav.Render() should render nav-list even with empty items, got: %s", got)
	}
}
