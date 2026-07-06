# Session: browser playground + interactive tutorial (yaegi-in-wasm)

- **Date:** 2026-07-06 13:57
- **Repo:** `~/projs/go/element/` (module `github.com/rohanthewiz/element`, branch `master`)
- **Sibling work:** mirrors go-styl's playground/tutorial (see go-styl session
  `2026-0706-1250-playground-highlighting-tutorial.md`); reuses its UI bones,
  overlay-editor invariant, and verification recipes.

> This session: a client-side **playground** and 14-lesson **interactive
> tutorial** for element, deployed to GitHub Pages. Since element's API *is*
> Go, the playground embeds the **yaegi interpreter compiled to wasm** with
> element + serr **sources** embedded — user programs are interpreted against
> the real library (generics, components package, debug mode included).
> Commit `a687a17` (43 files). No library changes.

---

## 1. Architecture — how Go runs in the browser

`playground/wasm/` is **its own Go module** (keeps yaegi out of element's
dependency tree; `replace github.com/rohanthewiz/element => ../..`).

- `runner/runner.go` — builds a virtual GOPATH (`fstest.MapFS`,
  `src/github.com/rohanthewiz/{element,serr}/...`) from sources embedded via
  `//go:embed all:srcfs`, then `interp.New` + `EvalWithContext` per run
  (fresh interpreter each time; 5s timeout). Stdout = the playground output.
- `main_js.go` (js&&wasm) — global `eleGo`: `run(src) -> {html, stderr, ms} |
  {error, line, col, stderr, ms}`, `examples()`, `version`; calls
  `eleGoReady()` when instantiated.
- `main_native.go` (!js) — the same runner as a CLI printing one JSON object;
  the lesson harness shells out to it. `-examples` lists snippets.
- `snippets/*/main.go` — 8 examples-menu programs, each a **compilable main
  package** (`go build ./...` in wasm/ keeps them honest against the real
  library) while the playground interprets the same text.
- `symbols/` — trimmed stdlib symbol extracts (18 packages, `gen.sh` to
  regenerate). Full `yaegi/stdlib.Symbols` = 38 MB wasm; trimmed = **12 MB,
  2.9 MB gzipped**, ~15–70 ms per run.
- `build.sh` — `stage` copies element root + components + debug assets +
  serr (from module cache via `go list -m -f '{{.Dir}}'`) into
  `wasm/runner/srcfs/` (gitignored), then builds `ele.wasm` + copies
  `wasm_exec.js`. `build.sh stage` alone is enough for the native runner.

## 2. yaegi facts (hard-won)

- **v0.16.1 is unusable**: a composite literal assigned to a *named return
  value* (`func Mk(...) (e El) { e = El{...}; return e }`) comes back
  **zeroed**. element.New uses exactly that shape → every Element lost its
  buffer pointer. **Fixed on yaegi master**; pinned
  `v0.16.2-0.20260209085605-fcb76d1ece0c` in `playground/wasm/go.mod`.
- yaegi cannot process `//go:embed` → runner inlines
  `assets/debug_table.{js,css}` into `element_debug.go` as quoted string
  literals at FS-build time.
- yaegi lacks Go 1.21 builtins `clear`, `min`, `max` (even on master):
  `clear(con.cmap)` rewritten to a delete loop; `components` gets a
  `zz_yaegi_compat.go` with int `min`/`max` (package-level shims legally
  shadow the predeclared names — all four call sites are int).
- Every rewrite **fails the virtual-FS build loudly** if the target pattern
  disappears, and `elementCompat` rejects any *other* element file that
  starts using embed/clear/min/max — future API drift surfaces immediately.
- `for {}` in user code **is cancelled** by `EvalWithContext` even
  single-threaded in wasm (yaegi checks its done channel at loop
  back-edges) — verified in headless Chrome (~5 s, then a proper error).
- SourcecodeFilesystem + `GoPath: "."` (fs.FS paths can't start with `/`;
  GoPath `/` silently fails to resolve `src/...`).
- Native interp of a full element program: ~13–18 ms. Interpreted components
  (Table with `[][]any` cells, Alert, Card, ProgressBar) and interpreted
  generics (`ForEach`, `ForEach2`) all work.

## 3. Playground UI (`index.html`, `highlight.js`)

Adapted from go-styl's page (same chrome, storage keys `go-ele-*`):

- Left pane Go editor; right pane **html out / preview** sub-tabs — preview
  is a `sandbox="allow-scripts"` iframe fed `srcdoc = stdout`.
- `highlight.js` exposes `eleHi = {go, html, css, escape, editor}`. New Go
  tokenizer (keywords/types/builtins/func-names/package-qualifiers, raw
  strings + block comments carried across lines) and HTML tokenizer
  (doctype, comments, tags/attrs; `<style>` bodies handed to the CSS
  tokenizer — ported from go-styl — with hex swatches in output panes only).
  **Invariant preserved: highlighted HTML text content is character-identical
  to the input** (overlay editor alignment).
- Tab key inserts a real tab (Go); Cmd/Ctrl+Enter forces a run; 300 ms
  debounce (interpretation is heavier than go-styl compiles).
- Gotcha fixed: `switchTab(tab)` must run unconditionally at boot —
  go-styl's `if (tab === 'tutorial')` pattern left `body[data-tab]` unset on
  the play tab.

## 4. Tutorial (`tutorial.js`)

Same engine as go-styl's (nav ✓ ticks, prose column with highlighted
snippets, live editor + output bench with its own html/preview toggle,
reset/solution/open-in-playground, drafts + progress in localStorage,
checkless lessons complete on visit). Checks take `(html, flat)`.

The 14 lessons: hello (execution-order model) → attributes (`*Class`
variants) → single tags → text (`T`/`F`, no escaping caveat) →
ForEach/ForEach2 → conditionals (`Wrap` vs IIFE) → components (interface +
`RenderComps`) → full pages (`Html`/`HtmlPage`) → components.Table →
UI-kit tour (Card/Alert/ProgressBar) → debug mode (fix bugs until
"No element concerns found.") → pooling & caching (`Cached` render counter:
"footer rendered 1 time(s)") → pretty output (`b.Pretty()`) → where-to-next
(rweb + go-styl cross-links).

**JS-template-literal trap:** `\n` inside lesson Go strings must be `\\n`
(it bit the perf lesson's Printf); backticks and `${` can't appear in lesson
code at all.

## 5. Verification (all green)

- **Lesson harness** (scratchpad `verify_lessons.js`, node → native CLI):
  per lesson, starter runs / starter does NOT pre-pass / solution runs AND
  passes. 14/14 after fixing: `_ = score,` in an argument list (statement,
  not expression — moved out) and the `\\n` escape above.
- **Highlighter round-trip** (`verify_highlight.js`): strip-tags+unescape of
  `eleHi.go/html(±swatch)` output == input exactly, for every lesson
  starter/solution/prose block, snippet, and interpreter output — 92/92.
- **Browser drive** (playwright-core + system Chrome, 29 steps): wasm boot,
  default run, overlay identity, tokenized panes, retype→rerun, error
  line:col, **infinite-loop timeout**, 8 examples, table in preview iframe,
  highlight toggle, debug report, 14-lesson nav, live task completion +
  tick, solution button, tutorial preview, open-in-playground, reload
  persistence, `#tutorial` hash, 480 px bench height, dark+light screenshots,
  no console errors.
- Native `-examples` JSON needed json tags on `Example` (Go default caps).

## 6. Deploy & repo notes

- `.github/workflows/pages.yml` — element's **first** workflow: build.sh →
  stage `index.html`, `*.js`, `ele.wasm` → deploy-pages. Trigger branch is
  **master** (element's default — not main like go-styl; initially wrong and
  amended). `go-version-file: playground/wasm/go.mod` (go 1.26 — the symbols
  extracts carry a go1.26 build tag).
- One-time setting (done): Settings → Pages → Source: GitHub Actions.
- `.claude` is **gitignored in this repo** — the new
  `.claude/skills/verify/SKILL.md` (build/stage/drive/harness recipe) stays
  local.
- Artifacts gitignored: `playground/{ele.wasm,wasm_exec.js}`,
  `playground/wasm/runner/srcfs/`.
- `playground/serve` — tiny static server (root module), `-addr :8080`.

## 7. Next-session pointers

- Possible follow-ups: run the interpreter in a **web worker** (kill hung
  runs instantly instead of 5 s timeout; keeps typing responsive), editor
  gutter markers from `{line, col}`, shareable URLs (source in hash),
  lesson deep-links (`#tut/<id>`), a forms lesson (FormField/SelectField),
  gofmt-on-save (yaegi has no formatter — would need a wasm gofmt).
- If element root ever adopts `min`/`max`/`clear`/`//go:embed` outside
  element_debug.go, `runner.elementCompat` must learn the new pattern (it
  will fail the build of the virtual FS with a pointed error).
- Symbols package list ↔ tutorial imports: adding a lesson that imports a
  new stdlib package requires regenerating `symbols/` (gen.sh) — wasm size
  is the tax.
