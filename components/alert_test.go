package components

import (
	"strings"
	"testing"

	"github.com/rohanthewiz/element"
)

func TestAlert_Render(t *testing.T) {
	tests := []struct {
		name     string
		alert    Alert
		contains []string
	}{
		{
			name: "basic info alert",
			alert: Alert{
				Type:    AlertInfo,
				Message: "This is an info message",
			},
			contains: []string{
				`class="alert alert-info"`,
				`role="alert"`,
				"This is an info message",
			},
		},
		{
			name: "success alert",
			alert: Alert{
				Type:    AlertSuccess,
				Message: "Operation successful",
			},
			contains: []string{
				`alert-success`,
				"Operation successful",
			},
		},
		{
			name: "warning alert",
			alert: Alert{
				Type:    AlertWarning,
				Message: "Warning message",
			},
			contains: []string{
				`alert-warning`,
				"Warning message",
			},
		},
		{
			name: "error alert",
			alert: Alert{
				Type:    AlertError,
				Message: "Error occurred",
			},
			contains: []string{
				`alert-error`,
				"Error occurred",
			},
		},
		{
			name: "alert with title",
			alert: Alert{
				Type:    AlertInfo,
				Title:   "Notice",
				Message: "Please read this",
			},
			contains: []string{
				"<strong>Notice</strong>",
				"Please read this",
			},
		},
		{
			name: "dismissible alert",
			alert: Alert{
				Type:        AlertInfo,
				Message:     "Dismissible message",
				Dismissible: true,
			},
			contains: []string{
				`alert-dismissible`,
				`class="alert-close"`,
				`type="button"`,
				`aria-label="Close"`,
				"×",
			},
		},
		{
			name: "default type is info",
			alert: Alert{
				Message: "Default alert",
			},
			contains: []string{
				`alert-info`,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := element.NewBuilder()
			tt.alert.Render(b)
			got := b.String()

			for _, want := range tt.contains {
				if !strings.Contains(got, want) {
					t.Errorf("Alert.Render() missing %q\ngot: %s", want, got)
				}
			}
		})
	}
}

func TestAlert_Render_NoTitle(t *testing.T) {
	b := element.NewBuilder()
	alert := Alert{
		Type:    AlertInfo,
		Message: "No title",
	}
	alert.Render(b)
	got := b.String()

	if strings.Contains(got, "<strong>") {
		t.Errorf("Alert without title should not render strong element, got: %s", got)
	}
}

func TestAlert_Render_NotDismissible(t *testing.T) {
	b := element.NewBuilder()
	alert := Alert{
		Type:    AlertInfo,
		Message: "Not dismissible",
	}
	alert.Render(b)
	got := b.String()

	if strings.Contains(got, "alert-close") {
		t.Errorf("Non-dismissible alert should not have close button, got: %s", got)
	}
	if strings.Contains(got, "alert-dismissible") {
		t.Errorf("Non-dismissible alert should not have dismissible class, got: %s", got)
	}
}

func TestAlertType_Constants(t *testing.T) {
	// Verify alert type constants have expected values
	if AlertInfo != "info" {
		t.Errorf("AlertInfo = %q, want %q", AlertInfo, "info")
	}
	if AlertSuccess != "success" {
		t.Errorf("AlertSuccess = %q, want %q", AlertSuccess, "success")
	}
	if AlertWarning != "warning" {
		t.Errorf("AlertWarning = %q, want %q", AlertWarning, "warning")
	}
	if AlertError != "error" {
		t.Errorf("AlertError = %q, want %q", AlertError, "error")
	}
}
