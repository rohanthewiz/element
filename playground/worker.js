/* worker.js — hosts the Go interpreter off the main thread.
 *
 * Why a worker: eleGo.run is synchronous, and an interpreted infinite loop
 * (`for {}`) spins forever. The wasm-side EvalWithContext timeout cannot
 * fire in a browser — Go's timer goroutine never gets scheduled while the
 * interpreter goroutine spins without yielding (it works natively, where
 * real threads exist). In a worker the main thread stays responsive and,
 * on timeout, runner.js simply terminate()s this worker and spawns a fresh
 * one — true preemption no in-wasm mechanism can provide.
 *
 * Protocol (see runner.js):
 *   -> {id, src}                              run request
 *   <- {type:'ready', version, examples}      wasm instantiated, eleGo installed
 *   <- {type:'result', id, r}                 r = eleGo.run(src) result object
 *   <- {type:'fatal', error}                  wasm failed to load
 */
importScripts('wasm_exec.js');

var go = new Go();

self.eleGoReady = function () {
	postMessage({
		type: 'ready',
		version: (self.eleGo && self.eleGo.version) || 'dev',
		examples: self.eleGo ? self.eleGo.examples() : [],
	});
};

// Prefer streaming; fall back for servers that mislabel .wasm's MIME type.
var load = WebAssembly.instantiateStreaming
	? WebAssembly.instantiateStreaming(fetch('ele.wasm'), go.importObject)
		.catch(function () {
			return fetch('ele.wasm').then(function (r) { return r.arrayBuffer(); })
				.then(function (buf) { return WebAssembly.instantiate(buf, go.importObject); });
		})
	: fetch('ele.wasm').then(function (r) { return r.arrayBuffer(); })
		.then(function (buf) { return WebAssembly.instantiate(buf, go.importObject); });

load.then(function (r) { go.run(r.instance); })
	.catch(function (err) { postMessage({ type: 'fatal', error: String(err) }); });

self.onmessage = function (e) {
	// Synchronous run: while this executes, further messages queue in the
	// worker's event loop (FIFO), so results always arrive in send order.
	postMessage({ type: 'result', id: e.data.id, r: self.eleGo.run(e.data.src) });
};
