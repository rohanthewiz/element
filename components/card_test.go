package components

import (
	"strings"
	"testing"

	"github.com/rohanthewiz/element"
)

func TestCard_Render(t *testing.T) {
	tests := []struct {
		name     string
		card     Card
		contains []string
	}{
		{
			name: "basic card with title and body",
			card: Card{
				Title: "Card Title",
				Body:  "Card body content",
			},
			contains: []string{
				`class="card"`,
				`class="card-header"`,
				`class="card-title"`,
				"Card Title",
				`class="card-body"`,
				"Card body content",
			},
		},
		{
			name: "card with footer",
			card: Card{
				Title:  "Title",
				Body:   "Body",
				Footer: "Footer text",
			},
			contains: []string{
				`class="card-footer"`,
				"Footer text",
			},
		},
		{
			name: "card with custom class",
			card: Card{
				Title: "Title",
				Body:  "Body",
				Class: "custom-card highlighted",
			},
			contains: []string{
				`class="card custom-card highlighted"`,
			},
		},
		{
			name: "card without title",
			card: Card{
				Body: "Just body content",
			},
			contains: []string{
				`class="card"`,
				`class="card-body"`,
				"Just body content",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := element.NewBuilder()
			tt.card.Render(b)
			got := b.String()

			for _, want := range tt.contains {
				if !strings.Contains(got, want) {
					t.Errorf("Card.Render() missing %q\ngot: %s", want, got)
				}
			}
		})
	}
}

func TestCard_Render_NoTitle(t *testing.T) {
	b := element.NewBuilder()
	card := Card{Body: "Body only"}
	card.Render(b)
	got := b.String()

	if strings.Contains(got, "card-header") {
		t.Errorf("Card without title should not render header, got: %s", got)
	}
}

func TestCard_Render_NoFooter(t *testing.T) {
	b := element.NewBuilder()
	card := Card{Title: "Title", Body: "Body"}
	card.Render(b)
	got := b.String()

	if strings.Contains(got, "card-footer") {
		t.Errorf("Card without footer should not render footer, got: %s", got)
	}
}

// TestBodyComponent implements element.Component for testing BodyComponent
type TestBodyComponent struct {
	Content string
}

func (c TestBodyComponent) Render(b *element.Builder) any {
	b.Span().T(c.Content)
	return nil
}

func TestCard_Render_WithBodyComponent(t *testing.T) {
	b := element.NewBuilder()
	card := Card{
		Title:         "Title",
		BodyComponent: TestBodyComponent{Content: "Component content"},
	}
	card.Render(b)
	got := b.String()

	if !strings.Contains(got, "<span>Component content</span>") {
		t.Errorf("Card.Render() should render BodyComponent, got: %s", got)
	}
}
