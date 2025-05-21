package element

import (
	"regexp"
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

func TestBuilder_WriteString(t *testing.T) {
	const want = `<span.*>Testing, testing</span>`
	t.Run("Builder WriteString()", func(t *testing.T) {
		b := NewBuilder()
		b.Ele("span").R(
			b.WriteString("Testing, testing"),
		)
		got := b.String()
		if !regexp.MustCompile(want).MatchString(got) {
			t.Errorf("String() = %v, want %v", got, want)
		}
	})
}

func TestBuilder_WriteBytes(t *testing.T) {
	const want = `<span.*>MyDoc: Test this</span>`
	t.Run("Builder WriteBytes()", func(t *testing.T) {
		b := NewBuilder()
		b.Ele("span").R(
			b.WriteBytes([]byte("MyDoc: Test this")),
		)

		got := b.String()
		if !regexp.MustCompile(want).MatchString(got) {
			t.Errorf("String() = %v, want %v", got, want)
		}
	})
}

func TestBuilder_F(t *testing.T) {
	tests := []struct {
		name         string
		formatString string
		args         []any
		expected     string
	}{
		{
			name:         "simple string",
			formatString: "Hello",
			args:         nil,
			expected:     "Hello",
		},
		{
			name:         "string with one argument",
			formatString: "Hello, %s!",
			args:         []any{"world"},
			expected:     "Hello, world!",
		},
		{
			name:         "multiple arguments",
			formatString: "%s %d %s",
			args:         []any{"test", 123, "abc"},
			expected:     "test 123 abc",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := NewBuilder()
			b.F(tt.formatString, tt.args...)

			result := b.String() // Assuming Builder has a String() method to get the content
			if result != tt.expected {
				t.Errorf("F() = %q, want %q", result, tt.expected)
			}
		})
	}
}

// BodyComponent implements Component interface for testing
type BodyComponent struct{}

func (tc BodyComponent) Render(b *Builder) any {
	b.Div().T("Test Content")
	return nil
}

func TestBuilder_HtmlPage(t *testing.T) {
	tests := []struct {
		name             string
		styles           string
		headWithoutStyle string
		body             Component
		want             string
	}{
		{
			name:             "Basic HTML page",
			styles:           "body { color: black; }",
			headWithoutStyle: "<title>Test</title>",
			body:             BodyComponent{},
			want:             "<!DOCTYPE html><html><head><title>Test</title><style>body { color: black; }</style></head><body><div>Test Content</div></body></html>",
		},
		{
			name:             "HTML page without styles",
			styles:           "",
			headWithoutStyle: "<title>Test</title>",
			body:             BodyComponent{},
			want:             "<!DOCTYPE html><html><head><title>Test</title></head><body><div>Test Content</div></body></html>",
		},
		{
			name:             "HTML page without head content",
			styles:           "body { color: tan; }",
			headWithoutStyle: "",
			body:             BodyComponent{},
			want:             "<!DOCTYPE html><html><head><style>body { color: tan; }</style></head><body><div>Test Content</div></body></html>",
		},
		{
			name:             "HTML page without styles and head content",
			styles:           "",
			headWithoutStyle: "",
			body:             BodyComponent{},
			want:             "<!DOCTYPE html><html><head></head><body><div>Test Content</div></body></html>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := NewBuilder()
			got := b.HtmlPage(tt.styles, tt.headWithoutStyle, tt.body)
			if got != tt.want {
				t.Errorf("HtmlPage() = %v, want %v", got, tt.want)
			}
		})
	}
}
