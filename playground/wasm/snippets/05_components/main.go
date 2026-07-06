// Components are types with Render(b *element.Builder) any — compose them
// like values, render them anywhere in a tree.
package main

import (
	"fmt"

	"github.com/rohanthewiz/element"
)

type UserBadge struct {
	Name  string
	Admin bool
}

func (u UserBadge) Render(b *element.Builder) (x any) {
	b.SpanClass("user").R(
		b.Strong().T(u.Name),
		b.Wrap(func() {
			if u.Admin {
				b.SmallClass("tag").T(" (admin)")
			}
		}),
	)
	return
}

func main() {
	b := element.NewBuilder()

	team := []UserBadge{
		{Name: "Ada", Admin: true},
		{Name: "Grace"},
		{Name: "Linus"},
	}

	b.DivClass("team").R(
		b.H3().T("Team"),
		b.Ul().R(
			element.ForEach(team, func(u UserBadge) {
				b.Li().R(
					b.RenderComps(u),
				)
			}),
		),
	)

	fmt.Println(b.String())
}
