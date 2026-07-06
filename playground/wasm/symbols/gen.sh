#!/bin/sh
# Regenerates the yaegi stdlib symbol extracts in this directory.
# Keep the package list in sync with what element, serr, and the tutorial
# lessons need — every package here grows the wasm binary.
set -eu
cd "$(dirname "$0")"

for pkg in bytes errors fmt math math/rand path/filepath regexp runtime \
	sort strconv strings sync sync/atomic time unicode unicode/utf8 \
	encoding/base64 encoding/json; do
	go run github.com/traefik/yaegi/cmd/yaegi extract "$pkg"
done
