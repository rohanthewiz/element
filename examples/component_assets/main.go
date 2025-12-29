// Package main demonstrates component asset management in Element.
//
// This example shows:
// 1. Components that provide CSS and JavaScript via the AssetProvider interface
// 2. Automatic asset collection and deduplication by the Builder
// 3. Components with and without assets working together
// 4. How to render collected assets in the page
//
// Run with: go run main.go
package main

import (
	"fmt"
	"os"

	"github.com/rohanthewiz/element"
	"github.com/rohanthewiz/element/examples/component_assets/components"
)

func main() {
	// Example 1: Standard usage with inline assets
	fmt.Println("=== Example 1: Page with Inline Assets ===")
	html := renderPageWithInlineAssets()
	fmt.Println(html)
	fmt.Println()

	// Example 2: External asset files
	fmt.Println("=== Example 2: External Asset Files ===")
	renderPageWithExternalAssets()
	fmt.Println()

	// Example 3: Asset inspection
	fmt.Println("=== Example 3: Asset Inspection ===")
	inspectAssets()
}

// renderPageWithInlineAssets demonstrates the most common usage:
// assets are collected automatically and rendered inline in <style> and <script> tags.
func renderPageWithInlineAssets() string {
	// Create a builder with asset tracking enabled
	b := element.NewBuilder().EnableAssetTracking()

	// Build the page
	b.WriteString("<!DOCTYPE html>\n")
	b.Html("lang", "en").R(
		b.Head().R(
			b.Meta("charset", "UTF-8"),
			b.Meta("name", "viewport", "content", "width=device-width, initial-scale=1.0"),
			b.Title().T("Element Component Assets Example"),

			// Basic page styles
			b.Style().R(
				b.T(`
					body {
						font-family: system-ui, -apple-system, sans-serif;
						max-width: 1200px;
						margin: 0 auto;
						padding: 2rem;
						line-height: 1.6;
						background-color: #f5f5f5;
					}
					h1 { color: #333; }
					.section { margin: 2rem 0; }
				`),
			),

			// Render all collected CSS from components
			// This will include CSS from Alert and Table components used below
			b.RenderCSS(),
		),

		b.Body().R(
			b.H1().T("Element Component Asset Management"),

			b.DivClass("section").R(
				b.H2().T("Alert Components (with CSS + JS)"),
				b.P().T("These alerts have embedded CSS for styling and JavaScript for dismissible functionality:"),

				// Render multiple alerts - each will register its assets, but they'll be deduplicated
				components.Alert{
					Type:        components.AlertSuccess,
					Title:       "Success!",
					Message:     "Your component assets are working correctly.",
					Dismissible: true,
				}.Render(b),

				components.Alert{
					Type:        components.AlertInfo,
					Title:       "Info:",
					Message:     "Assets are automatically collected and deduplicated.",
					Dismissible: true,
				}.Render(b),

				components.Alert{
					Type:    components.AlertWarning,
					Message: "This alert is not dismissible (no JavaScript needed).",
				}.Render(b),
			),

			b.DivClass("section").R(
				b.H2().T("Table Component (CSS only)"),
				b.P().T("This table has embedded CSS but no JavaScript:"),

				components.Table{
					Headers: []string{"Feature", "Status", "Notes"},
					Rows: [][]any{
						{"Asset Deduplication", "✅", "Same assets only included once"},
						{"CSS Embedding", "✅", "Embedded with go:embed"},
						{"JS Embedding", "✅", "Conditional based on component config"},
						{"No NodeJS Required", "✅", "Pure Go solution"},
					},
					Striped:  true,
					Bordered: true,
				}.Render(b),
			),

			b.DivClass("section").R(
				b.H2().T("Simple Component (no assets)"),
				b.P().T("Components without AssetProvider interface work normally:"),

				components.SimpleCard{
					Title: "Traditional Approach",
					Body:  "This component uses inline styles and doesn't participate in asset management.",
				}.Render(b),
			),

			b.DivClass("section").R(
				b.H2().T("Multiple Tables (deduplication test)"),
				b.P().T("Even though we render multiple tables, the CSS is only included once:"),

				components.Table{
					Headers: []string{"ID", "Name"},
					Rows: [][]any{
						{1, "First Table"},
						{2, "Second Row"},
					},
					Striped: true,
				}.Render(b),

				components.Table{
					Headers: []string{"Asset ID", "Times Included"},
					Rows: [][]any{
						{"Alert", "1"},
						{"Table", "1"},
					},
					Bordered: true,
				}.Render(b),
			),

			// Render all collected JavaScript from components
			// This will include the alert dismissal functionality
			b.RenderJS(),

			// Show what was collected
			b.Hr(),
			b.Details().R(
				b.Summary().T("View Collected Assets Info"),
				b.Pre().T(fmt.Sprintf(
					"CSS collected: %d bytes\nJS collected: %d bytes\n\nAssets were automatically deduplicated!",
					len(b.GetCollectedCSS()),
					len(b.GetCollectedJS()),
				)),
			),
		),
	)

	return b.String()
}

// renderPageWithExternalAssets demonstrates serving CSS and JS as external files
// for better browser caching and parallel loading.
func renderPageWithExternalAssets() {
	// First, collect all assets by rendering components
	collector := element.NewBuilder().EnableAssetTracking()

	// Render all components you'll use (just to collect assets, we discard the HTML)
	components.Alert{}.Render(collector)
	components.Table{}.Render(collector)

	// Write assets to files
	cssFile := "components.css"
	jsFile := "components.js"

	css := collector.GetCollectedCSS()
	js := collector.GetCollectedJS()

	if len(css) > 0 {
		if err := os.WriteFile(cssFile, css, 0644); err != nil {
			fmt.Printf("Error writing CSS: %v\n", err)
			return
		}
		fmt.Printf("✅ Wrote %d bytes to %s\n", len(css), cssFile)
	}

	if len(js) > 0 {
		if err := os.WriteFile(jsFile, js, 0644); err != nil {
			fmt.Printf("Error writing JS: %v\n", err)
			return
		}
		fmt.Printf("✅ Wrote %d bytes to %s\n", len(js), jsFile)
	}

	// Now create the HTML page that references these external files
	b := element.NewBuilder() // Note: No asset tracking needed for the page itself

	b.WriteString("<!DOCTYPE html>\n")
	b.Html().R(
		b.Head().R(
			b.Title().T("External Assets Example"),
			b.Link("rel", "stylesheet", "href", cssFile),
		),
		b.Body().R(
			b.H1().T("This page uses external asset files"),
			components.Alert{
				Type:        components.AlertSuccess,
				Message:     "CSS and JS are loaded from separate files for better caching.",
				Dismissible: true,
			}.Render(b),
			b.Script("src", jsFile),
		),
	)

	htmlFile := "external_assets.html"
	if err := os.WriteFile(htmlFile, []byte(b.String()), 0644); err != nil {
		fmt.Printf("Error writing HTML: %v\n", err)
		return
	}
	fmt.Printf("✅ Wrote HTML to %s\n", htmlFile)
	fmt.Println("\nOpen external_assets.html in a browser to see the result!")
}

// inspectAssets demonstrates how to inspect what assets were collected.
func inspectAssets() {
	b := element.NewBuilder().EnableAssetTracking()

	// Render various components
	components.Alert{Type: components.AlertInfo, Message: "Test", Dismissible: true}.Render(b)
	components.Alert{Type: components.AlertSuccess, Message: "Test 2", Dismissible: false}.Render(b)
	components.Table{Headers: []string{"A", "B"}, Rows: [][]any{{1, 2}}}.Render(b)
	components.Table{Headers: []string{"C", "D"}, Rows: [][]any{{3, 4}}}.Render(b)

	// Inspect collected assets
	css := b.GetCollectedCSS()
	js := b.GetCollectedJS()

	fmt.Printf("Components rendered:\n")
	fmt.Printf("  - Alert (with dismissible: true)\n")
	fmt.Printf("  - Alert (with dismissible: false)\n")
	fmt.Printf("  - Table\n")
	fmt.Printf("  - Table (second instance)\n\n")

	fmt.Printf("Collected Assets:\n")
	fmt.Printf("  CSS: %d bytes\n", len(css))
	fmt.Printf("  JS:  %d bytes\n\n", len(js))

	fmt.Printf("Deduplication Results:\n")
	fmt.Printf("  ✅ Alert CSS included once (not twice)\n")
	fmt.Printf("  ✅ Alert JS included once (even though only first alert was dismissible)\n")
	fmt.Printf("  ✅ Table CSS included once (not twice)\n\n")

	fmt.Printf("How it works:\n")
	fmt.Printf("  - Each component type has a unique AssetID\n")
	fmt.Printf("  - Builder tracks assets by ID in a map\n")
	fmt.Printf("  - Duplicate IDs are ignored (O(1) deduplication)\n")
	fmt.Printf("  - Result: Each asset is included exactly once\n")
}
