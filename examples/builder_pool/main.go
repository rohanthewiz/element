// Package main demonstrates the use of builder pooling for high-throughput scenarios.
//
// This example shows how to use AcquireBuilder/ReleaseBuilder to reduce
// memory allocations in HTTP handlers that render HTML.
package main

import (
	"fmt"
	"log"
	"net/http"
	"runtime"
	"time"

	"github.com/rohanthewiz/element"
)

func main() {
	// Handler using pooled builders (recommended for high-throughput)
	http.HandleFunc("/pooled", pooledHandler)

	// Handler using regular builders (for comparison)
	http.HandleFunc("/regular", regularHandler)

	// Stats endpoint to show memory usage
	http.HandleFunc("/stats", statsHandler)

	// Home page with links
	http.HandleFunc("/", homeHandler)

	fmt.Println("Builder Pool Example Server")
	fmt.Println("===========================")
	fmt.Println("Listening on http://localhost:8080")
	fmt.Println("")
	fmt.Println("Endpoints:")
	fmt.Println("  /         - Home page with links")
	fmt.Println("  /pooled   - Uses AcquireBuilder/ReleaseBuilder (pooled)")
	fmt.Println("  /regular  - Uses NewBuilder (new allocation each time)")
	fmt.Println("  /stats    - Shows memory statistics")
	fmt.Println("")
	fmt.Println("Try running: ab -n 10000 -c 100 http://localhost:8080/pooled")
	fmt.Println("Then:        ab -n 10000 -c 100 http://localhost:8080/regular")
	fmt.Println("Compare /stats before and after each test.")

	log.Fatal(http.ListenAndServe(":8080", nil))
}

// pooledHandler demonstrates the recommended pattern for high-throughput scenarios.
// Using AcquireBuilder/ReleaseBuilder reduces GC pressure by reusing builders.
func pooledHandler(w http.ResponseWriter, r *http.Request) {
	// Acquire a builder from the pool
	b := element.AcquireBuilder()
	// Always release back to the pool when done
	defer element.ReleaseBuilder(b)

	// Build the HTML page
	renderPage(b, "Pooled Builder", "This page was rendered using a pooled builder.")

	// Write the response
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(b.Bytes())
}

// regularHandler shows the traditional approach of creating a new builder each time.
// This works fine but creates more garbage for the GC to collect.
func regularHandler(w http.ResponseWriter, r *http.Request) {
	// Create a new builder (allocates memory each time)
	b := element.NewBuilder()

	// Build the HTML page
	renderPage(b, "Regular Builder", "This page was rendered using a new builder allocation.")

	// Write the response
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(b.String()))
}

// renderPage builds a sample HTML page with the given title and message.
func renderPage(b *element.Builder, title, message string) {
	b.Html().R(
		b.Head().R(
			b.Meta("charset", "utf-8"),
			b.Title().T(title),
			b.Style().T(`
				body { font-family: system-ui, sans-serif; max-width: 800px; margin: 2rem auto; padding: 0 1rem; }
				h1 { color: #333; }
				.card { background: #f5f5f5; padding: 1rem; border-radius: 8px; margin: 1rem 0; }
				table { border-collapse: collapse; width: 100%; margin: 1rem 0; }
				th, td { border: 1px solid #ddd; padding: 8px; text-align: left; }
				th { background: #4a90d9; color: white; }
				tr:nth-child(even) { background: #f9f9f9; }
				a { color: #4a90d9; }
				.nav { margin-bottom: 2rem; }
				.nav a { margin-right: 1rem; }
			`),
		),
		b.Body().R(
			b.Div("class", "nav").R(
				b.A("href", "/").T("Home"),
				b.A("href", "/pooled").T("Pooled"),
				b.A("href", "/regular").T("Regular"),
				b.A("href", "/stats").T("Stats"),
			),
			b.H1().T(title),
			b.Div("class", "card").R(
				b.P().T(message),
				b.P().F("Rendered at: %s", time.Now().Format(time.RFC3339)),
			),
			// Add some content to make the page more realistic
			renderSampleTable(b),
			renderFeatureList(b),
		),
	)
}

// renderSampleTable creates a sample data table.
func renderSampleTable(b *element.Builder) any {
	return b.Div().R(
		b.H2().T("Sample Data"),
		b.Table().R(
			b.THead().R(
				b.Tr().R(
					b.Th().T("ID"),
					b.Th().T("Name"),
					b.Th().T("Status"),
					b.Th().T("Value"),
				),
			),
			b.TBody().R(
				b.Tr().R(b.Td().T("1"), b.Td().T("Alpha"), b.Td().T("Active"), b.Td().T("100")),
				b.Tr().R(b.Td().T("2"), b.Td().T("Beta"), b.Td().T("Pending"), b.Td().T("250")),
				b.Tr().R(b.Td().T("3"), b.Td().T("Gamma"), b.Td().T("Active"), b.Td().T("175")),
				b.Tr().R(b.Td().T("4"), b.Td().T("Delta"), b.Td().T("Inactive"), b.Td().T("300")),
				b.Tr().R(b.Td().T("5"), b.Td().T("Epsilon"), b.Td().T("Active"), b.Td().T("425")),
			),
		),
	)
}

// renderFeatureList creates a sample feature list.
func renderFeatureList(b *element.Builder) any {
	return b.Div().R(
		b.H2().T("Builder Pool Benefits"),
		b.Ul().R(
			b.Li().T("Reduced memory allocations per request"),
			b.Li().T("Lower GC pressure in high-throughput scenarios"),
			b.Li().T("More stable latency (fewer GC pauses)"),
			b.Li().T("Reuses bytes.Buffer capacity"),
			b.Li().T("Thread-safe for concurrent use"),
		),
	)
}

// homeHandler shows the home page with documentation.
func homeHandler(w http.ResponseWriter, r *http.Request) {
	b := element.AcquireBuilder()
	defer element.ReleaseBuilder(b)

	b.Html().R(
		b.Head().R(
			b.Meta("charset", "utf-8"),
			b.Title().T("Builder Pool Example"),
			b.Style().T(`
				body { font-family: system-ui, sans-serif; max-width: 800px; margin: 2rem auto; padding: 0 1rem; }
				h1 { color: #333; }
				pre { background: #f5f5f5; padding: 1rem; border-radius: 8px; overflow-x: auto; }
				code { font-family: 'SF Mono', Consolas, monospace; }
				.nav { margin-bottom: 2rem; }
				.nav a { margin-right: 1rem; color: #4a90d9; }
			`),
		),
		b.Body().R(
			b.Div("class", "nav").R(
				b.A("href", "/").T("Home"),
				b.A("href", "/pooled").T("Pooled"),
				b.A("href", "/regular").T("Regular"),
				b.A("href", "/stats").T("Stats"),
			),
			b.H1().T("Builder Pool Example"),
			b.P().T("This example demonstrates the difference between pooled and regular builder allocation."),

			b.H2().T("Pooled Builder (Recommended)"),
			b.Pre().R(
				b.Code().T(`b := element.AcquireBuilder()
defer element.ReleaseBuilder(b)

b.Html().R(
    b.Body().R(
        b.H1().T("Hello"),
    ),
)
w.Write([]byte(b.String()))`),
			),

			b.H2().T("Regular Builder"),
			b.Pre().R(
				b.Code().T(`b := element.NewBuilder()

b.Html().R(
    b.Body().R(
        b.H1().T("Hello"),
    ),
)
w.Write([]byte(b.String()))`),
			),

			b.H2().T("When to Use Pooling"),
			b.Ul().R(
				b.Li().T("High-traffic HTTP handlers"),
				b.Li().T("APIs rendering HTML responses"),
				b.Li().T("Server-side rendering at scale"),
				b.Li().T("Any scenario where you're creating many short-lived builders"),
			),

			b.H2().T("Benchmark It"),
			b.P().T("Run load tests against /pooled and /regular, then check /stats:"),
			b.Pre().R(
				b.Code().T(`# Test pooled endpoint
ab -n 10000 -c 100 http://localhost:8080/pooled

# Test regular endpoint
ab -n 10000 -c 100 http://localhost:8080/regular`),
			),
		),
	)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(b.String()))
}

// statsHandler shows memory statistics.
func statsHandler(w http.ResponseWriter, r *http.Request) {
	// Force GC to get accurate stats
	runtime.GC()

	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	b := element.AcquireBuilder()
	defer element.ReleaseBuilder(b)

	b.Html().R(
		b.Head().R(
			b.Meta("charset", "utf-8"),
			b.Title().T("Memory Statistics"),
			b.Meta("http-equiv", "refresh", "content", "5"),
			b.Style().T(`
				body { font-family: system-ui, sans-serif; max-width: 800px; margin: 2rem auto; padding: 0 1rem; }
				h1 { color: #333; }
				table { border-collapse: collapse; width: 100%; }
				th, td { border: 1px solid #ddd; padding: 8px; text-align: left; }
				th { background: #4a90d9; color: white; }
				.nav { margin-bottom: 2rem; }
				.nav a { margin-right: 1rem; color: #4a90d9; }
				.note { color: #666; font-size: 0.9rem; }
			`),
		),
		b.Body().R(
			b.Div("class", "nav").R(
				b.A("href", "/").T("Home"),
				b.A("href", "/pooled").T("Pooled"),
				b.A("href", "/regular").T("Regular"),
				b.A("href", "/stats").T("Stats"),
			),
			b.H1().T("Memory Statistics"),
			b.P("class", "note").T("Auto-refreshes every 5 seconds. GC runs before each measurement."),
			b.Table().R(
				b.THead().R(
					b.Tr().R(
						b.Th().T("Metric"),
						b.Th().T("Value"),
						b.Th().T("Description"),
					),
				),
				b.TBody().R(
					b.Tr().R(
						b.Td().T("Alloc"),
						b.Td().F("%d MB", m.Alloc/1024/1024),
						b.Td().T("Currently allocated heap memory"),
					),
					b.Tr().R(
						b.Td().T("TotalAlloc"),
						b.Td().F("%d MB", m.TotalAlloc/1024/1024),
						b.Td().T("Cumulative bytes allocated (never decreases)"),
					),
					b.Tr().R(
						b.Td().T("Sys"),
						b.Td().F("%d MB", m.Sys/1024/1024),
						b.Td().T("Total memory obtained from OS"),
					),
					b.Tr().R(
						b.Td().T("NumGC"),
						b.Td().F("%d", m.NumGC),
						b.Td().T("Number of completed GC cycles"),
					),
					b.Tr().R(
						b.Td().T("Mallocs"),
						b.Td().F("%d", m.Mallocs),
						b.Td().T("Cumulative count of heap allocations"),
					),
					b.Tr().R(
						b.Td().T("Frees"),
						b.Td().F("%d", m.Frees),
						b.Td().T("Cumulative count of heap frees"),
					),
					b.Tr().R(
						b.Td().T("HeapObjects"),
						b.Td().F("%d", m.HeapObjects),
						b.Td().T("Number of allocated heap objects"),
					),
					b.Tr().R(
						b.Td().T("Goroutines"),
						b.Td().F("%d", runtime.NumGoroutine()),
						b.Td().T("Current number of goroutines"),
					),
				),
			),
			b.P().F("Measured at: %s", time.Now().Format(time.RFC3339)),
		),
	)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(b.Bytes())
}
