package main

import (
	"github.com/rohanthewiz/element"
	"github.com/rohanthewiz/element/components/form_field"
)

// EMPTY STRUCT for FORM COMPONENT
// ContactForm is a stateless component - no data fields needed
// All form structure and attributes are defined in the Render method
type ContactForm struct {
	Recipient string
}

// METHOD with POINTER PARAMETER and NAMED RETURN
// (cf ContactForm) - value receiver for the empty struct
// (b *element.Builder) - POINTER parameter to avoid copying the builder
// (dontCare any) - named return with 'any' type (we return nil via naked return)
func (cf ContactForm) Render(b *element.Builder) (x any) {
	// FORM ELEMENT: Build an HTML form with attributes
	// ATTRIBUTE PAIRS: "name", "value", "name", "value" pattern
	// action="/contact" - where to send form data (POST request to /contact endpoint)
	// method="POST" - HTTP method for form submission (POST for data modification)
	b.Form("action", "/contact", "method", "POST").R(
		// INPUT ELEMENT: Text input field
		// MULTIPLE ATTRIBUTES demonstrated:
		//   type="text" - standard text input (single line)
		//   name="name" - field name used when submitting form data
		//   placeholder="Name" - hint text shown when field is empty
		b.Input("type", "text", "name", "sender", "placeholder", "Your name").R(),

		form_field.FormField{
			Label:       "Email Address",
			Name:        "email",
			Type:        "email",
			Placeholder: "you@example.com",
			Required:    true,
			HelpText:    "We'll never share your email.",
		}.Render(b),

		/*		// INPUT ELEMENT: Email input field
				// type="email" - HTML5 input type that validates email format
				// Browser will enforce basic email validation before submission
				b.Input("type", "email", "name", "sender_email", "placeholder", "Your email").R(),
		*/
		b.Input("type", "hidden", "name", "recipient", "value", cf.Recipient).R(),

		// TEXTAREA ELEMENT: Multi-line text input
		// name="message" - field identifier for form submission
		// .R() with no arguments creates an empty textarea (no child elements)
		// TextArea is different from Input - it's a paired tag (<textarea></textarea>)
		b.TextArea("name", "message", "placeholder", "Message").R(),

		// BUTTON ELEMENT: Submit button
		// type="submit" - clicking this button submits the form
		// .T("Send") adds text content to the button
		// When clicked, browser sends POST request to /contact with form data
		b.Button("type", "submit").T("Send"),
		b.P("style", "color:red").T(" *Required fields"),
	)

	// NAKED RETURN: Returns the named value 'dontCare' (which is nil by default)
	// We don't need to return anything meaningful, so we use a naked return
	return
}
