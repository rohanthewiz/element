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
