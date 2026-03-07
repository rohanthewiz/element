# Element Library Architecture Review

> A deep dive into the architecture of Element, a Go library for programmatic HTML generation.

## Executive Summary

Element is a Go library that generates HTML programmatically without templates. Its core innovation is leveraging Go's natural function execution order (AST) to create an intuitive API that mirrors HTML's tree structure. With 7+ years of production use and ~3,800 lines of code, Element provides a type-safe, high-performance alternative to traditional templating engines.

---

## Core Architecture

### Design Philosophy

1. **Code as Templates** - HTML structure is expressed as Go code, not template strings
2. **Natural Nesting** - Go's function execution order matches HTML's tree structure
3. **Single-Pass Rendering** - HTML written to buffer immediately; no multi-pass compilation
4. **Zero Reflection** - Pure compiled Go for maximum performance
5. **Minimal Surface Area** - Small, focused API that's easy to learn

### The Execution Order Insight

The fundamental insight of Element is that Go's function argument evaluation order naturally creates parent-child relationships:

```go
b.Div().R(           // 1. Opens <div>
    b.P().T("text"), // 2. Opens <p>, renders "text", closes </p>
)                    // 3. Closes </div>

// Result: <div><p>text</p></div>
```

This eliminates the need for template syntax while maintaining readability.

---

## Module Architecture

```
element/
├── Core Layer
│   ├── element.go          - Element struct & rendering logic
│   ├── builder.go          - Builder container & string accumulator
│   ├── component.go        - Component interface for composition
│   └── tags.go             - Single-tag element definitions
│
├── API Layer
│   ├── builder_elements.go           - 84 element creation methods
│   ├── builder_elements_with_class.go - 90+ class convenience methods
│   ├── builder_funcs.go              - Functional helpers (Wrap, HtmlPage)
│   └── element_funcs.go              - Global helpers (ForEach, Vars)
│
├── Support Layer
│   ├── helpers.go          - Utility functions
│   ├── element_debug.go    - Debug system & visual reporting
│   └── assets/             - Embedded CSS/JS for debug UI
│
└── Tests & Examples
    ├── *_test.go           - Comprehensive test coverage
    └── example/            - Integration examples
```

### Layer Responsibilities

| Layer | Purpose | Key Principle |
|-------|---------|---------------|
| **Core** | Element lifecycle and rendering | Minimal, stable foundation |
| **API** | Developer-facing convenience methods | Ergonomic, discoverable |
| **Support** | Debug tools and utilities | Optional, non-intrusive |

---

## Core Components

### Element Struct

The `Element` struct represents an HTML element during construction:

```go
type Element struct {
    name       string            // Tag name (div, p, span, etc.)
    id         string            // Debug-mode tracking ID
    attrs      map[string]string // Element attributes
    sb         *bytes.Buffer  // Reference to builder's buffer
    issues     []string          // Debug-mode issues
    function   string            // Creation context (debug)
    location   string            // File:line (debug)
}
```

**Lifecycle:**
1. `New()` - Creates element, writes opening tag immediately
2. Attributes written as part of opening tag
3. Children rendered via termination methods
4. `R()`/`T()`/`F()` - Writes closing tag

### Builder Struct

The `Builder` is the primary API entry point:

```go
type Builder struct {
    Ele  elementFunc         // Element creation function
    Text textFunc            // Text node function
    s    *bytes.Buffer    // Internal HTML accumulator
}
```

**Key Methods:**
- `NewBuilder()` - Factory for fresh builder
- `String()` - Returns accumulated HTML
- `Reset()` - Clears buffer for reuse
- 84+ element methods (`Div()`, `P()`, `A()`, etc.)

### Component Interface

Simple interface enabling composition:

```go
type Component interface {
    Render(b *Builder) (x any)
}
```

Components can be rendered inline: `b.Div().R(myComponent.Render(b))`

---

## Termination Pattern

Every non-single element must be terminated. This is enforced by the API design:
Single tags are not required to be terminated, but for consistency just terminate with an R().

| Method | Use Case | Example |
|--------|----------|---------|
| `R(children...)` | Mixed/multiple children | `b.Div().R(b.P(), b.Span())` |
| `T(texts...)` | Text-only content | `b.P().T("Hello")` |
| `F(fmt, args)` | Formatted text | `b.Span().F("Count: %d", 42)` |
| `R()` (recommended) | Single-tag elements | `b.Br().R()`, `b.Img("src", "x.png").R()` |

This pattern prevents unclosed tags and provides clear semantics.

---

## Debug System

The debug system (`element_debug.go`) provides runtime issue detection:

**Tracked Issues:**
- Unclosed tags
- Unwrapped text content
- Children on single-tag elements
- Malformed attribute pairs

**Debug UI Features:**
- Tabbed HTML/Markdown views
- Copy-to-clipboard for element locations
- Issue deduplication by location
- Terminal output for CI integration

---

*Document last updated: March 2026*
