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

func (s SelectComponent) Render(b *element.Builder) (x any) {
	_, t := b.Vars()

	b.Select().R(
		func() (x any) {
			for _, color := range s.Items {
				params := []string{"value", color}
				if color == s.Selected {
					params = append(params, "selected", "selected")
				}
				b.Option(params...).R(t(color))
			}
			return
		}(),
	)
	return
}

func generateHTML(animals []string, colors []string) string {
	b, e, t := element.Vars()

	selComp := SelectComponent{Selected: "blue", Items: colors}

	b.Html().R(
		b.Head().R(
			b.Style().R(
				t(`
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
                }
            `),
			),
		),
		b.Body().R(
			b.Div("id", "page-container").R(
				b.H1().R(
					t("This is my heading"),
				),
				e("div", "class", "intro").R(
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
					func() (x any) {
						out := ""
						for i := 0; i < 10; i++ {
							if i > 0 {
								out += ","
							}
							out += strconv.Itoa(i)
						}
						return t(out)
					}(),
				),
				b.Div().R(
					t("Lorem Ipsum Lorem Ipsum Lorem<br>Ipsum Lorem Ipsum "),
					e("p").R(t("Finally...")),
				),
				// Iterate over a slice with a built-in function
				// You can actually do more with an inline anonymous function, and components,
				// so consider the For method here deprecated
				e("ul", "class", "list").For(animals, "li"),

				element.RenderComponents(b, selComp),
				b.P().R(), // quick spacer :-)
				e("div", "class", "footer").R(
					t("About | Privacy | Logout")),
			),
		),
	)

	return b.String()
}
