//go:build !(js && wasm)

// Native build of the playground runner — the same interpreter the browser
// uses, as a CLI. The tutorial verification harness shells out to this to
// batch-check every lesson without a browser:
//
//	go run . [-timeout 5s] [file.go]   # source from the file, or stdin
//	go run . -examples                 # list embedded snippets as JSON
//
// Output is one JSON object mirroring eleGo.run's result:
// {"html", "stderr", "ms"} on success, {"error", "line", "col", ...} on failure.
package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/rohanthewiz/element/playground/wasm/runner"
)

func main() {
	timeout := flag.Duration("timeout", 5*time.Second, "interpretation time limit")
	examples := flag.Bool("examples", false, "list embedded snippets as JSON and exit")
	flag.Parse()

	enc := json.NewEncoder(os.Stdout)
	if *examples {
		list, err := Examples()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		_ = enc.Encode(list)
		return
	}

	var src []byte
	var err error
	if flag.NArg() > 0 {
		src, err = os.ReadFile(flag.Arg(0))
	} else {
		src, err = io.ReadAll(os.Stdin)
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	ctx, cancel := context.WithTimeout(context.Background(), *timeout)
	defer cancel()
	res := runner.Run(ctx, string(src))

	out := map[string]any{"ms": res.Ms, "stderr": res.Stderr}
	if res.Err != nil {
		out["error"] = res.Err.Error()
		if res.Line > 0 {
			out["line"] = res.Line
			out["col"] = res.Col
		}
		_ = enc.Encode(out)
		os.Exit(2)
	}
	out["html"] = res.Stdout
	_ = enc.Encode(out)
}
