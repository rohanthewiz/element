package main

import (
	"fmt"
	"github.com/rohanthewiz/element"
)

func main() {
	e := element.New  // to keep things unobtrusive
	animals := []string{"cat", "mouse", "dog"}  // just an ordinary Go slice

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
					e("p").R("I've got plenty to say here ",
						e("span", "class", "highlight").R("important phrase!"),
						" More intro text",
					),
				),
				e("div").R(
					"Lorem Ipsum Lorem Ipsum Lorem<br>Ipsum Lorem Ipsum ",
					e("p").R("Finally..."),
				),
				e("ul", "class", "list").I(animals, e("li")), // Iterate my slice - move over Angular!
				e("div", "class", "footer").R("About | Privacy | Logout"),
			),
		),
	)

	fmt.Println(str)  // Use a good html viewer to see formatted result
	// -- I suggest JetBrains Goland ( Code | Reformat Code)
}
