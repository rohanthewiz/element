package main


import (
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/rohanthewiz/element"
)

func main() {
	r := fiber.New()
	r.Get("/", rootHandler)
	r.Listen(":8000")
}

func rootHandler(c *fiber.Ctx) error {
	animals := []string{"cat", "mouse", "dog"}
	colors := []string{"red", "blue", "green", "indigo", "violet"}
	c.Set("Content-Type", "text/html")
	err := c.SendString(generateTemplate(animals, colors))
	return err
}

func generateTemplate(animals []string, colors []string) string {
	e := element.New                           // to keep things unobtrusive
	t := element.Text
	s := &strings.Builder{}
	s.WriteString("<!DOCTYPE html>\n")
	e(s, "html", "lang", "en").R(
		e(s, "head").R(
			e(s, "style").R(
				t(s, `
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
		e(s, "body").R(
			e(s, "div", "id", "page-container").R(
				e(s, "h1").R(
					t(s, "This is my heading"),
				),
				e(s, "div", "class", "intro").R(
					e(s, "p").R(
						t(s, "I've got plenty to say here "),
						e(s, "span", "class", "highlight").R(
							t(s, "important phrase!", " More intro text"),
						),
					),
				),
				e(s, "p").R(
					t(s, "ABC Company"), e(s, "br").R(),
					func() (el element.Element) {
						out := ""
						for i := 0; i < 10; i++ {
							out += strconv.Itoa(i) + ","
						}
						return t(s, out)
					}(),
				),
				e(s, "div").R(
					t(s, "Lorem Ipsum Lorem Ipsum Lorem<br>Ipsum Lorem Ipsum "),
					e(s, "p").R(t(s, "Finally...")),
				),
				// Iterate over a slice with some built-in functions
				e(s, "ul", "class", "list").For(animals, "li").R(),

				// Iterate over a slice with an anonymous function - this is very versatile!
				e(s, "select").R(
					func() element.Element {
						for _, color := range colors {
							el :=  e(s, "option", "value", color)
							if color == "blue" {
								el = e(s, "option", "value", color, "selected", "selected")
							}
							el.R(t(s, color))
						}
						return t(s, "") // bogus as element is not used in R()
					}(),
				),
				e(s, "p").R(), // quick spacer :-)
				e(s, "div", "class", "footer").R(t(s, "About | Privacy | Logout")),
			),
		),
	)

	return s.String()
}
