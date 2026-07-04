package components

import (
	"strings"
	"testing"

	"github.com/rohanthewiz/element"
)

func TestDropdown_Render(t *testing.T) {
	b := element.NewBuilder()
	Dropdown{
		Label: "Account",
		Items: []DropdownItem{
			{Label: "Profile", Href: "/profile"},
			{Divider: true},
			{Label: "Sign out", Href: "/logout"},
		},
	}.Render(b)
	got := b.String()

	for _, want := range []string{
		`<details class="dropdown">`,
		`<summary class="dropdown-toggle">Account</summary>`,
		`<ul class="dropdown-menu">`,
		`<a class="dropdown-link" href="/profile">Profile</a>`,
		`<li class="dropdown-divider" role="separator">`,
		`<a class="dropdown-link" href="/logout">Sign out</a>`,
	} {
		if !strings.Contains(got, want) {
			t.Errorf("Dropdown.Render() missing %q\ngot: %s", want, got)
		}
	}
}
