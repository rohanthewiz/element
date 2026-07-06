// tutorial.js — the interactive element tutorial (the "Tutorial" tab of the
// playground). Lesson data + the tutorial UI.
//
// Each lesson: { id, title, nav, prose, code, task, check, solution }
//   prose    — array of segments: raw-HTML strings, or {code, lang} blocks
//              that render as highlighted snippets (lang: go|html|css|plain)
//   code     — starter Go program for the lesson editor
//   task     — one-line exercise shown above the editor
//   check    — (html, flat) => bool; flat is the program's output with
//              whitespace collapsed to single spaces. Passing marks the
//              lesson complete. Checks always receive the RAW output.
//   solution — source the "solution" button loads
//   rawOut   — show the program's raw output instead of the pretty-printed
//              form (for lessons where the raw text IS the point)
//
// The data half is plain JS so the verification harness can run every
// starter and solution through the real interpreter in CI/dev.
(() => {
'use strict';

const code = (src, lang) => ({ code: src, lang: lang || 'go' });

const LESSONS = [

// ---------------------------------------------------------------------------
{
  id: 'hello',
  title: 'Hello, element',
  nav: 'Hello, element',
  prose: [
`<h2>Hello, element</h2>
<p><a href="https://github.com/rohanthewiz/element" target="_blank" rel="noopener">element</a> generates HTML from plain Go — no templates, no reflection. The editor on the right holds a complete Go program; this page interprets it (with <a href="https://github.com/traefik/yaegi" target="_blank" rel="noopener">yaegi</a> compiled to WebAssembly) and shows whatever it prints. <strong>Stdout is the output pane.</strong></p>
<p>Everything starts with a <code>Builder</code>, which wraps a byte buffer. Each element method (<code>b.Div()</code>, <code>b.H1()</code>, …) writes its opening tag immediately, then one of three methods finishes it:</p>
<ul>
<li><code>.R(children…)</code> — render children, then the closing tag</li>
<li><code>.T(text…)</code> — text-only content, then the closing tag</li>
<li><code>.F(format, args…)</code> — <code>fmt.Sprintf</code>-style text content</li>
</ul>`,
    code(`b := element.NewBuilder()

b.Div().R(
    b.H1().T("Hi!"),
)

fmt.Println(b.String())`),
`<p>Why does this work? Go evaluates a call's <em>arguments before the call itself</em> — so <code>b.H1()</code> writes <code>&lt;h1&gt;</code>, <code>.T("Hi!")</code> writes the text and <code>&lt;/h1&gt;</code>, and only then does the outer <code>.R()</code> close the <code>&lt;/div&gt;</code>. The tree structure of your Go code <em>is</em> the tree structure of the HTML.</p>
<div class="tip">Every lesson is live — the program reruns as you type. Complete the task above the editor to earn a check mark; progress is saved in your browser. Use the <em>preview</em> toggle on the output pane to see the HTML rendered.</div>`,
  ],
  code: `package main

import (
	"fmt"

	"github.com/rohanthewiz/element"
)

func main() {
	b := element.NewBuilder()

	b.Div().R(
		b.H1().T("Hello from element"),
	)

	fmt.Println(b.String())
}
`,
  task: 'Add a <p> below the h1 with the text: Built with Go',
  check: (html, flat) => /<p>Built with Go<\/p>/.test(flat),
  solution: `package main

import (
	"fmt"

	"github.com/rohanthewiz/element"
)

func main() {
	b := element.NewBuilder()

	b.Div().R(
		b.H1().T("Hello from element"),
		b.P().T("Built with Go"),
	)

	fmt.Println(b.String())
}
`,
},

// ---------------------------------------------------------------------------
{
  id: 'attributes',
  title: 'Attributes',
  nav: 'Attributes',
  prose: [
`<h2>Attributes</h2>
<p>Attributes are passed as <code>key, value</code> string pairs after the element name — always an even count, in the order you want them emitted:</p>`,
    code(`b.Div("id", "hero", "class", "wide").R(...)
// <div id="hero" class="wide">...`),
`<p>Because <code>class</code> is by far the most common attribute, every element has a <code>*Class</code> variant that takes the class first, then normal pairs:</p>`,
    code(`b.PClass("note").T("classy")
// <p class="note">classy</p>

b.DivClass("card wide", "id", "c1").R(...)
// <div class="card wide" id="c1">...`),
`<div class="tip">There is no attribute escaping or validation — element writes exactly what you pass. That is a feature (speed, control) and a responsibility: quote-free strings only.</div>`,
  ],
  code: `package main

import (
	"fmt"

	"github.com/rohanthewiz/element"
)

func main() {
	b := element.NewBuilder()

	b.DivClass("links", "id", "resources").R(
		b.H3().T("Resources"),
	)

	fmt.Println(b.String())
}
`,
  task: 'Inside the div, add a link: b.AClass("ext", "href", "https://go.dev").T("Go") — read the html out pane.',
  check: (html, flat) => /<a class="ext" href="https:\/\/go\.dev">Go<\/a>/.test(flat),
  solution: `package main

import (
	"fmt"

	"github.com/rohanthewiz/element"
)

func main() {
	b := element.NewBuilder()

	b.DivClass("links", "id", "resources").R(
		b.H3().T("Resources"),
		b.AClass("ext", "href", "https://go.dev").T("Go"),
	)

	fmt.Println(b.String())
}
`,
},

// ---------------------------------------------------------------------------
{
  id: 'single-tags',
  title: 'Single tags',
  nav: 'Single tags',
  prose: [
`<h2>Single tags</h2>
<p>Void elements — <code>&lt;img&gt;</code>, <code>&lt;br&gt;</code>, <code>&lt;hr&gt;</code>, <code>&lt;input&gt;</code>, <code>&lt;meta&gt;</code>, <code>&lt;link&gt;</code>, and friends — have no closing tag and take no children. element knows which is which, so you don't have to track it: <strong>close everything with <code>.R()</code></strong> and single tags simply won't emit a closing tag.</p>`,
    code(`b.Img("src", "gopher.png", "alt", "a gopher").R()
b.Hr().R()
b.Input("type", "text", "name", "q", "placeholder", "search…").R()`),
    code(`<img src="gopher.png" alt="a gopher">
<hr>
<input type="text" name="q" placeholder="search…">`, 'html'),
`<p>Passing children to a single tag is a bug element will flag for you in debug mode (a later lesson) — the children render <em>after</em> the tag rather than inside it, since inside doesn't exist.</p>`,
  ],
  code: `package main

import (
	"fmt"

	"github.com/rohanthewiz/element"
)

func main() {
	b := element.NewBuilder()

	b.DivClass("figure").R(
		b.H3().T("A fine gopher"),
		b.Hr().R(),
	)

	fmt.Println(b.String())
}
`,
  task: 'Add an <img> after the <hr> with src "gopher.png" and alt "gopher".',
  check: (html, flat) => /<img [^>]*alt="gopher"/.test(flat) && /src="gopher\.png"/.test(flat),
  solution: `package main

import (
	"fmt"

	"github.com/rohanthewiz/element"
)

func main() {
	b := element.NewBuilder()

	b.DivClass("figure").R(
		b.H3().T("A fine gopher"),
		b.Hr().R(),
		b.Img("src", "gopher.png", "alt", "gopher").R(),
	)

	fmt.Println(b.String())
}
`,
},

// ---------------------------------------------------------------------------
{
  id: 'text',
  title: 'Text: T, F, and the builder',
  nav: 'Text & formatting',
  prose: [
`<h2>Text: <code>T</code>, <code>F</code>, and the builder</h2>
<p>Three ways to put text in the tree:</p>
<ul>
<li><code>el.T("a", "b", …)</code> — text children of an element, concatenated</li>
<li><code>el.F("x = %d", n)</code> — formatted text child (<code>fmt.Sprintf</code> semantics)</li>
<li><code>b.T(…)</code> / <code>b.F(…)</code> — text written <em>directly</em> at the current position, no wrapping element. This is how you mix text and inline elements inside <code>.R()</code>:</li>
</ul>`,
    code(`b.P().R(
    b.T("The answer is "),
    b.Strong().F("%d", 42),
    b.T(", obviously."),
)`),
    code(`<p>The answer is <strong>42</strong>, obviously.</p>`, 'html'),
`<div class="tip"><strong>element does not escape text.</strong> What you pass is what ships — ideal for trusted content and pre-rendered fragments. For user-supplied strings in real applications, escape first (e.g. <code>html.EscapeString</code> from the Go stdlib).</div>`,
  ],
  code: `package main

import (
	"fmt"

	"github.com/rohanthewiz/element"
)

func main() {
	b := element.NewBuilder()

	score := 42
	_ = score // use the variable in your task, then delete this line

	b.DivClass("stats").R(
		b.H3().T("Game over"),
		b.P().R(
			b.T("Thanks for playing, "),
			b.Em().T("gopher"),
			b.T("!"),
		),
	)

	fmt.Println(b.String())
}
`,
  task: 'Add a <p> rendering "Score: 42" from the score variable, using .F().',
  check: (html, flat) => /<p>Score: 42<\/p>/.test(flat),
  solution: `package main

import (
	"fmt"

	"github.com/rohanthewiz/element"
)

func main() {
	b := element.NewBuilder()

	score := 42

	b.DivClass("stats").R(
		b.H3().T("Game over"),
		b.P().R(
			b.T("Thanks for playing, "),
			b.Em().T("gopher"),
			b.T("!"),
		),
		b.P().F("Score: %d", score),
	)

	fmt.Println(b.String())
}
`,
},

// ---------------------------------------------------------------------------
{
  id: 'lists',
  title: 'Slices: ForEach',
  nav: 'Slices & ForEach',
  prose: [
`<h2>Slices: <code>ForEach</code></h2>
<p>The render tree is Go, so a slice becomes markup with an ordinary loop. <code>element.ForEach</code> packages the pattern as a child-position helper (it is a generic function — any element type works):</p>`,
    code(`fruits := []string{"apple", "banana", "cherry"}

b.Ul().R(
    element.ForEach(fruits, func(f string) {
        b.Li().T(f)
    }),
)`),
`<p><code>ForEach2</code> also hands you the index:</p>`,
    code(`element.ForEach2(fruits, func(f string, i int) {
    b.Li().F("%d. %s", i+1, f)
})`),
`<div class="tip">Both return a throwaway <code>any</code> so they can sit directly in an <code>.R(…)</code> argument list — the same convention every element method follows.</div>`,
  ],
  code: `package main

import (
	"fmt"

	"github.com/rohanthewiz/element"
)

func main() {
	b := element.NewBuilder()

	planets := []string{"Mercury", "Venus", "Earth"}

	b.Div().R(
		b.H3().T("Planets"),
		b.Ul().R(
			element.ForEach(planets, func(p string) {
				b.Li().T(p)
			}),
		),
	)

	fmt.Println(b.String())
}
`,
  task: 'Change the <ul> to an <ol> whose items read "1. Mercury", "2. Venus", "3. Earth" — use ForEach2.',
  check: (html, flat) => /<ol><li>1\. Mercury<\/li><li>2\. Venus<\/li><li>3\. Earth<\/li><\/ol>/.test(flat),
  solution: `package main

import (
	"fmt"

	"github.com/rohanthewiz/element"
)

func main() {
	b := element.NewBuilder()

	planets := []string{"Mercury", "Venus", "Earth"}

	b.Div().R(
		b.H3().T("Planets"),
		b.Ol().R(
			element.ForEach2(planets, func(p string, i int) {
				b.Li().F("%d. %s", i+1, p)
			}),
		),
	)

	fmt.Println(b.String())
}
`,
},

// ---------------------------------------------------------------------------
{
  id: 'conditionals',
  title: 'Conditionals',
  nav: 'Conditionals',
  prose: [
`<h2>Conditionals</h2>
<p>Two idioms put an <code>if</code> inside a render tree.</p>
<h3>1. <code>b.Wrap</code> — run statements</h3>`,
    code(`b.Div().R(
    b.Wrap(func() {
        if isEvening {
            b.H3().T("Good evening!")
        } else {
            b.H3().T("Good day!")
        }
    }),
)`),
`<h3>2. An immediately-invoked func — return a child</h3>
<p>When a branch should <em>be</em> the child value, use a small function literal returning <code>(x any)</code>, called in place:</p>`,
    code(`b.Div().R(
    func() (x any) {
        if count == 0 {
            return b.P().T("empty!")
        }
        return b.P().F("%d items", count)
    }(),
)`),
`<div class="tip">Both are plain Go — no template mini-language to learn. Pick <code>Wrap</code> for side-effect style, the IIFE for value style; they render identically.</div>`,
  ],
  code: `package main

import (
	"fmt"

	"github.com/rohanthewiz/element"
)

func main() {
	b := element.NewBuilder()

	loggedIn := false

	b.DivClass("nav").R(
		b.Wrap(func() {
			if loggedIn {
				b.SpanClass("user").T("Welcome back!")
			}
		}),
	)

	fmt.Println(b.String())
}
`,
  task: 'Add an else branch that renders a sign-in link: <a href="/login">Sign in</a>.',
  check: (html, flat) => /<a href="\/login">Sign in<\/a>/.test(flat),
  solution: `package main

import (
	"fmt"

	"github.com/rohanthewiz/element"
)

func main() {
	b := element.NewBuilder()

	loggedIn := false

	b.DivClass("nav").R(
		b.Wrap(func() {
			if loggedIn {
				b.SpanClass("user").T("Welcome back!")
			} else {
				b.A("href", "/login").T("Sign in")
			}
		}),
	)

	fmt.Println(b.String())
}
`,
},

// ---------------------------------------------------------------------------
{
  id: 'components',
  title: 'Components',
  nav: 'Components',
  prose: [
`<h2>Components</h2>
<p>A component is any type implementing one method:</p>`,
    code(`type Component interface {
    Render(b *Builder) (x any)
}`),
`<p>Fields are your props; <code>Render</code> writes through the shared builder. Render one or many anywhere in a tree with <code>b.RenderComps(…)</code>:</p>`,
    code(`type Greeting struct {
    Name string
}

func (g Greeting) Render(b *element.Builder) (x any) {
    b.PClass("greet").F("Hello, %s!", g.Name)
    return
}

b.Div().R(
    b.RenderComps(Greeting{Name: "Ada"}, Greeting{Name: "Grace"}),
)`),
`<div class="tip">Components compose: a page component renders section components which render widget components — all sharing one builder, one buffer, zero copying.</div>`,
  ],
  code: `package main

import (
	"fmt"

	"github.com/rohanthewiz/element"
)

type Todo struct {
	Title    string
	Priority string
}

func (t Todo) Render(b *element.Builder) (x any) {
	b.Li().R(
		b.T(t.Title),
	)
	return
}

func main() {
	b := element.NewBuilder()

	todos := []Todo{
		{Title: "Write docs", Priority: "high"},
		{Title: "Ship v1", Priority: "medium"},
	}

	b.Ul().R(
		element.ForEach(todos, func(t Todo) {
			b.RenderComps(t)
		}),
	)

	fmt.Println(b.String())
}
`,
  task: 'Extend Todo.Render to append the priority as <span class="badge">high</span> (etc.) inside each <li>.',
  check: (html, flat) => /<span class="badge">high<\/span>/.test(flat) && /<span class="badge">medium<\/span>/.test(flat),
  solution: `package main

import (
	"fmt"

	"github.com/rohanthewiz/element"
)

type Todo struct {
	Title    string
	Priority string
}

func (t Todo) Render(b *element.Builder) (x any) {
	b.Li().R(
		b.T(t.Title),
		b.SpanClass("badge").T(t.Priority),
	)
	return
}

func main() {
	b := element.NewBuilder()

	todos := []Todo{
		{Title: "Write docs", Priority: "high"},
		{Title: "Ship v1", Priority: "medium"},
	}

	b.Ul().R(
		element.ForEach(todos, func(t Todo) {
			b.RenderComps(t)
		}),
	)

	fmt.Println(b.String())
}
`,
},

// ---------------------------------------------------------------------------
{
  id: 'page',
  title: 'A full page',
  nav: 'Full pages',
  prose: [
`<h2>A full page</h2>
<p><code>b.Html()</code> writes the doctype and opens <code>&lt;html&gt;</code>; from there, head and body are just elements. For the common shape there is a one-call helper:</p>`,
    code(`page := b.HtmlPage(
    "body { font-family: sans-serif }",   // inner <style>
    "<title>Hi</title>",                  // rest of <head>
    Body{},                               // a Component for <body>
)`),
`<p>Or spell it out when you want full control:</p>`,
    code(`b.Html("lang", "en").R(
    b.Head().R(
        b.Title().T("Hi"),
        b.Style().T("body { margin: 0 }"),
    ),
    b.Body().R(
        b.H1().T("Hello"),
    ),
)`),
`<div class="tip">Switch the output pane to <em>preview</em> — a full page carries its own styles, so this is where the tutorial gets visual.</div>`,
  ],
  code: `package main

import (
	"fmt"

	"github.com/rohanthewiz/element"
)

func main() {
	b := element.NewBuilder()

	b.Html("lang", "en").R(
		b.Head().R(
			b.Style().T("body { font-family: sans-serif; background: #fdf6e3; padding: 2rem }"),
		),
		b.Body().R(
			b.H1().T("A whole page"),
			b.P().T("doctype, head, body — all from one builder"),
		),
	)

	fmt.Println(b.String())
}
`,
  task: 'Add a <title> of "My Page" to the head.',
  check: (html, flat) => /<title>My Page<\/title>/.test(flat),
  solution: `package main

import (
	"fmt"

	"github.com/rohanthewiz/element"
)

func main() {
	b := element.NewBuilder()

	b.Html("lang", "en").R(
		b.Head().R(
			b.Title().T("My Page"),
			b.Style().T("body { font-family: sans-serif; background: #fdf6e3; padding: 2rem }"),
		),
		b.Body().R(
			b.H1().T("A whole page"),
			b.P().T("doctype, head, body — all from one builder"),
		),
	)

	fmt.Println(b.String())
}
`,
},

// ---------------------------------------------------------------------------
{
  id: 'table',
  title: 'The components package: Table',
  nav: 'Table component',
  prose: [
`<h2>The components package: <code>Table</code></h2>
<p><code>github.com/rohanthewiz/element/components</code> ships ready-made UI pieces built on the same builder. <code>Table</code> turns headers plus rows into a full <code>&lt;table&gt;</code>:</p>`,
    code(`t := components.Table{
    Caption: "Inventory",
    Headers: []string{"Item", "Qty"},
    Rows: [][]any{
        {"carrots", 12},
        {"turnips", 7},
    },
    Striped: true,
}
b.RenderComps(t)`),
`<p>Cells are <code>any</code>: strings, numbers, and booleans are written directly; a cell that is itself a <code>Component</code> (a Button, a Badge, a link…) renders in place. <code>Columns</code> gives per-column alignment, <code>FooterRows</code> a <code>&lt;tfoot&gt;</code>, <code>EmptyMessage</code> a friendly empty state.</p>
<div class="tip">Styling flags (<code>Striped</code>, <code>Bordered</code>, <code>Hover</code>) only add classes like <code>table-striped</code> — bring your own CSS, or none.</div>`,
  ],
  code: `package main

import (
	"fmt"

	"github.com/rohanthewiz/element"
	"github.com/rohanthewiz/element/components"
)

func main() {
	b := element.NewBuilder()

	t := components.Table{
		Caption: "Garden inventory",
		Headers: []string{"Item", "Qty"},
		Rows: [][]any{
			{"carrots", 12},
			{"turnips", 7},
		},
	}

	b.RenderComps(t)

	fmt.Println(b.String())
}
`,
  task: 'Make the table striped (hint: one field) and confirm the class table-striped appears in the output.',
  check: (html, flat) => /table-striped/.test(flat),
  solution: `package main

import (
	"fmt"

	"github.com/rohanthewiz/element"
	"github.com/rohanthewiz/element/components"
)

func main() {
	b := element.NewBuilder()

	t := components.Table{
		Caption: "Garden inventory",
		Headers: []string{"Item", "Qty"},
		Rows: [][]any{
			{"carrots", 12},
			{"turnips", 7},
		},
		Striped: true,
	}

	b.RenderComps(t)

	fmt.Println(b.String())
}
`,
},

// ---------------------------------------------------------------------------
{
  id: 'uikit',
  title: 'More components: Card, Alert & friends',
  nav: 'UI kit tour',
  prose: [
`<h2>More components: Card, Alert &amp; friends</h2>
<p>The kit covers the usual suspects — each one a plain struct, each rendering semantic markup with predictable class names:</p>
<ul>
<li><code>Card</code> — title / body / footer container</li>
<li><code>Alert</code> — info, success, warning, error notices</li>
<li><code>Badge</code>, <code>Button</code>, <code>ProgressBar</code>, <code>Breadcrumb</code>, <code>Nav</code>, <code>Pagination</code>…</li>
<li><code>Tabs</code>, <code>Accordion</code>, <code>Modal</code>, <code>Dropdown</code> — CSS-only interactivity, no JavaScript</li>
<li>Form controls: <code>FormField</code>, <code>SelectField</code>, <code>CheckboxField</code>, <code>RadioGroup</code>…</li>
</ul>`,
    code(`components.Alert{
    Type:    components.AlertSuccess,
    Title:   "Deployed",
    Message: "All 3 services are green.",
}`),
`<div class="tip">Components are ordinary values — build slices of them, pass them as props (<code>Card.BodyComponent</code>, <code>Table</code> cells), or wrap them in your own types.</div>`,
  ],
  code: `package main

import (
	"fmt"

	"github.com/rohanthewiz/element"
	"github.com/rohanthewiz/element/components"
)

func main() {
	b := element.NewBuilder()

	b.DivClass("dashboard").R(
		b.RenderComps(
			components.Card{
				Title: "Deploy status",
				Body:  "3 of 3 services healthy.",
			},
			components.ProgressBar{Value: 82, ShowValue: true, Label: "rollout"},
		),
	)

	fmt.Println(b.String())
}
`,
  task: 'Add a success Alert with the message "Saved!" to the dashboard.',
  check: (html, flat) => /alert-success/.test(flat) && /Saved!/.test(flat),
  solution: `package main

import (
	"fmt"

	"github.com/rohanthewiz/element"
	"github.com/rohanthewiz/element/components"
)

func main() {
	b := element.NewBuilder()

	b.DivClass("dashboard").R(
		b.RenderComps(
			components.Card{
				Title: "Deploy status",
				Body:  "3 of 3 services healthy.",
			},
			components.ProgressBar{Value: 82, ShowValue: true, Label: "rollout"},
			components.Alert{
				Type:    components.AlertSuccess,
				Message: "Saved!",
			},
		),
	)

	fmt.Println(b.String())
}
`,
},

// ---------------------------------------------------------------------------
{
  id: 'debug',
  title: 'Debug mode',
  nav: 'Debug mode',
  prose: [
`<h2>Debug mode</h2>
<p>Because the tree is built by execution order, the classic mistakes are <em>structural</em>: a bare string where a child should be, children handed to a single tag, an element never closed. <code>element.DebugSet()</code> turns on a checker that catches them as the tree renders:</p>`,
    code(`element.DebugSet()
defer element.DebugClear()

// ... build ...

fmt.Println(element.DebugShow(element.DebugOptions{TextOnly: true}))`),
`<p>Concerns are collected per element (with the source function and file:line when compiled natively) and <code>DebugShow</code> prints a report — as HTML by default, or a text table with <code>TextOnly</code>. When the tree is clean it prints <code>No element concerns found.</code></p>
<div class="tip">Debug mode adds <code>data-ele-id</code> attributes and stack lookups — development only. The zero-cost path is the default.</div>`,
  ],
  code: `package main

import (
	"fmt"

	"github.com/rohanthewiz/element"
)

func main() {
	element.DebugSet()
	defer element.DebugClear()

	b := element.NewBuilder()

	b.Div().R(
		b.H1().T("Spot the bugs"),
		"a bare string — should be b.T(...)",
		b.Br().R(
			b.P().T("a <br> cannot hold children"),
		),
	)

	fmt.Println(element.DebugShow(element.DebugOptions{TextOnly: true}))
}
`,
  task: 'Fix both issues (wrap the string in b.T(), move the <p> out of the <br>) so the report prints: No element concerns found.',
  rawOut: true, // the report is text — show it exactly as printed
  check: (html, flat) => /No element concerns found\./.test(flat),
  solution: `package main

import (
	"fmt"

	"github.com/rohanthewiz/element"
)

func main() {
	element.DebugSet()
	defer element.DebugClear()

	b := element.NewBuilder()

	b.Div().R(
		b.H1().T("Spot the bugs"),
		b.T("a bare string — should be b.T(...)"),
		b.Br().R(),
		b.P().T("a <br> cannot hold children"),
	)

	fmt.Println(element.DebugShow(element.DebugOptions{TextOnly: true}))
}
`,
},

// ---------------------------------------------------------------------------
{
  id: 'perf',
  title: 'Performance: pooling & caching',
  nav: 'Pooling & caching',
  prose: [
`<h2>Performance: pooling &amp; caching</h2>
<p>element is already allocation-light (one buffer, no reflection). Two helpers go further in hot paths like HTTP handlers:</p>
<h3>Builder pooling</h3>`,
    code(`b := element.AcquireBuilder()
defer element.ReleaseBuilder(b)   // resets and returns it to the pool
// ... build, then copy the output out before release:
out := b.String()`),
`<h3>Component caching</h3>
<p><code>element.Cached(comp)</code> wraps any component; the first render runs <code>Render</code>, every later render replays the recorded bytes — for headers, footers, nav bars that don't change per request:</p>`,
    code(`footer := element.Cached(Footer{})   // renders once, replays after`),
`<div class="tip">The editor's program renders the footer on three pages and counts how many times <code>Render</code> actually ran. Watch the last line of output.</div>`,
  ],
  code: `package main

import (
	"fmt"

	"github.com/rohanthewiz/element"
)

var renders = 0

type Footer struct{}

func (Footer) Render(b *element.Builder) (x any) {
	renders++
	b.Footer().T("© 2026 gopher industries")
	return
}

func main() {
	var footer element.Component = Footer{}

	for page := 1; page <= 3; page++ {
		b := element.NewBuilder()
		b.Div().R(
			b.H3().F("page %d", page),
			b.RenderComps(footer),
		)
		fmt.Println(b.String())
	}

	fmt.Printf("footer rendered %d time(s)\\n", renders)
}
`,
  task: 'Wrap the footer in element.Cached(...) so it renders only 1 time for all three pages.',
  check: (html, flat) => /footer rendered 1 time\(s\)/.test(flat),
  solution: `package main

import (
	"fmt"

	"github.com/rohanthewiz/element"
)

var renders = 0

type Footer struct{}

func (Footer) Render(b *element.Builder) (x any) {
	renders++
	b.Footer().T("© 2026 gopher industries")
	return
}

func main() {
	footer := element.Cached(Footer{})

	for page := 1; page <= 3; page++ {
		b := element.NewBuilder()
		b.Div().R(
			b.H3().F("page %d", page),
			b.RenderComps(footer),
		)
		fmt.Println(b.String())
	}

	fmt.Printf("footer rendered %d time(s)\\n", renders)
}
`,
},

// ---------------------------------------------------------------------------
{
  id: 'pretty',
  title: 'Pretty output',
  nav: 'Pretty output',
  prose: [
`<h2>Pretty output</h2>
<p>element emits compact, single-line HTML — ideal to ship, hard to read. For debugging (or teaching!), <code>b.Pretty()</code> returns the same document indented:</p>`,
    code(`b.Div().R(b.P().T("hi"))

b.String()  // <div><p>hi</p></div>
b.Pretty()  // <div>
            //   <p>hi</p>
            // </div>`),
`<p>Formatting happens on the way out — the builder still stores the compact form, and <code>PrettyHTML(s)</code> is available as a standalone function for HTML from anywhere.</p>
<div class="tip">Pretty output is for eyes, not for diffing or storage: whitespace inside <code>&lt;pre&gt;</code> and inline text can matter in HTML, so ship <code>b.String()</code>.</div>
<p><em>Note:</em> this tutorial's output pane normally pretty-prints for readability (it runs <code>PrettyHTML</code> for you, like the Playground tab's <code>pretty</code> checkbox) — <strong>this lesson alone shows your program's raw output</strong>, so you can see exactly what <code>b.String()</code> vs <code>b.Pretty()</code> produce.</p>`,
  ],
  rawOut: true, // showing raw vs pretty IS the lesson
  code: `package main

import (
	"fmt"

	"github.com/rohanthewiz/element"
)

func main() {
	b := element.NewBuilder()

	b.DivClass("post").R(
		b.H2().T("On formatting"),
		b.P().R(
			b.T("Readable "),
			b.Em().T("and"),
			b.T(" shippable."),
		),
		b.Ul().R(
			b.Li().T("compact by default"),
			b.Li().T("pretty on demand"),
		),
	)

	fmt.Println(b.String())
}
`,
  task: 'Print the pretty form instead — the output pane should show indented, multi-line HTML.',
  check: (html, flat) => /\n\s+<li>/.test(html),
  solution: `package main

import (
	"fmt"

	"github.com/rohanthewiz/element"
)

func main() {
	b := element.NewBuilder()

	b.DivClass("post").R(
		b.H2().T("On formatting"),
		b.P().R(
			b.T("Readable "),
			b.Em().T("and"),
			b.T(" shippable."),
		),
		b.Ul().R(
			b.Li().T("compact by default"),
			b.Li().T("pretty on demand"),
		),
	)

	fmt.Println(b.Pretty())
}
`,
},

// ---------------------------------------------------------------------------
{
  id: 'next',
  title: 'Where to next',
  nav: 'Where to next',
  prose: [
`<h2>Where to next</h2>
<p>You have the whole model: a builder writing through execution order, <code>.R/.T/.F</code>, attributes, single tags, slices, conditionals, components, full pages, the components kit, debug mode, and the performance helpers. The editor holds a small showcase — tinker freely.</p>
<p>Beyond the browser:</p>
<ul>
<li><strong>Serve it</strong> — element pairs naturally with any Go web framework; with <a href="https://github.com/rohanthewiz/rweb" target="_blank" rel="noopener">rweb</a> a handler is <code>return s.WriteHTML(b.String())</code>.</li>
<li><strong>Style it</strong> — <a href="https://github.com/rohanthewiz/go-styl" target="_blank" rel="noopener">go-styl</a> compiles Stylus to CSS in-process (it has a playground and tutorial just like this one).</li>
<li><strong>Examples</strong> — the repo's <code>examples/</code> directory covers forms, interfaces, pooling, and cached components; the picker in the Playground tab has runnable tours.</li>
<li><strong>README</strong> — <a href="https://github.com/rohanthewiz/element#readme" target="_blank" rel="noopener">github.com/rohanthewiz/element</a> for the full API and benchmarks.</li>
</ul>
<div class="tip">This tutorial runs element itself, interpreted in your browser — the same source that renders your production HTML. No mock, no subset.</div>`,
  ],
  code: `package main

import (
	"fmt"

	"github.com/rohanthewiz/element"
	"github.com/rohanthewiz/element/components"
)

type Feature struct {
	Name string
	Done bool
}

func (f Feature) Render(b *element.Builder) (x any) {
	b.Li().R(
		b.Wrap(func() {
			if f.Done {
				b.SpanClass("done").T("✓ ")
			} else {
				b.SpanClass("todo").T("• ")
			}
		}),
		b.T(f.Name),
	)
	return
}

func main() {
	b := element.NewBuilder()

	feats := []Feature{
		{"builder & components", true},
		{"debug mode", true},
		{"your next app", false},
	}

	b.Html().R(
		b.Head().R(
			b.Title().T("element tour"),
			b.Style().T("body{font-family:sans-serif;padding:2rem} .done{color:green} .todo{color:#999}"),
		),
		b.Body().R(
			b.H1().T("A little of everything"),
			b.Ul().R(
				element.ForEach(feats, func(f Feature) {
					b.RenderComps(f)
				}),
			),
			b.RenderComps(components.ProgressBar{Value: 2, Max: 3, ShowValue: true, Label: "progress"}),
		),
	)

	fmt.Println(b.String())
}
`,
  task: 'No task here — you made it. Edit freely, or head to the Playground tab.',
  check: null,
  solution: null,
},
];

// ---------------------------------------------------------------------------
// Tutorial UI
// ---------------------------------------------------------------------------
function init() {
  const $ = id => document.getElementById(id);
  const doc = $('tut-doc'), navEl = $('tut-navlist');
  const ta = $('tsrc'), hlCode = $('tsrcHl');
  const outEl = $('tout'), errEl = $('terr'), previewEl = $('tpreview');
  const taskEl = $('tut-task'), statEl = $('tut-check');
  const posEl = $('tut-pos');

  const store = {
    read(k, dflt) { try { return localStorage.getItem(k) ?? dflt; } catch (_) { return dflt; } },
    write(k, v) { try { localStorage.setItem(k, v); } catch (_) {} },
    del(k) { try { localStorage.removeItem(k); } catch (_) {} },
  };

  let done;
  try { done = new Set(JSON.parse(store.read('go-ele-tut-done', '[]'))); }
  catch (_) { done = new Set(); }
  let cur = Math.min(LESSONS.length - 1,
    Math.max(0, parseInt(store.read('go-ele-tut-cur', '0'), 10) || 0));
  let lastRaw = '', lastShown = '';

  const hlOn = () => !document.body.classList.contains('nohl');
  const repaint = eleHi.editor(ta, hlCode, hlOn);

  // html / preview toggle on the lesson output pane
  let view = store.read('go-ele-tut-view', 'html');
  function setView(v) {
    view = v;
    outEl.hidden = v !== 'html';
    previewEl.hidden = v !== 'preview';
    $('tview-html').classList.toggle('on', v === 'html');
    $('tview-preview').classList.toggle('on', v === 'preview');
    if (v === 'preview') previewEl.srcdoc = lastRaw;
    store.write('go-ele-tut-view', v);
  }
  $('tview-html').addEventListener('click', () => setView('html'));
  $('tview-preview').addEventListener('click', () => setView('preview'));

  function saveDone() { store.write('go-ele-tut-done', JSON.stringify([...done])); }

  function block(seg) {
    if (typeof seg === 'string') return seg;
    const lang = seg.lang || 'go';
    const body = lang === 'go' ? eleHi.go(seg.code)
               : lang === 'html' ? eleHi.html(seg.code)
               : lang === 'css' ? eleHi.css(seg.code)
               : eleHi.escape(seg.code);
    return '<pre class="snip lang-' + lang + '"><code>' + body + '</code></pre>';
  }

  function renderNav() {
    navEl.innerHTML = '';
    LESSONS.forEach((l, i) => {
      const li = document.createElement('li');
      if (i === cur) li.className = 'cur';
      li.innerHTML = '<span class="n">' + (i + 1) + '</span>' + eleHi.escape(l.nav) +
        (done.has(l.id) ? '<span class="tick">✓</span>' : '');
      li.addEventListener('click', () => open(i));
      navEl.appendChild(li);
    });
    posEl.textContent = (cur + 1) + ' / ' + LESSONS.length;
    $('tut-prev').disabled = cur === 0;
    $('tut-next').disabled = cur === LESSONS.length - 1;
    $('tut-solution').style.display = LESSONS[cur].solution ? '' : 'none';
  }

  function open(i) {
    cur = i;
    store.write('go-ele-tut-cur', String(i));
    const l = LESSONS[i];
    doc.innerHTML = l.prose.map(block).join('');
    doc.scrollTop = 0;
    ta.value = store.read('go-ele-tut-draft-' + l.id, null) ?? l.code;
    taskEl.textContent = l.task || '';
    if (!l.check && !done.has(l.id)) { done.add(l.id); saveDone(); }
    renderNav();
    repaint();
    run();
  }

  function renderOut() {
    outEl.innerHTML = '';
    if (hlOn()) outEl.innerHTML = eleHi.html(lastShown, true);
    else outEl.textContent = lastShown;
  }

  function run() {
    if (!window.eleGo) return;
    const l = LESSONS[cur];
    const r = eleGo.run(ta.value);
    if (r.error !== undefined) {
      errEl.textContent = (r.line ? 'line ' + r.line + ':' + r.col + ' — ' : '') + r.error
        + (r.stderr ? '\n' + r.stderr : '');
      errEl.style.display = 'block';
      outEl.style.opacity = '0.45';
    } else {
      errEl.style.display = 'none';
      outEl.style.opacity = '';
      lastRaw = r.html;
      // The pane shows the pretty form for readability; checks and the
      // preview always work from the raw output.
      lastShown = (l.rawOut || r.pretty === undefined) ? r.html : r.pretty;
      renderOut();
      if (view === 'preview') previewEl.srcdoc = lastRaw;
      if (l.check) {
        const flat = r.html.replace(/\s+/g, ' ');
        if (l.check(r.html, flat)) {
          if (!done.has(l.id)) { done.add(l.id); saveDone(); renderNav(); }
          statEl.textContent = '✓ task complete';
          statEl.className = 'stat ok';
          return;
        }
      }
    }
    if (l.check) {
      statEl.textContent = done.has(l.id) ? '✓ solved earlier' : '○ not yet';
      statEl.className = done.has(l.id) ? 'stat ok' : 'stat';
    } else {
      statEl.textContent = '';
      statEl.className = 'stat';
    }
  }

  let timer = 0;
  ta.addEventListener('input', () => {
    store.write('go-ele-tut-draft-' + LESSONS[cur].id, ta.value);
    clearTimeout(timer);
    timer = setTimeout(run, 300);
  });
  // Tab inserts a real tab, like the playground editor (it's Go).
  ta.addEventListener('keydown', e => {
    if (e.key !== 'Tab') return;
    e.preventDefault();
    const { selectionStart: s, selectionEnd: t, value } = ta;
    ta.value = value.slice(0, s) + '\t' + value.slice(t);
    ta.selectionStart = ta.selectionEnd = s + 1;
    ta.dispatchEvent(new Event('input'));
  });

  $('tut-prev').addEventListener('click', () => cur > 0 && open(cur - 1));
  $('tut-next').addEventListener('click', () => cur < LESSONS.length - 1 && open(cur + 1));
  $('tut-reset').addEventListener('click', () => {
    store.del('go-ele-tut-draft-' + LESSONS[cur].id);
    ta.value = LESSONS[cur].code;
    repaint();
    run();
  });
  $('tut-solution').addEventListener('click', () => {
    const l = LESSONS[cur];
    if (!l.solution) return;
    ta.value = l.solution;
    store.write('go-ele-tut-draft-' + l.id, l.solution);
    repaint();
    run();
  });
  $('tut-toplay').addEventListener('click', () => {
    if (window.playHooks) window.playHooks.openInPlayground(ta.value);
  });

  setView(view);
  open(cur);
  return { run, repaint: () => { repaint(); renderOut(); } };
}

const api = { LESSONS, init };
if (typeof window !== 'undefined') window.eleTutorial = api;
if (typeof module !== 'undefined') module.exports = api;
})();
