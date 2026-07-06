module github.com/rohanthewiz/element/playground/wasm

go 1.26

require (
	github.com/rohanthewiz/element v0.5.6
	github.com/traefik/yaegi v0.16.2-0.20260209085605-fcb76d1ece0c
)

require github.com/rohanthewiz/serr v1.3.0 // indirect

// The playground always runs against the element checkout it lives in.
replace github.com/rohanthewiz/element => ../..
