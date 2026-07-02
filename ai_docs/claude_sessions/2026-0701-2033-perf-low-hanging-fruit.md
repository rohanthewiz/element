# Session: Performance Low-Hanging Fruit

**Date:** 2026-07-01 20:33 | **Session ID:** 6603ff72-99b3-40ee-8d0f-3147ae3dfa3e

## Goal

Identify and implement "low-hanging fruit" improvements in the Element library. The user suggested component caching (a simple on-instance cache); investigation surfaced several more wins, including one that dwarfed it.

## Findings (in order of value)

1. **Stack walks on every element creation** — `New()` called `serr.FunctionName()`/`serr.FunctionLoc()` unconditionally (two `runtime.Caller` walks + string splits per element), even though `function`/`location` are only read by debug-mode code. Gating them behind `IsDebugMode()` alone gave a **4.8× speedup** (6,623 → 1,368 ns/op on a div-with-two-spans benchmark).

2. **Debug-mode concurrency bugs**
   - `elementConcerns.UpsertConcern`/`.Clear` used **value receivers** on a struct containing a `sync.Mutex` (`go vet` copylocks). Each call locked a *copy* of the mutex, so the shared `cmap` had no protection — concurrent debug-mode rendering could panic with concurrent map writes.
   - `debugMode` was a plain global bool, and `DebugShow` toggled it off/on mid-call — a data race.
   - The race detector also caught a third: `seededRand` (`*rand.Rand`) is not safe for concurrent use and `genRandomId` is called per element in debug mode.

3. **`fmt.Sprintf` per attribute** in `writeOpeningTag`, plus string concatenation in tag open/close.

4. **Nondeterministic attribute order** — attributes lived in a `map[string]string`, so the same code produced differently-ordered HTML run to run (old tests even had comments working around map randomization).

5. **Component caching** — the user's original idea; straightforward since components write to the builder.

## Changes Made

- `element.go`
  - `serr.FunctionName`/`FunctionLoc` now only run when `IsDebugMode()` is true.
  - `Element.attrs map[string]string` → `attrPairs []string` (ordered key/value pairs; deterministic output).
  - `writeOpeningTag`/`close` write directly (`WriteByte`/`WriteString`) — no `fmt.Sprintf`, no concatenation.
  - `HasAttribute` scans the pair slice.
- `helpers.go`
  - `stringlistToMap` → `normalizeAttrPairs` + `upsertPair`: preserves insertion order, duplicate keys keep first position with last value winning, odd trailing arg still dropped (with debug concern), `data-ele-id` still appended in debug mode.
  - `lowerName` fast path: skips the `strings.ToLower` allocation when the element name is already lowercase (the common case — all builder methods pass lowercase literals).
  - `seededRand` now guarded by a mutex.
- `element_debug.go`
  - `debugMode` is an `atomic.Bool`.
  - `elementConcerns` methods use pointer receivers; `concerns` is a pointer; `DebugShow` reads via a locked `snapshot()`.
  - Simplified an issue-list copy loop (staticcheck S1001).
- `component_cached.go` (new)
  - `element.Cached(comp)` wraps a `Component`: renders once via `sync.Once` into cached bytes, replays with `WriteBytes` thereafter. Safe for concurrent use. **Bypasses the cache in debug mode** so concern tracking still works.

## Tests Added

- `element_attrs_test.go` — attribute order determinism (50 runs), duplicate last-wins, odd-count drop, `HasAttribute`, `*Class` ordering, uppercase name lowering, debug `data-ele-id`.
- `component_cached_test.go` — renders exactly once, works inside a render tree, 16-goroutine concurrent render, debug-mode bypass.
- `element_concurrent_test.go` — 16-goroutine rendering in debug mode (exercises the shared concerns map) and with the builder pool; meaningful under `-race`.
- `benchmarks_test.go` — `BenchmarkDivRender[Pooled]`, `BenchmarkComponent[Cached]`.

## Results

`go vet ./...` clean; `go test -race ./...` all pass; all example apps build.

| Benchmark (Apple M1 Pro) | Before | After |
|---|---|---|
| DivRender | 6,623 ns/op, 57 allocs | 573 ns/op, 12 allocs (**11.5×**) |
| Component vs Cached | 879 ns/op, 13 allocs | 44 ns/op, 1 alloc (**20×**) |

(Gating the stack walks alone accounted for 6,623 → 1,368; the attr-slice + direct-write changes took it to 573.)

## Behavior Notes / Compatibility

- Public API unchanged; `Cached()` is additive.
- Attribute output order is now the order passed in (previously random map order). Duplicate attribute keys: last value still wins (same as map behavior), rendered at the key's first position.
- Old test comments about map randomization are now obsolete — exact-string assertions on multi-attribute elements are safe.

## Deferred / Possible Follow-ups

- `Element` (~9 fields) is passed by value into `R()`/`T()` etc.; pointer receivers would change the API feel — not worth it yet.
- `Cached` has no invalidation; add `Invalidate()` if dynamic-but-mostly-static components need it.
- staticcheck QF1012 hints (Fprintf vs WriteString(Sprintf)) remain in a couple of debug-path spots — cosmetic.

## Session Events

- The auto-mode permission classifier briefly went down mid-session (Bash calls blocked); recovered by using Edit/Write tools for reverts and retrying Bash later.
