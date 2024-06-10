package element

import (
	"testing"
)

func TestNewBuilderWithAlias(t *testing.T) {
	tests := []struct {
		name  string
		wantB *Builder
	}{
		{name: "New Builder - with alias", wantB: &Builder{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotB := NB(); gotB == nil {
				t.Errorf("NewBuilder() = %v, want not nil", gotB)
			}
		})
	}
}

func TestBuilderWithAlias_String(t *testing.T) {
	const want = `<html><body><p>Hello World!</p></body></html>`
	t.Run("Builder String()  - with alias", func(t *testing.T) {
		b := NB()
		b.E("html").R(
			b.E("body").R(
				b.E("p").R(b.T("Hello World!"))))
		if got := b.S(); got != want {
			t.Errorf("String() = %v, want %v", got, want)
		}
	})
}

func TestBuildElementWithDelayedRenderWithAlias(t *testing.T) {
	const want = `<div><p style="color:red">Hello World!</p></div>`

	t.Run("Builder - Create element with delayed render - with alias", func(t *testing.T) {
		b := NB()

		pDelayedRender := b.ENR("p")

		// We could conditionally add a color style
		pDelayedRender.AA("style", "color:red")

		b.E("div").R(pDelayedRender.RenderOpeningTag().R(
			b.T("Hello World!")))

		if got := b.S(); got != want {
			t.Errorf("String() = %v, want %v", got, want)
		}
	})
}

func TestBuilder_WriteStringWithAlias(t *testing.T) {
	const want = `<span>Testing, testing</span>`
	t.Run("Builder WriteString() - with alias", func(t *testing.T) {
		b := NB()
		b.E("span").R(
			b.WS("Testing, testing"),
		)
		if got := b.S(); got != want {
			t.Errorf("String() = %v, want %v", got, want)
		}
	})
}
