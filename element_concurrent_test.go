package element

import (
	"strings"
	"sync"
	"testing"
)

// Rendering on separate builders from many goroutines must be safe,
// including in debug mode where concern tracking shares state.
// Run with -race to verify.
func TestConcurrentRenderDebugMode(t *testing.T) {
	DebugSet()
	defer DebugClear()

	var wg sync.WaitGroup
	for i := 0; i < 16; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 50; j++ {
				b := NewBuilder()
				// Every open/close in debug mode upserts into the shared concerns map
				b.Div("class", "row").R(
					b.Span().T("hello"),
				)
				if got := b.String(); !strings.Contains(got, "hello") {
					t.Errorf("unexpected output %q", got)
				}
			}
		}()
	}
	wg.Wait()
}

func TestConcurrentRenderWithPool(t *testing.T) {
	var wg sync.WaitGroup
	for i := 0; i < 16; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				b := AcquireBuilder()
				b.Div("class", "row").R(
					b.Span().T("hello"),
				)
				want := `<div class="row"><span>hello</span></div>`
				if got := b.String(); got != want {
					t.Errorf("got %q, want %q", got, want)
				}
				ReleaseBuilder(b)
			}
		}()
	}
	wg.Wait()
}
