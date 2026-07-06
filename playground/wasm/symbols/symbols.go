// Package symbols holds a trimmed set of stdlib symbols for the yaegi
// interpreter — only the packages exposed to playground user code, plus
// what element and serr themselves import. Using the full
// yaegi/stdlib.Symbols would triple the wasm binary size.
//
// Regenerate the extracted files with gen.sh after changing the package list.
package symbols

import "reflect"

// Symbols is the map the generated files populate via init().
var Symbols = map[string]map[string]reflect.Value{}
