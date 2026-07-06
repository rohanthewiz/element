package main

import (
	"io/fs"
	"runtime/debug"
	"sort"
	"strings"

	"embed"
)

// Each snippet is a standalone, compilable main program — CI builds them
// against the real element package, and the playground interprets the same
// text, so the examples menu can never drift from the library.
//
//go:embed snippets/*/main.go
var snippetsFS embed.FS

// Example is one playground examples-menu entry.
type Example struct {
	Name   string `json:"name"` // directory name minus its ordering prefix, e.g. "hello"
	Source string `json:"source"`
}

// Examples lists the embedded snippets in directory order.
func Examples() ([]Example, error) {
	dirs, err := fs.Glob(snippetsFS, "snippets/*/main.go")
	if err != nil {
		return nil, err
	}
	sort.Strings(dirs)
	list := make([]Example, 0, len(dirs))
	for _, path := range dirs {
		data, err := fs.ReadFile(snippetsFS, path)
		if err != nil {
			return nil, err
		}
		name := strings.TrimSuffix(strings.TrimPrefix(path, "snippets/"), "/main.go")
		// "01_hello" -> "hello"
		if _, rest, ok := strings.Cut(name, "_"); ok {
			name = rest
		}
		list = append(list, Example{Name: name, Source: string(data)})
	}
	return list, nil
}

// version reports the element vcs revision baked into the build, falling
// back to "dev" for local (uncommitted) builds.
func version() string {
	if bi, ok := debug.ReadBuildInfo(); ok {
		for _, s := range bi.Settings {
			if s.Key == "vcs.revision" && len(s.Value) >= 7 {
				return s.Value[:7]
			}
		}
	}
	return "dev"
}
