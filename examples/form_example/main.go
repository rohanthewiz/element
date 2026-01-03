package main

import (
	"log"

	"github.com/rohanthewiz/element"
	"github.com/rohanthewiz/rweb"
	"github.com/rohanthewiz/serr"
)

func main() {
	// Initialize the rweb server on port 8000 with verbose logging
	// This creates a high-performance web server with built-in middleware support
	s := rweb.NewServer(rweb.ServerOptions{
		Address: ":8000",
		Verbose: true,
	})

	// Add middleware to log request information (method, path, duration, etc.)
	s.Use(rweb.RequestInfo)
	// s.ElementDebugRoutes()

	// Define route handlers
	// The root handler demonstrates Element's core features including components and iteration
	s.Get("/", rootHandler)
	// Start the server and listen for requests
	// log.Fatal ensures we see any startup errors
	log.Fatal(s.Run())
}

// rootHandler serves the main page demonstrating various Element features:
// - CSS styling, components, iteration, and nested element structures
func rootHandler(c rweb.Context) error {
	// Generate the complete HTML page and send it as the response
	err := c.WriteHTML(HomePage{}.Render())
	if err != nil {
		// Wrap errors with context for better debugging
		return serr.Wrap(err)
	}
	return nil
}

type HomePage struct{}

func (hp HomePage) Render() string {
	b := element.NewBuilder()

	b.Html().R(
		b.Head().R(
			b.Style().T(`
/* Form */
.form-field {
    margin-bottom: 1.25rem;
}

.form-label {
    display: block;
    margin-bottom: 0.5rem;
    font-weight: 500;
}

.required { color: #e74c3c; }

.form-input {
    width: 100%;
    padding: 0.625rem;
    border: 1px solid #ddd;
    border-radius: 4px;
    font-size: 1rem;
}

.form-input:focus {
    outline: none;
    border-color: #3498db;
    box-shadow: 0 0 0 3px rgba(52,152,219,0.1);
}

.has-error .form-input {
    border-color: #e74c3c;
}

.form-error {
    display: block;
    color: #e74c3c;
    font-size: 0.875rem;
    margin-top: 0.25rem;
}

.form-help {
    display: block;
    color: #666;
    font-size: 0.875rem;
    margin-top: 0.25rem;
}

.form-actions {
    margin-top: 1.5rem;
}

.btn {
    padding: 0.625rem 1.25rem;
    border: none;
    border-radius: 4px;
    cursor: pointer;
    font-size: 1rem;
}

.btn-primary {
    background: #3498db;
    color: white;
}

.btn-primary:hover {
    background: #2980b9;
}
`),
		),
		b.Body().R(
			b.Div("class", "main-content", "style", "text-align:center").R(
				b.H1("style", "color:blue").T("Contact Form Example"),
				b.P("style", "color:green").T("Fill out the form below to send a message."),
				ContactForm{Recipient: "support@example.com"}.Render(b),
			),
		),
	)

	return b.String()
}
