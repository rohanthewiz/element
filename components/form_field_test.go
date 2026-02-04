package components

import (
	"strings"
	"testing"

	"github.com/rohanthewiz/element"
)

func TestFormField_Render(t *testing.T) {
	tests := []struct {
		name     string
		field    FormField
		contains []string
	}{
		{
			name: "basic text field",
			field: FormField{
				Label: "Username",
				Name:  "username",
			},
			contains: []string{
				`class="form-field"`,
				`class="form-label"`,
				`for="username"`,
				"Username",
				`type="text"`, // Default type
				`id="username"`,
				`name="username"`,
				`class="form-input"`,
			},
		},
		{
			name: "email field",
			field: FormField{
				Label: "Email",
				Name:  "email",
				Type:  "email",
			},
			contains: []string{
				`type="email"`,
			},
		},
		{
			name: "password field",
			field: FormField{
				Label: "Password",
				Name:  "password",
				Type:  "password",
			},
			contains: []string{
				`type="password"`,
			},
		},
		{
			name: "field with placeholder",
			field: FormField{
				Label:       "Name",
				Name:        "name",
				Placeholder: "Enter your name",
			},
			contains: []string{
				`placeholder="Enter your name"`,
			},
		},
		{
			name: "field with value",
			field: FormField{
				Label: "Name",
				Name:  "name",
				Value: "John Doe",
			},
			contains: []string{
				`value="John Doe"`,
			},
		},
		{
			name: "required field",
			field: FormField{
				Label:    "Email",
				Name:     "email",
				Required: true,
			},
			contains: []string{
				`required="required"`,
				`class="required"`,
				" *",
			},
		},
		{
			name: "field with help text",
			field: FormField{
				Label:    "Username",
				Name:     "username",
				HelpText: "Choose a unique username",
			},
			contains: []string{
				`class="form-help"`,
				"Choose a unique username",
			},
		},
		{
			name: "field with error",
			field: FormField{
				Label: "Email",
				Name:  "email",
				Error: "Invalid email format",
			},
			contains: []string{
				`class="form-field has-error"`,
				`class="form-error"`,
				"Invalid email format",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := element.NewBuilder()
			tt.field.Render(b)
			got := b.String()

			for _, want := range tt.contains {
				if !strings.Contains(got, want) {
					t.Errorf("FormField.Render() missing %q\ngot: %s", want, got)
				}
			}
		})
	}
}

func TestFormField_Render_DefaultType(t *testing.T) {
	b := element.NewBuilder()
	field := FormField{
		Label: "Test",
		Name:  "test",
		// Type not specified
	}
	field.Render(b)
	got := b.String()

	if !strings.Contains(got, `type="text"`) {
		t.Errorf("FormField should default to type=text, got: %s", got)
	}
}

func TestFormField_Render_NoPlaceholder(t *testing.T) {
	b := element.NewBuilder()
	field := FormField{
		Label: "Test",
		Name:  "test",
	}
	field.Render(b)
	got := b.String()

	if strings.Contains(got, "placeholder=") {
		t.Errorf("FormField without placeholder should not have placeholder attr, got: %s", got)
	}
}

func TestFormField_Render_NoValue(t *testing.T) {
	b := element.NewBuilder()
	field := FormField{
		Label: "Test",
		Name:  "test",
	}
	field.Render(b)
	got := b.String()

	if strings.Contains(got, "value=") {
		t.Errorf("FormField without value should not have value attr, got: %s", got)
	}
}

func TestFormField_Render_ErrorHidesHelpText(t *testing.T) {
	b := element.NewBuilder()
	field := FormField{
		Label:    "Test",
		Name:     "test",
		HelpText: "Help text",
		Error:    "Error message",
	}
	field.Render(b)
	got := b.String()

	// When there's an error, help text should not be shown
	if strings.Contains(got, "form-help") {
		t.Errorf("FormField with error should not show help text, got: %s", got)
	}
	if !strings.Contains(got, "form-error") {
		t.Errorf("FormField with error should show error, got: %s", got)
	}
}

func TestFormField_Render_NotRequired(t *testing.T) {
	b := element.NewBuilder()
	field := FormField{
		Label: "Test",
		Name:  "test",
	}
	field.Render(b)
	got := b.String()

	if strings.Contains(got, "required=") {
		t.Errorf("Non-required FormField should not have required attr, got: %s", got)
	}
	if strings.Contains(got, `class="required"`) {
		t.Errorf("Non-required FormField should not have required span, got: %s", got)
	}
}
