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
//   - Tick shown ticks every box the filter matches (all of them with no
//     filter), adding to what is already ticked. Clear unticks every box.
//   - The summary counts the ticked boxes and, with CountNoun set, totals
//     their Counts ("2 folders ticked · 37 DAGs"). Past MaxCount it adds
//     "(at most N)" and the class mcb-over. It is aria-live, so screen readers
//     hear the count change.
//   - Enter in the filter ticks the box it names exactly (value or label), or
//     failing that its only match, and never submits the surrounding form.
//     To make Enter also "go", add a keydown listener to .mcb-filter with
//     addEventListener: it runs after this one, so it sees the new tick.
//   - Every change, including the buttons', ends in one bubbling "change"
//     event on the component's root element, so a caller needs one listener.
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
//	normalize = function (text) → text   the filter text as matched, e.g. strip a
//	                                     pasted prefix; used by the filter and Enter
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
// Style hooks: .multi-checkbox, .mcb-toolbar, .mcb-filter-label, .mcb-filter,
// .mcb-summary (.mcb-over past MaxCount), .mcb-btn (.mcb-tick-shown,
// .mcb-clear), .mcb-body, .mcb-label, .mcb-grid, .mcb-opt, .mcb-cb, .mcb-name,
// .mcb-count, .mcb-empty.
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
				attrs := []string{"type", "search", "id", id + "-filter", "class", "mcb-filter mcb-js",
					"hidden", "hidden", "autocomplete", "off", "aria-controls", id + "-grid"}
				if m.FilterPlaceholder != "" {
					attrs = append(attrs, "placeholder", m.FilterPlaceholder)
				}
				b.Input(attrs...)
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
// Flow:
//
//	filter input ──► applyFilter (show = ticked || label matches; a note hides all)
//	box click ─────────────────────────────┐
//	Tick shown / Clear / Enter / setValues ├─► "change" on root ─► renderSummary + save
//	setOptions ────────────────────────────┘
//
// The hooks live on the instance object (api), which exists before any
// handler runs, so a caller assigning api.total after init is seen by the
// next summary without re-wiring anything.
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
    function needle() {
      if (!filter) return '';
      var t = typeof api.normalize === 'function' ? api.normalize(filter.value) : filter.value;
      return String(t == null ? '' : t).trim().toLowerCase();
    }
    function hits(cb, f) { return !f || nameOf(cb).indexOf(f) >= 0; }
    function changed() { root.dispatchEvent(new Event('change', { bubbles: true })); }
    function asSet(list) {
      var s = Object.create(null);
      (Array.isArray(list) ? list : []).forEach(function (v) { s[String(v)] = true; });
      return s;
    }

    // A ticked box is always shown, so a choice can't be hidden by the filter.
    function applyFilter() {
      var f = needle(), total = 0, shown = 0;
      each(function (cb) {
        var show = note === null && (cb.checked || hits(cb, f));
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

    root.addEventListener('change', function () { renderSummary(); save(); });

    if (filter) {
      filter.addEventListener('input', applyFilter);
      filter.addEventListener('keydown', function (ev) {
        if (ev.key !== 'Enter' || ev.isComposing) return;
        ev.preventDefault();
        var f = needle();
        if (!f || note !== null) return;  // nothing to pick from while a note replaces the boxes
        var exact = null, only = null, matches = 0;
        each(function (cb) {
          if (exact || cb.disabled) return;
          if (cb.value.toLowerCase() === f || nameOf(cb) === f) { exact = cb; return; }
          if (hits(cb, f)) { matches++; only = cb; }
        });
        var pick = exact || (matches === 1 ? only : null);
        if (pick && !pick.checked) { pick.checked = true; changed(); }
      });
    }

    var shownBtn = root.querySelector('.mcb-tick-shown');
    if (shownBtn) shownBtn.addEventListener('click', function () {
      if (note !== null) return;  // a note hides every box, so none is "shown"
      var f = needle(), any = false;
      each(function (cb) {
        if (!cb.disabled && !cb.checked && hits(cb, f)) { cb.checked = true; any = true; }
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
