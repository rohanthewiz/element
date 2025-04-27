package element

import (
	"fmt"
	"strings"
	"sync"
)

const (
	concernOther     = "other"
	concernOpenTag   = "open_tag"
	concernClosedTag = "closed_tag"
)

var debugMode = false // we will set to true to enable debug mode

// DebugSet sets the debug mode to true.
func DebugSet() {
	debugMode = true
}

func IsDebugMode() bool {
	return debugMode
}

func DebugClear() {
	debugMode = false
	concerns.Clear() // Clear the concerns map
}

var concerns = elementConcerns{cmap: make(map[string]Element, 8)}

type elementConcerns struct {
	cmap map[string]Element
	lock sync.Mutex // mutex to protect the map
}

// UpsertConcern adds a concern to the concerns map in the form of an element and a possible issue.
// If the concernType is "open_tag", it will just add the element to the map.
// If the concernType is "closed_tag", it will check if the element exists in the map and remove it.
// If the concernType is anything else, it will append the issue to the element's issues.
func (con elementConcerns) UpsertConcern(concernType string, el Element) {
	if !debugMode {
		return
	}

	key := buildConcernKey(concernType, el.id)

	// Turn on lock so we can safely modify the map
	con.lock.Lock()
	defer con.lock.Unlock() // Ensure the lock is released after the function returns

	if concernType == concernOpenTag {
		if el.IsSingleTag() {
			return // Single tags should not be added as open tags
		}
		con.cmap[key] = el
		return
	}

	if concernType == concernClosedTag {
		if el.IsSingleTag() {
			return // Single tags should not be closed
		}

		key = buildConcernKey(concernOpenTag, el.id) // we need to check for the open tag key
		if _, ok := con.cmap[key]; ok {              // if the element exists in the map
			delete(con.cmap, key) // remove it from the map
			// fmt.Println("Removed", key, "from concerns map for element", el.details())
		} else {
			fmt.Println("No open tag found for element (that's weird)", el.details())
		}
		return
	}

	if len(el.issues) == 0 { // if there are no issues provided, we don't need to do anything
		fmt.Println("No issues provided for concern", concernType, "for element", el.details())
		return
	}

	con.cmap[key] = el // update the map
	// fmt.Println("**-> Concern updated for ", el.details(), "with issues:", strings.Join(el.issues, ", "))
}

// Clear clears the concerns map
func (con elementConcerns) Clear() {
	con.lock.Lock()
	defer con.lock.Unlock() // Ensure the lock is released after the function returns
	clear(con.cmap)
}

type DebugOptions struct {
	TextOnly bool // Default output is HTML
}

func DebugShow(opts ...DebugOptions) (out string) {
	if !debugMode {
		b, _, _ := Vars()
		b.Body().R(
			b.P("style", "font-weight:bold").
				T("Debug mode is not enabled. Set debug mode to true to see element concerns."),
		)
		return b.String()
	}

	// Pause Debug mode so we can check the current concerns
	debugMode = false
	defer func() {
		debugMode = true // Restore debug mode after checking
	}()

	if len(concerns.cmap) <= 0 {
		msg := "No element concerns found."
		fmt.Println(msg)
		return msg // No concerns to report
	}

	if len(opts) > 0 && opts[0].TextOnly {
		fmt.Printf("\nELEMENT CONCERNS: %d issues ------------------\n\n", len(concerns.cmap))
		for key, el := range concerns.cmap {
			fmt.Printf("- key: %s\n", key)
			fmt.Println(el.details())
			if strings.HasPrefix(key, concernOpenTag) {
				fmt.Printf("appears to be not closed with R()\n\n")
			} else {
				fmt.Printf("### issues\n* %s\n\n", strings.Join(el.issues, "\n  * "))
			}
		}

	} else {
		b, _, _ := Vars()

		b.Html().R(
			b.Head().R(
				b.Title().T("Element Concerns Report"),
				b.Style().T(`
            .tbl-element-concerns {
                width: 100%;
                border-collapse: collapse;
                margin: 25px 0;
                font-size: 0.9em;
                font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Oxygen-Sans, Ubuntu, Cantarell, "Helvetica Neue", sans-serif;
                box-shadow: 0 0 20px rgba(0, 0, 0, 0.15);
                border-radius: 5px;
                overflow: hidden;
            }
            .tbl-element-concerns thead tr {
                background-color: #3498db;
                color: #ffffff;
                text-align: left;
                font-weight: bold;
            }
            .tbl-element-concerns th,
            .tbl-element-concerns td {
                padding: 12px 15px;
            }
            .tbl-element-concerns tbody tr {
                border-bottom: 1px solid #dddddd;
            }
            .tbl-element-concerns tbody tr:nth-of-type(even) {
                background-color: #f3f3f3;
            }
            .tbl-element-concerns tbody tr:last-of-type {
                border-bottom: 2px solid #3498db;
            }
            .tbl-element-concerns tbody tr:hover {
                background-color: #e6f7ff;
                cursor: pointer;
            }
            .tbl-element-concerns strong {
                color: #e74c3c;
            }
            .tbl-element-concerns ul {
                margin: 0;
                padding-left: 20px;
            }
            .tbl-element-concerns li {
                margin-bottom: 5px;
            }
            .tbl-element-concerns li:last-child {
                margin-bottom: 0;
            }
        `),
			),
			b.Body().R(
				b.H2().T("Element Concerns"),
				b.P().R(
					b.F("Total issues: %d", len(concerns.cmap)),
				),

				b.Table("class", "tbl-element-concerns").R(
					b.THead().R(
						b.Tr().R(
							b.Th().T("Key"),
							b.Th().T("Details"),
							b.Th().T("Issues"),
						),
					),
					b.TBody().R(
						b.Wrap(func() {
							for key, el := range concerns.cmap {
								b.Tr().R(
									b.Td().T(key),
									b.Td().T(el.detailsHtml()),
									b.Td().R(
										b.Wrap(func() {
											if strings.HasPrefix(key, concernOpenTag) {
												b.F("<strong>%s</strong> tag not closed", el.name)
											} else {
												b.Ul().R(
													ForEach(el.issues, func(issue string) {
														b.Li().T(issue)
													}),
												)
											}
										}),
									),
								)
							}
						}),
					),
				),
			),
		)

		out = b.String()
	}

	return
}

func buildConcernKey(concernType, id string) string {
	return concernType + "-" + id
}
