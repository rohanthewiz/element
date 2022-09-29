# Simple HTML generator
Forget templates, and having to learn some half-baked templating language to generate decent HTML pages.
This proves that HTML can be generated nicely from Go code. All is explicit and compiler-checked.
Though I may have good test coverage, this is beta. You can deploy to production, however you should verify results with a good html linter/checker (see JetBrains products - I *highly* recommend Goland! BTW, this is in production use)

## Usage
Simply create an element and render it: `e("span").R(t("Inner text")) // -> "<span>Inner Text</span>`
(Please see the full example below)

We use short method names and some aliases to keep the code as unobtrusive as possible.
**See the example:** https://github.com/rohanthewiz/element/tree/master/example/simple_element_example for a full, ready-to-compile example app.

```go
package main

import (
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/rohanthewiz/element"
)

func main() {
	r := fiber.New()
	r.Get("/", rootHandler)
	log.Fatal(r.Listen(":8000"))
}

func rootHandler(c *fiber.Ctx) error {
	animals := []string{"cat", "mouse", "dog"}
	colors := []string{"red", "blue", "green", "indigo", "violet"}
	c.Set("Content-Type", "text/html")
	err := c.SendString(generateHTML(animals, colors))
	return err
}

func generateHTML(animals []string, colors []string) string {
	b := element.NewBuilder()
	e := b.Ele
	t := b.Text

	_ = b.WriteString("<!DOCTYPE html>\n")
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
				e("div").R(
					t("Lorem Ipsum Lorem Ipsum Lorem<br>Ipsum Lorem Ipsum "),
					e("p").R(t("Finally...")),
				),
				// Iterate over a slice with a built-in function
				// You can actually do more with an inline anonymous function
				e("ul", "class", "list").For(animals, "li"),

				// Iterate over a slice with an anonymous function - this is very versatile!
				e("select").R(
					func() (x any) {
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

	return b.String()
}
```

Produces this:

![element_generated_page](https://user-images.githubusercontent.com/1130495/32986574-dc894b08-cc9a-11e7-82eb-f62fffb84895.png)

The bonus is that our HTML is already somewhat minified to one line so it's very efficient.
Here's what the formatted output can look like:

```html
<!-- Formatted with JetBrains' Goland (Code | Reformat Code) -->
<!DOCTYPE html>
<html lang="en">
<head>
    <style>
        #page-container {
            padding: 4rem;
            height: 100vh;
            background-color: rgb(232, 230, 228);
        }

        .intro {
            font-style: italic;
            font-size: 0.9rem;
            padding-left: 3em;
        }

        .highlight {
            background-color: yellow;
        }

        .footer {
            text-align: center;
            font-size: 0.8rem;
            border-top: 1px solid #ccc;
            padding: 1em;
        }
    </style>
</head>
<body>
<div id="page-container"><h1>This is my heading</h1>
    <div class="intro"><p>I've got plenty to say here <span class="highlight">important phrase! More intro text</span>
    </p></div>
    <p>ABC Company<br>0,1,2,3,4,5,6,7,8,9,</p>
    <div>Lorem Ipsum Lorem Ipsum Lorem<br>Ipsum Lorem Ipsum <p>Finally...</p></div>
    <ul class="list">
        <li>cat</li>
        <li>mouse</li>
        <li>dog</li>
    </ul>
    <select>
        <option value="red">red</option>
        <option value="blue" selected="selected">blue</option>
        <option value="green">green</option>
        <option value="indigo">indigo</option>
        <option value="violet">violet</option>
    </select>
    <p></p>
    <div class="footer">About | Privacy | Logout</div>
</div>
</body>
</html>
```

## Hints
- Single tag elements (like `br`) don't need to call `.R()`, however most other elements are dual tag and so must call `.R()`
- Use `go fmt` to format go code as normal
- Enjoy the full power and freedom of Go, while generating HTML responses!

## Contributing
If you have ideas, let me know. PRs are welcome, but keep the below in mind.
The idea is to keep this as *light* and unobtrusive as possible. Thanks!
Also, if possible try to maintain at least 95% coverage -- again Goland has all the tools needed for test coverage.