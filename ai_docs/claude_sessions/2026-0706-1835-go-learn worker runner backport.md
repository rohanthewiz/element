# Session: worker runner backported to element playground

- **Date:** 2026-07-06 18:35
- **Session ID:** `df41691e-e1c2-4246-ae55-b7216cc91a26` (same session as
  `2026-0706-1800-go-learn platform build and deploy.md` — this segment is
  the element backport that followed go-learn's deploy)
- **Repo:** `~/projs/go/element` (branch `master`)
- **Commit:** `0878fc5` playground: run the interpreter in a web worker
- **Live & verified:** https://rohanthewiz.github.io/element/

> This segment: ported go-learn's web-worker interpreter pattern into the
> element playground, fixing the tab-freezing `for {}` bug discovered during
> go-learn's final drive. Both sites now share the same runner architecture.

## 1. The bug (recap — full analysis in the 1800 session doc)

`interp.EvalWithContext` timeouts NEVER fire in browser wasm: Go's timer
goroutine starves while the interpreter goroutine spins synchronously on
the main thread (no async preemption on wasm; control never returns to the
JS event loop). Works natively — real threads. An interpreted `for {}`
wedged the element tab permanently; the 2026-0706-1357 session's claim that
for{} was cancelled in wasm was disproven by direct probe against the live
site. Memory saved: `wasm-eval-timeout-needs-worker.md` (now marked FIXED).

## 2. What changed in element (4 files)

- **`playground/worker.js` (new)** — importScripts('wasm_exec.js'),
  instantiates `ele.wasm`, `eleGoReady` → postMessage
  `{type:'ready', version, examples}` (examples must ride the handshake —
  the page can no longer call `eleGo.examples()` directly). `{id, src}` →
  `{type:'result', id, r}`; runs are synchronous in the worker so results
  are FIFO. Paths are worker-relative (flat layout: worker.js sits next to
  wasm_exec.js/ele.wasm both locally and on Pages).
- **`playground/runner.js` (new)** — global **`eleRun`**:
  `run(src) -> Promise`, `isReady()`, `onBoot(cb)` (fires ONCE, first boot
  only), `version()`, `examples()`. 6s watchdog → `worker.terminate()` +
  silent respawn; pending runs behind a wedge resolve with a restart error;
  runs requested while (re)spawning queue and flush on ready. **No
  auto-rerun after respawn** — a still-broken draft would re-wedge in a
  kill/respawn loop.
- **`playground/index.html`** — dropped the `wasm_exec.js` script tag and
  the main-thread `WebAssembly.instantiate` boot; playground `run()` is
  async with a `runSeq` superseded-guard; boot chrome (version, examples
  dropdown, loading fade, `eleTutorial.init()`) moved into
  `eleRun.onBoot`. Also added a data-URI SVG favicon (kills the perpetual
  console 404).
- **`playground/tutorial.js`** — `run()` async: `eleRun.run(...).then(...)`
  guarded by `seq !== runSeq || LESSONS[cur] !== l` (superseded run or
  lesson switched mid-flight). `window.eleGo` readiness check →
  `window.eleRun`.

No Go changes: `main_js.go` still installs `eleGo` + calls `eleGoReady` —
those globals now simply live inside the worker's scope.

## 3. Verification

14-check playwright drive (`drive_element.mjs` pattern in scratchpad; system
Chrome headless), run BOTH locally (`go run ./playground/serve -addr :8091`)
and against the live Pages site after deploy:
boot via worker · default program renders · examples menu ≥8 · version ·
pretty toggle round-trip · **playground for{} killed with "timed out" error
· worker respawned · next run recovers** · tutorial lesson 1 runs · solution
completes + tick · **tutorial for{} killed · reset recovers** · no console
errors. ALL PASS both environments. Deploy workflow green on push (the
`playground/*.js` glob in pages.yml stages runner.js + worker.js
automatically — no workflow change needed).

## 4. Notes for future sessions

- go-learn and element now share the worker-runner architecture but the
  code is copied, not shared (different globals: goRun/eleRun; different
  handshake payloads). If a third playground appears, consider extracting.
- Watchdog is 6s against a wasm-side 5s EvalWithContext that never fires in
  browsers — keep the wasm-side timeout anyway: it's what bounds the NATIVE
  runner (verification harnesses, element lesson checks).
- If element's lesson harness (node) is ever rebuilt: it uses the native
  CLI, unaffected by any of this.
- Element playground follow-ups still open from earlier sessions: gutter
  error markers, shareable URLs, lesson deep-links, forms lesson.
