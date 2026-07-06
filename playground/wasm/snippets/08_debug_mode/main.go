// Debug mode catches malformed trees: children passed to single tags,
// text not wrapped in b.T(), and more. Run it and read the report.
package main

import (
	"fmt"

	"github.com/rohanthewiz/element"
)

func main() {
	element.DebugSet()
	defer element.DebugClear()

	b := element.NewBuilder()

	b.Div().R(
		b.H1().T("Spot the bugs"),
		"oops — a bare string, not wrapped in b.T()",
		b.Br().R(
			b.P().T("a <br> cannot have children"),
		),
	)

	fmt.Println(element.DebugShow(element.DebugOptions{TextOnly: true}))
}
