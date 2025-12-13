package element

import (
	"bytes"
	_ "embed"
	"fmt"
	"strings"
	"sync"
)

const (
	concernOther     = "other"
	concernOpenTag   = "open_tag"
	concernClosedTag = "closed_tag"
)

//go:embed assets/debug_table.js
var tableJS string

//go:embed assets/debug_table.css
var tableCSS string

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

// DebugClearIssues clears only the concerns map without turning off debug mode
func DebugClearIssues() {
	concerns.Clear()
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
		msg := "Debug mode is not enabled. Set debug mode to true to see element concerns."
		fmt.Println(msg)

		// Still return HTML for backward compatibility
		b, _, _ := Vars()
		b.Body().R(
			b.P("style", "font-weight:bold").T(msg),
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

	// Deduplicate concerns
	seenIssues := make(map[string]bool)
	dedupedConcerns := make(map[string]Element)

	for key, el := range concerns.cmap {
		dedupKey := buildDedupKey(el, key)
		if !seenIssues[dedupKey] {
			seenIssues[dedupKey] = true
			dedupedConcerns[key] = el
		}
	}

	// Always output markdown to terminal
	fmt.Printf("\n## ELEMENT CONCERNS: %d issues\n\n", len(dedupedConcerns))
	fmt.Println("| Details | Issues |")
	fmt.Println("|---------|--------|")

	for key, el := range dedupedConcerns {
		details := fmt.Sprintf("**%s** tag %s - %s (%s)", el.name, el.id, el.function, el.location)

		var issues string
		if strings.HasPrefix(key, concernOpenTag) {
			issues = fmt.Sprintf("**%s** tag not closed", el.name)
		} else {
			if len(el.issues) > 0 {
				issueList := make([]string, len(el.issues))
				for i, issue := range el.issues {
					issueList[i] = "• " + issue
				}
				issues = strings.Join(issueList, "; ")
			}
		}

		fmt.Printf("| %s | %s |\n", details, issues)
	}
	fmt.Println()

	if len(opts) > 0 && opts[0].TextOnly {
		// TextOnly mode - just return empty string since we already printed to terminal
		return ""
	} else {
		b, _, _ := Vars()

		// Build markdown content
		var markdownContent bytes.Buffer
		markdownContent.WriteString("## Element Concerns\n\n")
		markdownContent.WriteString(fmt.Sprintf("Total issues: %d\n\n", len(dedupedConcerns)))
		markdownContent.WriteString("| Key | Details | Issues |\n")
		markdownContent.WriteString("|-----|---------|--------|\n")

		for key, el := range dedupedConcerns {
			var issuesText string
			if strings.HasPrefix(key, concernOpenTag) {
				issuesText = fmt.Sprintf("**%s** tag not closed", el.name)
			} else {
				if len(el.issues) > 0 {
					issueList := make([]string, len(el.issues))
					for i, issue := range el.issues {
						issueList[i] = issue
					}
					issuesText = strings.Join(issueList, ", ")
				}
			}
			markdownContent.WriteString(fmt.Sprintf("| %s | %s | %s |\n", key, el.detailsHtml(), issuesText))
		}

		b.Html().R(
			b.Head().R(
				b.Title().T("Element Concerns Report"),
				b.Style().T(tableCSS),
			),
			b.Body().R(
				b.Script().T(tableJS),
				b.H2().T("Element Concerns"),
				b.P().R(
					b.F("Total issues: %d", len(dedupedConcerns)),
				),
				b.ButtonClass("clear-button", "onclick", "clearIssues()").T("Clear Issues"),

				// Tab navigation
				b.DivClass("tabs").R(
					b.DivClass("tab active", "id", "html-tab", "onclick", "switchTab('html')").T("HTML"),
					b.DivClass("tab", "id", "markdown-tab", "onclick", "switchTab('markdown')").T("Markdown"),
				),

				// HTML content (default active)
				b.DivClass("tab-content active", "id", "html-content").R(
					b.TableClass("tbl-element-concerns").R(
						b.THead().R(
							b.Tr().R(
								b.Th().T("Key"),
								b.Th().T("Details"),
								b.Th().T("Issues"),
							),
						),
						b.TBody().R(
							b.Wrap(func() {
								for key, el := range dedupedConcerns {
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

				// Markdown content
				b.DivClass("tab-content", "id", "markdown-content").R(
					b.DivClass("markdown-view").R(
						b.ButtonClass("markdown-copy-button", "onclick", "copyMarkdownContent()").R(
							b.F(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24"><path d="M16 1H4c-1.1 0-2 .9-2 2v14h2V3h12V1zm3 4H8c-1.1 0-2 .9-2 2v14c0 1.1.9 2 2 2h11c1.1 0 2-.9 2-2V7c0-1.1-.9-2-2-2zm0 16H8V7h11v14z"/></svg>`),
							b.F("Copy"),
						),
						b.Pre("id", "markdown-content-pre").T(markdownContent.String()),
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

// buildDedupKey creates a composite key for deduplication based on location, element name, and issue
func buildDedupKey(el Element, concernKey string) string {
	// For open tag issues
	if strings.HasPrefix(concernKey, concernOpenTag) {
		return el.location + "|" + el.name + "|open_tag_not_closed"
	}

	// For other issues, concatenate all issue texts
	if len(el.issues) > 0 {
		issueText := strings.Join(el.issues, ";")
		return el.location + "|" + el.name + "|" + issueText
	}

	// Fallback
	return el.location + "|" + el.name + "|unknown"
}
