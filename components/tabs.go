package components

import (
	"strconv"
	"strings"

	"github.com/rohanthewiz/element"
)

// -----------------------------------------------------------------------------
// Tabs Component
// -----------------------------------------------------------------------------

// Tab is a single tab: a label and its panel content.
type Tab struct {
	Label   string            // Tab button text
	Content string            // Simple text content
	Body    element.Component // Alternative: render a component as the panel
}

// Tabs renders a tabbed interface using the CSS-only radio-button pattern —
// no JavaScript required. Each instance emits a small scoped <style> block
// with the show/hide rules for its panels.
//
// Style hooks (global, work for every instance):
//
//	.tabs .tab-label                          — tab buttons
//	.tabs > .tab-radio:checked + .tab-label   — the active tab button
//	.tabs .tab-panel                          — panel container
//
// When rendering more than one Tabs on a page, give each a distinct ID.
type Tabs struct {
	ID     string // Unique id for this instance on the page (default "tabs")
	Items  []Tab
	Active int    // Index of the initially active tab (default 0)
	Class  string // Additional CSS classes for the wrapper
}

// Render implements the element.Component interface.
func (t Tabs) Render(b *element.Builder) (x any) {
	if len(t.Items) == 0 {
		return
	}

	id := t.ID
	if id == "" {
		id = "tabs"
	}
	active := t.Active
	if active < 0 || active >= len(t.Items) {
		active = 0
	}
	wrapClass := "tabs"
	if t.Class != "" {
		wrapClass += " " + t.Class
	}

	b.DivClass(wrapClass, "id", id).R(
		// Radio + label pairs form the tab bar
		element.ForEach2(t.Items, func(item Tab, i int) {
			tabID := id + "-t" + strconv.Itoa(i)
			radioAttrs := []string{"type", "radio", "class", "tab-radio", "name", id, "id", tabID}
			if i == active {
				radioAttrs = append(radioAttrs, "checked", "checked")
			}
			b.Input(radioAttrs...)
			b.LabelClass("tab-label", "for", tabID).T(item.Label)
		}),
		// Panels
		b.DivClass("tab-panels").R(
			element.ForEach(t.Items, func(item Tab) {
				b.DivClass("tab-panel").R(
					b.Wrap(func() {
						if item.Body != nil {
							item.Body.Render(b)
						} else {
							b.T(item.Content)
						}
					}),
				)
			}),
		),
		// Scoped structural rules: hide radios and inactive panels
		b.Style().T(t.styleRules(id)),
	)
	return
}

// styleRules generates the per-instance show/hide CSS
func (t Tabs) styleRules(id string) string {
	var sb strings.Builder
	sb.WriteString("#" + id + " > .tab-radio{display:none}\n")
	sb.WriteString("#" + id + " > .tab-panels > .tab-panel{display:none}\n")
	for i := range t.Items {
		tabID := id + "-t" + strconv.Itoa(i)
		nth := strconv.Itoa(i + 1)
		sb.WriteString("#" + tabID + ":checked ~ .tab-panels > .tab-panel:nth-child(" + nth + "){display:block}\n")
	}
	return sb.String()
}
