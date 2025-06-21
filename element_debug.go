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

	// Always output markdown to terminal
	fmt.Printf("\n## ELEMENT CONCERNS: %d issues\n\n", len(concerns.cmap))
	fmt.Println("| Details | Issues |")
	fmt.Println("|---------|--------|")
	
	for key, el := range concerns.cmap {
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
		var markdownContent strings.Builder
		markdownContent.WriteString("## Element Concerns\n\n")
		markdownContent.WriteString(fmt.Sprintf("Total issues: %d\n\n", len(concerns.cmap)))
		markdownContent.WriteString("| Key | Details | Issues |\n")
		markdownContent.WriteString("|-----|---------|--------|\n")
		
		for key, el := range concerns.cmap {
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
				b.Style().T(`
            /* Tab styles */
            .tabs {
                display: flex;
                gap: 10px;
                margin-bottom: 20px;
                border-bottom: 2px solid #e0e0e0;
            }
            .tab {
                padding: 10px 20px;
                cursor: pointer;
                background-color: #f5f5f5;
                border: 1px solid #e0e0e0;
                border-bottom: none;
                border-radius: 5px 5px 0 0;
                font-weight: 500;
                transition: all 0.3s;
            }
            .tab:hover {
                background-color: #e8e8e8;
            }
            .tab.active {
                background-color: white;
                border-color: #3498db;
                border-top: 3px solid #3498db;
                color: #3498db;
                margin-bottom: -2px;
                padding-bottom: 12px;
            }
            .tab-content {
                display: none;
            }
            .tab-content.active {
                display: block;
            }
            
            /* Markdown view styles */
            .markdown-view {
                font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, monospace;
                background-color: #f8f9fa;
                padding: 20px;
                border-radius: 5px;
                overflow-x: auto;
            }
            .markdown-view pre {
                margin: 0;
                white-space: pre-wrap;
            }
            
            /* Existing table styles */
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
            .copy-icon {
                display: inline-block;
                width: 16px;
                height: 16px;
                margin-left: 5px;
                cursor: pointer;
                vertical-align: text-bottom;
                transition: opacity 0.2s;
                stroke: #3498db;
            }
            .copy-icon:hover {
                opacity: 0.7;
            }
            .notification {
                position: fixed;
                top: 20px;
                right: 20px;
                background-color: #2ecc71;
                color: white;
                padding: 10px 20px;
                border-radius: 5px;
                box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
                opacity: 0;
                transition: opacity 0.3s;
                z-index: 1000;
            }
            .notification.show {
                opacity: 1;
            }
            .clear-button {
                background-color: #6c757d;
                color: white;
                border: none;
                padding: 10px 20px;
                font-size: 14px;
                border-radius: 5px;
                cursor: pointer;
                margin: 10px 0;
                transition: background-color 0.3s;
            }
            .clear-button:hover {
                background-color: #5a6268;
            }
            .clear-button:active {
                transform: scale(0.98);
            }
        `),
			),
			b.Body().R(
				b.Script().T(`
					function switchTab(tabName) {
						// Remove active class from all tabs and contents
						const tabs = document.querySelectorAll('.tab');
						const contents = document.querySelectorAll('.tab-content');
						
						tabs.forEach(tab => tab.classList.remove('active'));
						contents.forEach(content => content.classList.remove('active'));
						
						// Add active class to selected tab and content
						document.getElementById(tabName + '-tab').classList.add('active');
						document.getElementById(tabName + '-content').classList.add('active');
					}
					
					function copyToClipboard(text) {
						navigator.clipboard.writeText(text).then(function() {
							showNotification('Copied: ' + text);
						}).catch(function(err) {
							console.error('Failed to copy: ', err);
							// Fallback for older browsers
							const textArea = document.createElement("textarea");
							textArea.value = text;
							textArea.style.position = "fixed";
							textArea.style.left = "-999999px";
							document.body.appendChild(textArea);
							textArea.focus();
							textArea.select();
							try {
								document.execCommand('copy');
								showNotification('Copied: ' + text);
							} catch (err) {
								console.error('Fallback copy failed: ', err);
							}
							document.body.removeChild(textArea);
						});
					}
					
					function showNotification(message) {
						const notification = document.createElement('div');
						notification.className = 'notification';
						notification.textContent = message;
						document.body.appendChild(notification);
						
						// Trigger reflow to enable transition
						notification.offsetHeight;
						notification.classList.add('show');
						
						setTimeout(function() {
							notification.classList.remove('show');
							setTimeout(function() {
								document.body.removeChild(notification);
							}, 300);
						}, 2000);
					}
					
					function clearIssues() {
						fetch('/debug/clear-issues')
							.then(response => {
								if (response.ok) {
									showNotification('Issues cleared successfully');
									setTimeout(function() {
										window.location.reload();
									}, 1000);
								} else {
									showNotification('Failed to clear issues');
								}
							})
							.catch(error => {
								console.error('Error clearing issues:', error);
								showNotification('Error clearing issues');
							});
					}
				`),
				b.H2().T("Element Concerns"),
				b.P().R(
					b.F("Total issues: %d", len(concerns.cmap)),
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
				
				// Markdown content
				b.DivClass("tab-content", "id", "markdown-content").R(
					b.DivClass("markdown-view").R(
						b.Pre().T(markdownContent.String()),
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
