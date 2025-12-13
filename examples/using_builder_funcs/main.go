package main

import (
	_ "embed"
	"log"
	"time"

	"github.com/rohanthewiz/element"
	"github.com/rohanthewiz/rweb"
)

// Pickup our styles from the styles.css file
//
//go:embed styles.css
var styleCSS string

const siteName = "Simple Site"

const localStyles = `
body {background-color: #333; color: #fff; font-family: sans-serif;}

div.title {
	font-size: 1.3rem; font-weight: bold; margin-bottom: 1rem; color: #fff;
}

table {
	border-collapse: collapse;
	width: 100%;
	font-family: sans-serif;
	background-color: #222;
	color: #fff;
}

th, td {
	border: 1px solid #666;
	padding: 8px;
	text-align: left;
}

tr.table-head {
	background-color: green;
	font-weight: bold;
}

tr:nth-child(even) {
	background-color: #444;
}
`

func main() {
	s := rweb.NewServer(rweb.ServerOptions{
		Address: ":8000",
		Verbose: true,
	})

	// Middleware
	s.Use(rweb.RequestInfo)

	s.Get("/", rootHandler)
	log.Fatal(s.Run())
}

type ListOfThings struct {
	Name    string
	Things  []string
	Numbers []float64
}

func (l ListOfThings) Render(b *element.Builder) (x any) {
	t := b.Text

	b.Div("class", "card").R( // wrapper
		b.H3().R(t(l.Name)),

		b.Wrap(func() {
			b.Ul().R(
				element.ForEach(l.Things,
					func(item string) {
						b.Li().R(t(item))
					}),
			)
		}),

		b.Wrap(func() {
			if len(l.Numbers) > 0 {
				b.Ul().R(
					element.ForEach(l.Numbers,
						func(number float64) {
							b.Li().R(b.F("%0.2f", number))
						}),
				)
			}
		}),
	)

	return // nil
}

func rootHandler(c rweb.Context) error {
	list := ListOfThings{Name: "Items", Things: []string{"Item A", "Item B", "Item C"}}
	list2 := ListOfThings{Name: "More Items", Things: []string{"Item 1", "Item 2", "Item 3", "Item 4", "Item 5"}}
	list3 := ListOfThings{Name: "Just Numbers", Numbers: []float64{17.0, 5.1, 98.7, 3.1415927}}

	b := element.AcquireBuilder() // get from pool
	defer element.ReleaseBuilder(b)

	b.Html().R(
		b.Head().R(
			b.Title().F("%s", siteName),
			b.Style().T(localStyles, styleCSS), // include styles from local and our styles.css file
		),

		b.Body().R(
			b.Div("class", "title").R(
				b.F("%s", siteName),
			),

			b.DivClass("container").R(
				// Render multiple components
				element.RenderComponents(b, list, list2, list3),
				b.Hr().R(),

				b.P("style", "font-weight:bold").R(
					b.T("Hello there big world!"), b.Br(),
					b.F("%s", time.Now().String()),
				),
				b.Wrap(func() {
					if 2 > 1 {
						b.Span().T("Yup. Two is greater than 1.")
					}
				}),

				b.Aside("style", "display:inline-block;float:right;padding:1rem").
					T("This is an aside!"),
			),

			b.Table().R(
				b.THead().R(
					b.Tr("class", "table-head").R(
						b.Th().T("Header 1"),
						b.Th().T("Header 2"),
						b.Th().T("Header 3"),
					),
				),
				b.TBody().R(
					b.Tr().R(
						b.Td().T("Row 1, Col 1"),
						b.Td().T("Row 1, Col 2"),
						b.Td().T("Row 1, Col 3"),
					),
					b.Tr().R(
						b.Td().T("Row 2, Col 1"),
						b.Td().T("Row 2, Col 2"),
						b.Td().T("Row 2, Col 3"),
					),
				),
			),

			b.Pre().T(`This is a preformatted block of text.
				It will be rendered as a block of text with  line breaks. 
				This is a preformatted block of text. It will be rendered as a block of text with  line breaks.`),

			b.DivClass("footer").F("Copyright &copy; %s &mdash; %s",
				time.Now().Format("2006"), "GodsCoders"),
		),
	)

	return c.WriteHTMLBytes(b.Bytes())
}
