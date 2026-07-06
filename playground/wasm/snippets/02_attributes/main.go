// Attributes are passed as key, value pairs. Single-tag elements like
// <img> and <hr> need no children — close them with .R().
package main

import (
	"fmt"

	"github.com/rohanthewiz/element"
)

func main() {
	b := element.NewBuilder()

	b.DivClass("profile", "id", "card-1").R(
		b.Img("src", "https://go.dev/images/gophers/motorcycle.svg",
			"alt", "Gopher on a motorcycle", "width", "180").R(),
		b.Hr().R(),
		b.A("href", "https://github.com/rohanthewiz/element", "target", "_blank").
			T("element on GitHub"),
		// The *Class builder variants put the class first, then pairs as usual.
		b.PClass("note", "data-lang", "go").T("class-first convenience methods"),
	)

	fmt.Println(b.String())
}
