# Session: go-learn — pluggable tutorial platform (LeetCode-in-Go) build

- **Date:** 2026-07-06 18:00
- **Session ID:** `df41691e-e1c2-4246-ae55-b7216cc91a26`
- **Repo built:** `~/projs/go/go-learn` (NEW, sibling of element; not yet git-initialized at save time)
- **Builds on:** `2026-0706-1357-playground-yaegi-wasm-tutorial.md` and
  `2026-0706-1540-playground-deploy-splitter-pretty.md` (element playground
  architecture is the donor)
- **Plan file:** `~/.claude/plans/piped-swimming-lake.md` (approved)

> New product: **go-learn** — an interactive tutorial platform, first track =
> 15 classic LeetCode problems solved in Go in the browser (yaegi-in-wasm),
> with per-test results, solutions, explanations, and inline SVG diagrams.
> Platform is pluggable: generic engine + track/kind/runner registries;
> a 3-lesson go-basics stub track proves the boundary. Target:
> github.com/rohanthewiz/go-learn → https://rohanthewiz.github.io/go-learn/

## 1. Locked decisions

15 curated classics · create repo + deploy this session (gh authed as repo
owner) · inline SVG diagrams themed by CSS vars · go-basics stub track ships ·
free Playground tab (Go → stdout). 8 categories (not 7): Arrays & Hashing,
Stack, Two Pointers, Sliding Window, Binary Search, Linked List, Trees, DP.

## 2. Architecture (what makes it pluggable)

- **wasm runner is generic**: `goRun.run(src) -> {stdout, stderr, ms} |
  {error, line, col, stderr, ms}`. No notion of tests. Single root Go module
  `github.com/rohanthewiz/go-learn` (no separate wasm module needed — nothing
  to keep yaegi out of). Runner = element's minus srcfs/gopathFS/compat
  (~120 lines deleted); pure stdlib. Symbols list (16 pkgs): + **reflect**
  (DeepEqual in harnesses), container/heap, container/list, math/bits;
  − regexp, runtime, sync, path/filepath, base64. 12 MB wasm / 2.9 MB gz.
- **Test model — JS-side import merge** (`engine/assemble.js`, UMD):
  user edits a full Go file (package/imports/types/func, NO main); track JS
  merges it with the problem's harness file (head-parse both: package clause +
  import decls only — Go grammar guarantees imports precede all other decls —
  dedupe by (alias,path), emit one file). `mapErrorLine` maps interpreter
  errors back to editor lines (regions: user/imports/harness).
  `parseSentinel` splits stdout on the LAST `__GOLEARN_RESULTS__` /
  `__GOLEARN_END__` line pair → user prints (console pane) + JSON
  `[{pass,input,want,got}]`. Spoof-proof because harness prints last.
- **Registries** (`engine/engine.js`): registerRunner / registerKind /
  registerTrack / registerItem / init / wasmReady. Kinds own the output pane:
  `lesson` (stdout + check(stdout, flat)) and `problem` (results table +
  console + solved-on-all-pass) in `engine/kinds.js`. Runner plugin
  `engine/runner-go.js`. Tracks = plain script tags, no build step.
- **Harness Go conventions** (in `tracks/leetcode/track.js` as shared JS
  constants spliced into harnesses): HARNESS_RT (runCase = per-case
  defer/recover so panics still emit full results; emitResults), LIST_HELPERS
  (sliceToList/listToSlice, cycle-capped), TREE_HELPERS (treeFromLevel/
  treeToLevel, []any with nil = LeetCode null encoding). ListNode/TreeNode are
  declared in the USER-visible starter; harness only uses them. Defensive
  copies of slice inputs. Order-insensitive compares normalize both sides
  (group-anagrams, three-sum). Starters always compile (stub returns).
- **Storage**: `golearn:<track>:done|cur|draft:<id>` + global
  `golearn:track|tab|hl|split|tut-split|play-src`.

## 3. Yaegi gotchas (beyond prior sessions')

- Same pin as element: `v0.16.2-0.20260209085605-fcb76d1ece0c`.
- Still NO `min`/`max`/`clear` builtins — all solutions/harnesses use
  explicit ifs (agents were told; held).
- `return nil;;;` is legal Go (empty statements) — don't use it as a
  "broken code" test fixture; use an undefined identifier.
- reflect/encoding/json/sort all work interpreted; full problem run ≈0.3–1 ms
  native.

## 4. Verification (all green at save time)

- `verify/verify.mjs` (runs in CI before deploy): loads tracks by parsing
  index.html script tags (catches missing tags), static shape checks, then
  per problem via mergeProgram + native runner: **starter compiles AND fails
  ≥1 test** (anti-vacuous), **solution passes all, stderr clean**; lessons:
  starter doesn't pre-pass check, solution does. 15 problems + 3 lessons PASS.
- `verify/one.mjs <problem.js>` — single-problem authoring check (agents used
  it; keeps index.html out of the authoring loop).
- assemble.js round-trip suite (scratch): 13 checks incl. alias imports,
  sentinel spoofing, panic recovery, error-line mapping. ALL PASS.
- Browser drives (playwright-core @ scratchpad + system Chrome, server
  `go run ./serve -addr :8090`):
  - quick drive 16/16 (boot, starter-fails table, solution → tick +
    explanation auto-open, error line mapping, console pane, track switch,
    per-track progress, playground stdout, favicon added to kill 404).
  - full drive (in flight at save): 15×(svg+table+failing verdict) ok,
    8 category groups (my check said 7 — content was right), solution-button
    flow ok; remaining steps: for{} timeout, reload persistence, mobile,
    screenshots.
- 5 parallel subagents authored the 14 remaining problems from the two-sum.js
  template; every one self-verified before finishing.

## 5. Deploy plan (next actions — user approved "create repo and deploy")

1. `git init` + initial commit in ~/projs/go/go-learn (master).
2. `gh repo create rohanthewiz/go-learn --public --source . --push`.
3. Pages: `gh api repos/rohanthewiz/go-learn/pages -X POST -f build_type=workflow`.
4. Workflow `.github/workflows/pages.yml` already written: verify job
   (node verify/verify.mjs) → deploy job (build.sh, stage index.html +
   highlight.js + go-learn.wasm + wasm_exec.js + engine/ + tracks/ → _site).
5. Poll `https://api.github.com/repos/rohanthewiz/go-learn/actions/runs`,
   then rerun full drive with URL=https://rohanthewiz.github.io/go-learn/.

## 6. File map (go-learn)

index.html (shell+CSS+boot, ~620 lines) · highlight.js (element's, eleHi→goHi)
· engine/{assemble,engine,kinds,runner-go}.js · tracks/go-basics/track.js ·
tracks/leetcode/track.js + problems/*.js (15) · wasm/{main_js,main_native,
version}.go + runner/runner.go + symbols/ (gen.sh, 16 extracts) ·
verify/{verify,one}.mjs · serve/main.go · build.sh · README.md ·
.github/workflows/pages.yml · .gitignore (go-learn.wasm, wasm_exec.js, _site/,
.claude/).

## 7. ADDENDUM — deployed; web-worker runner became REQUIRED

The full drive did NOT pass on the first architecture: an interpreted
`for {}` wedged the tab permanently. Probes proved **EvalWithContext
timeouts never fire in browser wasm** (timer goroutine starves while the
interpreter spins on the main thread; native is fine) — and the LIVE element
playground has the same bug (its session note claiming for{} cancellation in
wasm was verified wrong by direct probe; fix element later with this same
pattern).

Fix shipped in v1: interpreter moved into a web worker.
- `engine/worker.js` — importScripts wasm_exec.js, instantiates wasm, posts
  {ready, version}; {id,src} → synchronous goRun.run → {result} (FIFO).
- `engine/runner-go.js` — run(src) → Promise; 6s main-thread watchdog →
  worker.terminate() + silent respawn (NO auto-rerun after respawn, or a
  still-broken draft re-wedges in a kill loop); queued runs behind a wedge
  resolve with a restart error; runs requested while (re)spawning queue.
- kinds/engine/playground went async (stale() guard on ctx; runSeq in
  playground); index.html no longer loads wasm on the main thread.

Deployed: repo `rohanthewiz/go-learn` (initial commit, master), Pages
build_type=workflow, `site` workflow (verify → deploy) green, and the full
26-check drive passed **against the live site** — including
infinite-loop-killed + worker-respawn, reload persistence, mobile, no
console errors. https://rohanthewiz.github.io/go-learn/

## 8. Open ideas (post-v1)

Blind-75 expansion (content-only) · hard problems · hidden test cases ·
web-worker runner (kill hung runs instantly) · shareable URLs · lesson
deep-links · quiz kind for a system-design track (runner:'none' — the seam
exists) · gofmt-on-save (needs wasm gofmt).
