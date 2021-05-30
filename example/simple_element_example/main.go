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
	err := c.SendString(generateHTML(animals, colors))
	return err
}

func generateHTML(animals []string, colors []string) string {
	s := &strings.Builder{} // This allows us to reduce memory allocations as we build our HTML

	// These anonymous functions nicely wrap our string builder
	// so we don't have to explicitly pass it in to every element call below
	e := func(el string, p ...string) element.Element {
		return element.New(s, el, p...)
	}
	t := func(p ...string) int {
		return element.Text(s, p...)
	}

	/*
		// *The below is a perfect candidate for saving as a snippet / Live Template in your editor / IDE*
		// Place at the top of every function rendering HTML with Element
		s := &strings.Builder{}
		e := func(el string, p ...string) element.Element {
			return element.New(s, el, p...)
		}
		t := func(p ...string) int {
			return element.Text(s, p...)
		}
	*/

	s.WriteString("<!DOCTYPE html>\n")
	e("html", "lang", "en").R(
		e("head").R(
			e("style").R(
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
		e("body").R(
			e("div", "id", "page-container").R(
				e("h1").R(
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
				e("p").R(
					t("ABC Company"),
					e("br"), // single tags don't need to call `.R()`
					func() (x int) {
						out := ""
						for i := 0; i < 10; i++ {
							out += strconv.Itoa(i) + ","
						}
						return t(out)
					}(),
				),
				e("div").R(
					t("Lorem Ipsum Lorem Ipsum Lorem<br>Ipsum Lorem Ipsum "),
					e("p").R(t("Finally...")),
				),
				// Iterate over a slice with a built-in function
				// You can actually do more with an inline anonymous function
				e("ul", "class", "list").For(animals, "li"),

				// Iterate over a slice with an anonymous function - this is very versatile!
				e("select").R(
					func() (x int) {
						for _, color := range colors {
							var el element.Element
							if color == "blue" {
								el = e("option", "value", color, "selected", "selected")
							} else {
								el = e("option", "value", color)
							}
							el.R(t(color))
						}
						return
					}(),
				),
				e("p").R(), // quick spacer :-)
				e("div", "class", "footer").R(t("About | Privacy | Logout")),
			),
		),
	)

	return s.String()
}
