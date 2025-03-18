package main

import (
	"log"

	"github.com/rohanthewiz/element"
	"github.com/rohanthewiz/rweb"
)

func main() {
	s := rweb.NewServer(rweb.ServerOptions{
		Address: ":8000",
		Verbose: true,
	})

	// Middleware
	s.Use(rweb.RequestInfo)

	// Proxy root request to the target server
	s.Get("/", rootHandler)
	log.Fatal(s.Run())
}

func rootHandler(c rweb.Context) error {
	b, e, t := element.Vars()

	b.Html().R(
		b.Head().R(
			b.Title().R(t("Element Check")),
			b.Style().R(
				t(`body {background-color: #333; color: #fff; font-family: sans-serif;}
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
`),
			),
			b.Body().R(
				b.Div("class", "container").R(
					b.H1().R(t("Element Check")),
					b.P("style", "font-weight:bold").R(
						t("Hello there big world!"),
					),
					e("aside", "style", "display:inline-block;float:right").R(t("Not sure what you put in an aside!")),
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
				),
			),
		),
	)

	return c.WriteHTML(b.String())
}
