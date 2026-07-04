package components

import (
	"strings"
	"testing"

	"github.com/rohanthewiz/element"
)

func TestAccordion_Render(t *testing.T) {
	b := element.NewBuilder()
	Accordion{
		Items: []AccordionItem{
			{Title: "Section One", Content: "First content", Open: true},
			{Title: "Section Two", Body: Badge{Text: "42"}},
		},
		Exclusive: true,
		Name:      "faq",
	}.Render(b)
	got := b.String()

	for _, want := range []string{
		`<div class="accordion">`,
		`name="faq"`,
		`open="open"`,
		`<summary class="accordion-title">Section One</summary>`,
		"First content",
		`<span class="badge badge-neutral">42</span>`,
		"</details>",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("Accordion.Render() missing %q\ngot: %s", want, got)
		}
	}

	if strings.Count(got, "<details") != 2 {
		t.Errorf("expected 2 details elements, got: %s", got)
	}
	// Only the first item is open
	if strings.Count(got, `open="open"`) != 1 {
		t.Errorf("expected exactly one open section, got: %s", got)
	}
}

func TestAccordion_NonExclusiveHasNoName(t *testing.T) {
	b := element.NewBuilder()
	Accordion{Items: []AccordionItem{{Title: "A", Content: "a"}}}.Render(b)
	got := b.String()

	if strings.Contains(got, "name=") {
		t.Errorf("non-exclusive accordion should not set name, got: %s", got)
	}
}
