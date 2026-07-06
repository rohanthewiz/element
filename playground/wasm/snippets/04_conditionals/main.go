// Two ways to branch inside a render tree: b.Wrap for statements,
// or an immediately-invoked func() (x any) when you need a child value.
package main

import (
	"fmt"

	"github.com/rohanthewiz/element"
)

func main() {
	b := element.NewBuilder()

	hour := 20
	loggedIn := true

	b.Div().R(
		b.Wrap(func() {
			if hour >= 18 {
				b.H2().T("Good evening!")
			} else {
				b.H2().T("Good day!")
			}
		}),
		func() (x any) {
			if !loggedIn {
				return b.A("href", "/login").T("Sign in")
			}
			return b.P().T("Welcome back.")
		}(),
	)

	fmt.Println(b.String())
}
