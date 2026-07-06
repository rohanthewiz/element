// Hello, element! Build HTML with plain Go — no templates.
// Whatever this program prints becomes the playground's output.
package main

import (
	"fmt"

	"github.com/rohanthewiz/element"
)

func main() {
	b := element.NewBuilder()

	b.DivClass("greeting").R(
		b.H1().T("Hello from element"),
		b.P().R(
			b.T("HTML is built by "),
			b.Em().T("running Go"),
			b.T(" — arguments evaluate in order, so the tree just falls out."),
		),
	)

	fmt.Println(b.String())
}
