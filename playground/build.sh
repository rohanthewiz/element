#!/bin/sh
# Builds the WASM playground into this directory:
#   playground/ele.wasm      (yaegi runner + embedded element source, js/wasm)
#   playground/wasm_exec.js  (Go's JS support shim, copied from GOROOT)
#
# `build.sh stage` only refreshes wasm/runner/srcfs (the element + serr
# sources the interpreter reads) — enough for native `go run ./wasm`.
#
# Serve the directory with any static file server, e.g.:
#   go run ./playground/serve       # tiny bundled server on :8080
set -eu
cd "$(dirname "$0")"

# --- stage: copy interpretable sources into the embed tree ------------------
srcfs=wasm/runner/srcfs
rm -rf "$srcfs"
mkdir -p "$srcfs/element/assets" "$srcfs/element/components" "$srcfs/serr"

for f in ../*.go; do
  case "$f" in *_test.go) continue ;; esac
  cp "$f" "$srcfs/element/"
done
for f in ../components/*.go; do
  case "$f" in *_test.go) continue ;; esac
  cp "$f" "$srcfs/element/components/"
done
cp ../assets/debug_table.js ../assets/debug_table.css "$srcfs/element/assets/"

serr_dir=$(cd wasm && go list -m -f '{{.Dir}}' github.com/rohanthewiz/serr)
for f in "$serr_dir"/*.go; do
  case "$f" in *_test.go) continue ;; esac
  cp "$f" "$srcfs/serr/"
done
chmod -R u+w "$srcfs"   # module cache files are read-only

[ "${1:-}" = "stage" ] && { echo "staged $srcfs"; exit 0; }

# --- build -------------------------------------------------------------------
(cd wasm && GOOS=js GOARCH=wasm go build -trimpath -ldflags='-s -w' -o ../ele.wasm .)

goroot=$(go env GOROOT)
if [ -f "$goroot/lib/wasm/wasm_exec.js" ]; then     # Go >= 1.24
  cp "$goroot/lib/wasm/wasm_exec.js" .
else                                                # Go <= 1.23
  cp "$goroot/misc/wasm/wasm_exec.js" .
fi

ls -lh ele.wasm
