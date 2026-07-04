package components

import (
	"strings"
	"testing"

	"github.com/rohanthewiz/element"
)

func TestButton_Render(t *testing.T) {
	tests := []struct {
		name     string
		button   Button
		contains []string
	}{
		{
			name:   "default button",
			button: Button{Text: "Save"},
			contains: []string{
				`<button class="btn btn-primary" type="button">Save</button>`,
			},
		},
		{
			name:   "submit with variant and size",
			button: Button{Text: "Delete", Type: "submit", Variant: "danger", Size: "sm"},
			contains: []string{
				`class="btn btn-danger btn-sm"`, `type="submit"`,
			},
		},
		{
			name:   "disabled button",
			button: Button{Text: "Wait", Disabled: true},
			contains: []string{
				`disabled="disabled"`,
			},
		},
		{
			name:   "link styled as button",
			button: Button{Text: "Docs", Href: "/docs", Variant: "secondary"},
			contains: []string{
				`<a class="btn btn-secondary" href="/docs" role="button">Docs</a>`,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := element.NewBuilder()
			tt.button.Render(b)
			got := b.String()
			for _, want := range tt.contains {
				if !strings.Contains(got, want) {
					t.Errorf("Button.Render() missing %q\ngot: %s", want, got)
				}
			}
		})
	}
}

func TestBadge_Render(t *testing.T) {
	b := element.NewBuilder()
	Badge{Text: "3 new", Type: AlertSuccess, Pill: true}.Render(b)
	got := b.String()

	if !strings.Contains(got, `<span class="badge badge-success badge-pill">3 new</span>`) {
		t.Errorf("Badge.Render() unexpected output: %s", got)
	}
}

func TestProgressBar_Render(t *testing.T) {
	b := element.NewBuilder()
	ProgressBar{Value: 30, Max: 60, Label: "Upload", ShowValue: true}.Render(b)
	got := b.String()

	for _, want := range []string{
		`role="progressbar"`,
		`aria-valuenow="30"`, `aria-valuemax="60"`,
		`aria-label="Upload"`,
		`style="width:50%"`,
		">50%</div>",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("ProgressBar.Render() missing %q\ngot: %s", want, got)
		}
	}
}

func TestProgressBar_ClampsValue(t *testing.T) {
	b := element.NewBuilder()
	ProgressBar{Value: 150}.Render(b)
	got := b.String()

	if !strings.Contains(got, `aria-valuenow="100"`) || !strings.Contains(got, "width:100%") {
		t.Errorf("ProgressBar should clamp value to max, got: %s", got)
	}
}
