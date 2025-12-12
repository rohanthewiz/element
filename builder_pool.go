package element

import "sync"

// builderPool is a pool of reusable Builder instances.
// Using a pool reduces allocation pressure in high-throughput scenarios
// by reusing Builder instances across requests.
var builderPool = sync.Pool{
	New: func() any {
		return NewBuilder()
	},
}

// AcquireBuilder gets a Builder from the pool, or creates a new one if the pool is empty.
// The returned Builder is reset and ready for use.
//
// Usage:
//
//	b := element.AcquireBuilder()
//	defer element.ReleaseBuilder(b)
//	// ... use builder ...
//	html := b.String()
//
// For high-throughput HTTP handlers, using AcquireBuilder/ReleaseBuilder
// reduces GC pressure compared to creating new builders with NewBuilder().
func AcquireBuilder() *Builder {
	return builderPool.Get().(*Builder)
}

// ReleaseBuilder returns a Builder to the pool for reuse.
// After calling ReleaseBuilder, the Builder must not be used again.
//
// The builder is automatically reset before being returned to the pool.
// In debug mode, builders are not pooled to prevent cross-request debug state confusion.
func ReleaseBuilder(b *Builder) {
	if b == nil {
		return
	}

	b.Reset()

	// In debug mode, don't pool builders to avoid debug state leaking between requests
	if IsDebugMode() {
		return
	}

	builderPool.Put(b)
}
