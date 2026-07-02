package element

import (
	"strings"
	"sync"
	"sync/atomic"
	"testing"
)

// countingComp tracks how many times Render is actually invoked
type countingComp struct {
	renders atomic.Int32
}

func (c *countingComp) Render(b *Builder) (x any) {
	c.renders.Add(1)
	b.DivClass("cached-widget").R(
		b.Span().T("static content"),
	)
	return
}

func TestCachedRendersOnce(t *testing.T) {
	comp := &countingComp{}
	cached := Cached(comp)

	want := `<div class="cached-widget"><span>static content</span></div>`

	for i := 0; i < 5; i++ {
		b := NewBuilder()
		cached.Render(b)
		if got := b.String(); got != want {
			t.Fatalf("render %d: got %q, want %q", i, got, want)
		}
	}

	if n := comp.renders.Load(); n != 1 {
		t.Errorf("expected underlying component to render once, rendered %d times", n)
	}
}

func TestCachedInRenderTree(t *testing.T) {
	cached := Cached(&countingComp{})

	b := NewBuilder()
	b.Body().R(
		b.H1().T("Title"),
		cached.Render(b),
	)

	want := `<body><h1>Title</h1><div class="cached-widget"><span>static content</span></div></body>`
	if got := b.String(); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestCachedConcurrent(t *testing.T) {
	comp := &countingComp{}
	cached := Cached(comp)

	want := `<div class="cached-widget"><span>static content</span></div>`

	var wg sync.WaitGroup
	for i := 0; i < 16; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 50; j++ {
				b := AcquireBuilder()
				cached.Render(b)
				if got := b.String(); got != want {
					t.Errorf("got %q, want %q", got, want)
				}
				ReleaseBuilder(b)
			}
		}()
	}
	wg.Wait()

	if n := comp.renders.Load(); n != 1 {
		t.Errorf("expected underlying component to render once, rendered %d times", n)
	}
}

// In debug mode the cache is bypassed so concern tracking still works
func TestCachedBypassesCacheInDebugMode(t *testing.T) {
	DebugSet()
	defer DebugClear()

	comp := &countingComp{}
	cached := Cached(comp)

	for i := 0; i < 3; i++ {
		b := NewBuilder()
		cached.Render(b)
		if got := b.String(); !strings.Contains(got, "static content") {
			t.Fatalf("render %d: unexpected output %q", i, got)
		}
	}

	if n := comp.renders.Load(); n != 3 {
		t.Errorf("expected 3 renders in debug mode, got %d", n)
	}
}
