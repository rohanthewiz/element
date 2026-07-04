package components

import "github.com/rohanthewiz/element"

// -----------------------------------------------------------------------------
// Accordion Component
// -----------------------------------------------------------------------------

// AccordionItem is one expandable section of an Accordion.
type AccordionItem struct {
	Title   string            // Section title (the clickable summary)
	Content string            // Simple text content
	Body    element.Component // Alternative: render a component as the section body
	Open    bool              // Whether the section starts expanded
}

// Accordion renders expandable sections using native <details>/<summary> —
// no JavaScript required. With Exclusive set, all sections share a name
// attribute so the browser keeps at most one open at a time.
type Accordion struct {
	Items     []AccordionItem
	Exclusive bool   // Only one section open at a time (native <details name=...>)
	Name      string // Group name used when Exclusive (default "accordion")
	Class     string // Additional CSS classes for each section
}

// Render implements the element.Component interface.
func (a Accordion) Render(b *element.Builder) (x any) {
	itemClass := "accordion-item"
	if a.Class != "" {
		itemClass = a.Class + " " + itemClass
	}
	groupName := a.Name
	if groupName == "" {
		groupName = "accordion"
	}

	b.DivClass("accordion").R(
		element.ForEach(a.Items, func(item AccordionItem) {
			attrs := []string{"class", itemClass}
			if a.Exclusive {
				attrs = append(attrs, "name", groupName)
			}
			if item.Open {
				attrs = append(attrs, "open", "open")
			}

			b.Details(attrs...).R(
				b.SummaryClass("accordion-title").T(item.Title),
				b.DivClass("accordion-body").R(
					b.Wrap(func() {
						if item.Body != nil {
							item.Body.Render(b)
						} else {
							b.T(item.Content)
						}
					}),
				),
			)
		}),
	)
	return
}
