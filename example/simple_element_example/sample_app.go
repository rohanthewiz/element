package main

import (
	"github.com/labstack/echo"
	"github.com/rohanthewiz/element"
)

func main() {
	router := echo.New()
	router.GET("/", rootHandler)
	router.Logger.Fatal(router.Start(":8000"))
}

func rootHandler(c echo.Context) error {
	animals := []string{"cat", "mouse", "dog"} // just an ordinary Go slice
	c.HTMLBlob(200, generateTemplate(animals))
	return nil
}

func generateTemplate(animals []string) []byte {
	e := element.New                           // to keep things unobtrusive
	food := []string{"chicken", "bread", "cheese"}

	str := e("html").R(
		e("head").R(
			e("style").R(`
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
		e("body").R(
			e("div", "id", "page-container").R(
				e("h1").R("This is my heading"),
				e("div", "class", "intro").R(
					e("p").R(
						"I've got plenty to say here ",
						e("span", "class", "highlight").R("important phrase!"),
						" More intro text",
					),
				),
				e("div").R(
					"Lorem Ipsum Lorem Ipsum Lorem<br>Ipsum Lorem Ipsum ",
					e("p").R("Finally..."),
				),
				e("ul", "class", "list").For(animals, "li"), // Iterate my slice - move over Angular!
				e("ul", "class", "conditional-list-false").ForIf(false, food, "li"),
				e("ul", "class", "conditional-list-true").ForIf(true, food, "li"),
				e("div", "class", "footer").R("About | Privacy | Logout"),
			),
		),
	)

	return []byte(str)
}
