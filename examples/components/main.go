// Package main demonstrates the reusable components in
// github.com/rohanthewiz/element/components.
// Run with: go run .
// Then visit: http://localhost:8080
package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/rohanthewiz/element"
	"github.com/rohanthewiz/element/components"
)

func main() {
	http.HandleFunc("/", homeHandler)

	fmt.Println("Components Example Server")
	fmt.Println("=========================")
	fmt.Println("Visit: http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// cellLink shows components rendering inside table cells
type cellLink struct {
	Href, Label string
}

func (c cellLink) Render(b *element.Builder) (x any) {
	b.A("href", c.Href).T(c.Label)
	return
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	b := element.AcquireBuilder()
	defer element.ReleaseBuilder(b)

	navItems := []components.NavItem{
		{Label: "Home", Href: "/", Active: true},
		{Label: "Products", Href: "/products"},
		{Label: "About", Href: "/about"},
	}

	breadcrumbs := []components.BreadcrumbItem{
		{Label: "Home", Href: "/"},
		{Label: "Components", Href: "/components"},
		{Label: "Demo"}, // Current page
	}

	definitions := []components.Definition{
		{Term: "Element", Definition: "A Go library for programmatic HTML generation without templates."},
		{Term: "Builder", Definition: "The primary API entry point that accumulates HTML in an internal buffer."},
		{Term: "Component", Definition: "A reusable HTML fragment implementing the Render(b *Builder) any interface."},
	}

	// Build the page
	b.Html().R(
		b.Head().R(
			b.Meta("charset", "utf-8").R(),
			b.Meta("name", "viewport", "content", "width=device-width, initial-scale=1").R(),
			b.Title().T("Element Components Demo"),
			b.Style().T(getStyles()),
		),
		b.Body().R(
			// Navigation with a dropdown menu
			components.Nav{Items: navItems, Brand: "Element", Class: "main-nav"}.Render(b),

			b.DivClass("container").R(
				// Breadcrumb
				section(b, "Breadcrumb",
					components.Breadcrumb{Items: breadcrumbs},
				),

				// Alerts
				section(b, "Alerts",
					components.Alert{Type: components.AlertSuccess, Title: "Success!", Message: "Your changes have been saved.", Dismissible: true},
					components.Alert{Type: components.AlertInfo, Message: "This is an informational message."},
					components.Alert{Type: components.AlertWarning, Title: "Warning:", Message: "Your session will expire in 5 minutes."},
					components.Alert{Type: components.AlertError, Title: "Error!", Message: "Failed to process your request.", Dismissible: true},
				),

				// Badges and buttons
				section(b, "Badges & Buttons",
					components.Badge{Text: "neutral"},
					components.Badge{Text: "3 new", Type: components.AlertInfo, Pill: true},
					components.Badge{Text: "active", Type: components.AlertSuccess},
					components.Badge{Text: "degraded", Type: components.AlertWarning},
					components.Badge{Text: "down", Type: components.AlertError, Pill: true},
					element.CompFunc(func(b *element.Builder) (x any) {
						b.DivClass("btn-row").R(
							components.Button{Text: "Primary"}.Render(b),
							components.Button{Text: "Secondary", Variant: "secondary"}.Render(b),
							components.Button{Text: "Danger", Variant: "danger"}.Render(b),
							components.Button{Text: "Small", Size: "sm"}.Render(b),
							components.Button{Text: "Disabled", Disabled: true}.Render(b),
							components.Button{Text: "Link Button", Href: "/docs", Variant: "secondary"}.Render(b),
						)
						return
					}),
				),

				// Table with columns, footer, and component cells
				section(b, "Table",
					components.Table{
						Caption: "Q3 user activity (component cells, alignment, footer)",
						Columns: []components.Column{
							{Header: "ID", Align: "right"},
							{Header: "Name"},
							{Header: "Visits", Align: "right"},
							{Header: "Status", Align: "center"},
							{Header: "Action", Align: "center"},
						},
						Rows: [][]any{
							{1, "Alice Smith", 1042, components.Badge{Text: "Active", Type: components.AlertSuccess}, cellLink{"/users/1", "Edit"}},
							{2, "Bob Johnson", 233, components.Badge{Text: "Pending", Type: components.AlertWarning}, cellLink{"/users/2", "Edit"}},
							{3, "Carol White", 3970, components.Badge{Text: "Active", Type: components.AlertSuccess}, cellLink{"/users/3", "Edit"}},
						},
						FooterRows: [][]any{
							{"", "Total", 5245, "", ""},
						},
						Striped:    true,
						Bordered:   true,
						Hover:      true,
						Responsive: true,
						RowClass: func(i int, row []any) string {
							if v, ok := row[2].(int); ok && v > 3000 {
								return "row-highlight"
							}
							return ""
						},
					},
					components.Table{
						Headers:      []string{"A", "B", "C"},
						EmptyMessage: "No records found — the EmptyMessage field in action.",
						Bordered:     true,
					},
				),

				// Tabs — CSS-only, no JavaScript
				section(b, "Tabs (no JavaScript)",
					components.Tabs{
						ID: "demo-tabs",
						Items: []components.Tab{
							{Label: "Overview", Content: "Tabs are driven by hidden radio inputs and a few generated CSS rules — no JS."},
							{Label: "Details", Body: components.DefinitionList{Items: definitions}},
							{Label: "Data", Body: components.Table{
								Headers: []string{"Metric", "Value"},
								Rows:    [][]any{{"Users", 1234}, {"Revenue", "$5,678"}},
							}},
						},
					},
				),

				// Accordion — native details/summary
				section(b, "Accordion (native <details>, exclusive open)",
					components.Accordion{
						Exclusive: true,
						Name:      "faq",
						Items: []components.AccordionItem{
							{Title: "What is Element?", Content: "A Go library for generating HTML without templates.", Open: true},
							{Title: "Does the accordion need JavaScript?", Content: "No — it uses native <details>/<summary>. The Exclusive option uses the details name attribute so the browser keeps one section open."},
							{Title: "Can sections hold components?", Body: components.Badge{Text: "Yes — any element.Component", Type: components.AlertSuccess}},
						},
					},
				),

				// Modal — native Popover API
				section(b, "Modal (native Popover API)",
					components.Modal{
						ID:           "demo-modal",
						Title:        "Hello from a no-JS modal",
						Content:      "The browser handles open, close, Esc, and click-outside via the Popover API.",
						TriggerText:  "Open Modal",
						TriggerClass: "btn btn-primary",
					},
				),

				// Dropdown
				section(b, "Dropdown (native <details>)",
					components.Dropdown{
						Label: "Account ▾",
						Items: []components.DropdownItem{
							{Label: "Profile", Href: "/profile"},
							{Label: "Settings", Href: "/settings"},
							{Divider: true},
							{Label: "Sign out", Href: "/logout"},
						},
					},
				),

				// Progress
				section(b, "Progress",
					components.ProgressBar{Value: 35, Label: "Storage used", ShowValue: true},
					components.ProgressBar{Value: 80, Label: "Upload", ShowValue: true, Class: "progress-lg"},
				),

				// Form controls
				section(b, "Form Controls",
					element.CompFunc(func(b *element.Builder) (x any) {
						b.Form("action", "#", "method", "post").R(
							components.FormField{
								Label: "Email Address", Name: "email", Type: "email",
								Placeholder: "you@example.com", Required: true,
								HelpText: "We'll never share your email.",
							}.Render(b),
							components.FormField{
								Label: "Username", Name: "username",
								Value: "invalid user!", Error: "Username can only contain letters and numbers.",
							}.Render(b),
							components.SelectField{
								Label: "Country", Name: "country", Prompt: "Choose a country…", Required: true,
								Options: []components.SelectOption{
									{Value: "us", Label: "United States"},
									{Value: "ca", Label: "Canada"},
									{Value: "mx", Label: "Mexico"},
								},
							}.Render(b),
							components.TextAreaField{
								Label: "Bio", Name: "bio", Rows: 4,
								Placeholder: "Tell us about yourself",
								HelpText:    "Markdown is supported.",
							}.Render(b),
							components.RadioGroup{
								Legend: "Plan", Name: "plan", Selected: "pro", Inline: true,
								Options: []components.SelectOption{
									{Value: "free", Label: "Free"},
									{Value: "pro", Label: "Pro"},
									{Value: "team", Label: "Team"},
								},
							}.Render(b),
							components.CheckboxField{
								Label: "Subscribe to the newsletter", Name: "subscribe", Checked: true,
							}.Render(b),
							b.DivClass("form-actions").R(
								components.Button{Text: "Submit", Type: "submit"}.Render(b),
							),
						)
						return
					}),
				),

				// Pagination
				section(b, "Pagination",
					components.Pagination{
						CurrentPage: 3, TotalPages: 10,
						BaseURL:   "/page/%d",
						ShowFirst: true, ShowLast: true,
					},
				),
			),

			// Footer
			b.FooterClass("footer").R(
				b.P().T("Element Components Demo — Built with github.com/rohanthewiz/element"),
			),
		),
	)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(b.Bytes())
}

// section renders a titled demo section containing the given components
func section(b *element.Builder, title string, comps ...element.Component) (x any) {
	b.DivClass("section").R(
		b.H2().T(title),
		element.RenderComponents(b, comps...),
	)
	return
}
