# Simple HTML generator
Forget templates, and having to learn some half-baked templating language to generate decent HTML pages.
This proves that HTML can be generated nicely from Go code. All is explicit and compiler-checked.
Though I may have 100% test coverage, this is alpha, so you probably shouldn't deploy to production without a good html linter (see JetBrains products - I *highly* recommend Goland)

## Usage
Simply create an element and render it: `e("span").R("Inner text")`
We use short method names to keep the code as unobtrusive as possible.
See the example (example/main) for some working code.

```go
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
                e("div", "class", "footer").R("About | Privacy | Logout"),
            ),
        ),
    )

    fmt.Println(str)  // Use a good html viewer to see formatted result
}
```

Produces this:

![element_generated_page](https://user-images.githubusercontent.com/1130495/32986574-dc894b08-cc9a-11e7-82eb-f62fffb84895.png)

The bonus is that our HTML is already somewhat minified to one line so it's very efficient.
Here's what the formatted output looks like:

```html
<!-- Formatted with JetBrains' Goland (Code | Reformat Code) -->
<html>
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
    <div class="intro"><p>I've got plenty to say here <span class="highlight">important phrase!</span> More intro text
    </p></div>
    <div>Lorem Ipsum Lorem Ipsum Lorem<br>Ipsum Lorem Ipsum <p>Finally...</p></div>
    <ul class="list">
        <li>cat</li>
        <li>mouse</li>
        <li>dog</li>
    </ul>
    <div class="footer">About | Privacy | Logout</div>
</div>
</body>
</html>
```

## Style Hints
It Go code man `go fmt` as you please. I do suggest couple things though.

* If you are rendering a short inner text for the element, keep that on one line: `e("span").R("please note")`
* If you are rendering multiple items, especially nested elements, break the render into multiple lines

```go
e("div", "class", "wrapper").R(
	e("div", "class", "inner").R(
		e("p").R("Here is some information")
	)
)
```

## Contributing
Give me some ideas, code, and time :-) if you'd like to see this become better.
The idea is to keep this as **light** and unobtrusive as possible. Thanks!
Also, if possible try to maintain 100% coverage -- again Goland has all the tools needed for test coverage