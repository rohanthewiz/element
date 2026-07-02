package element

import (
	"bytes"
	"sync"
)

// cachedComponent wraps a Component, rendering it once
// and replaying the cached bytes on subsequent renders
type cachedComponent struct {
	comp Component
	once sync.Once
	out  []byte
}

// Cached wraps a component so it is rendered once, with the output bytes
// reused on every subsequent Render. Use it for static components
// (nav bars, footers, icon sets) that are expensive to build but don't change.
//
// The wrapper is safe for concurrent use.
// In debug mode the underlying component is rendered fresh each time
// so element concerns are still tracked.
//
// Example:
//
//	var footer = element.Cached(Footer{})
//	// in a handler
//	b.Body().R(
//		content.Render(b),
//		footer.Render(b),
//	)
func Cached(comp Component) Component {
	return &cachedComponent{comp: comp}
}

func (c *cachedComponent) Render(b *Builder) (x any) {
	if IsDebugMode() {
		return c.comp.Render(b)
	}

	c.once.Do(func() {
		cb := NewBuilder()
		c.comp.Render(cb)
		c.out = bytes.Clone(cb.Bytes())
	})

	_ = b.WriteBytes(c.out)
	return
}
