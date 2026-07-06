// Package runner executes user Go programs that import
// github.com/rohanthewiz/element (and .../element/components) using the
// yaegi interpreter. The element and serr sources are embedded at build time
// (see playground/build.sh, which stages them into srcfs/) and exposed to
// the interpreter through a virtual GOPATH filesystem, so the playground
// needs no network and no Go toolchain at runtime — it works the same
// compiled to WebAssembly in a browser as it does natively.
package runner

import (
	"bytes"
	"context"
	"embed"
	"fmt"
	"io/fs"
	"strconv"
	"strings"
	"sync"
	"testing/fstest"
	"time"

	"github.com/traefik/yaegi/interp"

	"github.com/rohanthewiz/element/playground/wasm/symbols"
)

// srcfs is staged by playground/build.sh:
//
//	srcfs/element/*.go + assets/debug_table.{js,css}
//	srcfs/element/components/*.go
//	srcfs/serr/*.go
//
//go:embed all:srcfs
var rawFS embed.FS

const (
	elementMod    = "github.com/rohanthewiz/element"
	componentsMod = "github.com/rohanthewiz/element/components"
	serrMod       = "github.com/rohanthewiz/serr"
)

var (
	gopathOnce sync.Once
	gopath     fstest.MapFS
	gopathErr  error
)

// gopathFS lazily assembles the virtual GOPATH the interpreter resolves
// imports against. Built once; yaegi only reads from it.
func gopathFS() (fstest.MapFS, error) {
	gopathOnce.Do(func() {
		gopath, gopathErr = buildGopathFS()
	})
	return gopath, gopathErr
}

func buildGopathFS() (fstest.MapFS, error) {
	m := fstest.MapFS{}
	stage := func(from, mod string, transform func(name string, data []byte) ([]byte, error)) error {
		entries, err := fs.ReadDir(rawFS, from)
		if err != nil {
			return fmt.Errorf("srcfs missing %s (run playground/build.sh stage): %w", from, err)
		}
		for _, e := range entries {
			name := e.Name()
			if e.IsDir() || !strings.HasSuffix(name, ".go") || strings.HasSuffix(name, "_test.go") {
				continue
			}
			data, err := fs.ReadFile(rawFS, from+"/"+name)
			if err != nil {
				return err
			}
			if transform != nil {
				if data, err = transform(name, data); err != nil {
					return fmt.Errorf("%s/%s: %w", from, name, err)
				}
			}
			m["src/"+mod+"/"+name] = &fstest.MapFile{Data: data}
		}
		return nil
	}

	if err := stage("srcfs/element", elementMod, elementCompat); err != nil {
		return nil, err
	}
	if err := stage("srcfs/element/components", componentsMod, nil); err != nil {
		return nil, err
	}
	// yaegi (through v0.16.x master) has no min/max builtins (Go 1.21).
	// The components package uses them on ints only; a package-level shim
	// shadows the predeclared names, which is legal Go.
	m["src/"+componentsMod+"/zz_yaegi_compat.go"] = &fstest.MapFile{Data: []byte(
		"package components\n\n" +
			"func min(a, b int) int {\n\tif a < b {\n\t\treturn a\n\t}\n\treturn b\n}\n\n" +
			"func max(a, b int) int {\n\tif a > b {\n\t\treturn a\n\t}\n\treturn b\n}\n")}
	if err := stage("srcfs/serr", serrMod, nil); err != nil {
		return nil, err
	}
	return m, nil
}

// elementCompat rewrites element sources for interpretation. yaegi cannot
// process //go:embed directives or the Go 1.21 `clear` builtin, both used
// only in element_debug.go. Each rewrite must apply — a silent no-op would
// surface later as a confusing interpreter error, so unmatched patterns fail
// the build of the virtual filesystem instead.
func elementCompat(name string, data []byte) ([]byte, error) {
	if name != "element_debug.go" {
		if bytes.Contains(data, []byte("//go:embed")) || containsBuiltin(data) {
			return nil, fmt.Errorf("uses //go:embed, clear, min, or max — extend elementCompat")
		}
		return data, nil
	}
	s := string(data)

	replace := func(old, new, what string) error {
		replaced := strings.Replace(s, old, new, 1)
		if replaced == s {
			return fmt.Errorf("expected %s not found — element_debug.go changed, update elementCompat", what)
		}
		s = replaced
		return nil
	}

	if err := replace("\t_ \"embed\"\n", "", `import _ "embed"`); err != nil {
		return nil, err
	}
	for _, embedded := range []struct{ file, varName string }{
		{"assets/debug_table.js", "tableJS"},
		{"assets/debug_table.css", "tableCSS"},
	} {
		asset, err := fs.ReadFile(rawFS, "srcfs/element/"+embedded.file)
		if err != nil {
			return nil, fmt.Errorf("embedded asset %s: %w", embedded.file, err)
		}
		err = replace(
			"//go:embed "+embedded.file+"\nvar "+embedded.varName+" string",
			"var "+embedded.varName+" = "+strconv.Quote(string(asset)),
			"//go:embed "+embedded.file)
		if err != nil {
			return nil, err
		}
	}
	if err := replace("\tclear(con.cmap)",
		"\tfor k := range con.cmap {\n\t\tdelete(con.cmap, k)\n\t}",
		"clear(con.cmap)"); err != nil {
		return nil, err
	}
	return []byte(s), nil
}

// containsBuiltin reports whether source calls a builtin yaegi lacks.
// Line-based and crude, but only guards library sources we control.
func containsBuiltin(data []byte) bool {
	for _, pat := range [][]byte{[]byte("\tclear("), []byte(" clear("), []byte("= min("), []byte("= max(")} {
		if bytes.Contains(data, pat) {
			return true
		}
	}
	return false
}

// Result is the outcome of interpreting one user program.
type Result struct {
	Stdout string  // whatever the program printed — the playground treats it as HTML
	Stderr string  // interpreter panic traces and any user writes to stderr
	Ms     float64 // wall-clock interpretation time in milliseconds
	Err    error   // nil on success
	// Best-effort position of Err within the user's source (1-based; 0 if unknown).
	Line, Col int
}

// Run interprets src (a complete Go main program) with a fresh interpreter.
// The context bounds runaway programs — pass a timeout for untrusted input.
func Run(ctx context.Context, src string) Result {
	gofs, err := gopathFS()
	if err != nil {
		return Result{Err: err}
	}

	var stdout, stderr bytes.Buffer
	i := interp.New(interp.Options{
		GoPath:               ".",
		SourcecodeFilesystem: gofs,
		Stdout:               &stdout,
		Stderr:               &stderr,
	})
	if err := i.Use(symbols.Symbols); err != nil {
		return Result{Err: err}
	}

	start := time.Now()
	_, err = i.EvalWithContext(ctx, src)
	res := Result{
		Stdout: stdout.String(),
		Stderr: stderr.String(),
		Ms:     float64(time.Since(start).Microseconds()) / 1000,
		Err:    err,
	}
	if err != nil {
		res.Line, res.Col = errPosition(err.Error())
	}
	return res
}

// errPosition extracts a leading "line:col:" position from a yaegi error
// message. Errors pointing inside library sources ("src/github.com/...")
// carry no user position.
func errPosition(msg string) (line, col int) {
	head, _, ok := strings.Cut(msg, ": ")
	if !ok {
		return 0, 0
	}
	l, c, ok := strings.Cut(head, ":")
	if !ok {
		return 0, 0
	}
	ln, err1 := strconv.Atoi(l)
	cn, err2 := strconv.Atoi(c)
	if err1 != nil || err2 != nil || ln < 1 || cn < 1 {
		return 0, 0
	}
	return ln, cn
}
