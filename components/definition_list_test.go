package components

import (
	"strings"
	"testing"

	"github.com/rohanthewiz/element"
)

func TestDefinitionList_Render(t *testing.T) {
	tests := []struct {
		name     string
		dl       DefinitionList
		contains []string
	}{
		{
			name: "basic definition list",
			dl: DefinitionList{
				Items: []Definition{
					{Term: "HTML", Definition: "HyperText Markup Language"},
					{Term: "CSS", Definition: "Cascading Style Sheets"},
				},
			},
			contains: []string{
				"<dl",
				`class="definition-list"`,
				`class="definition-term"`,
				"<dt",
				"HTML",
				`class="definition-desc"`,
				"<dd",
				"HyperText Markup Language",
				"CSS",
				"Cascading Style Sheets",
				"</dl>",
			},
		},
		{
			name: "definition list with custom class",
			dl: DefinitionList{
				Class: "glossary",
				Items: []Definition{
					{Term: "Term", Definition: "Def"},
				},
			},
			contains: []string{
				`class="definition-list glossary"`,
			},
		},
		{
			name: "single definition",
			dl: DefinitionList{
				Items: []Definition{
					{Term: "API", Definition: "Application Programming Interface"},
				},
			},
			contains: []string{
				"API",
				"Application Programming Interface",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := element.NewBuilder()
			tt.dl.Render(b)
			got := b.String()

			for _, want := range tt.contains {
				if !strings.Contains(got, want) {
					t.Errorf("DefinitionList.Render() missing %q\ngot: %s", want, got)
				}
			}
		})
	}
}

func TestDefinitionList_Render_Empty(t *testing.T) {
	b := element.NewBuilder()
	dl := DefinitionList{
		Items: []Definition{},
	}
	dl.Render(b)
	got := b.String()

	// Should still render dl element
	if !strings.Contains(got, "<dl") {
		t.Errorf("DefinitionList.Render() should render dl even when empty, got: %s", got)
	}
	// But should not have dt or dd
	if strings.Contains(got, "<dt") {
		t.Errorf("Empty DefinitionList should not have dt elements, got: %s", got)
	}
	if strings.Contains(got, "<dd") {
		t.Errorf("Empty DefinitionList should not have dd elements, got: %s", got)
	}
}

func TestDefinitionList_Render_OrderPreserved(t *testing.T) {
	b := element.NewBuilder()
	dl := DefinitionList{
		Items: []Definition{
			{Term: "First", Definition: "1st"},
			{Term: "Second", Definition: "2nd"},
			{Term: "Third", Definition: "3rd"},
		},
	}
	dl.Render(b)
	got := b.String()

	// Check order by finding positions
	firstPos := strings.Index(got, "First")
	secondPos := strings.Index(got, "Second")
	thirdPos := strings.Index(got, "Third")

	if firstPos > secondPos || secondPos > thirdPos {
		t.Errorf("DefinitionList.Render() should preserve order, got: %s", got)
	}
}
