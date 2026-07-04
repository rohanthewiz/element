package components

import (
	"fmt"
	"strconv"

	"github.com/rohanthewiz/element"
)

// -----------------------------------------------------------------------------
// Pagination Component
// -----------------------------------------------------------------------------

// Pagination renders page navigation controls.
type Pagination struct {
	CurrentPage int    // Current page number (1-based)
	TotalPages  int    // Total number of pages
	BaseURL     string // URL pattern with %d placeholder for page number
	ShowFirst   bool   // Show "First" link
	ShowLast    bool   // Show "Last" link
	Class       string // Additional CSS classes
}

// Render implements the element.Component interface.
func (p Pagination) Render(b *element.Builder) (x any) {
	if p.TotalPages <= 1 {
		return // Don't render if only one page
	}

	navClass := "pagination"
	if p.Class != "" {
		navClass += " " + p.Class
	}

	b.NavClass(navClass, "aria-label", "Page navigation").R(
		b.UlClass("pagination-list").R(
			// First page link
			func() (x any) {
				if p.ShowFirst && p.CurrentPage > 1 {
					b.Li().R(
						b.AClass("pagination-link", "href", fmt.Sprintf(p.BaseURL, 1)).T("« First"),
					)
				}
				return
			}(),
			// Previous page link
			func() (x any) {
				if p.CurrentPage > 1 {
					b.Li().R(
						b.AClass("pagination-link", "href", fmt.Sprintf(p.BaseURL, p.CurrentPage-1), "rel", "prev").T("‹ Prev"),
					)
				}
				return
			}(),
			// Page numbers (show current and neighbors)
			func() (x any) {
				start := max(1, p.CurrentPage-2)
				end := min(p.TotalPages, p.CurrentPage+2)

				for i := start; i <= end; i++ {
					liClass := "pagination-item"
					if i == p.CurrentPage {
						liClass += " active"
						b.LiClass(liClass).R(
							b.SpanClass("pagination-current", "aria-current", "page").T(strconv.Itoa(i)),
						)
					} else {
						b.LiClass(liClass).R(
							b.AClass("pagination-link", "href", fmt.Sprintf(p.BaseURL, i)).T(strconv.Itoa(i)),
						)
					}
				}
				return
			}(),
			// Next page link
			func() (x any) {
				if p.CurrentPage < p.TotalPages {
					b.Li().R(
						b.AClass("pagination-link", "href", fmt.Sprintf(p.BaseURL, p.CurrentPage+1), "rel", "next").T("Next ›"),
					)
				}
				return
			}(),
			// Last page link
			func() (x any) {
				if p.ShowLast && p.CurrentPage < p.TotalPages {
					b.Li().R(
						b.AClass("pagination-link", "href", fmt.Sprintf(p.BaseURL, p.TotalPages)).T("Last »"),
					)
				}
				return
			}(),
		),
	)
	return
}
