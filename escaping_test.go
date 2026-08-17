package element

import (
	"strings"
	"testing"
)

// element is a raw writer by design, with exactly one exception: an attribute
// value cannot contain the double quote that delimits it. These tests pin both
// halves of that — the exception, and the fact that it is the only one.

// A value that closes its attribute and opens a new one is the whole reason
// writeAttrValue exists. Everything after the quote is read by the parser as
// further attributes on the element, which is how a data field becomes an event
// handler.
const attrBreakout = `x" onmouseover="alert(1)`

func TestAttributeValueCannotCloseItsAttribute(t *testing.T) {
	b := NewBuilder()
	b.Div("data-id", attrBreakout).R()
	got := b.String()

	// The raw sequence must not survive. Checking for the sequence rather than
	// for the absence of a quote is deliberate: the output is full of legitimate
	// quotes, and only this one is a breakout.
	if strings.Contains(got, `" onmouseover="`) {
		t.Errorf("rendered %s — the value closed its attribute and injected a handler", got)
	}

	// It must still be present, escaped. Without this a renderer that dropped
	// the attribute entirely would pass the check above.
	if !strings.Contains(got, "&#34;") {
		t.Errorf("rendered %s — the quote was neither escaped nor emitted", got)
	}

	// Exactly one attribute should have been produced.
	if n := strings.Count(got, `="`); n != 1 {
		t.Errorf("rendered %s — expected 1 attribute, found %d", got, n)
	}
}

// Values that need no escaping must come through byte for byte. This is the
// overwhelmingly common case and the one most likely to be disturbed by a change
// to the escaping path.
func TestAttributeValuesAreOtherwiseUntouched(t *testing.T) {
	cases := []string{
		"btn btn-primary",
		"/path/to/thing?a=1&b=2",               // bare & stays bare
		"event.stopPropagation(); copyID('x')", // single quotes are safe in a double-quoted value
		"a < b > c",                            // angle brackets do not delimit anything here
		"&#39;",                                // a caller's own character reference is not re-encoded
		"",
	}

	for _, val := range cases {
		b := NewBuilder()
		b.Div("data-x", val).R()
		want := `<div data-x="` + val + `"></div>`
		if got := b.String(); got != want {
			t.Errorf("value %q rendered %s, want %s", val, got, want)
		}
	}
}

// Escaping an already escaped value has to stay correct — html.EscapeString
// leaves no double quotes behind, so writeAttrValue finds nothing to do and the
// value is not encoded twice.
func TestPreEscapedAttributeValuesAreNotDoubleEncoded(t *testing.T) {
	b := NewBuilder()
	b.Div("data-x", "&#34;already escaped&#34;").R()

	want := `<div data-x="&#34;already escaped&#34;"></div>`
	if got := b.String(); got != want {
		t.Errorf("rendered %s, want %s", got, want)
	}
}

// T must keep writing verbatim. This is not a style preference — inlining a
// stylesheet or a script depends on it, and those break silently and completely
// if T ever starts escaping.
func TestTStaysRaw(t *testing.T) {
	const css = `.a > .b { content: "x"; }`

	b := NewBuilder()
	b.Style().T(css)

	want := "<style>" + css + "</style>"
	if got := b.String(); got != want {
		t.Errorf("rendered %s, want %s — T must not escape", got, want)
	}
}

func TestTEEscapes(t *testing.T) {
	b := NewBuilder()
	b.Div().TE(`<img src=x onerror=alert(1)>`)
	got := b.String()

	if strings.Contains(got, "<img") {
		t.Errorf("rendered %s — TE emitted a live tag", got)
	}
	if !strings.Contains(got, "&lt;img") {
		t.Errorf("rendered %s — TE did not emit the escaped text", got)
	}
}

// TE takes a list, like T, and every member is escaped — not just the first.
func TestTEEscapesEveryArgument(t *testing.T) {
	b := NewBuilder()
	b.Div().TE("safe", "<script>", "also safe")

	if got := b.String(); strings.Contains(got, "<script>") {
		t.Errorf("rendered %s — a later argument escaped unescaped", got)
	}
}

func TestBuilderTEEscapes(t *testing.T) {
	b := NewBuilder()
	b.TE(`<b>`)

	if got := b.String(); got != "&lt;b&gt;" {
		t.Errorf("rendered %q, want %q", got, "&lt;b&gt;")
	}
}

// The escaping path must not cost anything when there is nothing to escape,
// since it runs for every attribute of every element.
func BenchmarkAttrValueNoQuote(b *testing.B) {
	for i := 0; i < b.N; i++ {
		bld := NewBuilder()
		bld.Div("class", "btn btn-primary", "data-id", "batch-1234-abcd").R()
	}
}

func BenchmarkAttrValueWithQuote(b *testing.B) {
	for i := 0; i < b.N; i++ {
		bld := NewBuilder()
		bld.Div("class", `btn "primary"`, "data-id", `batch-"1234"`).R()
	}
}
