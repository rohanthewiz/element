# Element Library Implementation Proposals

> Detailed implementation proposals for high-value architectural improvements.

This document expands on the improvement opportunities identified in [ARCHITECTURE.md](../ARCHITECTURE.md), providing implementation details, code examples, and considerations for three key enhancements.

---

## Table of Contents

1. [Builder Pool for High-Throughput Scenarios](#1-builder-pool-for-high-throughput-scenarios)
2. [Automatic Attribute Escaping](#2-automatic-attribute-escaping)
3. [Streaming/Chunked Rendering](#3-streamingchunked-rendering)

---

## 1. Builder Pool for High-Throughput Scenarios

**Priority:** P1 (Quick Win)
**Effort:** Low
**Risk:** Low

### Problem Statement

Every HTTP request that renders HTML creates a new `Builder`, which internally allocates a fresh `bytes.Buffer`. In high-traffic applications, this creates allocation pressure:

```go
func handler(w http.ResponseWriter, r *http.Request) {
    b := element.NewBuilder()  // New allocation every request
    // ... render page ...
    w.Write([]byte(b.String()))
    // Builder becomes garbage after request
}
```

With thousands of requests per second, the garbage collector works harder to reclaim these short-lived allocations, leading to:
- Increased GC pause times
- Higher tail latencies (p99, p999)
- Unnecessary CPU cycles spent on allocation/deallocation

### Proposed Solution

Use Go's `sync.Pool` to reuse `Builder` instances across requests.

#### New API

```go
// pool.go (new file) or add to builder.go

package element

import "sync"

var builderPool = sync.Pool{
    New: func() any {
        return NewBuilder()
    },
}

// AcquireBuilder gets a Builder from the pool, or creates a new one if the pool is empty.
// The returned Builder is reset and ready for use.
func AcquireBuilder() *Builder {
    return builderPool.Get().(*Builder)
}

// ReleaseBuilder returns a Builder to the pool for reuse.
// After calling ReleaseBuilder, the Builder must not be used again.
func ReleaseBuilder(b *Builder) {
    b.Reset()
    builderPool.Put(b)
}
```

#### Reset Method Enhancement

Ensure the `Reset()` method clears all state:

```go
// In builder.go

// Reset clears the builder's internal buffer, preparing it for reuse.
// After Reset, the builder can be used to generate new HTML.
func (b *Builder) Reset() {
    b.s.Reset()
    // Note: Ele and Text functions reference b.s, so they remain valid
}
```

#### Usage Pattern

```go
func handler(w http.ResponseWriter, r *http.Request) {
    b := element.AcquireBuilder()
    defer element.ReleaseBuilder(b)

    b.Html().R(
        b.Head().R(
            b.Title().T("My Page"),
        ),
        b.Body().R(
            b.H1().T("Hello, World!"),
        ),
    )

    w.Write([]byte(b.String()))
}
```

### Implementation Details

#### Debug Mode Interaction

Pooled builders should not carry debug state between requests:

```go
func ReleaseBuilder(b *Builder) {
    b.Reset()

    // In debug mode, we may want to avoid pooling entirely
    // to prevent cross-request debug state confusion
    if IsDebugMode() {
        return  // Don't return to pool in debug mode
    }

    builderPool.Put(b)
}
```

Alternatively, ensure debug state is cleared:

```go
func (b *Builder) Reset() {
    b.s.Reset()
    // If builder tracks any debug state, clear it here
}
```

#### Benchmark Suite

Add benchmarks to measure the improvement:

```go
// builder_bench_test.go

package element

import "testing"

func BenchmarkBuilderAllocation(b *testing.B) {
    b.ReportAllocs()
    for i := 0; i < b.N; i++ {
        builder := NewBuilder()
        builder.Div("class", "container").R(
            builder.P().T("Hello, World!"),
            builder.Ul().R(
                builder.Li().T("Item 1"),
                builder.Li().T("Item 2"),
                builder.Li().T("Item 3"),
            ),
        )
        _ = builder.String()
    }
}

func BenchmarkBuilderPooled(b *testing.B) {
    b.ReportAllocs()
    for i := 0; i < b.N; i++ {
        builder := AcquireBuilder()
        builder.Div("class", "container").R(
            builder.P().T("Hello, World!"),
            builder.Ul().R(
                builder.Li().T("Item 1"),
                builder.Li().T("Item 2"),
                builder.Li().T("Item 3"),
            ),
        )
        _ = builder.String()
        ReleaseBuilder(builder)
    }
}

func BenchmarkBuilderPooledParallel(b *testing.B) {
    b.ReportAllocs()
    b.RunParallel(func(pb *testing.PB) {
        for pb.Next() {
            builder := AcquireBuilder()
            builder.Div("class", "container").R(
                builder.P().T("Hello, World!"),
            )
            _ = builder.String()
            ReleaseBuilder(builder)
        }
    })
}
```

### Testing Strategy

```go
// pool_test.go

package element

import (
    "sync"
    "testing"
)

func TestAcquireReleaseBuilder(t *testing.T) {
    b := AcquireBuilder()
    if b == nil {
        t.Fatal("AcquireBuilder returned nil")
    }

    b.Div().T("test")
    result := b.String()
    if result != "<div>test</div>" {
        t.Errorf("unexpected output: %s", result)
    }

    ReleaseBuilder(b)
}

func TestPooledBuilderIsReset(t *testing.T) {
    // Get a builder and use it
    b1 := AcquireBuilder()
    b1.Div().T("first")
    ReleaseBuilder(b1)

    // Get another builder (likely the same one from pool)
    b2 := AcquireBuilder()
    b2.Span().T("second")
    result := b2.String()

    // Should NOT contain content from previous use
    if result != "<span>second</span>" {
        t.Errorf("builder not properly reset, got: %s", result)
    }
    ReleaseBuilder(b2)
}

func TestPoolConcurrentAccess(t *testing.T) {
    var wg sync.WaitGroup
    iterations := 1000

    for i := 0; i < iterations; i++ {
        wg.Add(1)
        go func(n int) {
            defer wg.Done()
            b := AcquireBuilder()
            b.Div().F("Item %d", n)
            _ = b.String()
            ReleaseBuilder(b)
        }(i)
    }

    wg.Wait()
}
```

### Expected Impact

| Metric | Before | After (Expected) |
|--------|--------|------------------|
| Allocations per request | 2-3 | 0-1 |
| GC pressure | Higher | Lower |
| P99 latency | Variable | More stable |
| Memory churn | High | Low |

### Migration Guide

The pooled API is opt-in and backwards compatible:

```go
// Existing code continues to work
b := element.NewBuilder()

// New high-performance pattern
b := element.AcquireBuilder()
defer element.ReleaseBuilder(b)
```

---

## 2. Automatic Attribute Escaping

**Priority:** P1 (Security Critical)
**Effort:** Low-Medium
**Risk:** Medium (Breaking Change)

### Problem Statement

Currently, attribute values are written directly without HTML escaping:

```go
// User input could contain malicious content
userInput := `" onclick="alert('xss')`
b.Div("data-name", userInput).R()

// Output: <div data-name="" onclick="alert('xss')"></div>
// This is an XSS vulnerability!
```

This puts the burden on developers to remember to escape user input, which is error-prone and has led to countless security vulnerabilities across the industry.

### Proposed Solution

Automatically escape attribute values by default, with an explicit opt-out for trusted content.

#### New Types and Functions

```go
// escape.go (new file)

package element

import (
    "html"
    "strings"
)

// Raw represents a string that should not be HTML-escaped.
// Use this for trusted content like inline JavaScript handlers.
// Security: Only use Raw with content you fully control, never with user input.
type Raw string

// escapeAttrValue escapes a string for safe use in HTML attributes.
// Returns the original string unchanged if no escaping is needed (fast path).
func escapeAttrValue(s string) string {
    // Fast path: if no special characters, return as-is
    if !strings.ContainsAny(s, `"'&<>`) {
        return s
    }
    return html.EscapeString(s)
}

// escapeText escapes a string for safe use in HTML text content.
func escapeText(s string) string {
    if !strings.ContainsAny(s, `&<>`) {
        return s
    }
    return html.EscapeString(s)
}
```

#### Modified Element Creation

Update the `New()` function in `element.go`:

```go
// In element.go - modify attribute handling

func New(s *bytes.Buffer, el string, attrs ...any) Element {
    // ... existing setup code ...

    s.WriteString("<")
    s.WriteString(el)

    // Process attributes in pairs
    for i := 0; i < len(attrs); i += 2 {
        key := fmt.Sprintf("%v", attrs[i])

        if i+1 < len(attrs) {
            s.WriteString(` `)
            s.WriteString(key)
            s.WriteString(`="`)

            // Check if value is Raw (trusted) or needs escaping
            switch v := attrs[i+1].(type) {
            case Raw:
                s.WriteString(string(v))  // No escaping for Raw
            case string:
                s.WriteString(escapeAttrValue(v))  // Escape strings
            default:
                s.WriteString(escapeAttrValue(fmt.Sprintf("%v", v)))
            }

            s.WriteString(`"`)
        }
    }

    // ... rest of function ...
}
```

#### Modified Text Methods

Update the `T()` method to escape text content:

```go
// In element.go

func (e Element) T(texts ...string) any {
    for _, text := range texts {
        e.sb.WriteString(escapeText(text))
    }
    // ... closing tag logic ...
}
```

### API Design

#### Basic Usage (Auto-Escaped)

```go
// User input is automatically escaped
userInput := `<script>alert("xss")</script>`
b.P().T(userInput)
// Output: <p>&lt;script&gt;alert("xss")&lt;/script&gt;</p>

// Attribute values are also escaped
userName := `"; onclick="steal()"`
b.Div("data-user", userName).R()
// Output: <div data-user="&#34;; onclick=&#34;steal()&#34;"></div>
```

#### Trusted Content (Raw)

```go
// For trusted JavaScript handlers
b.Button("onclick", element.Raw("handleClick(event)")).T("Click Me")
// Output: <button onclick="handleClick(event)">Click Me</button>

// For pre-escaped HTML content
trustedHTML := fetchSanitizedHTML()
b.Div().R(element.Raw(trustedHTML))
```

#### URL Attribute Safety

Add additional protection for URL attributes:

```go
// In escape.go

// sanitizeURL checks for dangerous URL schemes
func sanitizeURL(s string) string {
    lower := strings.ToLower(strings.TrimSpace(s))

    // Block javascript: and data: URLs
    if strings.HasPrefix(lower, "javascript:") ||
       strings.HasPrefix(lower, "data:text/html") ||
       strings.HasPrefix(lower, "vbscript:") {
        return "#blocked"  // Safe fallback
    }

    return escapeAttrValue(s)
}

// isURLAttribute returns true for attributes that contain URLs
func isURLAttribute(key string) bool {
    switch strings.ToLower(key) {
    case "href", "src", "action", "formaction", "poster", "data":
        return true
    }
    return false
}
```

### Character Escaping Reference

| Character | Escaped Form | Context |
|-----------|--------------|---------|
| `"` | `&#34;` or `&quot;` | Attributes |
| `'` | `&#39;` | Attributes |
| `&` | `&amp;` | Both |
| `<` | `&lt;` | Both |
| `>` | `&gt;` | Both |

### Performance Considerations

#### Benchmark the Fast Path

```go
func BenchmarkEscapeNoSpecialChars(b *testing.B) {
    s := "hello world class-name"  // No special chars
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _ = escapeAttrValue(s)
    }
}

func BenchmarkEscapeWithSpecialChars(b *testing.B) {
    s := `hello "world" & <friends>`
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _ = escapeAttrValue(s)
    }
}
```

#### Optimization: Avoid Allocation When Possible

```go
func escapeAttrValue(s string) string {
    // Fast path: check if escaping is needed
    needsEscape := false
    for i := 0; i < len(s); i++ {
        switch s[i] {
        case '"', '\'', '&', '<', '>':
            needsEscape = true
            break
        }
    }

    if !needsEscape {
        return s  // Return original string, no allocation
    }

    return html.EscapeString(s)
}
```

### Testing Strategy

```go
// escape_test.go

package element

import "testing"

func TestAttributeEscaping(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        expected string
    }{
        {"no escape needed", "hello", "hello"},
        {"double quote", `say "hello"`, `say &#34;hello&#34;`},
        {"single quote", "it's", "it&#39;s"},
        {"ampersand", "a & b", "a &amp; b"},
        {"less than", "a < b", "a &lt; b"},
        {"greater than", "a > b", "a &gt; b"},
        {"xss attempt", `" onclick="alert(1)`, `&#34; onclick=&#34;alert(1)`},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := escapeAttrValue(tt.input)
            if result != tt.expected {
                t.Errorf("escapeAttrValue(%q) = %q, want %q",
                    tt.input, result, tt.expected)
            }
        })
    }
}

func TestRawBypassesEscaping(t *testing.T) {
    b := NewBuilder()
    b.Button("onclick", Raw("handleClick()")).T("Click")

    result := b.String()
    expected := `<button onclick="handleClick()">Click</button>`

    if result != expected {
        t.Errorf("Raw not working: got %q, want %q", result, expected)
    }
}

func TestTextContentEscaping(t *testing.T) {
    b := NewBuilder()
    b.P().T("<script>alert('xss')</script>")

    result := b.String()
    if strings.Contains(result, "<script>") {
        t.Error("script tag was not escaped in text content")
    }
}

func TestURLSanitization(t *testing.T) {
    tests := []struct {
        input    string
        blocked  bool
    }{
        {"https://example.com", false},
        {"javascript:alert(1)", true},
        {"JAVASCRIPT:alert(1)", true},
        {"data:text/html,<script>", true},
        {"/relative/path", false},
    }

    for _, tt := range tests {
        result := sanitizeURL(tt.input)
        isBlocked := result == "#blocked"
        if isBlocked != tt.blocked {
            t.Errorf("sanitizeURL(%q) blocked=%v, want blocked=%v",
                tt.input, isBlocked, tt.blocked)
        }
    }
}
```

### Migration Strategy

This is a breaking change. Recommended approach:

#### Phase 1: Opt-In (v1.x)

```go
// Add builder option for escaping
func NewBuilder(opts ...Option) *Builder {
    b := &Builder{...}
    for _, opt := range opts {
        opt(b)
    }
    return b
}

func WithAutoEscape(enabled bool) Option {
    return func(b *Builder) {
        b.autoEscape = enabled
    }
}

// Usage
b := element.NewBuilder(element.WithAutoEscape(true))
```

#### Phase 2: Default On (v2.0)

```go
// Escaping enabled by default
// Raw type available for opt-out
```

### Security Audit Checklist

After implementation, verify:

- [ ] All attribute values are escaped by default
- [ ] Text content in `T()` is escaped
- [ ] Text content in `F()` is escaped
- [ ] `Raw` type bypasses escaping correctly
- [ ] URL attributes block `javascript:` schemes
- [ ] Nested quotes are properly escaped
- [ ] Unicode characters are handled correctly

---

## 3. Streaming/Chunked Rendering

**Priority:** P3 (Strategic)
**Effort:** High
**Risk:** High (Architectural Change)

### Problem Statement

Currently, the entire HTML document must be generated in memory before any bytes can be sent to the client:

```go
func handler(w http.ResponseWriter, r *http.Request) {
    b := element.NewBuilder()

    // Entire page built in memory
    b.Html().R(
        b.Head().R(/* styles, scripts */),
        b.Body().R(
            renderHeader(b),
            renderLargeDataTable(b),  // Could be megabytes
            renderFooter(b),
        ),
    )

    // Only NOW can we send anything
    w.Write([]byte(b.String()))
}
```

This causes:
- **High TTFB** (Time To First Byte) - User sees blank screen until entire page is ready
- **High memory usage** - Entire page held in `bytes.Buffer`
- **Poor perceived performance** - No progressive rendering

### Proposed Solution

Support writing directly to an `io.Writer` to enable chunked transfer encoding and progressive rendering.

### Implementation Approach

Given the complexity, we propose a phased approach:

#### Phase 1: WriteTo Method (Simple)

Add ability to write buffered content to an `io.Writer`:

```go
// In builder.go

import "io"

// WriteTo writes the accumulated HTML to w.
// Implements io.WriterTo interface.
func (b *Builder) WriteTo(w io.Writer) (int64, error) {
    n, err := io.WriteString(w, b.String())
    return int64(n), err
}

// FlushTo writes accumulated HTML to w and resets the buffer.
// Useful for sending content in chunks.
func (b *Builder) FlushTo(w io.Writer) error {
    _, err := io.WriteString(w, b.String())
    b.Reset()
    return err
}
```

**Usage:**

```go
func handler(w http.ResponseWriter, r *http.Request) {
    b := element.NewBuilder()

    // Render and flush head immediately
    b.WriteString("<!DOCTYPE html>")
    b.Html().Open()
    b.Head().R(
        b.Meta("charset", "utf-8"),
        b.Link("rel", "stylesheet", "href", "/styles.css"),
        b.Script("src", "/app.js"),
    )
    b.FlushTo(w)  // Browser starts fetching CSS/JS NOW

    // Render body
    b.Body().R(
        renderContent(b),
    )
    b.Html().Close()
    b.FlushTo(w)
}
```

**New Element Methods Needed:**

```go
// Open writes the opening tag without requiring immediate children
func (e Element) Open() {
    // Opening tag already written in New()
    // This is just for semantic clarity
}

// Close writes the closing tag
func (e Element) Close() {
    if !e.IsSingleTag() {
        e.sb.WriteString("</")
        e.sb.WriteString(e.name)
        e.sb.WriteString(">")
    }
}
```

#### Phase 2: Streaming Builder (Advanced)

Create a new builder type that writes directly to an `io.Writer`:

```go
// streaming.go (new file)

package element

import (
    "fmt"
    "io"
    "sync"
)

// StreamBuilder writes HTML directly to an io.Writer.
// Unlike Builder, content is not buffered in memory.
type StreamBuilder struct {
    w      io.Writer
    err    error      // First error encountered
    mu     sync.Mutex // Protects concurrent writes
}

// NewStreamBuilder creates a builder that writes directly to w.
func NewStreamBuilder(w io.Writer) *StreamBuilder {
    return &StreamBuilder{w: w}
}

// Err returns the first error encountered during writing, if any.
func (sb *StreamBuilder) Err() error {
    return sb.err
}

// writeString writes s to the underlying writer.
// Stops writing after first error.
func (sb *StreamBuilder) writeString(s string) {
    if sb.err != nil {
        return
    }
    _, sb.err = io.WriteString(sb.w, s)
}

// Flush flushes the underlying writer if it supports flushing.
func (sb *StreamBuilder) Flush() error {
    if sb.err != nil {
        return sb.err
    }
    if f, ok := sb.w.(interface{ Flush() error }); ok {
        return f.Flush()
    }
    if f, ok := sb.w.(http.Flusher); ok {
        f.Flush()
    }
    return nil
}
```

**Streaming Element:**

```go
// StreamElement represents an element being written to a stream.
type StreamElement struct {
    sb   *StreamBuilder
    name string
}

// R renders children and writes the closing tag.
func (e *StreamElement) R(children ...any) error {
    // Children are already rendered (executed before R is called)
    if !isSingleTag(e.name) {
        e.sb.writeString("</")
        e.sb.writeString(e.name)
        e.sb.writeString(">")
    }
    return e.sb.err
}

// T writes text content and the closing tag.
func (e *StreamElement) T(texts ...string) error {
    for _, text := range texts {
        e.sb.writeString(escapeText(text))
    }
    if !isSingleTag(e.name) {
        e.sb.writeString("</")
        e.sb.writeString(e.name)
        e.sb.writeString(">")
    }
    return e.sb.err
}
```

**Element Methods:**

```go
// Div creates a div element and writes its opening tag.
func (sb *StreamBuilder) Div(attrs ...string) *StreamElement {
    sb.writeString("<div")
    writeAttrs(sb, attrs)
    sb.writeString(">")
    return &StreamElement{sb: sb, name: "div"}
}

// Helper to write attributes
func writeAttrs(sb *StreamBuilder, attrs []string) {
    for i := 0; i < len(attrs); i += 2 {
        if i+1 < len(attrs) {
            sb.writeString(` `)
            sb.writeString(attrs[i])
            sb.writeString(`="`)
            sb.writeString(escapeAttrValue(attrs[i+1]))
            sb.writeString(`"`)
        }
    }
}
```

**Usage:**

```go
func handler(w http.ResponseWriter, r *http.Request) {
    // Enable chunked encoding
    w.Header().Set("Transfer-Encoding", "chunked")

    sb := element.NewStreamBuilder(w)

    sb.Html().R(
        sb.Head().R(
            sb.Title().T("My Page"),
        ),
        sb.Body().R(
            sb.H1().T("Welcome"),  // Sent immediately
            sb.Div().R(
                renderLargeTable(sb),  // Streamed as generated
            ),
        ),
    )

    if err := sb.Err(); err != nil {
        log.Printf("render error: %v", err)
    }
}
```

#### Phase 3: Unified Interface (Future)

Create a common interface so components work with both builders:

```go
// html_writer.go

package element

// HTMLWriter is the interface for writing HTML content.
// Both Builder and StreamBuilder implement this interface.
type HTMLWriter interface {
    WriteString(s string) error
    WriteRaw(s string) error  // Write without escaping
}

// Component can render to any HTMLWriter.
type Component interface {
    Render(w HTMLWriter) error
}
```

### Architecture Considerations

#### Challenge: Return Value Pattern

Current API relies on return values for nesting:

```go
b.Div().R(
    b.P().T("hello"),  // P().T() returns value used by R()
)
```

With streaming, we might need callback pattern:

```go
sb.Div(func() {
    sb.P().T("hello")
})
```

**Proposed Solution:** Keep return values but make them no-ops for nesting:

```go
// StreamElement.R returns nil but provides nesting syntax
func (e *StreamElement) R(children ...any) any {
    // Children already executed and written
    e.sb.writeString("</" + e.name + ">")
    return nil
}
```

#### Challenge: Error Handling

Streaming can fail mid-render:

```go
type StreamBuilder struct {
    w   io.Writer
    err error  // Capture first error, stop writing
}

// All writes check for prior error
func (sb *StreamBuilder) writeString(s string) {
    if sb.err != nil {
        return  // Silently skip after first error
    }
    _, sb.err = io.WriteString(sb.w, s)
}
```

#### Challenge: Component Compatibility

Existing components need adaptation:

```go
// Current component
type MyComponent struct{}

func (c *MyComponent) Render(b *Builder) any {
    return b.Div().T("hello")
}

// Streaming-compatible component
type MyStreamComponent struct{}

func (c *MyStreamComponent) Render(w HTMLWriter) error {
    w.WriteString("<div>hello</div>")
    return nil
}
```

### Testing Strategy

```go
// streaming_test.go

package element

import (
    "bytes"
    "errors"
    "testing"
)

func TestStreamBuilder(t *testing.T) {
    var buf bytes.Buffer
    sb := NewStreamBuilder(&buf)

    sb.Div("class", "test").R(
        sb.P().T("Hello"),
    )

    if err := sb.Err(); err != nil {
        t.Fatalf("unexpected error: %v", err)
    }

    expected := `<div class="test"><p>Hello</p></div>`
    if buf.String() != expected {
        t.Errorf("got %q, want %q", buf.String(), expected)
    }
}

func TestStreamBuilderErrorHandling(t *testing.T) {
    // Writer that fails after N bytes
    fw := &failingWriter{failAfter: 10}
    sb := NewStreamBuilder(fw)

    sb.Div().R(
        sb.P().T("This is a long text that will cause failure"),
    )

    if sb.Err() == nil {
        t.Error("expected error from failing writer")
    }
}

type failingWriter struct {
    written   int
    failAfter int
}

func (fw *failingWriter) Write(p []byte) (int, error) {
    if fw.written >= fw.failAfter {
        return 0, errors.New("write failed")
    }
    fw.written += len(p)
    return len(p), nil
}

func TestStreamBuilderFlush(t *testing.T) {
    var buf bytes.Buffer
    sb := NewStreamBuilder(&buf)

    sb.Html().Open()
    sb.Head().R(
        sb.Title().T("Test"),
    )
    sb.Flush()

    // Should have partial content
    if !bytes.Contains(buf.Bytes(), []byte("<head>")) {
        t.Error("head not flushed")
    }
}
```

### Performance Comparison

| Aspect | Buffered Builder | Streaming Builder |
|--------|------------------|-------------------|
| Memory | O(page size) | O(1) |
| TTFB | After full render | Immediate |
| Throughput | Higher (no syscalls) | Lower (many syscalls) |
| Complexity | Simple | Complex |
| Error handling | Easy | Requires care |

### Recommendation

1. **Implement Phase 1 immediately** - `WriteTo()` and `FlushTo()` are simple additions with clear value
2. **Prototype Phase 2** - Build streaming builder as separate type, don't modify existing API
3. **Defer Phase 3** - Unified interface can wait until streaming patterns are proven

### Example: Progressive Page Load

```go
func handler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/html; charset=utf-8")

    b := element.NewBuilder()

    // Send document head immediately
    b.WriteString("<!DOCTYPE html>")
    b.Html().Open()
    b.Head().R(
        b.Meta("charset", "utf-8"),
        b.Title().T("Dashboard"),
        b.Link("rel", "stylesheet", "href", "/app.css"),
        b.Script("src", "/app.js", "defer", "true"),
    )
    b.FlushTo(w)  // TTFB: browser starts loading CSS/JS

    // Send page shell
    b.Body().Open()
    b.Header("class", "site-header").R(
        b.Nav().R(/* navigation */),
    )
    b.FlushTo(w)  // User sees header

    // Render main content (potentially slow)
    b.Main().R(
        renderDashboard(b, r.Context()),
    )
    b.FlushTo(w)  // User sees content

    // Send footer and close
    b.Footer().R(/* footer content */)
    b.Body().Close()
    b.Html().Close()
    b.FlushTo(w)
}
```

---

## Summary

| Proposal | Status | Next Step |
|----------|--------|-----------|
| Builder Pool | Ready to implement | Add benchmark, implement, test |
| Attribute Escaping | Design complete | Decide on migration strategy |
| Streaming Rendering | Phased approach | Implement Phase 1 first |

### Implementation Order

1. **Builder Pool** (1-2 hours) - Immediate performance win
2. **Attribute Escaping** (4-8 hours) - Critical for security
3. **Streaming Phase 1** (2-4 hours) - `FlushTo()` method
4. **Streaming Phase 2** (1-2 days) - Full streaming builder

---

*Document created: December 2024*
*For: Element HTML Generation Library*
