# Simple HTML generator
Forget complexity of frontend frameworks and templates to generate decent HTML pages! 
Element generates HTML nicely, and simply from Go code. Everything is pure compiled Go without reflection, annotations or weird rules. This has been in use in production for several years.

- **This library includes AI-agent documentation at `ai_docs/SKILL.md`**.
- Can we make this a standard, so there is **one, easily discoverable, source of truth for a library's SKILL**?

## Benefits

1. You never leave the safety, speed, and familiarity of Go
    - The full power of Go at your finger tips
    - Compiles single-pass with the rest of your Go program -- no weird annotations or build steps
    - You don't have to learn some half-baked templating language
    - You rip at Go speed the whole time
2. Server-side by nature -- not an after thought
3. Zero dependencies, like Node, so you don't add vulnerabilities and complexities into your critical projects!
4. Buffer pools for super-high traffic situations
5. Easy to use with the natural formatting of an HTML tree structure

## Usage

```go
	b := element.NewBuilder()

	// Body(), like most builder methods, writes the *opening* tag to the underlying string builder
	// and is terminated with R(), which allows children to be rendered, then any closing tag appended,
	// or could be terminated with T(args...string), if the children are all literal text.
	b.Body().R(
		// DivClass is a convenience method for outputting and element with at least a "class" attribute
		b.DivClass("container").R( // -> Creates `<div class="container">...{span and paragraph children} ...</div>`
			b.Span().T("Some text"), // -> Adds "<span>Some text</span>"

			b.P().R( // -> Adds<p><a href="https://example.com">Example.com</a></p>"
				b.A("href", "http://example.com").T("Example.com"),
			), // ending </p> tag
		), // ending </div> tag
	) // </body>

	b.String()
	//-> <body><div class="container"><span>Some text</span><p><a href="http://example.com">Example.com</a></p></div></body>
```

See the full examples in the examples folder.

## How it works
- Element maintains an underlying `bytes.Buffer`, to which it appends HTML as you go.
- When you call `b.P()` an opening paragraph tag is immediately added to the string builder.
- `b.P()` returns an element.Element.
- Elements may have children, so we don't render their closing tags as yet. This is where `R()` comes in. 
- `R()` causes the calling function (of the element) to wait until all/any children (lexically, arguments) are resolved (rendered).
    before calling `close()` on the parent element.
- In other words, element leverages the natural order of function and argument execution (a tree, AST) to properly layout HTML elements (also a tree).
- Element therefore is natural Go, not a templating or pseudo language shimmed in, but pure, one-shot compiled Go!
- Also, as everything is written in a single pass with very little memory allocation, it runs at the full speed of Go!

## Core principles of Element
- Element takes a code-first approach to producing HTML.
- First we get a builder from Element: `b := element.NewBuilder()`.
- Note: Builder will render element tags to an internal string builder.
- Then we use the builder to first render the opening tag of the element internally, and return the Element. `b.P("id", "special")` // -> Element
- The element will then close in a few possible ways:
    - Most often we will use the Element method `R()` to
        - *allow any children to render* `b.P("id", "special").R(b.Span().T("We are saying something"))`
        - then close out the current element by writing its closing tag (if any) to the internal string builder
    - If the children of the current tag will be purely literal text, it is optimal to use the method `T()` to
        - *allow the pure text children to render*. Example: `b.P("id", "special").T("one", "two", "three")`
        - then close out the current element by writing its closing tag (if any) to the internal string builder
    - If the children of the current tag will be a single formatted piece text, we should use the method `F()` to
        - *allow the formatted text to render* similiar to `fmt.Sprintf()`. Example: `b.P("id", "special").F("There are %d items", count)`
        - then close out the current element by writing its closing tag (if any) to the internal string builder
- For elements that have no children prefer to close with `R()`, rather than `T()` or `T("")`
- Prefer `b.T()` over `b.Text()` since `b.Text()` is being deprecated.
- When rendering elements, where we include a "class" attribute make sure to use the corresponding builder function ending in "Class".
    - Example: use `b.PClass("intro")` instead of `b.P("class", "intro")`
- When we are done building elements we get the fully rendered HTML with `b.String()`

### Things to Note
- The actual values returned by children elements are ignored.
- `R()` receive arguments of `any` type, but they are discarded
- `T()` like `R()` can terminate an opened element, but `T()` is used when the children are literal text only.
- `F()` can also terminate an opened element with a single formatted string. Example: `b.Span().F("The value is: %0.1f", value)` 
- In debug mode, we do peek at the arguments passed to `R()` to help identify any issue in children elements.

## Examples

### Building with Multiple Functions
Simply pass the builder around (it is a pointer), and create elements in the desired order. Again, things are written immediately to the builder's buffer.

If you are calling directly from a render tree, as in the example below, you will need to return something to satisfy the calling function (in this case the container div's `R()`), but the calling function doesn't do anything with its arguments, they are there just to establish calling order.

```go
package main

func main() {
    b := element.B()
    
    b.DivClass("container").R(
        b.H2().T("Section 1"),
        generateList(b),
    )

    fmt.Println(b.Pretty())
    /* // Output:
    <div class="container">
      <h2>Section 1</h2>
      <ul>
        <li>Item 1</li>
        <li>Item 2</li>
      </ul>
    </div>
    */
}

func generateList(b *element.Builder) (x any) { // we don't care about the return
	// Everything is rendered in the builder immediately!
	b.Ul().R(
		b.Li().T("Item 1"),
		b.Li().T("Item 2"),
	)
	return
}
````

### A Full Example

We use short method names and some aliases to keep the code as unobtrusive as possible.
**See the examples folder** https://github.com/rohanthewiz/element/tree/master/examples for full, ready-to-compile examples.

```go
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

	s.Get("/other-page", otherPageHandler)

	// Set debug mode, then go to the page you want to check, refresh it,
	// then go to /debug/show to see any issues
	s.Get("/debug/set", func(c rweb.Context) error {
		element.DebugSet()
		return c.WriteHTML("<h3>Debug mode set.</h3> <a href='/'>Home</a>")
	})

	s.Get("/debug/show", func(c rweb.Context) error {
		err := c.WriteHTML(element.DebugShow())
		return err
	})

	s.Get("/debug/clear", func(c rweb.Context) error {
		element.DebugClear()
		return c.WriteHTML("<h3>Debug mode is off.</h3> <a href='/'>Home</a>")
	})

	// Run Server
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

// SelectComponent is an example of a component
// note that a component is anything that has a Render method
// taking an `*element.Builder` and returning `any`
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
	b := element.NewBuilder()

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
				b.DivClass("intro").R(), // this should not show any issues
				b.Div("class", "intro", "unpaired").R( // testing bad pairs
					b.P().R(
						b.T("I've got plenty to say here "),
						b.SpanClass("highlight").R(
							b.T("important phrase!", " More intro text"),
						),
					),
				),
				b.P().R(
					b.T("ABC Company"),
					b.Br(), // single tags don't need to call `.R()` or `.T()`, but no harm in calling `.R()` on single tags
					b.Wrap(func() {
						out := ""
						for i := 0; i < 10; i++ {
							if i > 0 {
								out += ","
							}
							out += strconv.Itoa(i)
						}
						b.T(out)
					}),
				),
				b.Div().R(
					b.T("Lorem Ipsum Lorem Ipsum Lorem<br>Ipsum Lorem Ipsum "),
					b.T("Finally..."),
				),
				// Iterate over a slice with a built-in function
				// You can actually do more with an inline anonymous function, and components,
				// so consider the For method here deprecated
				b.UlClass("list").R(
					element.ForEach(animals, func(animal string) {
						b.Li().T("This is a ", animal, " in a list item")
					}),
				),

				// Render a select component
				element.RenderComponents(b, selComp),
				b.P().R(), // quick spacer :-)
				b.DivClass("footer").T("About | Privacy | Logout"),
			),
		),
	)

	out := b.String()
	// Enable below to see the generated HTML
	// fmt.Println(strings.Repeat("-", 60))
	// fmt.Println(out)
	// fmt.Println(strings.Repeat("=", 60))
	return out
}

// ----- OTHER PAGE -----
// demonstrates an alternative (I think better) page construction technique

func otherPageHandler(c rweb.Context) error {
	return c.WriteHTML(otherHTMLPage())
}

func otherHTMLPage() (out string) {
	b := element.NewBuilder()
	b.HtmlPage("body {background-color:#eee;}", "<title>My Other Page</title>", otherBody{})
	return b.String()
}

type otherBody struct{}

func (ob otherBody) Render(b *element.Builder) (x any) {
	b.H1().T("This is my other page")
	b.P().R(
		b.T("This is a simple example of using the Element library to generate HTML."),
	)
	b.DivClass("footer").T("About | Privacy | Logout")
	return
}
```

The bonus is that our HTML is already somewhat minified to one line so it's very efficient.
Here's what the formatted output can look like:
(Note that Element no longer generates data-ele-id's in normal mode.)

```html
<!-- Formatted with JetBrains' Goland (Code | Reformat Code) -->
<!DOCTYPE html>
<html data-ele-id="html-gUr_ZVN3">
<head data-ele-id="head-jXV6jyD-">
    <style data-ele-id="style-jaP1Z7RD">
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
        }</style>
</head>
<body data-ele-id="body-q0-_xsVj">
<div id="page-container" data-ele-id="div-OWVE-Cdh"><h1 data-ele-id="h1-pa6mczX9">This is my heading</h1>
    <div class="intro" data-ele-id="div-yJWOSiuA"></div>
    <div class="intro" data-ele-id="div-cNDBsWT0"><p data-ele-id="p-dqbVmkSN">I've got plenty to say here <span
            class="highlight" data-ele-id="span-_KNmmqqv">important phrase! More intro text</span></p></div>
    <p data-ele-id="p-PvKlD4JZ">ABC Company<br data-ele-id="br--sDnADbX">0,1,2,3,4,5,6,7,8,9</p>
    <div data-ele-id="div-p_jnZHug">Lorem Ipsum Lorem Ipsum Lorem<br>Ipsum Lorem Ipsum Finally...</div>
    <ul class="list" data-ele-id="ul-YP6wxCfB">
        <li data-ele-id="li-c4nMXicM">This is a cat in a list item</li>
        <li data-ele-id="li-XNgvBCVl">This is a mouse in a list item</li>
        <li data-ele-id="li-udTQgyXz">This is a dog in a list item</li>
    </ul>
    <select data-ele-id="select-RtGDjeiY">
        <option data-ele-id="option-z7CHgXy6" value="red">red</option>
        <option selected="selected" data-ele-id="option-u_JNNpLo" value="blue">blue</option>
        <option value="green" data-ele-id="option-b3A3WoJW">green</option>
        <option value="indigo" data-ele-id="option-OrBW4DJk">indigo</option>
        <option value="violet" data-ele-id="option-8nVTX0OZ">violet</option>
    </select>
    <p data-ele-id="p-aVhF4jEp"></p>
    <div class="footer" data-ele-id="div-8tQscer_">About | Privacy | Logout</div>
</div>
</body>
</html>
```

## Hints
- Use Builder to create elements -- this is the new way that comes with good benefits.
- You can create elements directly with Element, but there should be no need to do that now. Using Builder provides more features including convenience (less typing) and great debugging.
- Single tag elements (like `br`) don't need to call `.R()`, however most other elements are dual tag and so must call `.R()`
- Practically, just include `.R()` for all elements unless you are terminating an element with just pure text, in which case you can terminate with `.T()` or `.F()`.
- Use `go fmt` to format go code as normal
- Embrace the full power, safety and speed of Go. Say goodbye to the jungle of frontend frameworks!

## Enabling debugging
- Example uses rweb - `go get github.com/rohanthewiz/rweb`

### Turn Element debugging on

```go
    s := rweb.NewServer(rweb.ServerOptions{
        Address: ":8000",
        Verbose: true,
    })

	s.Get("/debug/set", func(c rweb.Context) error {
		element.DebugSet()
		return c.WriteHTML("<h3>Debug mode set.</h3> <a href='/'>Home</a>")
	})
```

###  Show any issues found
```go
	s.Get("/debug/show", func(c rweb.Context) error {
		err := c.WriteHTML(element.DebugShow())
		return err
	})
```

### Clear Debugging, so full performance is restored

```go
	s.Get("/debug/clear", func(c rweb.Context) error {
		element.DebugClear()
		return c.WriteHTML("<h3>Debug mode is off.</h3> <a href='/'>Home</a>")
	})
```

- See `example/simple_element_example/main.go` for a full example

## Contributing
If you have ideas, let me know, but keep the below in mind.
- The idea is to keep this as *light* and unobtrusive as possible - like Go. Thanks!
- Also, if possible, try to achieve full coverage of any new code added. Goland had all the tools needed for test / coverage long before AI!
