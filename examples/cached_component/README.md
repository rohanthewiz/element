# Cached Component Example

Demonstrates `element.Cached`, which wraps a `Component` so it renders **once**
and replays the cached bytes on every subsequent render.

```go
var cachedNav = element.Cached(NavBar{}) // package-level: shared across requests

// in a handler
b.Body().R(
    cachedNav.Render(b),   // first call renders; later calls copy cached bytes
    dynamicContent(b),
    cachedFooter.Render(b),
)
```

## When to use it

- Static components: nav bars, footers, icon sets, headers that don't change
- High-throughput handlers where rebuilding the same HTML per request is waste

## Behavior notes

- Safe for concurrent use (`sync.Once` inside)
- In **debug mode** the cache is bypassed so element concern tracking still works
- There is no invalidation — if the content can change, don't cache it

## Performance

From the library benchmarks (Apple M1 Pro):

| Benchmark | ns/op | allocs/op |
|---|---|---|
| Component (uncached) | ~879 | 13 |
| Component (cached) | ~44 | 1 |

## Run

```bash
go run .
# visit http://localhost:8080 and reload — the nav render count stays at 1
```
