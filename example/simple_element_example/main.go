package main

import (
	"log"
	"strconv"

	"github.com/rohanthewiz/element"
	"github.com/rohanthewiz/rweb"
	"github.com/rohanthewiz/serr"
)

func main() {
	s := rweb.NewServer(rweb.ServerOptions{
		Address: ":8000",
		Verbose: true,
	})

	s.Use(rweb.RequestInfo) // Stats Middleware

	s.Get("/", rootHandler)
	log.Fatal(s.Run())
}

func rootHandler(c rweb.Context) error {
	animals := []string{"cat", "mouse", "dog"}
	colors := []string{"red", "blue", "green", "indigo", "violet"}
	err := c.WriteHTML(generateHTML(animals, colors))
	if err != nil {
		return serr.Wrap(err)
	}
	return nil
}

type SelectComponent struct {
	Selected string
	Items    []string
}

// Render method signature matches the element.Component interface
func (s SelectComponent) Render(b *element.Builder) (x any) {
	b.Select().R(
		func() (x any) {
			for _, color := range s.Items {
				params := []string{"value", color}
				if color == s.Selected {
					params = append(params, "selected", "selected")
				}
				b.Option(params...).T(color)
			}
			return
		}(),
	)
	return
}

func generateHTML(animals []string, colors []string) string {
	// Set Debug mode and output debug into at end
	// element.DebugSet()
	// defer element.DebugCheck()

	b, e, t := element.Vars()

	selComp := SelectComponent{Selected: "blue", Items: colors}

	b.Html().R(
		b.Head().R(
			b.Style().T(`
                #page-container {
                    padding: 4rem; height: 100vh; background-color: rgb(232, 230, 228);
                }
                .intro {
                    font-style: italic; font-size: 0.9rem; padding-left: 3em;
                }
                .highlight {
                    background-color: yellow;
                }
                .footer {
                    text-align: center; font-size: 0.8rem; border-top: 1px solid #ccc; padding: 1em;
                }`,
			),
		),
		b.Body().R(
			b.Div("id", "page-container").R(
				b.H1().T("This is my heading"),
				e("div", "class", "intro", "unpaired").R( // testing bad pairs
					e("p").R(
						t("I've got plenty to say here "),
						e("span", "class", "highlight").R(
							t("important phrase!", " More intro text"),
						),
					),
				),
				b.P().R(
					t("ABC Company"),
					e("br"), // single tags don't need to call `.R()`
					b.Wrap(func() {
						out := ""
						for i := 0; i < 10; i++ {
							if i > 0 {
								out += ","
							}
							out += strconv.Itoa(i)
						}
						t(out)
					}),
				),
				b.Div().R(
					t("Lorem Ipsum Lorem Ipsum Lorem<br>Ipsum Lorem Ipsum "),
					e("p").T("Finally..."),
				),
				// Iterate over a slice with a built-in function
				// You can actually do more with an inline anonymous function, and components,
				// so consider the For method here deprecated
				e("ul", "class", "list").R(
					element.ForEach(animals, func(animal string) {
						b.Li().T("This is a ", animal, " in a list item")
					}),
				),

				// Render a select component
				element.RenderComponents(b, selComp),
				b.P().R(), // quick spacer :-)
				e("div", "class", "footer").
					T("About | Privacy | Logout"),
			),
		),
	)

	return b.String()
}
