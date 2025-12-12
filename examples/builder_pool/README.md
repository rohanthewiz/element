# Builder Pool Example

This example demonstrates how to use the builder pool for high-throughput HTTP scenarios.

## Overview

The `element` library provides a `sync.Pool`-based builder pool that reduces memory allocations in high-request-rate applications. Instead of creating a new `Builder` for each request, you can acquire one from the pool and release it when done.

## Quick Start

```bash
cd examples/builder_pool
go run main.go
```

Then visit http://localhost:8080

## API

### AcquireBuilder

Gets a `Builder` from the pool, or creates a new one if the pool is empty. The returned builder is reset and ready for use.

```go
b := element.AcquireBuilder()
```

### ReleaseBuilder

Returns a `Builder` to the pool for reuse. After calling `ReleaseBuilder`, the builder must not be used again.

```go
element.ReleaseBuilder(b)
```

## Usage Pattern

```go
func handler(w http.ResponseWriter, r *http.Request) {
    // Acquire a builder from the pool
    b := element.AcquireBuilder()
    // Always release back to the pool when done
    defer element.ReleaseBuilder(b)

    // Build your HTML
    b.Html().R(
        b.Head().R(
            b.Title().T("My Page"),
        ),
        b.Body().R(
            b.H1().T("Hello, World!"),
            b.P().T("This page uses a pooled builder."),
        ),
    )

    // Write the response
    w.Header().Set("Content-Type", "text/html; charset=utf-8")
    w.Write([]byte(b.String()))
}
```

## When to Use Pooling

Use `AcquireBuilder`/`ReleaseBuilder` when:

- **High-traffic HTTP handlers** - Thousands of requests per second
- **APIs rendering HTML** - Server-side rendered responses
- **Batch processing** - Generating many HTML documents in a loop
- **Memory-constrained environments** - Where reducing GC pressure matters

Use regular `NewBuilder()` when:

- **Low-traffic scenarios** - Occasional HTML generation
- **Debug mode** - Pooling is automatically disabled in debug mode
- **Simplicity preferred** - When the overhead doesn't matter

## Performance Benefits

The builder pool provides:

1. **Reduced allocations** - Reuses `Builder` and `strings.Builder` instances
2. **Lower GC pressure** - Fewer objects for the garbage collector to track
3. **More stable latency** - Fewer GC pauses means more consistent response times
4. **Capacity reuse** - The internal `strings.Builder` retains its grown capacity

### Benchmark Results

From the library's benchmark suite:

```
BenchmarkBuilderAllocation-16    56412    19684 ns/op   5541 B/op   80 allocs/op
BenchmarkBuilderPooled-16        62487    19716 ns/op   5457 B/op   76 allocs/op
```

The pooled version saves 4 allocations per operation (the `Builder` struct and `strings.Builder`).

## Load Testing

You can verify the benefits with load testing tools like `ab` (Apache Bench):

```bash
# Start the server
go run main.go

# In another terminal, test the pooled endpoint
ab -n 10000 -c 100 http://localhost:8080/pooled

# Test the regular endpoint
ab -n 10000 -c 100 http://localhost:8080/regular

# Check memory stats
curl http://localhost:8080/stats
```

Compare the `TotalAlloc` and `NumGC` values between runs.

## Thread Safety

The builder pool is safe for concurrent use. Multiple goroutines can call `AcquireBuilder()` and `ReleaseBuilder()` simultaneously without external synchronization.

```go
// Safe to use in concurrent handlers
http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
    b := element.AcquireBuilder()
    defer element.ReleaseBuilder(b)
    // ... each request gets its own builder instance
})
```

## Debug Mode

TODO verify this as I am feeling that the only debug state maintained is element concerns... // When debug mode is enabled (`element.DebugSet()`), builders are **not** returned to the pool. This prevents debug state from leaking between requests and makes debugging easier.

```go
// In debug mode, this effectively becomes NewBuilder()
b := element.AcquireBuilder()
defer element.ReleaseBuilder(b)  // Builder is discarded, not pooled
```

## Best Practices - Important!

1. **Always use defer** - Ensures the builder is released even if a panic occurs:
   ```go
   b := element.AcquireBuilder()
   defer element.ReleaseBuilder(b)
   ```

2. **Don't store references** - Don't keep references to a released builder:
   ```go
   b := element.AcquireBuilder()
   html := b.String()  // Get the string before releasing
   element.ReleaseBuilder(b)
   // Don't use 'b' after this point!
   ```

3. **Extract the string first** - Call `b.String()` before releasing:
   ```go
   b := element.AcquireBuilder()
   defer element.ReleaseBuilder(b)

   b.Div().T("Hello")
   return b.String()  // Safe: string is copied before release
   ```

4. **Use regular builders for stored results** - If you need to keep the builder around:
   ```go
   // For long-lived builders, use NewBuilder
   cache := element.NewBuilder()
   // ... build cached content ...
   cachedHTML := cache.String()
   ```

## Example Endpoints

This example provides several endpoints:

| Endpoint | Description |
|----------|-------------|
| `/` | Home page with documentation |
| `/pooled` | Renders a page using pooled builder |
| `/regular` | Renders a page using regular builder |
| `/stats` | Shows runtime memory statistics |

## Files

- `main.go` - Example HTTP server demonstrating builder pooling
- `README.md` - This documentation file
