// Render slices with ForEach / ForEach2 — ordinary Go generics,
// right inside the render tree.
package main

import (
	"fmt"

	"github.com/rohanthewiz/element"
)

func main() {
	b := element.NewBuilder()

	fruits := []string{"apple", "banana", "cherry"}
	langs := []string{"Go", "Zig", "Rust"}

	b.Div().R(
		b.H3().T("Fruits"),
		b.Ul().R(
			element.ForEach(fruits, func(f string) {
				b.Li().T(f)
			}),
		),
		b.H3().T("Languages, ranked"),
		b.Ol().R(
			element.ForEach2(langs, func(lang string, i int) {
				b.Li().F("#%d %s", i+1, lang)
			}),
		),
	)

	fmt.Println(b.String())
}
