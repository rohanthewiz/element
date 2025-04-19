package main

import (
	// _ "embed"
	"log"
	"time"

	"github.com/rohanthewiz/element"
	"github.com/rohanthewiz/element/comps"
	"github.com/rohanthewiz/rweb"
)

const siteName = "WORK IN PROGRESS - Data Table Example"

const localStyles = `
div.title {
	font-size: 1.3rem; font-weight: bold; margin-bottom: 1rem; color: #fff;
}`

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

func rootHandler(c rweb.Context) error {
	dt := comps.NewDataTable(comps.DataTableOptions{
		CustomStyles: localStyles,
	})

	b, e, t := element.Vars()

	b.Html().R(
		b.Head().R(
			b.Title().R(
				b.F("%s", siteName),
			), // this is the title of the page tab)), // this is the title of the page tab
			b.Style().R(
				// t can render a list of strings
				t(localStyles),
			),
		),

		b.Body().R(
			e("div", "class", "title").R(
				b.F("%s", siteName),
			),

			b.Div("class", "container").R(

				// Render the data table
				b.RenderComps(dt),

				b.Div("class", "footer").R(
					b.F("Copyright &copy; %s &mdash; %s",
						time.Now().Format("2006"), "GodsCoders"),
				),
			),
		),
	)

	return c.WriteHTML(b.String())
}
