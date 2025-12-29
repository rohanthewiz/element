# Component Asset Management in Element

## Problem Statement

Components in Element often need associated CSS and JavaScript to function properly. Currently:
- CSS is either embedded inline in Go code or manually managed in separate files
- No standard way for components to declare their asset dependencies
- No automatic asset deduplication (same CSS/JS might be included multiple times)
- No clear pattern for organizing component assets
- Manual burden on developers to remember which components need which assets

**Goals:**
1. Provide a simple, consistent way for components to declare CSS/JS dependencies
2. Automatically deduplicate assets across components
3. Embed all assets into the Go binary (no external files at runtime)
4. Avoid NodeJS tooling
5. Keep the API lightweight and maintain Element's philosophy
6. Support both component-scoped and shared assets

## Proposed Solution

### Architecture Overview

The solution consists of four complementary layers:

```
┌─────────────────────────────────────────────────────────┐
│  Layer 4: File Organization Convention (Developer UX)   │
│  - Colocate assets with components                      │
│  - Use go:embed for embedding                           │
└─────────────────────────────────────────────────────────┘
                          ↓
┌─────────────────────────────────────────────────────────┐
│  Layer 3: AssetProvider Interface (Component Level)     │
│  - Optional interface for components                    │
│  - Methods: CSS() []byte, JS() []byte                   │
└─────────────────────────────────────────────────────────┘
                          ↓
┌─────────────────────────────────────────────────────────┐
│  Layer 2: Builder Asset Collection (Runtime)            │
│  - Builder tracks used components                       │
│  - Collects and deduplicates assets                     │
│  - Methods: RenderCSS(), RenderJS()                     │
└─────────────────────────────────────────────────────────┘
                          ↓
┌─────────────────────────────────────────────────────────┐
│  Layer 1: Asset Registry (Optional, for shared assets)  │
│  - Global registry for commonly used assets             │
│  - Components register at init()                        │
│  - Useful for shared libraries (jQuery, etc.)           │
└─────────────────────────────────────────────────────────┘
```

---

## Implementation Details

### Layer 1: AssetProvider Interface

Add an optional interface that components can implement to provide assets:

```go
// AssetProvider is an optional interface that components can implement
// to declare their CSS and JavaScript dependencies.
type AssetProvider interface {
    // AssetID returns a unique identifier for deduplication.
    // Typically the component type name (e.g., "Alert", "Table").
    AssetID() string

    // CSS returns the component's CSS as bytes. Return nil if no CSS.
    CSS() []byte

    // JS returns the component's JavaScript as bytes. Return nil if no JS.
    JS() []byte
}
```

**Usage Example:**

```go
package components

import _ "embed"

//go:embed alert.css
var alertCSS []byte

//go:embed alert.js
var alertJS []byte

type Alert struct {
    Type    AlertType
    Message string
}

// AssetID implements AssetProvider
func (a Alert) AssetID() string {
    return "Alert"
}

// CSS implements AssetProvider
func (a Alert) CSS() []byte {
    return alertCSS
}

// JS implements AssetProvider
func (a Alert) JS() []byte {
    return alertJS
}

// Render implements Component
func (a Alert) Render(b *element.Builder) any {
    // Register this component's assets with the builder
    b.RegisterAssets(a)

    // ... render implementation
    return nil
}
```

### Layer 2: Builder Asset Management

Extend the Builder to track and manage component assets:

```go
// In builder.go

type Builder struct {
    buf         *bytes.Buffer
    // ... existing fields

    // Asset tracking
    assetRegistry map[string]assetEntry  // Key: AssetID, Value: asset data
    assetsEnabled bool                   // Whether to track assets
}

type assetEntry struct {
    id  string
    css []byte
    js  []byte
}

// EnableAssetTracking enables automatic asset collection.
// Call this when creating a builder if you want to use component assets.
func (b *Builder) EnableAssetTracking() *Builder {
    b.assetsEnabled = true
    if b.assetRegistry == nil {
        b.assetRegistry = make(map[string]assetEntry)
    }
    return b
}

// RegisterAssets registers a component's assets if it implements AssetProvider.
// This is called automatically when rendering components.
func (b *Builder) RegisterAssets(comp Component) {
    if !b.assetsEnabled {
        return
    }

    provider, ok := comp.(AssetProvider)
    if !ok {
        return // Component doesn't provide assets
    }

    id := provider.AssetID()
    if id == "" {
        return
    }

    // Check if already registered (deduplication)
    if _, exists := b.assetRegistry[id]; exists {
        return
    }

    // Register the asset
    b.assetRegistry[id] = assetEntry{
        id:  id,
        css: provider.CSS(),
        js:  provider.JS(),
    }
}

// RenderCSS renders all collected CSS in a <style> tag.
// Typically called in the <head> section.
func (b *Builder) RenderCSS() *Builder {
    if !b.assetsEnabled || len(b.assetRegistry) == 0 {
        return b
    }

    b.Style().R(func() any {
        for _, entry := range b.assetRegistry {
            if len(entry.css) > 0 {
                b.WriteBytes(entry.css)
                b.WriteByte('\n')
            }
        }
        return nil
    })

    return b
}

// RenderJS renders all collected JavaScript in a <script> tag.
// Typically called before </body>.
func (b *Builder) RenderJS() *Builder {
    if !b.assetsEnabled || len(b.assetRegistry) == 0 {
        return b
    }

    b.Script().R(func() any {
        for _, entry := range b.assetRegistry {
            if len(entry.js) > 0 {
                b.WriteBytes(entry.js)
                b.WriteByte('\n')
            }
        }
        return nil
    })

    return b
}

// GetCollectedCSS returns all CSS as a single byte slice.
// Useful for external stylesheets.
func (b *Builder) GetCollectedCSS() []byte {
    if !b.assetsEnabled {
        return nil
    }

    var buf bytes.Buffer
    for _, entry := range b.assetRegistry {
        if len(entry.css) > 0 {
            buf.Write(entry.css)
            buf.WriteByte('\n')
        }
    }
    return buf.Bytes()
}

// GetCollectedJS returns all JavaScript as a single byte slice.
// Useful for external script files.
func (b *Builder) GetCollectedJS() []byte {
    if !b.assetsEnabled {
        return nil
    }

    var buf bytes.Buffer
    for _, entry := range b.assetRegistry {
        if len(entry.js) > 0 {
            buf.Write(entry.js)
            buf.WriteByte('\n')
        }
    }
    return buf.Bytes()
}
```

### Layer 3: Asset Registry (Global, for shared assets)

For assets shared across many components (like normalize.css, utility libraries, etc.):

```go
// In asset_registry.go (new file)

package element

import (
    "sync"
)

// SharedAsset represents a reusable asset that can be shared across components
type SharedAsset struct {
    ID  string
    CSS []byte
    JS  []byte
}

var (
    sharedAssets = make(map[string]SharedAsset)
    assetMutex   sync.RWMutex
)

// RegisterSharedAsset registers a global asset that can be used by multiple components.
// This should typically be called in init() functions.
func RegisterSharedAsset(asset SharedAsset) {
    assetMutex.Lock()
    defer assetMutex.Unlock()

    if asset.ID == "" {
        return
    }
    sharedAssets[asset.ID] = asset
}

// GetSharedAsset retrieves a previously registered shared asset.
func GetSharedAsset(id string) (SharedAsset, bool) {
    assetMutex.RLock()
    defer assetMutex.RUnlock()

    asset, ok := sharedAssets[id]
    return asset, ok
}

// SharedAssetComponent is a helper type that provides shared assets.
// Use this when you want to include a shared asset in a page.
type SharedAssetComponent struct {
    AssetIDs []string
}

func (s SharedAssetComponent) AssetID() string {
    // Return a composite ID
    return "shared:" + strings.Join(s.AssetIDs, ",")
}

func (s SharedAssetComponent) CSS() []byte {
    var buf bytes.Buffer
    for _, id := range s.AssetIDs {
        if asset, ok := GetSharedAsset(id); ok && len(asset.CSS) > 0 {
            buf.Write(asset.CSS)
            buf.WriteByte('\n')
        }
    }
    return buf.Bytes()
}

func (s SharedAssetComponent) JS() []byte {
    var buf bytes.Buffer
    for _, id := range s.AssetIDs {
        if asset, ok := GetSharedAsset(id); ok && len(asset.JS) > 0 {
            buf.Write(asset.JS)
            buf.WriteByte('\n')
        }
    }
    return buf.Bytes()
}

func (s SharedAssetComponent) Render(b *Builder) any {
    b.RegisterAssets(s)
    return nil
}
```

**Usage Example:**

```go
package components

import (
    _ "embed"
    "github.com/rohanthewiz/element"
)

//go:embed normalize.css
var normalizeCSS []byte

//go:embed utility.css
var utilityCSS []byte

func init() {
    // Register commonly used assets
    element.RegisterSharedAsset(element.SharedAsset{
        ID:  "normalize",
        CSS: normalizeCSS,
    })

    element.RegisterSharedAsset(element.SharedAsset{
        ID:  "utilities",
        CSS: utilityCSS,
    })
}

// In your page rendering:
func RenderPage() string {
    b := element.NewBuilder().EnableAssetTracking()

    b.Html().R(
        b.Head().R(
            // Include shared assets
            element.SharedAssetComponent{
                AssetIDs: []string{"normalize", "utilities"},
            }.Render(b),

            // Render all collected CSS
            b.RenderCSS(),
        ),
        b.Body().R(
            // Your components
            Alert{Message: "Hello"}.Render(b),
            Table{...}.Render(b),

            // Render all collected JS
            b.RenderJS(),
        ),
    )

    return b.String()
}
```

### Layer 4: File Organization Convention

Recommended directory structure:

```
myproject/
├── components/
│   ├── alert/
│   │   ├── alert.go           # Component implementation
│   │   ├── alert.css          # Component-specific CSS
│   │   └── alert.js           # Component-specific JavaScript
│   ├── table/
│   │   ├── table.go
│   │   └── table.css          # No JS needed
│   ├── modal/
│   │   ├── modal.go
│   │   ├── modal.css
│   │   └── modal.js
│   └── shared/
│       ├── normalize.css      # Shared CSS
│       ├── utilities.css
│       └── shared.go          # Registers shared assets
└── main.go
```

**Alternative (simpler) for small projects:**

```
components/
├── components.go
├── alert.css
├── alert.js
├── table.css
├── modal.css
├── modal.js
└── shared.css
```

---

## Usage Patterns

### Pattern 1: Simple Page with Components

```go
func RenderPage(w http.ResponseWriter, r *http.Request) {
    b := element.NewBuilder().EnableAssetTracking()

    b.Doctype()
    b.Html().R(
        b.Head().R(
            b.Title().T("My Page"),
            b.RenderCSS(), // Automatically includes CSS from all components used below
        ),
        b.Body().R(
            Nav{Brand: "MyApp", Items: navItems}.Render(b),
            Alert{Type: AlertSuccess, Message: "Welcome!"}.Render(b),
            Card{Title: "Main Content", Body: "..."}.Render(b),

            b.RenderJS(), // Automatically includes JS from all components
        ),
    )

    w.Write(b.Bytes())
}
```

### Pattern 2: External Stylesheet/Script Files

For better caching, you can serve assets as separate files:

```go
// Handler for CSS file
func HandleCSS(w http.ResponseWriter, r *http.Request) {
    b := element.NewBuilder().EnableAssetTracking()

    // Render all components to collect their assets
    Alert{}.Render(b)
    Table{}.Render(b)
    Card{}.Render(b)
    // ... all components you use

    css := b.GetCollectedCSS()

    w.Header().Set("Content-Type", "text/css")
    w.Header().Set("Cache-Control", "public, max-age=31536000") // Cache for 1 year
    w.Write(css)
}

// In your HTML:
b.Head().R(
    b.Link("rel", "stylesheet", "href", "/static/components.css"),
)
```

### Pattern 3: Mixing Component and Page-Level Assets

```go
// Page-specific styles
//go:embed page_specific.css
var pageCSS []byte

type PageWrapper struct {
    Content element.Component
}

func (p PageWrapper) AssetID() string {
    return "PageWrapper"
}

func (p PageWrapper) CSS() []byte {
    return pageCSS
}

func (p PageWrapper) JS() []byte {
    return nil
}

func (p PageWrapper) Render(b *element.Builder) any {
    b.RegisterAssets(p)

    b.DivClass("page-wrapper").R(
        p.Content.Render(b),
    )
    return nil
}
```

---

## Advanced Features

### Asset Minification (Optional)

For production, you can minify assets at build time without NodeJS:

```go
// build.go (run with go:generate)

package main

import (
    "os"
    "path/filepath"

    // Go-based CSS minifier
    "github.com/tdewolff/minify/v2"
    "github.com/tdewolff/minify/v2/css"
    "github.com/tdewolff/minify/v2/js"
)

//go:generate go run build.go

func main() {
    m := minify.New()
    m.AddFunc("text/css", css.Minify)
    m.AddFunc("application/javascript", js.Minify)

    // Minify all CSS files
    filepath.Walk("./components", func(path string, info os.FileInfo, err error) error {
        if filepath.Ext(path) == ".css" {
            // Minify and overwrite
            // ... implementation
        }
        return nil
    })
}
```

### Content-Based Cache Busting

```go
import "crypto/sha256"

func (b *Builder) GetCSSHash() string {
    css := b.GetCollectedCSS()
    hash := sha256.Sum256(css)
    return fmt.Sprintf("%x", hash[:8]) // First 8 bytes as hex
}

// Usage:
func RenderPage(w http.ResponseWriter) {
    b := element.NewBuilder().EnableAssetTracking()

    // ... render components to collect assets

    cssHash := b.GetCSSHash()

    b.Head().R(
        b.Link("rel", "stylesheet", "href", "/static/components."+cssHash+".css"),
    )
}
```

### Conditional Asset Loading

```go
// Only load JS if the component actually has interactive features
func (b *Builder) RenderJSIfNeeded() *Builder {
    js := b.GetCollectedJS()
    if len(js) == 0 {
        return b // No JS to render
    }

    return b.RenderJS()
}
```

---

## Migration Guide

### For Existing Projects

**Step 1:** Enable asset tracking in your builder:

```go
// Before:
b := element.NewBuilder()

// After:
b := element.NewBuilder().EnableAssetTracking()
```

**Step 2:** Move inline CSS to separate files:

```go
// Before (in components.go):
const alertCSS = `
.alert {
    padding: 1rem;
    /* ... */
}
`

// After: Create alert.css file, then embed it:
//go:embed alert.css
var alertCSS []byte
```

**Step 3:** Implement AssetProvider for components that need assets:

```go
func (a Alert) AssetID() string { return "Alert" }
func (a Alert) CSS() []byte { return alertCSS }
func (a Alert) JS() []byte { return alertJS }
```

**Step 4:** Call RegisterAssets in Render:

```go
func (a Alert) Render(b *element.Builder) any {
    b.RegisterAssets(a) // Add this line

    // ... existing render code
    return nil
}
```

**Step 5:** Use RenderCSS() and RenderJS() in your page:

```go
b.Head().R(
    b.RenderCSS(), // Replaces manual style tag
)
b.Body().R(
    // ... components
    b.RenderJS(), // Replaces manual script tag
)
```

---

## Performance Considerations

### Memory Usage

- Asset deduplication prevents duplicate CSS/JS in memory
- Using `[]byte` instead of `string` for assets avoids extra allocations
- Builder's asset registry is a simple map (O(1) lookups)

### Build Time

- `//go:embed` happens at compile time (zero runtime cost)
- No reflection or runtime scanning
- Asset registration is simple map insertion

### Runtime Performance

- Asset collection: O(n) where n = number of unique components
- Deduplication: O(1) map lookup per component
- Rendering: Single pass through collected assets

### Benchmarks (Expected)

```
BenchmarkRenderWithAssets-8     100000    10234 ns/op    4096 B/op    12 allocs/op
BenchmarkRenderNoAssets-8       150000     6891 ns/op    2048 B/op     8 allocs/op
```

Overhead: ~50% time, ~100% memory (for asset tracking)
**Conclusion:** Use asset tracking when you need it, disable for simple pages.

---

## Tools and Ecosystem

### Without NodeJS

**CSS Processing:**
- [tdewolff/minify](https://github.com/tdewolff/minify) - Pure Go minifier
- [evanw/esbuild](https://esbuild.github.io/) - Fast bundler (written in Go)
- Manual CSS authoring (recommended for simplicity)

**JavaScript Processing:**
- [evanw/esbuild](https://esbuild.github.io/) - Can bundle and minify JS
- [tdewolff/minify](https://github.com/tdewolff/minify) - Can minify JS
- Vanilla JS (recommended to avoid tooling)

**Asset Verification:**
- `go:generate` scripts to validate CSS/JS syntax
- Standard Go tooling (go build, go test)

---

## Best Practices

### 1. Keep Components Self-Contained

```go
// Good: Component declares its own dependencies
type Alert struct { ... }
func (a Alert) CSS() []byte { return alertCSS }

// Bad: Expecting global CSS to be present
type Alert struct { ... }
// (relies on developer remembering to include alert.css separately)
```

### 2. Use Shared Assets for Common Dependencies

```go
// Good: Register normalize.css once, use everywhere
func init() {
    element.RegisterSharedAsset(element.SharedAsset{
        ID: "normalize",
        CSS: normalizeCSS,
    })
}

// Bad: Every component embeds normalize.css
```

### 3. Minimize JavaScript

```go
// Good: Pure CSS solution
type Accordion struct { ... }
func (a Accordion) CSS() []byte {
    // Use CSS :target or :checked pseudo-classes
}

// Consider: Only add JS if CSS can't do it
type Modal struct { ... }
func (m Modal) JS() []byte {
    // JavaScript for focus trap, ESC key handling
}
```

### 4. Version Your Assets

```go
const AlertVersion = "1.2.0"

func (a Alert) AssetID() string {
    return "Alert@" + AlertVersion
}
```

### 5. Document Asset Dependencies

```go
// Alert renders a styled notification message.
//
// Assets:
//   - CSS: Basic alert styling with color variants
//   - JS: Dismissible functionality (optional)
//
// Dependencies:
//   - Shared asset: "utilities" (for .sr-only class)
type Alert struct { ... }
```

---

## FAQ

**Q: Do I have to use asset tracking?**

A: No! It's completely optional. Continue using inline CSS or manual asset management if you prefer.

**Q: What if my component doesn't need assets?**

A: Simply don't implement the AssetProvider interface. The component will work normally.

**Q: Can I mix asset-providing components with regular components?**

A: Yes! The builder only collects assets from components that implement AssetProvider.

**Q: How do I handle CSS conflicts between components?**

A: Use BEM-style naming (`.alert__title`) or CSS modules pattern. Consider adding a namespace:

```css
/* alert.css */
.ele-alert { ... }
.ele-alert__title { ... }
```

**Q: Can I use this with HTMX or Alpine.js?**

A: Absolutely! Just register them as shared assets:

```go
//go:embed htmx.min.js
var htmxJS []byte

func init() {
    element.RegisterSharedAsset(element.SharedAsset{
        ID: "htmx",
        JS: htmxJS,
    })
}
```

**Q: What about TypeScript or SASS?**

A: Compile them to CSS/JS first, then embed the compiled output. You can use `go:generate` to automate this:

```go
//go:generate sass alert.scss alert.css
//go:generate esbuild alert.ts --outfile=alert.js

//go:embed alert.css
var alertCSS []byte

//go:embed alert.js
var alertJS []byte
```

**Q: How do I serve assets from a CDN?**

A: Don't use RenderCSS()/RenderJS(). Instead, write assets to files during build:

```go
// build.go
func main() {
    b := element.NewBuilder().EnableAssetTracking()
    // ... register all components

    os.WriteFile("dist/components.css", b.GetCollectedCSS(), 0644)
    os.WriteFile("dist/components.js", b.GetCollectedJS(), 0644)
}
```

Then upload `dist/` to your CDN and reference it:

```go
b.Link("rel", "stylesheet", "href", "https://cdn.example.com/components.css")
```

---

## Summary

This asset management system:

✅ **Simple** - Optional interfaces, minimal API surface
✅ **Efficient** - Automatic deduplication, embedded in binary
✅ **Flexible** - Works with any workflow (inline, external files, CDN)
✅ **No NodeJS** - Pure Go, standard library + optional minifiers
✅ **Performant** - Single-pass collection, O(1) deduplication
✅ **Backward Compatible** - Existing code works without changes
✅ **Element Philosophy** - Lightweight, no reflection, no magic

The system scales from simple pages (no asset tracking) to complex applications (full asset management) without forcing developers into a specific workflow.
