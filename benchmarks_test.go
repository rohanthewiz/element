package element

import "testing"

func BenchmarkDivRender(b *testing.B) {
	for i := 0; i < b.N; i++ {
		bld := NewBuilder()
		bld.Div("class", "row", "id", "main").R(
			bld.Span("class", "cell").T("hello"),
			bld.Span("class", "cell").T("world"),
		)
		_ = bld.String()
	}
}

func BenchmarkDivRenderPooled(b *testing.B) {
	for i := 0; i < b.N; i++ {
		bld := AcquireBuilder()
		bld.Div("class", "row", "id", "main").R(
			bld.Span("class", "cell").T("hello"),
			bld.Span("class", "cell").T("world"),
		)
		_ = bld.String()
		ReleaseBuilder(bld)
	}
}

type benchComp struct{}

func (benchComp) Render(b *Builder) (x any) {
	b.DivClass("footer").R(
		b.Ul().R(
			b.Li().R(b.A("href", "/about").T("About")),
			b.Li().R(b.A("href", "/contact").T("Contact")),
			b.Li().R(b.A("href", "/privacy").T("Privacy")),
		),
	)
	return
}

func BenchmarkComponent(b *testing.B) {
	comp := benchComp{}
	for i := 0; i < b.N; i++ {
		bld := AcquireBuilder()
		comp.Render(bld)
		_ = bld.String()
		ReleaseBuilder(bld)
	}
}

func BenchmarkComponentCached(b *testing.B) {
	cached := Cached(benchComp{})
	for i := 0; i < b.N; i++ {
		bld := AcquireBuilder()
		cached.Render(bld)
		_ = bld.String()
		ReleaseBuilder(bld)
	}
}
