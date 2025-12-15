// Package main demonstrates the use of reusable Element components.
// Run with: go run .
// Then visit: http://localhost:8080
package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/rohanthewiz/element"
)

func main() {
	http.HandleFunc("/", homeHandler)

	fmt.Println("Components Example Server")
	fmt.Println("=========================")
	fmt.Println("Visit: http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	b := element.AcquireBuilder()
	defer element.ReleaseBuilder(b)

	// Sample data
	tableHeaders := []string{"ID", "Name", "Email", "Status"}
	tableData := [][]any{
		{1, "Alice Smith", "alice@example.com", "Active"},
		{2, "Bob Johnson", "bob@example.com", "Pending"},
		{3, "Carol White", "carol@example.com", "Active"},
		{4, "David Brown", "david@example.com", "Inactive"},
	}

	navItems := []NavItem{
		{Label: "Home", Href: "/", Active: true},
		{Label: "Products", Href: "/products"},
		{Label: "About", Href: "/about"},
		{Label: "Contact", Href: "/contact"},
	}

	breadcrumbs := []BreadcrumbItem{
		{Label: "Home", Href: "/"},
		{Label: "Products", Href: "/products"},
		{Label: "Electronics", Href: "/products/electronics"},
		{Label: "Laptops", Href: ""}, // Current page
	}

	definitions := []Definition{
		{Term: "Element", Definition: "A Go library for programmatic HTML generation without templates."},
		{Term: "Builder", Definition: "The primary API entry point that accumulates HTML in an internal buffer."},
		{Term: "Component", Definition: "A reusable HTML fragment implementing the Render(b *Builder) any interface."},
	}

	// Build the page
	b.Html().R(
		b.Head().R(
			b.Meta("charset", "utf-8"),
			b.Meta("name", "viewport", "content", "width=device-width, initial-scale=1"),
			b.Title().T("Element Components Demo"),
			b.Style().T(getStyles()),
		),
		b.Body().R(
			// Navigation
			Nav{Items: navItems, Brand: "Element", Class: "main-nav"}.Render(b),

			b.DivClass("container").R(
				// Breadcrumb
				b.DivClass("section").R(
					b.H2().T("Breadcrumb Component"),
					Breadcrumb{Items: breadcrumbs}.Render(b),
				),

				// Alerts
				b.DivClass("section").R(
					b.H2().T("Alert Components"),
					Alert{Type: AlertSuccess, Title: "Success!", Message: "Your changes have been saved.", Dismissible: true}.Render(b),
					Alert{Type: AlertInfo, Message: "This is an informational message."}.Render(b),
					Alert{Type: AlertWarning, Title: "Warning:", Message: "Your session will expire in 5 minutes."}.Render(b),
					Alert{Type: AlertError, Title: "Error!", Message: "Failed to process your request.", Dismissible: true}.Render(b),
				),

				// Table
				b.DivClass("section").R(
					b.H2().T("Table Component"),
					b.P().T("A data table with headers and rows:"),
					Table{
						Headers:  tableHeaders,
						Rows:     tableData,
						Striped:  true,
						Bordered: true,
					}.Render(b),
				),

				// Cards
				b.DivClass("section").R(
					b.H2().T("Card Components"),
					b.DivClass("card-grid").R(
						Card{
							Title:  "Simple Card",
							Body:   "This is a basic card with just text content.",
							Footer: "Last updated: Today",
						}.Render(b),
						Card{
							Title: "Card with Table",
							BodyComponent: Table{
								Headers: []string{"Metric", "Value"},
								Rows: [][]any{
									{"Users", 1234},
									{"Revenue", "$5,678"},
									{"Growth", "12%"},
								},
							},
						}.Render(b),
						Card{
							Title: "Feature List",
							BodyComponent: featureList{},
						}.Render(b),
					),
				),

				// Definition List
				b.DivClass("section").R(
					b.H2().T("Definition List Component"),
					DefinitionList{Items: definitions}.Render(b),
				),

				// Pagination
				b.DivClass("section").R(
					b.H2().T("Pagination Component"),
					b.P().T("Page 3 of 10:"),
					Pagination{
						CurrentPage: 3,
						TotalPages:  10,
						BaseURL:     "/page/%d",
						ShowFirst:   true,
						ShowLast:    true,
					}.Render(b),
				),

				// Form Fields
				b.DivClass("section").R(
					b.H2().T("Form Field Components"),
					b.Form("action", "#", "method", "post").R(
						FormField{
							Label:       "Email Address",
							Name:        "email",
							Type:        "email",
							Placeholder: "you@example.com",
							Required:    true,
							HelpText:    "We'll never share your email.",
						}.Render(b),
						FormField{
							Label:       "Password",
							Name:        "password",
							Type:        "password",
							Required:    true,
						}.Render(b),
						FormField{
							Label: "Username",
							Name:  "username",
							Type:  "text",
							Value: "invalid user!",
							Error: "Username can only contain letters and numbers.",
						}.Render(b),
						b.DivClass("form-actions").R(
							b.ButtonClass("btn btn-primary", "type", "submit").T("Submit"),
						),
					),
				),

				// Usage examples
				b.DivClass("section").R(
					b.H2().T("Usage Example"),
					b.Pre().R(
						b.Code().T(`// Create a table component
table := Table{
    Headers: []string{"ID", "Name", "Status"},
    Rows: [][]any{
        {1, "Alice", "Active"},
        {2, "Bob", "Pending"},
    },
    Striped:  true,
    Bordered: true,
}

// Render it
b := element.AcquireBuilder()
defer element.ReleaseBuilder(b)

b.Div().R(
    table.Render(b),
)

// Or render multiple components
element.RenderComponents(b, table, card, alert)`),
					),
				),
			),

			// Footer
			b.FooterClass("footer").R(
				b.P().T("Element Components Demo - Built with github.com/rohanthewiz/element"),
			),
		),
	)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(b.Bytes())
}

// featureList is a simple component for demonstrating nested components
type featureList struct{}

func (f featureList) Render(b *element.Builder) (x any) {
	features := []string{
		"No external dependencies",
		"Type-safe HTML generation",
		"Single-pass rendering",
		"Builder pooling support",
	}

	b.Ul().R(
		element.ForEach(features, func(feature string) {
			b.Li().T(feature)
		}),
	)
	return
}

func getStyles() string {
	return `
* { box-sizing: border-box; margin: 0; padding: 0; }

body {
    font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
    line-height: 1.6;
    color: #333;
    background: #f5f5f5;
}

.container {
    max-width: 1200px;
    margin: 0 auto;
    padding: 2rem;
}

.section {
    background: white;
    border-radius: 8px;
    padding: 1.5rem;
    margin-bottom: 2rem;
    box-shadow: 0 2px 4px rgba(0,0,0,0.1);
}

.section h2 {
    margin-bottom: 1rem;
    padding-bottom: 0.5rem;
    border-bottom: 2px solid #eee;
}

/* Navigation */
.nav {
    background: #2c3e50;
    padding: 1rem 2rem;
    display: flex;
    align-items: center;
    gap: 2rem;
}

.nav-brand {
    color: white;
    font-size: 1.5rem;
    font-weight: bold;
    text-decoration: none;
}

.nav-list {
    list-style: none;
    display: flex;
    gap: 1rem;
}

.nav-link {
    color: rgba(255,255,255,0.8);
    text-decoration: none;
    padding: 0.5rem 1rem;
    border-radius: 4px;
}

.nav-link:hover, .nav-item.active .nav-link {
    color: white;
    background: rgba(255,255,255,0.1);
}

/* Breadcrumb */
.breadcrumb-list {
    list-style: none;
    display: flex;
    flex-wrap: wrap;
    padding: 0.75rem 0;
}

.breadcrumb-item {
    display: flex;
    align-items: center;
}

.breadcrumb-item a {
    color: #3498db;
    text-decoration: none;
}

.breadcrumb-item a:hover {
    text-decoration: underline;
}

.breadcrumb-current {
    color: #666;
}

.breadcrumb-separator {
    color: #999;
    margin: 0 0.25rem;
}

/* Alerts */
.alert {
    padding: 1rem;
    border-radius: 4px;
    margin-bottom: 1rem;
    display: flex;
    align-items: center;
    gap: 0.5rem;
}

.alert-info { background: #d1ecf1; color: #0c5460; border: 1px solid #bee5eb; }
.alert-success { background: #d4edda; color: #155724; border: 1px solid #c3e6cb; }
.alert-warning { background: #fff3cd; color: #856404; border: 1px solid #ffeeba; }
.alert-error { background: #f8d7da; color: #721c24; border: 1px solid #f5c6cb; }

.alert-close {
    margin-left: auto;
    background: none;
    border: none;
    font-size: 1.25rem;
    cursor: pointer;
    opacity: 0.5;
}

.alert-close:hover { opacity: 1; }

/* Table */
.table {
    width: 100%;
    border-collapse: collapse;
    margin: 1rem 0;
}

.table th, .table td {
    padding: 0.75rem;
    text-align: left;
}

.table th {
    background: #3498db;
    color: white;
    font-weight: 600;
}

.table-striped tbody tr:nth-child(even) {
    background: #f8f9fa;
}

.table-bordered th, .table-bordered td {
    border: 1px solid #dee2e6;
}

/* Cards */
.card-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
    gap: 1.5rem;
}

.card {
    background: white;
    border: 1px solid #ddd;
    border-radius: 8px;
    overflow: hidden;
}

.card-header {
    background: #f8f9fa;
    padding: 1rem;
    border-bottom: 1px solid #ddd;
}

.card-title {
    margin: 0;
    font-size: 1.1rem;
}

.card-body {
    padding: 1rem;
}

.card-footer {
    background: #f8f9fa;
    padding: 0.75rem 1rem;
    border-top: 1px solid #ddd;
    font-size: 0.875rem;
    color: #666;
}

/* Definition List */
.definition-list {
    display: grid;
    grid-template-columns: max-content 1fr;
    gap: 0.5rem 1.5rem;
}

.definition-term {
    font-weight: 600;
    color: #2c3e50;
}

.definition-desc {
    color: #555;
}

/* Pagination */
.pagination-list {
    list-style: none;
    display: flex;
    gap: 0.25rem;
}

.pagination-list li {
    display: flex;
}

.pagination-link, .pagination-current {
    padding: 0.5rem 0.75rem;
    border: 1px solid #ddd;
    border-radius: 4px;
    text-decoration: none;
    color: #3498db;
}

.pagination-link:hover {
    background: #f0f0f0;
}

.pagination-current {
    background: #3498db;
    color: white;
    border-color: #3498db;
}

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

/* Code */
pre {
    background: #2c3e50;
    color: #ecf0f1;
    padding: 1rem;
    border-radius: 4px;
    overflow-x: auto;
}

code {
    font-family: 'SF Mono', Consolas, monospace;
    font-size: 0.875rem;
}

/* Footer */
.footer {
    text-align: center;
    padding: 2rem;
    color: #666;
    background: white;
    margin-top: 2rem;
}
`
}
