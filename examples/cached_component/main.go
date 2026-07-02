// Package main demonstrates component caching with element.Cached.
//
// Static components (nav bars, footers, icon sets) are often rebuilt on every
// request even though their output never changes. Wrapping such a component
// with element.Cached renders it once and replays the cached bytes on every
// subsequent render -- a ~20x speedup with a single allocation.
//
// Run with: go run .
// Then visit: http://localhost:8080
package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/rohanthewiz/element"
)

// renderCount tracks how many times NavBar actually renders,
// proving the cache is doing its job
var renderCount atomic.Int64

// -----------------------------------------------------------------------------
// NavBar - a static component worth caching
// -----------------------------------------------------------------------------

type NavBar struct{}

func (n NavBar) Render(b *element.Builder) (x any) {
	renderCount.Add(1)

	b.NavClass("navbar").R(
		b.SpanClass("brand").T("Cached Components"),
		b.Ul().R(
			b.Li().R(b.A("href", "/").T("Home")),
			b.Li().R(b.A("href", "/about").T("About")),
			b.Li().R(b.A("href", "/contact").T("Contact")),
		),
	)
	return
}

// -----------------------------------------------------------------------------
// Footer - another static component
// -----------------------------------------------------------------------------

type Footer struct{}

func (f Footer) Render(b *element.Builder) (x any) {
	b.FooterClass("footer").R(
		b.P().T("Built with github.com/rohanthewiz/element"),
	)
	return
}

// Wrap the static components once at package level.
// Every handler reuses these; the underlying components render exactly once.
// (Cached is safe for concurrent use.)
var (
	cachedNav    = element.Cached(NavBar{})
	cachedFooter = element.Cached(Footer{})
)

func main() {
	http.HandleFunc("/", homeHandler)

	fmt.Println("Cached Component Example Server")
	fmt.Println("===============================")
	fmt.Println("Visit: http://localhost:8080")
	fmt.Println("Reload the page a few times -- the nav render count stays at 1.")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	b := element.AcquireBuilder()
	defer element.ReleaseBuilder(b)

	b.Html().R(
		b.Head().R(
			b.Meta("charset", "utf-8").R(),
			b.Title().T("Cached Component Example"),
			b.Style().T(pageCSS),
		),
		b.Body().R(
			cachedNav.Render(b), // replayed from cache after the first request

			b.DivClass("content").R(
				b.H1().T("Component Caching"),
				// The dynamic part of the page still renders fresh each request
				b.P().F("This page was generated at %s.", time.Now().Format(time.RFC1123)),
				b.P().F("The nav bar above has rendered %d time(s) since the server started, "+
					"no matter how many requests have been served.", renderCount.Load()),
				b.P().R(
					b.T("Wrap any static "),
					b.Code().T("Component"),
					b.T(" with "),
					b.Code().T("element.Cached(comp)"),
					b.T(" to render it once and reuse the bytes thereafter."),
				),
			),

			cachedFooter.Render(b),
		),
	)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(b.Bytes()) // Bytes() avoids a string conversion
}

const pageCSS = `
body { font-family: sans-serif; margin: 0; }
.navbar { display: flex; align-items: center; gap: 2rem; background: #2c3e50; color: #fff; padding: 0.75rem 1.5rem; }
.navbar .brand { font-weight: bold; }
.navbar ul { display: flex; gap: 1rem; list-style: none; margin: 0; padding: 0; }
.navbar a { color: #ecf0f1; text-decoration: none; }
.content { padding: 1.5rem; }
.footer { background: #ecf0f1; padding: 1rem 1.5rem; margin-top: 2rem; }
code { background: #f4f4f4; padding: 0.1rem 0.3rem; border-radius: 3px; }
`
