# element playground

A browser playground + interactive tutorial for
[element](https://github.com/rohanthewiz/element): write a Go program on the
left, see the HTML it prints on the right (as highlighted source or a live
preview) — no server, no toolchain, entirely client-side.

Since element's API *is* Go, the playground embeds the
[yaegi](https://github.com/traefik/yaegi) interpreter compiled to
WebAssembly. The element and serr **sources** ship inside the wasm binary and
are interpreted on demand, so what runs in the page is the real library —
generics, components package, debug mode and all.

## Pieces

- `wasm/` — its own Go module (keeps yaegi out of element's dependency tree).
  - `runner/` — builds a virtual GOPATH from embedded element + serr sources
    (staged into `runner/srcfs/` by `build.sh`; two mechanical rewrites make
    `element_debug.go` interpretable: `//go:embed` assets become string
    literals, and the Go 1.21 `clear`/`min`/`max` builtins yaegi lacks get
    shims) and interprets one user program per call with a 5s timeout.
  - `main_js.go` — the browser API: `eleGo.run(src)` →
    `{html, stderr, ms} | {error, line, col}`, `eleGo.examples()`,
    `eleGo.version`. Stdout is the output.
  - `main_native.go` — the same runner as a CLI (`go run ./wasm < prog.go`),
    used by the verification harness.
  - `snippets/*/main.go` — the examples menu; each is a standalone,
    *compilable* main package, so `go build ./...` in `wasm/` keeps the
    examples honest against the real library.
  - `symbols/` — trimmed stdlib symbol extracts for the interpreter
    (regenerate with `symbols/gen.sh`; the full yaegi stdlib would triple
    the binary).
- `index.html` — the page: Playground + Tutorial tabs, overlay editors,
  html/preview output panes, token palettes for dark & light themes.
- `highlight.js` — dependency-free tokenizers (Go, HTML, CSS) + the overlay
  editor. Invariant: highlighted HTML's text content is character-identical
  to the input, so the transparent-textarea overlay stays column-aligned.
- `tutorial.js` — 14 interactive lessons with live checks; lesson data is
  Node-requirable for batch verification.
- `serve/` — tiny static file server for local development.

## Build & run locally

```sh
./playground/build.sh          # stages srcfs, builds ele.wasm (~12 MB, ~3 MB gzipped)
go run ./playground/serve      # http://localhost:8080
```

`./playground/build.sh stage` refreshes only the embedded sources — enough
for the native runner: `cd playground/wasm && go run . < prog.go`.

## Deploy

`.github/workflows/pages.yml` builds and publishes to GitHub Pages on pushes
to `main` that touch the library or the playground. One-time setup:
repo **Settings → Pages → Build and deployment → Source: GitHub Actions**.

## Interpreter limits

- yaegi is pinned past v0.16.1 (a named-return bug there zeroed every
  `element.New` result); the pin is in `wasm/go.mod`.
- The Go 1.21+ builtins `min`, `max`, `clear` are unavailable in user code.
- Available stdlib in user programs: bytes, errors, fmt, math, math/rand,
  path/filepath, regexp, runtime, sort, strconv, strings, sync, sync/atomic,
  time, unicode, unicode/utf8, encoding/base64, encoding/json
  (see `symbols/gen.sh`).
- A program runs for at most 5 seconds — runaway loops are cancelled.
