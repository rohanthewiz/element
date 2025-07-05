// Package main demonstrates the Element library for programmatic HTML generation.
// This example shows how to use Element with the rweb web server to create
// dynamic HTML pages without templates, using Go's natural function execution order.
package main

import (
	"fmt"
	"log"
	"strconv"

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
	s.ElementDebugRoutes()

	// Define route handlers
	// The root handler demonstrates Element's core features including components and iteration
	s.Get("/", rootHandler)

	// The other-page handler shows an alternative approach using HtmlPage helper
	s.Get("/other-page", otherPageHandler)

	// Demonstrate standalone HTML generation (without HTTP context)
	// This shows Element can be used independently of web servers
	fmt.Println(TestRender())

	// Start the server and listen for requests
	// log.Fatal ensures we see any startup errors
	log.Fatal(s.Run())
}

// rootHandler serves the main page demonstrating various Element features:
// - CSS styling, components, iteration, and nested element structures
func rootHandler(c rweb.Context) error {
	// Sample data for demonstrating iteration and component rendering
	animals := []string{"cat", "mouse", "dog"}
	colors := []string{"red", "blue", "green", "indigo", "violet"}

	// Generate the complete HTML page and send it as the response
	err := c.WriteHTML(generateHTML(animals, colors))
	if err != nil {
		// Wrap errors with context for better debugging
		return serr.Wrap(err)
	}
	return nil
}

// SelectComponent demonstrates Element's component system for reusable UI elements.
// Components implement the element.Component interface by having a Render method
// that takes an *element.Builder and returns any. This enables composition
// of complex UI elements that can be reused across different pages.
type SelectComponent struct {
	Selected string   // The currently selected value
	Items    []string // All available options for the select dropdown
}

// Render implements the element.Component interface.
// It generates a <select> element with <option> children based on the component's data.
// The builder pattern ensures proper HTML structure with automatic tag closing.
func (s SelectComponent) Render(b *element.Builder) (x any) {
	// Create the select element and use R() to render its children
	b.Select().R(
		// Anonymous function allows us to generate multiple option elements
		// The function executes immediately (note the () at the end)
		func() (x any) {
			for _, color := range s.Items {
				// Build attribute pairs dynamically
				params := []string{"value", color}

				// Add 'selected' attribute for the currently selected item
				if color == s.Selected {
					params = append(params, "selected", "selected")
				}

				// Create option element with attributes and text content
				b.Option(params...).T(color)
			}
			return
		}(),
	)
	return
}

// generateHTML creates a complete HTML page demonstrating Element's features:
// - Nested element structure matching HTML's tree structure
// - CSS styling within the page
// - Dynamic content generation with loops
// - Component integration
// - Various element termination methods (R() vs T())
func generateHTML(animals []string, colors []string) string {
	// Create a new builder - this is the starting point for all HTML generation
	// The builder maintains an internal string builder for efficient concatenation
	b := element.NewBuilder()

	// Create a reusable select component with "blue" pre-selected
	selComp := SelectComponent{Selected: "blue", Items: colors}

	// Start building the HTML structure
	// Notice how the nesting matches HTML's structure naturally
	b.Html().R(
		// Head section contains metadata and styles
		b.Head().R(
			// Embedded CSS - T() is used since style content is pure text
			b.Style().T(`
                #page-container {
                    padding: 4rem; height: 100vh; background-color: rgb(232, 230, 228);
                }
                .intro {
                    font-style: italic; font-size: 0.9rem; padding-left: 3em;
                }
                .highlight {
                    background-color: yellow;
                }
                .footer {
                    text-align: center; font-size: 0.8rem; border-top: 1px solid #ccc; padding: 1em;
                }`,
			),
		),
		// Body contains the visible page content
		b.Body().R(
			// Main container div with id attribute
			// Attributes are passed as pairs: "name", "value"
			b.Div("id", "page-container").R(
				// Simple heading with text content
				b.H1().T("This is my heading"),

				// DivClass is a convenience method for elements with class attributes
				b.DivClass("intro").R(), // Empty div - R() with no children

				// Intentionally malformed attributes to test debug mode
				// "unpaired" has no value - debug mode will catch this
				b.Div("class", "intro", "unpaired").R(
					// Paragraph with mixed content
					b.P().R(
						// b.T() adds text nodes directly to the builder
						b.T("I've got plenty to say here "),
						// Inline span with highlighting class
						b.SpanClass("highlight").R(
							// Multiple strings can be passed to T()
							b.T("important phrase!", " More intro text"),
						),
					),
				),
				// Another paragraph demonstrating various features
				b.P().R(
					b.T("ABC Company"),
					// Self-closing tags like <br> don't require R() or T(), but if unsure if an element is self-closing use R() anyway and it will do the right thing.
					b.Br(),
					// Wrap allows arbitrary Go code to generate content
					b.Wrap(func() {
						// Generate comma-separated numbers programmatically
						out := ""
						for i := 0; i < 10; i++ {
							if i > 0 {
								out += ","
							}
							out += strconv.Itoa(i)
						}
						// Add the generated string as text content
						b.T(out)
					}),
				),
				// Div with multiple text nodes
				// Note: HTML in text content is not escaped
				b.Div().R(
					b.T("Lorem Ipsum Lorem Ipsum Lorem<br>Ipsum Lorem Ipsum "),
					b.T("Finally..."),
				),
				// Demonstrate iteration over a slice to generate list items
				// ForEach is a helper function, but inline anonymous functions
				// or components often provide more flexibility
				b.UlClass("list").R(
					// ForEach executes the callback for each item in the slice
					element.ForEach(animals, func(animal string) {
						// Each iteration creates a new list item
						b.Li().T("This is a ", animal, " in a list item")
					}),
				),

				// Render the reusable select component
				// RenderComponents calls the component's Render method
				element.RenderComponents(b, selComp),

				// Empty paragraph used as a spacer
				b.P().R(),

				// Footer div with simple text content
				b.DivClass("footer").T("About | Privacy | Logout"),
			),
		),
	)

	// Convert the builder's internal buffer to a string
	// This is the final HTML output
	out := b.String()

	// Uncomment these lines to debug the generated HTML
	// fmt.Println(strings.Repeat("-", 60))
	// fmt.Println(out)
	// fmt.Println(strings.Repeat("=", 60))

	return out
}

// TestRender demonstrates basic Element usage with a simple HTML fragment.
// This shows the fundamental builder pattern and how elements nest naturally.
func TestRender() string {
	// Always start with NewBuilder()
	b := element.NewBuilder()

	// Element creation follows this pattern:
	// 1. Call the element method (e.g., Body()) - this writes the opening tag
	// 2. Add children as arguments to R() or text with T()
	// 3. R() or T() writes the closing tag
	b.Body().R(
		// DivClass is a shorthand for Div("class", "container")
		// The R() method renders children and closes the tag
		b.DivClass("container").R(
			// Span with text-only content uses T()
			b.Span().T("Some text"),

			// Paragraph contains an anchor element
			b.P().R(
				// Anchor with href attribute and text content
				b.A("href", "http://example.com").T("Example.com"),
			), // P closing tag written here
		), // Div closing tag written here
	) // Body closing tag written here

	// Get the complete HTML string from the builder
	return b.String()
	// Output: <body><div class="container"><span>Some text</span><p><a href="http://example.com">Example.com</a></p></div></body>
}

// ----- OTHER PAGE -----

// otherPageHandler demonstrates an alternative page construction technique
// using HtmlPage helper method and component-based body rendering
func otherPageHandler(c rweb.Context) error {
	return c.WriteHTML(otherHTMLPage())
}

// otherHTMLPage uses the HtmlPage helper method which:
// 1. Creates a complete HTML document structure
// 2. Adds the provided CSS to a <style> tag in <head>
// 3. Adds any additional head content (like <title>)
// 4. Renders the body component
func otherHTMLPage() (out string) {
	b := element.NewBuilder()

	// HtmlPage is a convenience method that generates a complete HTML page
	// Parameters: CSS styles, additional head content, body component
	b.HtmlPage("body {background-color:#eee;}", "<title>My Other Page</title>", otherBody{})

	return b.String()
}

// otherBody implements the element.Component interface
// This demonstrates how to create page-specific body content as a component
type otherBody struct{}

// Render generates the body content for the other page.
// Notice that elements are rendered in sequence - there's no need to
// wrap everything in a container unless you want that in your HTML.
func (ob otherBody) Render(b *element.Builder) (x any) {
	// Add a heading
	b.H1().T("This is my other page")

	// Add a paragraph with descriptive text
	b.P().R(
		b.T("This is a simple example of using the Element library to generate HTML."),
	)

	// Input elements are self-closing in HTML, but Element allows R() for consistency
	// The span here demonstrates that it won't be rendered inside the input
	// This will show as an issue in Debug mode
	b.Input("type", "text").R(b.Span().T("I shouldn't be here"))

	// Footer section
	b.DivClass("footer").T("About | Privacy | Logout")

	return
}
