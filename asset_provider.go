package element

import (
	"bytes"
)

// AssetProvider is an optional interface that components can implement
// to declare their CSS and JavaScript dependencies.
//
// Components implementing this interface can have their assets automatically
// collected and deduplicated by the Builder when asset tracking is enabled.
//
// Example:
//
//	type Alert struct { ... }
//
//	func (a Alert) AssetID() string { return "Alert" }
//	func (a Alert) CSS() []byte { return alertCSS }
//	func (a Alert) JS() []byte { return alertJS }
type AssetProvider interface {
	// AssetID returns a unique identifier for this component's assets.
	// This is used for deduplication - components with the same AssetID
	// will only have their assets included once.
	//
	// Typically this should be the component type name (e.g., "Alert", "Table").
	// For versioning, you can include a version: "Alert@1.2.0"
	AssetID() string

	// CSS returns the component's CSS as bytes.
	// Return nil if the component doesn't need CSS.
	//
	// Typically loaded with //go:embed:
	//   //go:embed alert.css
	//   var alertCSS []byte
	CSS() []byte

	// JS returns the component's JavaScript as bytes.
	// Return nil if the component doesn't need JavaScript.
	//
	// Typically loaded with //go:embed:
	//   //go:embed alert.js
	//   var alertJS []byte
	JS() []byte
}

// assetEntry represents a collected asset from a component.
type assetEntry struct {
	id  string
	css []byte
	js  []byte
}

// RegisterAssets registers a component's assets if it implements AssetProvider.
// This is called by components in their Render() method to declare their dependencies.
//
// If asset tracking is not enabled on the builder, this is a no-op.
// If the component doesn't implement AssetProvider, this is a no-op.
// If the asset ID has already been registered, this is a no-op (deduplication).
//
// Example in component:
//
//	func (a Alert) Render(b *element.Builder) any {
//	    b.RegisterAssets(a)
//	    // ... render implementation
//	    return nil
//	}
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
		return // No ID provided
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

// EnableAssetTracking enables automatic asset collection for components.
// Call this when creating a builder if you want to use component asset management.
//
// When enabled, the builder will collect CSS and JavaScript from components
// that implement the AssetProvider interface. Assets are automatically deduplicated.
//
// Example:
//
//	b := element.NewBuilder().EnableAssetTracking()
//	// ... render components
//	b.Head().R(
//	    b.RenderCSS(), // Outputs all collected CSS
//	)
func (b *Builder) EnableAssetTracking() *Builder {
	b.assetsEnabled = true
	if b.assetRegistry == nil {
		b.assetRegistry = make(map[string]assetEntry)
	}
	return b
}

// RenderCSS renders all collected CSS in a <style> tag.
// This should typically be called in the <head> section.
//
// If asset tracking is not enabled or no CSS has been collected, this is a no-op.
//
// Example:
//
//	b.Head().R(
//	    b.Title().T("My Page"),
//	    b.RenderCSS(),
//	)
func (b *Builder) RenderCSS() *Builder {
	if !b.assetsEnabled || len(b.assetRegistry) == 0 {
		return b
	}

	hasCSS := false
	for _, entry := range b.assetRegistry {
		if len(entry.css) > 0 {
			hasCSS = true
			break
		}
	}

	if !hasCSS {
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
// This should typically be called before </body>.
//
// If asset tracking is not enabled or no JavaScript has been collected, this is a no-op.
//
// Example:
//
//	b.Body().R(
//	    // ... components
//	    b.RenderJS(),
//	)
func (b *Builder) RenderJS() *Builder {
	if !b.assetsEnabled || len(b.assetRegistry) == 0 {
		return b
	}

	hasJS := false
	for _, entry := range b.assetRegistry {
		if len(entry.js) > 0 {
			hasJS = true
			break
		}
	}

	if !hasJS {
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

// GetCollectedCSS returns all collected CSS as a single byte slice.
// This is useful for serving CSS as an external file or for inspection.
//
// If asset tracking is not enabled, this returns nil.
//
// Example (serving CSS as external file):
//
//	func HandleCSS(w http.ResponseWriter, r *http.Request) {
//	    b := element.NewBuilder().EnableAssetTracking()
//	    // ... render all components to collect assets
//	    Alert{}.Render(b)
//	    Table{}.Render(b)
//
//	    css := b.GetCollectedCSS()
//	    w.Header().Set("Content-Type", "text/css")
//	    w.Write(css)
//	}
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

	if buf.Len() == 0 {
		return nil
	}

	return buf.Bytes()
}

// GetCollectedJS returns all collected JavaScript as a single byte slice.
// This is useful for serving JavaScript as an external file or for inspection.
//
// If asset tracking is not enabled, this returns nil.
//
// Example (serving JS as external file):
//
//	func HandleJS(w http.ResponseWriter, r *http.Request) {
//	    b := element.NewBuilder().EnableAssetTracking()
//	    // ... render all components to collect assets
//	    Alert{}.Render(b)
//	    Table{}.Render(b)
//
//	    js := b.GetCollectedJS()
//	    w.Header().Set("Content-Type", "application/javascript")
//	    w.Write(js)
//	}
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

	if buf.Len() == 0 {
		return nil
	}

	return buf.Bytes()
}
