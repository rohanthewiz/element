# Simple HTML generator
Forget templates, and having to learn some half-baked templating language to generate decent HTML pages.
This proves that HTML can be generated nicely from Go code. All is explicit and compiler-checked.
Though I may have 100% test coverage, this is beta, so you can deploy to production, however you should verify results with a good html linter (see JetBrains products - I *highly* recommend Goland!)

## Usage
Simply create an element and render it: `e("span").R("Inner text")`
We use short method names to keep the code as unobtrusive as possible.
See the example https://github.com/rohanthewiz/element/tree/master/example/simple_element_example for a current, full example app.

```go
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

```

Produces this:

![element_generated_page](https://user-images.githubusercontent.com/1130495/32986574-dc894b08-cc9a-11e7-82eb-f62fffb84895.png)

The bonus is that our HTML is already somewhat minified to one line so it's very efficient.
Here's what the formatted output looks like:

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
        <option value="blue">
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

## Style Hints
It's Go code man, `go fmt` as you please. I do suggest a couple things though.

* If you are rendering a short inner text for the element, keep that on one line: `e(s, "span").R(t(s,"please note"))`
* If you are rendering multiple items, especially nested elements, break the render into multiple lines

## Contributing
If you have ideas, let me know. PRs are welcome, but keep the below in mind.
The idea is to keep this as **light** and unobtrusive as possible. Thanks!
Also, if possible try to maintain at least 95% coverage -- again Goland has all the tools needed for test coverage.