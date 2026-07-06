//go:build js && wasm

// Command wasm is the WebAssembly build of the element playground runner.
// It installs a global `eleGo` object with:
//
//	eleGo.run(src) ->
//	    {html, pretty, stderr, ms} | {error, line, col, stderr, ms}
//	eleGo.examples() -> [{name, source}]
//	eleGo.version -> short vcs revision or "dev"
//
// `pretty` is `html` re-indented with element.PrettyHTML — display sugar;
// previews and lesson checks should keep using the raw `html`.
//
// `src` is a complete Go main program; whatever it writes to stdout is
// returned as `html`. User code may import github.com/rohanthewiz/element
// and .../element/components — both are interpreted from embedded source.
package main

import (
	"context"
	"syscall/js"
	"time"

	"github.com/rohanthewiz/element"
	"github.com/rohanthewiz/element/playground/wasm/runner"
)

// runTimeout bounds a single interpretation so a runaway user loop
// (`for {}`) doesn't wedge the page forever.
const runTimeout = 5 * time.Second

func main() {
	g := js.Global()
	api := g.Get("Object").New()
	api.Set("run", js.FuncOf(run))
	api.Set("examples", js.FuncOf(listExamples))
	api.Set("version", version())
	g.Set("eleGo", api)
	// Tell the page the API is ready (wasm instantiation is async).
	if cb := g.Get("eleGoReady"); cb.Type() == js.TypeFunction {
		cb.Invoke()
	}
	select {} // keep the Go runtime alive for future JS calls
}

// run implements eleGo.run(src).
func run(_ js.Value, args []js.Value) any {
	if len(args) < 1 || args[0].Type() != js.TypeString {
		return map[string]any{"error": "run(src): src must be a string"}
	}
	ctx, cancel := context.WithTimeout(context.Background(), runTimeout)
	defer cancel()

	res := runner.Run(ctx, args[0].String())
	out := map[string]any{"ms": res.Ms, "stderr": res.Stderr}
	if res.Err != nil {
		out["error"] = res.Err.Error()
		if res.Line > 0 {
			out["line"] = res.Line
			out["col"] = res.Col
		}
		return out
	}
	out["html"] = res.Stdout
	out["pretty"] = element.PrettyHTML(res.Stdout)
	return out
}

// listExamples implements eleGo.examples().
func listExamples(_ js.Value, _ []js.Value) any {
	examples, err := Examples()
	if err != nil {
		return []any{}
	}
	list := make([]any, 0, len(examples))
	for _, ex := range examples {
		list = append(list, map[string]any{"name": ex.Name, "source": ex.Source})
	}
	return list
}
