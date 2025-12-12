package element

import (
	"sync"
	"testing"
)

func TestAcquireBuilder(t *testing.T) {
	b := AcquireBuilder()
	if b == nil {
		t.Fatal("AcquireBuilder returned nil")
	}
	defer ReleaseBuilder(b)

	// Verify the builder works correctly
	b.Div("class", "test").R(
		b.P().T("Hello"),
	)

	result := b.String()
	expected := `<div class="test"><p>Hello</p></div>`
	if result != expected {
		t.Errorf("unexpected output: got %q, want %q", result, expected)
	}
}

func TestReleaseBuilderNil(t *testing.T) {
	// Should not panic on nil
	ReleaseBuilder(nil)
}

func TestPooledBuilderIsReset(t *testing.T) {
	// Get a builder and use it
	b1 := AcquireBuilder()
	b1.Div().T("first content")
	content1 := b1.String()
	if content1 != "<div>first content</div>" {
		t.Errorf("first builder output wrong: %q", content1)
	}
	ReleaseBuilder(b1)

	// Get another builder (likely the same one from pool)
	b2 := AcquireBuilder()
	b2.Span().T("second")
	result := b2.String()
	defer ReleaseBuilder(b2)

	// Should NOT contain content from previous use
	expected := "<span>second</span>"
	if result != expected {
		t.Errorf("builder not properly reset, got %q, want %q", result, expected)
	}
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
			result := b.String()
			if result == "" {
				t.Error("empty result from concurrent builder")
			}
			ReleaseBuilder(b)
		}(i)
	}

	wg.Wait()
}

func TestPoolMultipleAcquireRelease(t *testing.T) {
	// Acquire and release multiple times to exercise the pool
	for i := 0; i < 100; i++ {
		b := AcquireBuilder()
		b.P().F("iteration %d", i)
		_ = b.String()
		ReleaseBuilder(b)
	}
}

func TestPooledBuilderFunctionality(t *testing.T) {
	b := AcquireBuilder()
	defer ReleaseBuilder(b)

	// Test various builder methods work correctly
	b.Html().R(
		b.Head().R(
			b.Title().T("Test Page"),
		),
		b.Body().R(
			b.Div("id", "main").R(
				b.H1().T("Hello"),
				b.P("class", "intro").T("Welcome"),
				b.Ul().R(
					b.Li().T("Item 1"),
					b.Li().T("Item 2"),
				),
			),
		),
	)

	result := b.String()

	// Verify key elements are present
	checks := []string{
		"<html>",
		"<head>",
		"<title>Test Page</title>",
		"<body>",
		`<div id="main">`,
		"<h1>Hello</h1>",
		`<p class="intro">Welcome</p>`,
		"<ul>",
		"<li>Item 1</li>",
		"<li>Item 2</li>",
		"</ul>",
		"</div>",
		"</body>",
		"</html>",
	}

	for _, check := range checks {
		if !contains(result, check) {
			t.Errorf("result missing %q\nfull result: %s", check, result)
		}
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// Benchmarks

// BenchmarkBuilderAllocation benchmarks creating new builders each time
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

// BenchmarkBuilderPooled benchmarks using pooled builders
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

// BenchmarkBuilderPooledParallel benchmarks pooled builders under concurrent load
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

// BenchmarkBuilderAllocationParallel benchmarks new builders under concurrent load
func BenchmarkBuilderAllocationParallel(b *testing.B) {
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			builder := NewBuilder()
			builder.Div("class", "container").R(
				builder.P().T("Hello, World!"),
			)
			_ = builder.String()
		}
	})
}

// BenchmarkBuilderLargePage benchmarks building a larger page
func BenchmarkBuilderLargePage(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		builder := NewBuilder()
		buildLargePage(builder)
		_ = builder.String()
	}
}

// BenchmarkBuilderLargePagePooled benchmarks building a larger page with pooling
func BenchmarkBuilderLargePagePooled(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		builder := AcquireBuilder()
		buildLargePage(builder)
		_ = builder.String()
		ReleaseBuilder(builder)
	}
}

// buildLargePage generates a more realistic page structure for benchmarking
func buildLargePage(b *Builder) {
	b.Html().R(
		b.Head().R(
			b.Title().T("Benchmark Page"),
			b.Meta("charset", "utf-8"),
			b.Meta("name", "viewport", "content", "width=device-width"),
		),
		b.Body().R(
			b.Header().R(
				b.Nav().R(
					b.A("href", "/").T("Home"),
					b.A("href", "/about").T("About"),
					b.A("href", "/contact").T("Contact"),
				),
			),
			b.Main().R(
				b.H1().T("Welcome"),
				b.P().T("This is a benchmark page with various elements."),
				b.Div("class", "content").R(
					buildTable(b),
					buildList(b),
				),
			),
			b.Footer().R(
				b.P().T("Copyright 2024"),
			),
		),
	)
}

func buildTable(b *Builder) any {
	return b.Table("class", "data-table").R(
		b.THead().R(
			b.Tr().R(
				b.Th().T("Name"),
				b.Th().T("Value"),
				b.Th().T("Status"),
			),
		),
		b.TBody().R(
			b.Tr().R(b.Td().T("Row 1"), b.Td().T("100"), b.Td().T("Active")),
			b.Tr().R(b.Td().T("Row 2"), b.Td().T("200"), b.Td().T("Pending")),
			b.Tr().R(b.Td().T("Row 3"), b.Td().T("300"), b.Td().T("Active")),
			b.Tr().R(b.Td().T("Row 4"), b.Td().T("400"), b.Td().T("Inactive")),
			b.Tr().R(b.Td().T("Row 5"), b.Td().T("500"), b.Td().T("Active")),
		),
	)
}

func buildList(b *Builder) any {
	return b.Ul("class", "feature-list").R(
		b.Li().T("Feature One"),
		b.Li().T("Feature Two"),
		b.Li().T("Feature Three"),
		b.Li().T("Feature Four"),
		b.Li().T("Feature Five"),
	)
}
