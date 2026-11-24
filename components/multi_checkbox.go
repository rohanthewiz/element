package components

import (
	"html"
	"strconv"
	"strings"

	"github.com/rohanthewiz/element"
)

// -----------------------------------------------------------------------------
// Multi Checkbox Component
// -----------------------------------------------------------------------------

// CheckOption is one box in a MultiCheckbox.
type CheckOption struct {
	Value    string // Submitted value, and what Selected and the script's values() hold
	Label    string // Display text (defaults to Value); escaped, since it is usually data
	Title    string // Tooltip (defaults to Label, so a name cut short by the ellipsis can still be read in full)
	Count    int    // Shown as "(n)" after the label and summed in the summary; only used when CountNoun is set
	Checked  bool   // Start ticked
	Disabled bool   // Cannot be changed; Tick shown and Clear leave it as it is
}

// MultiCheckbox renders a grid of checkboxes with a filter, a running count
// of what is ticked, and Tick shown / Clear buttons — for picking any number
// of items out of a long list (folders, tags, accounts).
//
//	Filter folders [acme        ]  [Extra…]  1 folder ticked · 1 DAG  [Tick shown] [Clear]
//	Folders: ┌──────────────────────────────────────────────────────────┐
//	         │ ☐ exports/db_to_db/acme/emea/prod/ (2)  ☑ acme/prod/ (1) │  ≤ Columns per row,
//	         └──────────────────────────────────────────────────────────┘  scrolls past MaxHeight
//
// Behaviour (a small inline script, see multiCheckboxJS):
//
//   - The filter narrows which boxes are shown, by a case-blind substring of
//     the label. It never changes what is ticked, and a ticked box stays shown
//     whatever the filter says, so a choice cannot be hidden and then
//     submitted by surprise.
//   - The filter takes a comma-separated list of terms ("q1, acme"), and a box
//     is shown when any one of them is in its label. Spaces around a term are
//     ignored, and so are empty terms, so "q1," is just "q1" and "," alone is
//     no filter at all. The box's tooltip says so.
//   - Tick shown ticks every box the filter matches (all of them with no
//     filter), adding to what is already ticked. Clear unticks every box.
//   - The summary counts the ticked boxes and, with CountNoun set, totals
//     their Counts ("2 folders ticked · 37 DAGs"). Past MaxCount it adds
//     "(at most N)" and the class mcb-over. It is aria-live, so screen readers
//     hear the count change.
//   - Enter in the filter ticks, for each term, the box it names exactly
//     (value or label), or failing that its only match; a term with several
//     matches and no exact one ticks nothing. So "sales/q1, acme/prod" ticks
//     both. The whole text is tried as one term first, so a label that itself
//     holds a comma ("Smith, Jo") is ticked by typing it in full rather than
//     split. Enter never submits the surrounding form. To make Enter also
//     "go", add a keydown listener to .mcb-filter with addEventListener: it
//     runs after this one, so it sees the new ticks.
//   - Every change, including the buttons', ends in one bubbling "change"
//     event on the component's root element, so a caller needs one listener.
//   - With FilterHistoryKey set, the filters used lately can be picked from
//     a dropdown beside the filter (see Filter history, below).
//
// Filter history. With FilterHistoryKey set, a ▾ button beside the filter
// drops down the filters used lately, newest first, kept in localStorage
// under that key:
//
//	Filter folders [acme        ][▾]  …
//	               ┌─────────────────┐
//	               │ sales/q1, acme  │  ← click it, or Enter on it, and it is put in the
//	               │ acme/prod       │    filter. It only filters: nothing is ticked
//	               ├─────────────────┤    until Enter or Tick shown, as with typed text
//	               │ Clear history   │
//	               └─────────────────┘
//
// A filter counts as used, and is remembered, on Enter, on Tick shown, and
// when a box is ticked or unticked by hand while it is in the box — not on
// every keystroke, which would fill the list with half-typed prefixes. A
// filter used again moves to the top rather than repeating. Two texts are the
// same filter when their terms are, as matched: "fm, om" and "FM,om" are one
// entry, and so are a pasted prefix and the bare name when the normalize hook
// strips it. The list keeps FilterHistorySize entries (default 10).
//
// Keys: ↓ (or Alt+↓) in the filter, or ↓ on ▾, opens it with the first entry
// focused; ↑/↓, Home and End move; Enter or Space uses an entry; Delete or
// Backspace forgets one; Escape closes it and returns to the filter; Tab
// closes it onto the ▾ button (Shift+Tab: onto the filter). A press anywhere
// outside it closes it, the filter included. The ▾ button is disabled while
// there is nothing to show.
//
// It is a real form control: every box carries Name, so a form submits the
// ticked values as Name=a&Name=b. Filtering hides boxes with the hidden
// attribute rather than removing them, so the checkboxes stay the only state
// — nothing parallel to drift — and a hidden box would still submit (though a
// ticked box is never hidden). Without JavaScript (disabled, or blocked by a
// Content-Security-Policy without 'unsafe-inline') it degrades to a plain
// checkbox grid: the filter, summary and buttons are rendered hidden and only
// the script reveals them, so they are never on screen and dead.
//
// Script API, for lists built in the browser: ElementMultiCheckbox.init(root)
// returns (or, the first time, creates) the instance, also kept on root.mcb:
//
//	values()       → the ticked values, in grid order
//	setValues(arr) → tick exactly these (disabled boxes keep their state)
//	setOptions(arr)→ replace the boxes with [{value, label, title, count, checked, disabled}];
//	                 ticks carry over for values still in the list; clears a note
//	setNote(text)  → show text in place of the boxes ("Loading…", an error) without
//	                 removing them, so their ticks survive a reload; null hides it
//	refresh()      → re-apply the filter and summary after changing boxes, a hook's
//	                 inputs, or data-max-count by hand
//
// Two optional hooks, assigned on the instance (then call refresh()):
//
//	normalize = function (text) → text   one filter term as matched, e.g. strip a
//	                                     pasted prefix; used by the filter and Enter.
//	                                     Called once per comma-separated term (and
//	                                     on the whole text for Enter's first try),
//	                                     so "repo/a, repo/b" strips both prefixes
//	total     = function (values) → n    the summary's total in place of the summed
//	                                     Counts, for totals that depend on more than
//	                                     the ticked boxes (a "subfolders" toggle);
//	                                     only used when CountNoun is set
//
// The summary re-reads data-max-count every time, so a limit that is only
// known later (from an API answer) can be set on the root followed by refresh().
//
// With StorageKey set, the ticks are remembered in the browser's localStorage
// and restored on load, where they win over Checked and Selected. When the
// grid starts empty and is filled by setOptions, they are restored by the
// first setOptions instead.
//
// Layout: a small scoped <style> holds the grid (Columns per row, stepping
// down to 6 / 4 / 2 as the component narrows, via container queries on
// .mcb-body), the ellipsis on long labels, and the toolbar row. Those rules
// sit inside :where(), which has zero specificity, so any rule in your
// stylesheet overrides them. The one rule that does not is
// "#ID [hidden]{display:none!important}", because a stylesheet's
// .mcb-opt{display:…} would otherwise un-hide filtered-out boxes. Container
// queries make .mcb-body take the width of its parent; it cannot shrink to
// fit its contents.
//
// The filter history's dropdown hangs under the filter (absolute, in a
// position:relative .mcb-filter-wrap), so it needs no room in the layout and
// scrolls with the page; an ancestor with overflow:hidden that ends right
// under the toolbar would cut it off. Its default colors are the system ones
// (Canvas, CanvasText, Highlight), which follow color-scheme; a page that
// darkens itself with a class should color .mcb-history itself.
//
// Style hooks: .multi-checkbox, .mcb-toolbar, .mcb-filter-label, .mcb-filter,
// .mcb-summary (.mcb-over past MaxCount), .mcb-btn (.mcb-tick-shown,
// .mcb-clear, .mcb-history-btn), .mcb-body, .mcb-label, .mcb-grid, .mcb-opt,
// .mcb-cb, .mcb-name, .mcb-count, .mcb-empty; with FilterHistoryKey,
// .mcb-filter-wrap, .mcb-history, .mcb-history-list, .mcb-history-item
// ([aria-selected=true] for the one the filter holds now), .mcb-history-clear.
//
// ID is used in CSS selectors and element ids, so keep it a plain identifier
// and distinct per instance on the page.
type MultiCheckbox struct {
	ID       string        // Unique id for this instance on the page (default Name, else "multi-checkbox")
	Name     string        // Form name shared by every box; empty submits nothing (script-only use)
	Label    string        // Text beside the grid, e.g. "Folders:"
	Options  []CheckOption // The boxes, in display order
	Selected []string      // Values to start ticked, as well as any Option.Checked (handy when re-rendering a submitted form)

	Columns   int    // Most boxes per row (default 8)
	MaxHeight string // CSS length the grid scrolls past (default "16rem")

	FilterLabel       string // Text before the filter box (default "Filter:")
	FilterPlaceholder string // Placeholder in the filter box

	// Recent filters, in a dropdown beside the filter (see the behaviour
	// notes above). Kept in the browser's localStorage, apart from
	// StorageKey's ticks: what was typed to find boxes is a habit of the
	// person, not part of the selection, and outlives it. Instances that
	// share a key share the list.
	FilterHistoryKey  string // localStorage key for the recent filters; empty = no history and no ▾ button
	FilterHistorySize int    // Most recent filters kept (default 10)

	// The summary's wording: "<n> <Noun|NounPlural> <TickedWord> · <total> <CountNoun|CountNounPlural>".
	Noun            string // One option, e.g. "folder" (default "item")
	NounPlural      string // Default Noun + "s"
	TickedWord      string // Default "ticked"
	CountNoun       string // What Count counts, e.g. "DAG"; empty leaves counts out of the labels and summary
	CountNounPlural string // Default CountNoun + "s"

	// MaxCount warns once the total passes it: the summed Counts when
	// CountNoun is set, otherwise the number of boxes ticked. A warning
	// only — enforce the limit on the server. 0 = no limit.
	MaxCount int

	TickShownText string // Tick shown button text (default "Tick shown"); custom text drops the default tooltip
	ClearText     string // Clear button text (default "Clear"); custom text drops the default tooltip

	Extra      element.Component // More controls for the toolbar, after the filter (e.g. an "include subfolders" checkbox)
	StorageKey string            // localStorage key to remember the ticks under; empty = not remembered
	HelpText   string            // Optional help text below the grid
	Error      string            // Error message to display
	Class      string            // Additional CSS classes for the wrapper
}

// Render implements the element.Component interface.
func (m MultiCheckbox) Render(b *element.Builder) (x any) {
	id := m.ID
	if id == "" {
		id = m.Name
	}
	if id == "" {
		id = "multi-checkbox"
	}

	noun := orDefault(m.Noun, "item")
	nouns := orDefault(m.NounPlural, noun+"s")
	counts := m.CountNoun != ""

	// The script reads its wording and limits from data attributes on the
	// root, rather than from values written into the script, so the script
	// text is identical for every instance and carries no caller data that
	// would need JS-string escaping.
	rootClass := "multi-checkbox"
	if m.Class != "" {
		rootClass += " " + m.Class
	}
	rootAttrs := []string{
		"class", fieldClass(m.Error, rootClass), "id", id,
		"data-noun", noun, "data-nouns", nouns,
		"data-ticked", orDefault(m.TickedWord, "ticked"),
	}
	if counts {
		rootAttrs = append(rootAttrs,
			"data-count-noun", m.CountNoun,
			"data-count-nouns", orDefault(m.CountNounPlural, m.CountNoun+"s"))
	}
	if m.MaxCount > 0 {
		rootAttrs = append(rootAttrs, "data-max-count", strconv.Itoa(m.MaxCount))
	}
	if m.StorageKey != "" {
		rootAttrs = append(rootAttrs, "data-storage-key", m.StorageKey)
	}
	if m.Name != "" {
		// For setOptions, which builds boxes in the browser and must name
		// them as these are named.
		rootAttrs = append(rootAttrs, "data-name", m.Name)
	}
	hist := m.FilterHistoryKey != ""
	if hist {
		rootAttrs = append(rootAttrs, "data-history-key", m.FilterHistoryKey)
		if m.FilterHistorySize > 0 {
			rootAttrs = append(rootAttrs, "data-history-size", strconv.Itoa(m.FilterHistorySize))
		}
	}

	// Selected as a set, so ticking is one lookup per option. Only built
	// when there is something in it; the common case allocates nothing.
	var selected map[string]bool
	if len(m.Selected) > 0 {
		selected = make(map[string]bool, len(m.Selected))
		for _, v := range m.Selected {
			selected[v] = true
		}
	}

	// The grid is a role=group named by the label, since a single <label>
	// can't name a set of controls. (A <fieldset>/<legend> would, but a
	// legend can't sit beside the grid the way this layout wants.)
	gridAttrs := []string{"class", "mcb-grid", "id", id + "-grid", "role", "group"}
	if m.Label != "" {
		gridAttrs = append(gridAttrs, "aria-labelledby", id+"-label")
	}
	gridAttrs = append(gridAttrs, describedByAttrs(id, m.HelpText, m.Error)...)

	// The "nothing to show" note. Shown from the start only when there are
	// no options at all; otherwise the script shows it when the filter
	// leaves nothing, and words it for that case.
	emptyAttrs := []string{"class", "mcb-empty"}
	emptyText := "No " + nouns + "."
	if len(m.Options) > 0 {
		emptyAttrs = append(emptyAttrs, "hidden", "hidden")
		emptyText = "No " + nouns + " match the filter."
	}

	shownText, shownTitle := m.TickShownText, ""
	if shownText == "" {
		shownText, shownTitle = "Tick shown", "Tick every "+noun+" the filter shows"
	}
	clearText, clearTitle := m.ClearText, ""
	if clearText == "" {
		clearText, clearTitle = "Clear", "Untick every "+noun
	}

	b.Div(rootAttrs...).R(
		// The toolbar. Everything in it that needs the script carries
		// mcb-js and starts hidden; the script reveals it. Extra is not
		// hidden: it is the caller's, and may work without the script.
		b.DivClass("mcb-toolbar").R(
			b.LabelClass("mcb-filter-label mcb-js", "for", id+"-filter", "hidden", "hidden").T(orDefault(m.FilterLabel, "Filter:")),
			b.Wrap(func() {
				// No name: the filter text is a view setting, not form data.
				// The title is the only place the comma lists are mentioned,
				// so they can be found without reading the placeholder (which
				// is the caller's, and gone once something is typed).
				attrs := []string{"type", "search", "id", id + "-filter", "class", "mcb-filter mcb-js",
					"hidden", "hidden", "autocomplete", "off", "aria-controls", id + "-grid",
					"title", "Separate several terms with commas"}
				if m.FilterPlaceholder != "" {
					attrs = append(attrs, "placeholder", m.FilterPlaceholder)
				}
				if !hist {
					b.Input(attrs...)
					return
				}
				// With history, the filter, its ▾ button and the dropdown
				// share one positioned box, so the dropdown can hang under
				// the filter whatever else the toolbar holds:
				//
				//   div.mcb-filter-wrap.mcb-js        position:relative, inline-flex
				//     input.mcb-filter                (as without history)
				//     button.mcb-history-btn          ▾, aria-expanded
				//     div.mcb-history[hidden]         absolute, top:100%
				//       div.mcb-history-list[role=listbox]   entries, built by the script
				//       button.mcb-history-clear      Clear history
				//
				// The wrap is script-only like the filter (mcb-js, hidden), so
				// with no script it adds no empty flex item to the toolbar.
				// The dropdown is not mcb-js: the script opens it, it does not
				// merely reveal it. The entries are left to the script, since
				// only the browser knows them. A div, not a span: the wrap
				// holds the dropdown's div.
				b.DivClass("mcb-filter-wrap mcb-js", "hidden", "hidden").R(
					b.Input(attrs...),
					// The glyph means nothing read aloud; aria-label names it.
					// The script swaps the title while there is nothing to
					// show, and only it knows whether there is.
					b.ButtonClass("mcb-btn mcb-history-btn", "type", "button", "id", id+"-history-btn",
						"aria-haspopup", "listbox", "aria-expanded", "false", "aria-controls", id+"-history-list",
						"aria-label", "Recent filters", "title", historyBtnTitle).T("▾"),
					b.Div("class", "mcb-history", "id", id+"-history", "hidden", "hidden").R(
						b.Div("class", "mcb-history-list", "id", id+"-history-list",
							"role", "listbox", "aria-label", "Recent filters").R(),
						b.ButtonClass("mcb-history-clear", "type", "button",
							"title", "Forget every recent filter").T("Clear history"),
					),
				)
			}),
			b.Wrap(func() {
				if m.Extra != nil {
					m.Extra.Render(b)
				}
			}),
			// Written by the script. One line (nowrap in the example
			// stylesheet) so the buttons after it don't jump as it grows.
			b.SpanClass("mcb-summary mcb-js", "id", id+"-summary", "aria-live", "polite", "hidden", "hidden").R(),
			b.Wrap(func() {
				attrs := []string{"type", "button", "hidden", "hidden"}
				if shownTitle != "" {
					attrs = append(attrs, "title", shownTitle)
				}
				b.ButtonClass("mcb-btn mcb-tick-shown mcb-js", attrs...).T(shownText)
			}),
			b.Wrap(func() {
				attrs := []string{"type", "button", "hidden", "hidden"}
				if clearTitle != "" {
					attrs = append(attrs, "title", clearTitle)
				}
				b.ButtonClass("mcb-btn mcb-clear mcb-js", attrs...).T(clearText)
			}),
		),

		b.DivClass("mcb-body").R(
			b.Wrap(func() {
				if m.Label != "" {
					b.SpanClass("mcb-label", "id", id+"-label").T(m.Label)
				}
			}),
			b.Div(gridAttrs...).R(
				// One <label> per box, wrapping its input, so the whole
				// label is the click target with no per-box id to invent.
				//
				//   label.mcb-opt[title]
				//     input.mcb-cb[name, value, data-count]
				//     span.mcb-name "acme/prod/"   ← ellipsis cuts this
				//     span.mcb-count "(1)"         ← never cut
				//
				// The name and count are separate spans so a long name is
				// cut while its count stays readable. setOptions in the
				// script builds this same shape; keep the two in step.
				//
				// Value and Title go through html.EscapeString: they are
				// data, and the attribute writer only escapes the double
				// quote (enough to stop a breakout, not enough for a value
				// like "a&amp" to survive the round trip unchanged).
				element.ForEach(m.Options, func(o CheckOption) {
					label := orDefault(o.Label, o.Value)
					title := orDefault(o.Title, label)

					attrs := []string{"type", "checkbox", "class", "mcb-cb", "value", html.EscapeString(o.Value)}
					if m.Name != "" {
						attrs = append(attrs, "name", m.Name)
					}
					if counts {
						attrs = append(attrs, "data-count", strconv.Itoa(o.Count))
					}
					if o.Checked || selected[o.Value] {
						attrs = append(attrs, "checked", "checked")
					}
					if o.Disabled {
						attrs = append(attrs, "disabled", "disabled")
					}

					b.LabelClass("mcb-opt", "title", html.EscapeString(title)).R(
						b.Input(attrs...),
						b.SpanClass("mcb-name").TE(label),
						b.Wrap(func() {
							if counts {
								b.SpanClass("mcb-count").T("(" + strconv.Itoa(o.Count) + ")")
							}
						}),
					)
				}),
				b.Span(emptyAttrs...).T(emptyText),
			),
		),

		b.Wrap(func() {
			renderFieldMessages(b, id, m.HelpText, m.Error)
		}),

		b.Style().T(m.styleRules(id)),
		// Last child of the root: the script finds its instance as
		// document.currentScript.parentNode, so it needs no id.
		b.Script().T(multiCheckboxJS),
	)
	return
}

// historyBtnTitle is the ▾ button's tooltip while it has filters to show.
// The script reads it off the button at start-up and puts it back whenever
// the list stops being empty, so the text lives only here.
const historyBtnTitle = "Recent filters (or ↓ in the filter)"

// orDefault returns s, or def when s is empty.
func orDefault(s, def string) string {
	if s == "" {
		return def
	}
	return s
}

// styleRules generates the per-instance layout CSS.
//
// Everything but the [hidden] rule is wrapped in :where(), which zeroes its
// specificity: these are defaults, and a plain ".mcb-grid{gap:2px}" in the
// page's stylesheet should beat them even though this <style> comes later
// in the document. The [hidden] rule is the exception — it must beat a
// stylesheet's display rule, or filtering would hide nothing.
//
// Columns: the grid's tracks are minmax(0, max-content), so each column is
// as wide as its widest label until the row runs out of room, and then the
// labels are cut with an ellipsis rather than the grid overflowing. A fixed
// track count can't adapt on its own, so container queries on .mcb-body step
// it down at fixed widths. Only the steps below Columns are emitted: with
// Columns 3, only the step to 2 applies.
func (m MultiCheckbox) styleRules(id string) string {
	cols := m.Columns
	if cols <= 0 {
		cols = 8
	}
	maxH := orDefault(m.MaxHeight, "16rem")
	scope := "#" + id + " "
	track := func(n int) string {
		return "grid-template-columns:repeat(" + strconv.Itoa(n) + ",minmax(0,max-content))"
	}

	var sb strings.Builder
	sb.WriteString(":where(" + scope + ".mcb-toolbar){display:flex;flex-wrap:wrap;align-items:center;gap:6px 10px}\n")
	sb.WriteString(":where(" + scope + ".mcb-body){display:flex;align-items:flex-start;gap:10px;container-type:inline-size}\n")
	sb.WriteString(":where(" + scope + ".mcb-grid){display:grid;" + track(cols) +
		";gap:4px 14px;flex:1;min-width:0;max-height:" + maxH + ";overflow-y:auto}\n")
	sb.WriteString(":where(" + scope + ".mcb-opt){display:inline-flex;align-items:center;gap:4px;min-width:0;white-space:nowrap;cursor:pointer}\n")
	sb.WriteString(":where(" + scope + ".mcb-name){overflow:hidden;text-overflow:ellipsis;min-width:0}\n")
	sb.WriteString(":where(" + scope + ".mcb-empty){grid-column:1/-1}\n")
	if m.FilterHistoryKey != "" {
		// The history dropdown. Only emitted when there is one.
		//  - The wrap is the containing block, so top:100% is "just under
		//    the filter" wherever the toolbar wraps it to.
		//  - width:max-content sizes the dropdown to its longest entry, as
		//    an absolute box's shrink-to-fit would only up to the wrap's
		//    width; min-width keeps it at least as wide as filter + ▾, and
		//    max-width stops a very long entry running off the page (it is
		//    cut with an ellipsis, and its title has it in full).
		//  - System colors, so it reads in light and dark color-schemes
		//    without a stylesheet; z-index lifts it over the grid below.
		sb.WriteString(":where(" + scope + ".mcb-filter-wrap){position:relative;display:inline-flex;align-items:center;gap:4px}\n")
		sb.WriteString(":where(" + scope + ".mcb-history){position:absolute;top:100%;left:0;z-index:10;margin-top:2px;" +
			"box-sizing:border-box;width:max-content;min-width:100%;max-width:min(40rem,90vw);max-height:16rem;overflow-y:auto;" +
			"padding:4px 0;background:Canvas;color:CanvasText;border:1px solid GrayText;border-radius:4px;" +
			"box-shadow:0 4px 12px rgba(0,0,0,.15)}\n")
		sb.WriteString(":where(" + scope + ".mcb-history-item){padding:3px 10px;cursor:pointer;white-space:nowrap;overflow:hidden;text-overflow:ellipsis}\n")
		sb.WriteString(":where(" + scope + ".mcb-history-item:hover," + scope + ".mcb-history-item:focus){background:Highlight;color:HighlightText;outline:none}\n")
		sb.WriteString(":where(" + scope + ".mcb-history-clear){display:block;box-sizing:border-box;width:100%;margin-top:4px;padding:3px 10px;" +
			"border:0;border-top:1px solid GrayText;background:none;color:inherit;font:inherit;text-align:left;cursor:pointer;opacity:.7}\n")
		sb.WriteString(":where(" + scope + ".mcb-history-clear:hover," + scope + ".mcb-history-clear:focus-visible){opacity:1}\n")
		sb.WriteString(":where(" + scope + ".mcb-history-btn:disabled){opacity:.4;cursor:default}\n")
	}
	sb.WriteString(scope + "[hidden]{display:none!important}\n")

	// Widths in px of .mcb-body, widest first, and the column count each
	// one steps down to.
	for _, step := range [...]struct{ width, cols int }{{1050, 6}, {800, 4}, {500, 2}} {
		if cols > step.cols {
			sb.WriteString("@container (max-width:" + strconv.Itoa(step.width) + "px){:where(" +
				scope + ".mcb-grid){" + track(step.cols) + "}}\n")
		}
	}
	return sb.String()
}

// multiCheckboxJS is the behaviour, inlined after each instance.
//
// Every instance emits the same text. The first run defines
// window.ElementMultiCheckbox; later runs find it and only call init, so a
// page with several instances holds one copy of the functions. init is
// idempotent (a second call returns the existing instance), so markup that is
// re-inserted and has its script re-run — htmx swaps, for one — is not wired
// twice.
//
// The checkboxes are the only state. Nothing mirrors the ticks in a JS
// object, so a form submit, values(), the summary and storage all read the
// same thing and cannot disagree.
//
// No innerHTML anywhere: setOptions builds nodes and sets textContent, so a
// label from data is only ever text.
//
// The filter text is a comma-separated list of terms. terms() splits it,
// normalizes and lowercases each part, and drops the empty ones; [] means no
// filter. A box hits when any term is a substring of its label, so the list
// is an OR. Enter resolves each term on its own (pickFor), after first trying
// the whole text as a single term so a label containing a comma can still be
// ticked by typing it in full:
//
//	"sales/q1, acme/prod, sales/"
//	   │ whole text names a box? ── yes ──► tick that one box
//	   no
//	   ├─ "sales/q1"   exact label ──────► tick
//	   ├─ "acme/prod"  exact value (beats the longer label that also holds it) ──► tick
//	   └─ "sales/"     several matches, none exact ──► nothing
//	all of a press's new ticks ──► one "change"
//
// Flow:
//
//	filter input ──► applyFilter (show = ticked || any term matches; a note hides all)
//	box click ─────────────────────────────┐
//	Tick shown / Clear / Enter / setValues ├─► "change" on root ─► renderSummary + save
//	setOptions ────────────────────────────┘
//
// The hooks live on the instance object (api), which exists before any
// handler runs, so a caller assigning api.total after init is seen by the
// next summary without re-wiring anything.
//
// Filter history (only when the root has data-history-key). localStorage is
// the list, as the checkboxes are the ticks: it is re-read whenever it is
// needed (opening the dropdown, focusing the filter) rather than held, so
// another tab's additions show up, and hMem is only a fallback for when
// storage can't be read at all. Entries are compared by histId, their terms
// as matched (norm per comma part), so the normalize hook decides what
// counts as "the same filter".
//
//	Enter / Tick shown / a box clicked ──► remember(filter text)
//	                                         drop the entry with the same histId,
//	                                         put this text first, keep hSize
//	▾ click, ↓ in the filter ──► openHistory ──► entries drawn as text nodes
//	entry click / Enter ──► useHistory: close, filter.value = text,
//	                        "input" on the filter (so applyFilter, and any
//	                        caller listening, run as for typing), remember
//
// Using an entry only filters. It ticks nothing, as typing the same text
// wouldn't, so picking one can't change a selection by surprise.
const multiCheckboxJS = `(function (script) {
  'use strict';
  var M = window.ElementMultiCheckbox;
  if (!M) {
    M = window.ElementMultiCheckbox = { init: init };
  }
  if (script && script.parentNode) M.init(script.parentNode);

  function init(root) {
    if (!root || !root.querySelector) return null;
    if (root.mcb) return root.mcb;
    var ds = root.dataset;
    var grid = root.querySelector('.mcb-grid');
    if (!grid) return null;
    var filter = root.querySelector('.mcb-filter');
    var summary = root.querySelector('.mcb-summary');
    var empty = grid.querySelector('.mcb-empty');
    var noun = ds.noun || 'item', nouns = ds.nouns || noun + 's';
    var ticked = ds.ticked || 'ticked';
    var cNoun = ds.countNoun || '', cNouns = ds.countNouns || cNoun + 's';
    var key = ds.storageKey || '';
    var api = { normalize: null, total: null };
    var note = null;     // setNote text; while set, it replaces the boxes
    var pending = null;  // remembered ticks awaiting the first setOptions (grid started empty)

    // Boxes are looked up each time rather than cached: setOptions replaces them.
    function each(fn) {
      var list = grid.querySelectorAll('.mcb-cb');
      for (var i = 0; i < list.length; i++) fn(list[i]);
    }
    function optOf(cb) { return cb.closest('.mcb-opt') || cb; }
    function nameOf(cb) {
      var n = optOf(cb).querySelector('.mcb-name');
      return (n ? n.textContent : cb.value).toLowerCase();
    }
    // norm is one term as matched: the caller's normalize hook, then trimmed
    // and lowercased (labels are compared lowercased too).
    function norm(text) {
      var t = typeof api.normalize === 'function' ? api.normalize(text) : text;
      return String(t == null ? '' : t).trim().toLowerCase();
    }
    // terms is the filter as a list of terms, split on commas; [] = no filter.
    // Each part is normalized on its own, so a hook that strips a pasted
    // prefix strips it from every term, not just the first.
    function terms() {
      var out = [];
      if (!filter) return out;
      var parts = filter.value.split(',');
      for (var i = 0; i < parts.length; i++) {
        var t = norm(parts[i]);
        if (t) out.push(t);
      }
      return out;
    }
    // hits: does any term occur in the box's label? No terms matches everything.
    function hits(cb, ts) {
      if (!ts.length) return true;
      var n = nameOf(cb);
      for (var i = 0; i < ts.length; i++) {
        if (n.indexOf(ts[i]) >= 0) return true;
      }
      return false;
    }
    // pickFor is the box one term names for Enter: an exact value or label,
    // else its only match, else null (no match, or several and none exact).
    // Disabled boxes are never picked.
    function pickFor(f) {
      var exact = null, only = null, matches = 0;
      each(function (cb) {
        if (exact || cb.disabled) return;
        var n = nameOf(cb);
        if (cb.value.toLowerCase() === f || n === f) { exact = cb; return; }
        if (n.indexOf(f) >= 0) { matches++; only = cb; }
      });
      return exact || (matches === 1 ? only : null);
    }
    function changed() { root.dispatchEvent(new Event('change', { bubbles: true })); }
    function asSet(list) {
      var s = Object.create(null);
      (Array.isArray(list) ? list : []).forEach(function (v) { s[String(v)] = true; });
      return s;
    }

    // A ticked box is always shown, so a choice can't be hidden by the filter.
    function applyFilter() {
      var ts = terms(), total = 0, shown = 0;
      each(function (cb) {
        var show = note === null && (cb.checked || hits(cb, ts));
        optOf(cb).hidden = !show;
        total++;
        if (show) shown++;
      });
      if (empty) {
        empty.textContent = note !== null ? note : total ? 'No ' + nouns + ' match the filter.' : 'No ' + nouns + '.';
        empty.hidden = shown > 0;
      }
    }

    function values() {
      var v = [];
      each(function (cb) { if (cb.checked) v.push(cb.value); });
      return v;
    }

    function renderSummary() {
      if (!summary) return;
      var n = 0, sum = 0, ticks = [];
      each(function (cb) {
        if (cb.checked) { n++; ticks.push(cb.value); sum += parseInt(cb.getAttribute('data-count'), 10) || 0; }
      });
      if (cNoun && typeof api.total === 'function') {
        var t = Number(api.total(ticks));
        if (isFinite(t)) sum = t;
      }
      var max = parseInt(ds.maxCount, 10) || 0;  // read each time: a caller may set it later
      var measured = cNoun ? sum : n;
      var over = max > 0 && measured > max;
      var t = n ? n + ' ' + (n === 1 ? noun : nouns) + ' ' + ticked : 'No ' + noun + ' ' + ticked;
      if (n && cNoun) t += ' · ' + sum.toLocaleString() + ' ' + (sum === 1 ? cNoun : cNouns);
      if (over) t += ' (at most ' + max.toLocaleString() + ')';
      summary.textContent = t;
      summary.classList.toggle('mcb-over', over);
    }

    // Best effort: private modes and full quotas throw, and the page works the same without it.
    function save() {
      if (!key) return;
      try { localStorage.setItem(key, JSON.stringify(values())); } catch (e) {}
    }
    function tickOnly(want) {
      each(function (cb) { if (!cb.disabled) cb.checked = !!want[cb.value]; });
    }

    function setValues(list) {
      tickOnly(asSet(list));
      applyFilter();
      changed();
    }

    // Builds the same shape as the Go renderer (see the diagram in MultiCheckbox.Render).
    function setOptions(list) {
      var keep = asSet(values()), name = ds.name || '';
      if (pending) {
        for (var k in pending) keep[k] = true;
        pending = null;
      }
      note = null;
      var frag = document.createDocumentFragment();
      (Array.isArray(list) ? list : []).forEach(function (o) {
        o = o || {};
        var value = String(o.value == null ? '' : o.value);
        var label = o.label == null || o.label === '' ? value : String(o.label);
        var cb = document.createElement('input');
        cb.type = 'checkbox';
        cb.className = 'mcb-cb';
        cb.value = value;
        if (name) cb.name = name;
        if (cNoun) cb.setAttribute('data-count', String(parseInt(o.count, 10) || 0));
        cb.disabled = !!o.disabled;
        cb.checked = !!o.checked || (!o.disabled && !!keep[value]);
        var lab = document.createElement('label');
        lab.className = 'mcb-opt';
        lab.title = o.title ? String(o.title) : label;
        var nm = document.createElement('span');
        nm.className = 'mcb-name';
        nm.textContent = label;
        lab.appendChild(cb);
        lab.appendChild(nm);
        if (cNoun) {
          var c = document.createElement('span');
          c.className = 'mcb-count';
          c.textContent = '(' + (parseInt(o.count, 10) || 0) + ')';
          lab.appendChild(c);
        }
        frag.appendChild(lab);
      });
      each(function (cb) { var o = optOf(cb); if (o.parentNode) o.parentNode.removeChild(o); });
      grid.insertBefore(frag, empty && empty.parentNode === grid ? empty : null);
      applyFilter();
      changed();
    }

    function setNote(text) {
      note = text == null || text === '' ? null : String(text);
      applyFilter();
    }

    // ----- Filter history -----
    // hKey is cleared when the markup lacks a part, so every function below
    // can rely on filter, hBtn, hPop and hList once hKey is set.
    var hKey = ds.historyKey || '';
    var hSize = parseInt(ds.historySize, 10);
    if (!(hSize > 0)) hSize = 10;
    var hWrap = hKey ? root.querySelector('.mcb-filter-wrap') : null;
    var hBtn = hWrap && hWrap.querySelector('.mcb-history-btn');
    var hPop = hWrap && hWrap.querySelector('.mcb-history');
    var hList = hPop && hPop.querySelector('.mcb-history-list');
    var hClear = hPop && hPop.querySelector('.mcb-history-clear');
    if (!filter || !hBtn || !hPop || !hList) hKey = '';
    var hTitle = hBtn ? hBtn.title : '';  // the server's wording, put back whenever there is a list
    var hMem = [];                        // the list when storage can't be read

    // histId is a filter text's identity: its terms as matched, so "fm, om"
    // and "FM,om" are one entry. '' = nothing worth remembering.
    function histId(text) {
      var parts = String(text == null ? '' : text).split(','), out = [];
      for (var i = 0; i < parts.length; i++) {
        var t = norm(parts[i]);
        if (t) out.push(t);
      }
      return out.join(',');
    }
    // readHistory returns the remembered filters, newest first. Anything
    // stored that isn't a usable string is skipped rather than trusted.
    function readHistory() {
      if (!hKey) return [];
      try {
        var s = JSON.parse(localStorage.getItem(hKey) || '[]');
        hMem = Array.isArray(s) ? s : [];
      } catch (e) {}  // no storage, or not JSON: keep this page's own copy
      var out = [];
      for (var i = 0; i < hMem.length && out.length < hSize; i++) {
        if (typeof hMem[i] === 'string' && histId(hMem[i])) out.push(hMem[i]);
      }
      return out;
    }
    function writeHistory(list) {
      hMem = list;
      try {
        if (list.length) localStorage.setItem(hKey, JSON.stringify(list));
        else localStorage.removeItem(hKey);
      } catch (e) {}
      syncHistoryBtn(list);
    }
    // remember puts text first in the list. The text is kept as typed
    // (trimmed), not rebuilt from its terms, so a label holding a comma
    // ("Smith, Jo") comes back exactly as Enter's whole-text try needs it.
    function remember(text) {
      if (!hKey) return;
      var t = String(text == null ? '' : text).trim(), id = histId(t);
      if (!id) return;
      var list = readHistory().filter(function (x) { return histId(x) !== id; });
      list.unshift(t);
      writeHistory(list.slice(0, hSize));
    }
    function syncHistoryBtn(list) {
      hBtn.disabled = !list.length;
      hBtn.title = list.length ? hTitle : 'No recent filters yet';
    }

    function historyOpen() { return !!hKey && !hPop.hidden; }
    // The entries are rebuilt on every open, from storage. The one matching
    // the filter now is aria-selected, a hook for styling it.
    function drawHistory(list) {
      while (hList.firstChild) hList.removeChild(hList.firstChild);
      var cur = histId(filter.value);
      list.forEach(function (t) {
        var it = document.createElement('div');
        it.className = 'mcb-history-item';
        it.setAttribute('role', 'option');
        it.setAttribute('aria-selected', histId(t) === cur ? 'true' : 'false');
        it.tabIndex = -1;  // focusable by script and click, not a Tab stop
        it.textContent = t;
        it.title = t;      // a long one is cut with an ellipsis
        hList.appendChild(it);
      });
    }
    function openHistory(focusFirst) {
      var list = readHistory();
      syncHistoryBtn(list);
      if (!list.length) { closeHistory(false); return; }
      drawHistory(list);
      if (hPop.hidden) {
        hPop.hidden = false;
        hBtn.setAttribute('aria-expanded', 'true');
        // Capture, so a handler that stops propagation can't keep it open.
        document.addEventListener('pointerdown', onOutside, true);
      }
      if (focusFirst) hList.firstChild.focus();
    }
    function closeHistory(refocus) {
      if (historyOpen()) {
        hPop.hidden = true;
        hBtn.setAttribute('aria-expanded', 'false');
        document.removeEventListener('pointerdown', onOutside, true);
      }
      if (refocus) filter.focus();
    }
    // A press outside the dropdown and its button closes it. The filter
    // counts as outside: clicking into it means typing, not picking.
    function onOutside(ev) {
      if (!hPop.contains(ev.target) && !hBtn.contains(ev.target)) closeHistory(false);
    }
    function useHistory(text) {
      closeHistory(true);
      filter.value = text;
      filter.dispatchEvent(new Event('input', { bubbles: true }));
      remember(text);
    }
    // forget drops one entry and keeps the focus at the same place in the
    // list (the next entry, or the new last one), so Delete can be pressed
    // again and again.
    function forget(item) {
      var at = Array.prototype.indexOf.call(hList.children, item);
      var id = histId(item.textContent);
      var list = readHistory().filter(function (x) { return histId(x) !== id; });
      writeHistory(list);
      if (!list.length) { closeHistory(true); return; }
      drawHistory(list);
      hList.children[Math.min(at, list.length - 1)].focus();
    }
    // The dropdown's arrow-key stops: each entry, then Clear history.
    function stops() {
      var s = Array.prototype.slice.call(hList.children);
      if (hClear) s.push(hClear);
      return s;
    }

    // Restore before the first draw, so the remembered ticks are what is shown.
    if (key) {
      try {
        var saved = JSON.parse(localStorage.getItem(key) || 'null');
        if (Array.isArray(saved)) {
          if (grid.querySelector('.mcb-cb')) tickOnly(asSet(saved));
          else pending = asSet(saved);
        }
      } catch (e) {}
    }

    root.addEventListener('change', function (ev) {
      // A box ticked or unticked by hand while a filter is in the box: the
      // filter was used to find it. (The root's own change, from changed(),
      // has the root as its target, and the filter's change is not a box.)
      if (ev.target !== root && ev.target.classList && ev.target.classList.contains('mcb-cb')) {
        remember(filter && filter.value);
      }
      renderSummary();
      save();
    });

    if (filter) {
      filter.addEventListener('input', applyFilter);
      filter.addEventListener('keydown', function (ev) {
        if (ev.isComposing) return;
        // ↓ (or Alt+↓, the usual "open the list" key) opens the history.
        // Only when it has entries, and not with Shift, Ctrl or Meta, which
        // move or select in the text; otherwise ↓ does what it always did.
        if (ev.key === 'ArrowDown' && hKey && !ev.shiftKey && !ev.ctrlKey && !ev.metaKey) {
          if (readHistory().length) {
            ev.preventDefault();
            openHistory(true);
          }
          return;
        }
        if (ev.key !== 'Enter') return;
        ev.preventDefault();
        remember(filter.value);  // ignores an empty or commas-only filter
        var ts = terms();
        if (!ts.length || note !== null) return;  // nothing to pick from while a note replaces the boxes
        // The whole text as one term first: with no comma it is the only
        // term anyway, and with one it lets a label that holds a comma be
        // named in full. Only when it names nothing is the list split.
        var whole = pickFor(norm(filter.value));
        var picks = whole ? [whole] : ts.map(pickFor);
        var any = false;
        picks.forEach(function (cb) {
          if (cb && !cb.checked) { cb.checked = true; any = true; }
        });
        if (any) changed();  // one change for the whole press, however many were ticked
      });
    }

    var shownBtn = root.querySelector('.mcb-tick-shown');
    if (shownBtn) shownBtn.addEventListener('click', function () {
      remember(filter && filter.value);
      if (note !== null) return;  // a note hides every box, so none is "shown"
      var ts = terms(), any = false;
      each(function (cb) {
        if (!cb.disabled && !cb.checked && hits(cb, ts)) { cb.checked = true; any = true; }
      });
      if (any) changed();
    });

    var clearBtn = root.querySelector('.mcb-clear');
    if (clearBtn) clearBtn.addEventListener('click', function () {
      var any = false;
      each(function (cb) {
        if (!cb.disabled && cb.checked) { cb.checked = false; any = true; }
      });
      applyFilter();
      if (any) changed();
    });

    if (hKey) {
      hBtn.addEventListener('click', function (ev) {
        if (historyOpen()) { closeHistory(false); return; }
        // detail is 0 for a click made with Enter or Space: the focus goes
        // into the list then, as a keyboard user expects. A mouse click
        // leaves it on the button.
        openHistory(ev.detail === 0);
      });
      hBtn.addEventListener('keydown', function (ev) {
        if (ev.key === 'ArrowDown') { ev.preventDefault(); openHistory(true); }
      });
      // Another tab may have added the first entry since the button was
      // disabled; focusing the filter is when it would be wanted.
      filter.addEventListener('focus', function () { syncHistoryBtn(readHistory()); });
      hList.addEventListener('click', function (ev) {
        var it = ev.target.closest ? ev.target.closest('.mcb-history-item') : null;
        if (it) useHistory(it.textContent);
      });
      if (hClear) hClear.addEventListener('click', function () {
        writeHistory([]);
        closeHistory(true);
      });
      hPop.addEventListener('keydown', function (ev) {
        var s = stops(), at = s.indexOf(ev.target);
        var item = ev.target.classList && ev.target.classList.contains('mcb-history-item') ? ev.target : null;
        switch (ev.key) {
          case 'ArrowDown':
          case 'ArrowUp':
            // Round the ends, so ↑ from the first entry reaches Clear history.
            ev.preventDefault();
            s[at < 0 ? 0 : (at + (ev.key === 'ArrowDown' ? 1 : s.length - 1)) % s.length].focus();
            break;
          case 'Home': ev.preventDefault(); s[0].focus(); break;
          case 'End': ev.preventDefault(); s[s.length - 1].focus(); break;
          case 'Enter':
          case ' ':
            // On Clear history these are the button's own click.
            if (item) { ev.preventDefault(); useHistory(item.textContent); }
            break;
          case 'Delete':
          case 'Backspace':  // a Mac's Delete key
            if (item) { ev.preventDefault(); forget(item); }
            break;
          case 'Escape':
            // Stopped, so a dialog the component sits in doesn't close too.
            ev.preventDefault();
            ev.stopPropagation();
            closeHistory(true);
            break;
          case 'Tab':
            // Onto a control that is still on screen. Left to the browser,
            // Tab would start from an entry being hidden under it.
            ev.preventDefault();
            closeHistory(false);
            (ev.shiftKey ? filter : hBtn).focus();
            break;
        }
      });
      syncHistoryBtn(readHistory());
    }

    var js = root.querySelectorAll('.mcb-js');
    for (var i = 0; i < js.length; i++) js[i].hidden = false;

    applyFilter();
    renderSummary();

    api.values = values;
    api.setValues = setValues;
    api.setOptions = setOptions;
    api.setNote = setNote;
    api.refresh = function () { applyFilter(); renderSummary(); };
    root.mcb = api;
    return api;
  }
})(document.currentScript);
`
