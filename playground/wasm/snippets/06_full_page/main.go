// A complete HTML page — styles, head, body — rendered by one builder.
// Switch the output pane to Preview to see it styled.
package main

import (
	"fmt"

	"github.com/rohanthewiz/element"
)

type Home struct{}

func (Home) Render(b *element.Builder) (x any) {
	b.DivClass("wrap").R(
		b.H1().T("element"),
		b.PClass("tagline").T("Generate HTML with pure Go — fast, simple, no templates."),
		b.Ul().R(
			element.ForEach([]string{
				"Function execution order builds the tree",
				"No reflection in the hot path",
				"Components, pooling, caching",
			}, func(point string) {
				b.Li().T(point)
			}),
		),
	)
	return
}

func main() {
	b := element.NewBuilder()

	page := b.HtmlPage(
		`body { font-family: sans-serif; background: #0e1116; color: #e6edf3;
		        display: grid; place-items: center; min-height: 95vh; }
		 .wrap { max-width: 42rem; padding: 2rem; border: 1px solid #30363d;
		        border-radius: 12px; }
		 .tagline { color: #7ee787; }`,
		"<title>element demo</title>",
		Home{},
	)

	fmt.Println(page)
}
