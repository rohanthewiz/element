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

	b := element.B()

	b.Html().R(
		b.Head().R(
			b.Title().R(
				b.F("%s", siteName),
			), // this is the title of the page tab)), // this is the title of the page tab
			b.Style().R(
				// t can render a list of strings
				t(localStyles, styleCSS), // include styles from local and our styles.css file
			),
		),

		b.Body().R(
			e("div", "class", "title").R(
				b.F("%s", siteName),
			),

			b.Div("class", "container").R(
				// Show some lists
				element.RenderComponents(b, list, list2, list3),
				b.P("style", "font-weight:bold").R(
					t("Hello there big world!"), b.Br().R(),
					b.F("%s", time.Now().String()),
				),
				b.Wrap(func() {
					if 2 > 1 {
						t("Yup. Two is greater than 1.")
					}
				}),
				b.Aside("style", "display:inline-block;float:right;padding:1rem").R(
					t("This is an aside!")),

				b.Table().R(
					b.THead().R(
						b.Tr("class", "table-head").R(
							b.Th().R(t("Header 1")),
							b.Th().R(t("Header 2")),
							b.Th().R(t("Header 3")),
						),
					),
					b.TBody().R(
						b.Tr().R(
							b.Td().R(t("Row 1, Col 1")),
							b.Td().R(t("Row 1, Col 2")),
							b.Td().R(t("Row 1, Col 3")),
						),
						b.Tr().R(
							b.Td().R(t("Row 2, Col 1")),
							b.Td().R(t("Row 2, Col 2")),
							b.Td().R(t("Row 2, Col 3")),
						),
					),
				),
				b.Pre().R(
					t(`This is a preformatted block of text.
    It will be rendered as a block of text with  line breaks. 
    This is a preformatted block of text. It will be rendered as a block of text with  line breaks.`),
				),
				b.Div("class", "footer").R(
					b.F("Copyright &copy; %s &mdash; %s",
						time.Now().Format("2006"), "GodsCoders"),
				),
			),
		),
	)

	return c.WriteHTML(b.String())
}
