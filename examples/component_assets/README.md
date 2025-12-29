# Component Asset Management Example

This example demonstrates Element's component asset management system - a simple, efficient way to manage CSS and JavaScript for reusable components without NodeJS.

## Features Demonstrated

1. **AssetProvider Interface** - Components declare their CSS/JS dependencies
2. **Automatic Asset Collection** - Builder automatically collects assets from components
3. **Deduplication** - Same assets only included once, even if component used multiple times
4. **Embedded Assets** - CSS/JS embedded in Go binary using `//go:embed`
5. **Flexible Output** - Assets can be inline or served as external files
6. **No NodeJS** - Pure Go solution with zero external dependencies

## Running the Example

```bash
# Run the demo
go run main.go

# This will:
# 1. Print HTML with inline assets to stdout
# 2. Generate external asset files (components.css, components.js)
# 3. Generate an HTML file that uses external assets
# 4. Show asset inspection and deduplication results
```

## Project Structure

```
component_assets/
├── main.go                    # Example demonstrations
├── components/
│   ├── components.go          # Component implementations
│   ├── alert.css             # Alert component styles
│   ├── alert.js              # Alert component JavaScript
│   └── table.css             # Table component styles
└── README.md
```

## How It Works

### 1. Components Implement AssetProvider

```go
type Alert struct {
    Type        AlertType
    Message     string
    Dismissible bool
}

// Unique identifier for deduplication
func (a Alert) AssetID() string {
    return "Alert"
}

// CSS embedded at compile time
//go:embed alert.css
var alertCSS []byte

func (a Alert) CSS() []byte {
    return alertCSS
}

// Conditional JavaScript
func (a Alert) JS() []byte {
    if a.Dismissible {
        return alertJS
    }
    return nil
}
```

### 2. Components Register Assets in Render

```go
func (a Alert) Render(b *element.Builder) any {
    // Register assets with builder (handles deduplication)
    b.RegisterAssets(a)

    // ... normal rendering code
    return nil
}
```

### 3. Builder Collects and Renders Assets

```go
// Create builder with asset tracking
b := element.NewBuilder().EnableAssetTracking()

// Render components (assets collected automatically)
Alert{Message: "Hello"}.Render(b)
Table{Headers: []string{"A", "B"}}.Render(b)
Alert{Message: "World"}.Render(b)  // Alert assets NOT duplicated!

// Render collected assets
b.Head().R(
    b.RenderCSS(),  // All CSS from components
)
b.Body().R(
    // ... components ...
    b.RenderJS(),   // All JavaScript from components
)
```

## Usage Patterns

### Pattern 1: Inline Assets (Simple)

Best for: Simple apps, development, single-page apps

```go
b := element.NewBuilder().EnableAssetTracking()

b.Html().R(
    b.Head().R(
        b.RenderCSS(),  // Inline <style> tag
    ),
    b.Body().R(
        Alert{}.Render(b),
        b.RenderJS(),    // Inline <script> tag
    ),
)
```

**Pros:** Simple, single file, no HTTP requests
**Cons:** No browser caching, larger HTML

### Pattern 2: External Files (Production)

Best for: Production apps, better caching, CDN deployment

```go
// Step 1: Collect assets
collector := element.NewBuilder().EnableAssetTracking()
Alert{}.Render(collector)
Table{}.Render(collector)

// Step 2: Write to files
os.WriteFile("assets/components.css", collector.GetCollectedCSS(), 0644)
os.WriteFile("assets/components.js", collector.GetCollectedJS(), 0644)

// Step 3: Reference in HTML
b := element.NewBuilder()  // No asset tracking needed
b.Head().R(
    b.Link("rel", "stylesheet", "href", "/assets/components.css"),
)
```

**Pros:** Browser caching, parallel loading, CDN-friendly
**Cons:** Requires build step, multiple files

### Pattern 3: Hybrid (Best of Both)

Best for: Server-rendered apps with frequent updates

```go
// Critical CSS inline, rest external
b := element.NewBuilder().EnableAssetTracking()

b.Head().R(
    // Inline critical/small CSS
    b.RenderCSS(),

    // External for larger assets
    b.Link("rel", "stylesheet", "href", "/assets/vendor.css"),
)
```

## Advanced Features

### Conditional Asset Loading

```go
// Only include JavaScript if component actually needs it
func (a Alert) JS() []byte {
    if a.Dismissible {
        return alertJS  // Only load JS for dismissible alerts
    }
    return nil
}
```

### Asset Versioning

```go
func (a Alert) AssetID() string {
    return "Alert@1.2.0"  // Version in ID
}
```

### Shared Assets

For assets used by many components (like normalize.css):

```go
// In shared/shared.go
//go:embed normalize.css
var normalizeCSS []byte

func init() {
    element.RegisterSharedAsset(element.SharedAsset{
        ID:  "normalize",
        CSS: normalizeCSS,
    })
}

// Use in pages
SharedAssetComponent{
    AssetIDs: []string{"normalize", "utilities"},
}.Render(b)
```

## Performance

### Deduplication

- O(1) lookup via map
- Each asset included exactly once
- No duplicate CSS/JS in output

### Memory

- Assets embedded at compile time (zero runtime file I/O)
- `[]byte` used to avoid string allocations
- Registry cleared when builder is garbage collected

### Build Time

- `//go:embed` happens at compile time
- No runtime asset discovery or scanning
- No reflection

## Migration from Inline CSS

If you have existing components with inline CSS:

```go
// Before:
const alertCSS = `
.alert {
    padding: 1rem;
    /* ... */
}
`

func (a Alert) Render(b *element.Builder) any {
    b.Style().T(alertCSS)  // Repeated on every render!
    // ...
}
```

```go
// After:
//go:embed alert.css
var alertCSS []byte

func (a Alert) AssetID() string { return "Alert" }
func (a Alert) CSS() []byte { return alertCSS }
func (a Alert) JS() []byte { return nil }

func (a Alert) Render(b *element.Builder) any {
    b.RegisterAssets(a)  // Included once, even if rendered 100 times
    // ...
}
```

## Best Practices

### 1. Keep Components Self-Contained

✅ Good:
```go
// alert.go, alert.css, alert.js in same directory
// Component embeds its own dependencies
```

❌ Bad:
```go
// Component expects global CSS to be present
// Developer must remember to include styles separately
```

### 2. Use Descriptive AssetIDs

✅ Good:
```go
func (a Alert) AssetID() string { return "Alert@1.0.0" }
func (t Table) AssetID() string { return "Table@2.1.0" }
```

❌ Bad:
```go
func (a Alert) AssetID() string { return "comp1" }
```

### 3. Minimize JavaScript

✅ Good:
```go
// Use CSS for animations, transitions, show/hide
// Only use JS when CSS can't do it
```

### 4. Document Dependencies

```go
// Alert renders a notification with optional dismiss button.
//
// Assets:
//   - CSS: Alert styling (info, success, warning, error variants)
//   - JS: Dismissible functionality (only if Dismissible=true)
//
// Dependencies:
//   - None (fully self-contained)
type Alert struct { ... }
```

## FAQ

**Q: Do I have to use asset management?**
A: No! It's completely optional. Continue using inline styles if you prefer.

**Q: Can I mix components with and without assets?**
A: Yes! Components without AssetProvider work normally.

**Q: What about Tailwind CSS?**
A: You can use Tailwind! Just compile it to CSS and embed it as a shared asset.

**Q: What about TypeScript/SASS?**
A: Compile to JS/CSS first, then embed the compiled output:

```go
//go:generate sass alert.scss alert.css
//go:generate esbuild alert.ts --outfile=alert.js

//go:embed alert.css
var alertCSS []byte
```

**Q: How do I handle CSS conflicts?**
A: Use namespaced class names (BEM, CSS modules pattern):

```css
/* alert.css */
.ele-alert { ... }
.ele-alert__title { ... }
.ele-alert--success { ... }
```

## See Also

- [COMPONENT_ASSET_MANAGEMENT.md](../../ai_docs/COMPONENT_ASSET_MANAGEMENT.md) - Full design document
- [Asset Provider Interface](../../asset_provider.go) - Implementation reference
- [Components Example](../components/) - Components without asset management
