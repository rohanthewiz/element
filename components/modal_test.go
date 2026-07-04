package components

import (
	"strings"
	"testing"

	"github.com/rohanthewiz/element"
)

func TestModal_Render(t *testing.T) {
	b := element.NewBuilder()
	Modal{
		ID:          "confirm",
		Title:       "Confirm Delete",
		Content:     "Are you sure?",
		TriggerText: "Delete",
	}.Render(b)
	got := b.String()

	for _, want := range []string{
		// Trigger
		`popovertarget="confirm"`,
		">Delete</button>",
		// Popover
		`id="confirm"`, `popover="auto"`, `role="dialog"`,
		`aria-labelledby="confirm-title"`,
		`<h3 class="modal-title" id="confirm-title">Confirm Delete</h3>`,
		`popovertargetaction="hide"`,
		`aria-label="Close"`,
		"Are you sure?",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("Modal.Render() missing %q\ngot: %s", want, got)
		}
	}
}

func TestModal_NoTrigger(t *testing.T) {
	b := element.NewBuilder()
	Modal{ID: "m1", Content: "body"}.Render(b)
	got := b.String()

	if strings.Contains(got, "modal-trigger") {
		t.Errorf("Modal without TriggerText should not render a trigger, got: %s", got)
	}
	// The close button still targets the popover
	if !strings.Contains(got, `popovertarget="m1"`) {
		t.Errorf("Modal close button should target the popover, got: %s", got)
	}
}
