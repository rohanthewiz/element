package components

import "github.com/rohanthewiz/element"

// -----------------------------------------------------------------------------
// Modal Component
// -----------------------------------------------------------------------------

// Modal renders a popup dialog using the native Popover API — no JavaScript
// required. The browser handles open/close, Esc-to-dismiss, and
// click-outside-to-dismiss (light dismiss).
//
// Style hooks: .modal (the popover), .modal-header, .modal-title,
// .modal-close, .modal-body, and .modal::backdrop for the overlay.
//
// To open the modal from elsewhere on the page, render any button with
// popovertarget set to this modal's ID.
type Modal struct {
	ID           string            // Unique id for this instance on the page (default "modal")
	Title        string            // Optional header title
	Content      string            // Simple text content
	Body         element.Component // Alternative: render a component as the modal body
	TriggerText  string            // Text for the built-in trigger button; empty renders no trigger
	TriggerClass string            // Class for the trigger button (default "btn modal-trigger")
	Class        string            // Additional CSS classes for the modal
}

// Render implements the element.Component interface.
func (m Modal) Render(b *element.Builder) (x any) {
	id := m.ID
	if id == "" {
		id = "modal"
	}
	modalClass := "modal"
	if m.Class != "" {
		modalClass += " " + m.Class
	}

	// Optional trigger button
	if m.TriggerText != "" {
		triggerClass := m.TriggerClass
		if triggerClass == "" {
			triggerClass = "btn modal-trigger"
		}
		b.ButtonClass(triggerClass, "type", "button", "popovertarget", id).T(m.TriggerText)
	}

	// The popover itself
	b.DivClass(modalClass, "id", id, "popover", "auto", "role", "dialog", "aria-labelledby", id+"-title").R(
		b.DivClass("modal-header").R(
			b.Wrap(func() {
				if m.Title != "" {
					b.H3Class("modal-title", "id", id+"-title").T(m.Title)
				}
			}),
			b.ButtonClass("modal-close",
				"type", "button",
				"popovertarget", id,
				"popovertargetaction", "hide",
				"aria-label", "Close").T("×"),
		),
		b.DivClass("modal-body").R(
			b.Wrap(func() {
				if m.Body != nil {
					m.Body.Render(b)
				} else {
					b.T(m.Content)
				}
			}),
		),
	)
	return
}
