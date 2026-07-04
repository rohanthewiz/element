package components

import (
	"strconv"

	"github.com/rohanthewiz/element"
)

// -----------------------------------------------------------------------------
// Progress Bar Component
// -----------------------------------------------------------------------------

// ProgressBar renders a horizontal progress indicator with proper ARIA
// attributes. It is div-based (rather than <progress>) so it can be fully
// styled with the .progress and .progress-fill classes.
type ProgressBar struct {
	Value     int    // Current value
	Max       int    // Maximum value (default 100)
	Label     string // Accessible label describing what is progressing
	ShowValue bool   // Show the percentage as text inside the bar
	Class     string // Additional CSS classes
}

// Render implements the element.Component interface.
func (p ProgressBar) Render(b *element.Builder) (x any) {
	maxVal := p.Max
	if maxVal <= 0 {
		maxVal = 100
	}
	val := min(max(p.Value, 0), maxVal)
	pct := val * 100 / maxVal

	barClass := "progress"
	if p.Class != "" {
		barClass += " " + p.Class
	}

	attrs := []string{
		"class", barClass,
		"role", "progressbar",
		"aria-valuenow", strconv.Itoa(val),
		"aria-valuemin", "0",
		"aria-valuemax", strconv.Itoa(maxVal),
	}
	if p.Label != "" {
		attrs = append(attrs, "aria-label", p.Label)
	}

	b.Div(attrs...).R(
		b.DivClass("progress-fill", "style", "width:"+strconv.Itoa(pct)+"%").R(
			b.Wrap(func() {
				if p.ShowValue {
					b.T(strconv.Itoa(pct) + "%")
				}
			}),
		),
	)
	return
}
