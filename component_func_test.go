package element

import (
	"strings"
	"testing"
)

func TestCompFunc(t *testing.T) {
	b := NewBuilder()
	comp := CompFunc(func(b *Builder) (x any) {
		b.H1().T("Hello!")
		return
	})

	// CompFunc satisfies the Component interface
	RenderComponents(b, comp)

	if got := b.String(); got != "<h1>Hello!</h1>" {
		t.Errorf("CompFunc render = %q, want %q", got, "<h1>Hello!</h1>")
	}
}

func TestForEach2(t *testing.T) {
	b := NewBuilder()
	items := []string{"alpha", "beta"}

	b.Ul().R(
		ForEach2(items, func(item string, i int) {
			b.Li().F("%d. %s", i+1, item)
		}),
	)

	got := b.String()
	for _, want := range []string{"<li>1. alpha</li>", "<li>2. beta</li>"} {
		if !strings.Contains(got, want) {
			t.Errorf("ForEach2 output missing %q, got: %s", want, got)
		}
	}
}
