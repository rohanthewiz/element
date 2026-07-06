// The components package ships ready-made UI pieces. Table takes headers
// and rows of any values — strings, numbers, or other components.
package main

import (
	"fmt"

	"github.com/rohanthewiz/element"
	"github.com/rohanthewiz/element/components"
)

func main() {
	b := element.NewBuilder()

	table := components.Table{
		Caption: "Go HTML builders, very scientifically compared",
		Headers: []string{"Library", "Approach", "Speed"},
		Rows: [][]any{
			{"element", "builder over bytes.Buffer", "fast"},
			{"html/template", "parsed templates", "fine"},
			{"string concat", "chaos", "yolo"},
		},
		Striped: true,
		Hover:   true,
	}

	b.RenderComps(table)

	fmt.Println(b.String())
}
