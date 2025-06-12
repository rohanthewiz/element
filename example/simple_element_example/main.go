package main

import (
	"fmt"
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

	// Set debug mode, then go to the page you want to check and refresh it,
	// then go to /debug/show to see any issues (or watch the console)
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

	// Render an HTML fragment -- just bc we can
	fmt.Println(TestRender())

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

func TestRender() string {
	b := element.NewBuilder()

	// Body() like most builder methods writes the *opening* tag to the underlying string builder
	// and are mostly terminated with R(), which allows children to be rendered, and any closing tag appended,
	// or can be terminated with T(args...string), if the children are text-only.
	b.Body().R(
		// DivClass is a convenience method for outputting and element with at least a "class" attribute
		b.DivClass("container").R( // -> Creates `<div class="container">...{span and paragraph children} ...</div>`
			b.Span().T("Some text"), // -> Adds "<span>Some text</span>"

			b.P().R( // -> Adds<p><a href="https://example.com">Example.com</a></p>"
				b.A("href", "http://example.com").T("Example.com"),
			), // ending </p> tag
		), // ending </div> tag
	) // </body>

	return b.String()
	// -> <body><div class="container"><span>Some text</span><p><a href="http://example.com">Example.com</a></p></div></body>
}

// ----- OTHER PAGE -----

// otherPageHandler demonstrates an alternative page construction technique
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
	b.Input("type", "text").R(b.Span().T("I shouldn't be here"))
	b.DivClass("footer").T("About | Privacy | Logout")
	return
}
