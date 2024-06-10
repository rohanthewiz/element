package element

import (
	"testing"
)

func TestNewBuilder(t *testing.T) {
	tests := []struct {
		name  string
		wantB *Builder
	}{
		{name: "New Builder", wantB: &Builder{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotB := NewBuilder(); gotB == nil {
				t.Errorf("NewBuilder() = %v, want not nil", gotB)
			}
		})
	}
}

func TestBuilder_String(t *testing.T) {
	const want = `<html><body><p>Hello World!</p></body></html>`
	t.Run("Builder String()", func(t *testing.T) {
		b := NewBuilder()
		b.Ele("html").R(
			b.Ele("body").R(
				b.Ele("p").R(b.Text("Hello World!"))))
		if got := b.String(); got != want {
			t.Errorf("String() = %v, want %v", got, want)
		}
	})
}

func TestBuildElementWithDelayedRender(t *testing.T) {
	const want = `<div><p style="color:red">Hello World!</p></div>`

	t.Run("Builder - Create element with delayed render", func(t *testing.T) {
		b := NewBuilder()

		pDelayedRender := b.EleNoRender("p")

		// We could conditionally add a color style
		pDelayedRender.AddAttributes("style", "color:red")

		b.Ele("div").R(pDelayedRender.RenderOpeningTag().R(
			b.Text("Hello World!")))

		if got := b.String(); got != want {
			t.Errorf("String() = %v, want %v", got, want)
		}
	})
}

func TestBuilder_WriteString(t *testing.T) {
	const want = `<span>Testing, testing</span>`
	t.Run("Builder WriteString()", func(t *testing.T) {
		b := NewBuilder()
		b.Ele("span").R(
			b.WriteString("Testing, testing"),
		)
		if got := b.String(); got != want {
			t.Errorf("String() = %v, want %v", got, want)
		}
	})
}

func TestBuilder_WriteBytes(t *testing.T) {
	const want = `<span>MyDoc: Test this</span>`
	t.Run("Builder WriteBytes()", func(t *testing.T) {
		b := NewBuilder()
		b.Ele("span").R(
			b.WriteBytes([]byte("MyDoc: Test this")),
		)

		if got := b.String(); got != want {
			t.Errorf("String() = %v, want %v", got, want)
		}
	})
}
