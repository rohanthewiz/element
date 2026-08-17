# Attribute-value escaping, and a TE counterpart to T

**Session:** `67789422-08e6-4652-b190-a5e675318f31`
**Date:** 2026-08-17
**Repo:** `~/projs/go/element` (branch `escape-attribute-values`)
**Origin:** an XSS audit in a downstream consumer (`~/cbre_projs/edp_dataflow`),
which is where the motivating bug was found

---

## How this came about

The change did not start here. A SAST scanner flagged two lines in the consuming
app; auditing outward from them turned up unescaped user data reaching the DOM in
that app's dashboard, rendered two ways — server-side through element, and
client-side in JS on live updates.

The initial read of the server side was that it was **safe**, on the assumption
that `Builder.T` escaped. It does not. `builder_funcs.go` writes its arguments
straight to the buffer, and the library performs no escaping anywhere. That is a
coherent design choice, not an oversight, but it means the consumer's server-
rendered rows carried the same injection as the client-rendered ones.

Two library-shaped gaps came out of fixing it, and they are opposite in kind.

---

## 1. Attribute values — a bug, not a policy

`element.go:176` was the single place every attribute value in the library is
written. It emitted the value raw between the quotes it had just opened.

```
b.Div("data-id", `x" onmouseover="alert(1)`)
  →  <div data-id="x" onmouseover="alert(1)"></div>
```

One attribute in, two attributes out, the second an event handler. This is the
whole reason `writeAttrValue` now exists.

**What it escapes: the double quote. Nothing else.** That is the point, and each
omission is a decision:

| Character | Left alone because |
| --- | --- |
| `<` `>` | Do not terminate a quoted attribute value — the tokenizer is looking for the closing quote and nothing else |
| `'` | Cannot close a value the renderer opened with `"` |
| `&` | Escaping it would double-encode every caller who already passes character references |

The ampersand is the one worth arguing about. Encoding it is stricter — a bare
`&` ought to be `&amp;` — but a careful caller building a JS string literal for
an inline handler passes `&#39;` deliberately, and would get `&amp;#39;`
rendered as visible text. A real regression traded for a cosmetic fix, and a
bare `&` cannot break out of an attribute value anyway.

A useful consequence: `html.EscapeString` on a value stays correct under this
change, because an already escaped string has no double quotes left to find.
Pinned by `TestPreEscapedAttributeValuesAreNotDoubleEncoded`.

**Why this is not a breaking change.** A raw `"` inside an attribute value is
broken output every time — there is no program that wants it. So escaping it
cannot alter what any *correct* caller renders, only what an incorrect one does.
The library's existing test suite, which is full of attribute assertions, passed
untouched.

Fast path preserved: `strings.IndexByte` finds no quote in the overwhelming
majority of values and the string is written straight through with no
allocation. 317 ns/op with no quote vs 327 with two, 7 allocs either way.

---

## 2. `T` stays raw — `TE` is the new one

The opposite conclusion for element *content*.

`T` cannot start escaping. `b.Style().T(css)` and `b.Script().T(js)` inline whole
stylesheets and scripts, and both break completely if it does. That is now
pinned by `TestTStaysRaw` so the reasoning survives the next person who notices
`T` is unescaped and decides to "fix" it.

So the safe path is a new one-character-away method rather than a changed
default:

```go
b.Td().T(css)       // verbatim — markup you authored
b.Td().TE(userName) // escaped — text the program did not author
```

The asymmetry is what makes the one-character difference defensible. Reaching
for `T` on a database value fails silently and invisibly until someone stores a
tag. Reaching for `TE` on markup you meant to inline fails instantly and
obviously — the tags show up on screen as text.

**No `FE`.** Escaping a format string is ambiguous: the caller means "escape the
arguments, not the template", and an API that guesses would be worse than one
that makes the intent explicit — `TE(fmt.Sprintf(...))`.

Added on both `Element` and `Builder`.

---

## What was rejected

**An `EscapeHTML` helper on the library.** The first instinct was to export the
escaper the consumer had hand-rolled. Benchmarking killed it: `html.EscapeString`
covers exactly the same five characters (spelling `&#34;` where a hand-rolled
version writes `&quot;` — parsers treat them identically), it is stdlib so it
costs no dependency, and it is ~12× faster than a `strings.Replacer` built per
call, which is the shape a local copy drifts into. A library function would be a
second name for something already in the standard library.

**Escaping by default with a `Raw()` escape hatch.** Correct-by-default is the
better policy in the abstract, but it inverts a published API and breaks every
existing `T(css)` call. The value is not worth the churn for a library whose
whole proposition is being a fast raw writer.

---

## Files

| File | Change |
| --- | --- |
| `element.go` | `writeAttrValue` + call site at the attribute writer; `Element.TE`; doc note on `Element.T` |
| `builder_funcs.go` | `Builder.TE`; doc note on `Builder.T` |
| `escaping_test.go` | **new** — 9 tests + 2 benchmarks |

---

## Verification

- Existing suite passes unchanged — the meaningful backward-compatibility signal,
  since it asserts heavily on attribute output.
- The breakout test was confirmed to catch the old behaviour by temporarily
  restoring the raw write: it rendered `<div data-id="x" onmouseover="alert(1)">`,
  failing on all three assertions (handler injected, quote not escaped, two
  attributes where one was expected).
- `gofmt`, `go vet`, `go build ./...` clean.

**Process note.** The first attempt at that regression check chained a
`git checkout element.go` after a patch step that silently failed to match, which
reverted the file and lost the edits. They were rebuilt and the check redone
against a file copy instead. Worth remembering: never chain a destructive restore
behind a patch whose match was not verified.

---

## Open items

1. **Nothing is committed on the branch yet** beyond this doc — review the diff
   before merging.
2. **`&` in attribute values stays unencoded.** Documented as deliberate in
   `writeAttrValue`. If element ever grows a strict mode, that is where it goes.
3. **`TE` does not help attributes.** Attribute values are now safe from
   breakout but there is no escaped-attribute affordance; a caller wanting full
   entity encoding still reaches for `html.EscapeString`. Possibly fine forever —
   worth revisiting only if the breakout fix proves insufficient in practice.
4. **README/CLAUDE.md not updated.** The T/TE distinction is the kind of thing
   that belongs in the docs, not only in doc comments.
