package main

import (
	"fmt"
	"log"

	"github.com/rohanthewiz/element"
	"github.com/rohanthewiz/rweb"
)

func main() {
	s := rweb.NewServer(rweb.ServerOptions{
		Address: ":8000",
		Verbose: true,
	})

	s.Use(rweb.RequestInfo) // Stats Middleware

	s.Get("/", func(ctx rweb.Context) (err error) {
		b, _, _ := element.Vars()

		bob := Cat{Color: "yellow", Lazy: true}
		phil := Dog{Color: "brown", Weight: 1.0}

		b.Body().R(
			b.H2().T("Rendering Animals"),
			element.RenderComponents(b, bob),

			b.Div("class", "footer").R(
				b.P().T("I am a footer"),
			),
		)

		AnimalMakeNoise([]Animal{bob, phil})

		return ctx.WriteHTML(b.String())
	})

	s.Get("/debug-set", func(c rweb.Context) error {
		element.DebugSet()
		return c.WriteHTML("<h3>Debug mode set.</h3> <a href='/'>Home</a>")
	})

	s.Get("/debug-clear", func(c rweb.Context) error {
		element.DebugClear()
		return c.WriteHTML("<h3>Debug mode is off.</h3> <a href='/'>Home</a>")
	})

	s.Get("/debug-show", func(c rweb.Context) error {
		err := c.WriteHTML(element.DebugShow())
		return err
	})

	log.Fatal(s.Run())
}

type Cat struct {
	Color  string
	Lazy   bool
	ImgSrc string
}

func (c Cat) MakeNoise() string {
	return "Meow"
}

func (c Cat) Render(b *element.Builder) (x any) {
	b.P().T("Yes I am a <b>cat</b>")
	return
}

// Define Dog
type Dog struct {
	Color  string
	Weight float64
}

func (d Dog) MakeNoise() string {
	return "Woof"
}

// --- End Dog Definition

type Animal interface {
	MakeNoise() string
}

// AnimalMakeNoise A function (not a method)
func AnimalMakeNoise(animal []Animal) {
	for _, a := range animal {
		fmt.Println(a.MakeNoise())
	}
}
