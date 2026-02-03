package definition_list

import "github.com/rohanthewiz/element"

// -----------------------------------------------------------------------------
// Definition List Component
// -----------------------------------------------------------------------------

// Definition represents a term and its definition.
type Definition struct {
	Term       string // The term being defined
	Definition string // The definition
}

// DefinitionList renders a definition list (<dl>) element.
type DefinitionList struct {
	Items []Definition // Term/definition pairs
	Class string       // Optional CSS class
}

// Render implements the element.Component interface.
func (dl DefinitionList) Render(b *element.Builder) (x any) {
	dlClass := "definition-list"
	if dl.Class != "" {
		dlClass += " " + dl.Class
	}

	b.DlClass(dlClass).R(
		func() (x any) {
			for _, item := range dl.Items {
				b.DtClass("definition-term").T(item.Term)
				b.DdClass("definition-desc").T(item.Definition)
			}
			return
		}(),
	)
	return
}
