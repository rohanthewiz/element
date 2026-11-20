package components

import (
	"bytes"
	"context"
	"encoding/json"
	"html"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/rohanthewiz/element"
)

// TestMultiCheckbox_Browser runs MultiCheckbox's inline script (multiCheckboxJS)
// in headless Chrome. The Go tests above can only see the markup; this is where
// the filter, Tick shown, Clear, Enter, the summary, storage, setOptions,
// setNote and the hooks are actually exercised.
//
// How it works, with no browser-automation dependency:
//
//	browserTestPage() ──► page.html ──► chrome --headless --dump-dom ──► DOM after load
//	   three instances                  (runs every <script>)             │
//	   + mcbBrowserTestJS                                                 ▼
//	     (drives them, writes                                  <pre id="results">[{name, pass, detail}…]</pre>
//	      JSON results into                                               │
//	      <pre id="results">)                                             ▼
//	                                                      one subtest per check
//
// The checks run synchronously at the end of <body>, so the DOM Chrome dumps
// already holds every result. The script's last entry is a sentinel; if it
// is missing, the script threw or stopped early, and the test says so rather
// than passing on the checks it did reach.
//
// It skips under -short and when no Chrome or Chromium can be found (set
// CHROME_BIN to point at one). Each run gets a fresh profile in t.TempDir(),
// so localStorage starts empty.
func TestMultiCheckbox_Browser(t *testing.T) {
	if testing.Short() {
		t.Skip("browser test skipped with -short")
	}
	chrome := findChrome()
	if chrome == "" {
		t.Skip("no Chrome or Chromium found; set CHROME_BIN to run the browser test")
	}

	dir := t.TempDir()
	page := filepath.Join(dir, "page.html")
	if err := os.WriteFile(page, []byte(browserTestPage()), 0o644); err != nil {
		t.Fatal(err)
	}

	args := []string{
		"--headless=new", "--disable-gpu", "--no-first-run", "--no-default-browser-check",
		"--user-data-dir=" + filepath.Join(dir, "profile"),
		// Lets the page settle before the dump; the checks need no timers,
		// so this only bounds how long Chrome may take.
		"--virtual-time-budget=5000",
		"--dump-dom", fileURL(page),
	}
	// Chrome refuses to start sandboxed as root, which is how containers
	// usually run CI. The page is this test's own file, nothing fetched.
	if os.Geteuid() == 0 {
		args = append([]string{"--no-sandbox"}, args...)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
	defer cancel()
	var stdout, stderr bytes.Buffer
	cmd := exec.CommandContext(ctx, chrome, args...)
	cmd.Stdout, cmd.Stderr = &stdout, &stderr
	if err := cmd.Run(); err != nil {
		t.Fatalf("running %s: %v\nstderr: %s", chrome, err, tail(stderr.String(), 2000))
	}

	results, err := parseBrowserResults(stdout.String())
	if err != nil {
		t.Fatalf("%v\nDOM tail: %s", err, tail(stdout.String(), 2000))
	}
	for _, r := range results {
		r := r
		t.Run(r.Name, func(t *testing.T) {
			if !r.Pass {
				t.Errorf("got %s", r.Detail)
			}
		})
	}
	if last := results[len(results)-1]; last.Name != browserTestSentinel {
		t.Fatalf("the test script stopped early; last result: %+v", last)
	}
}

// browserTestSentinel is the name of the script's last result.
const browserTestSentinel = "script reached the end"

// browserCheck is one result written by mcbBrowserTestJS.
type browserCheck struct {
	Name   string `json:"name"`
	Pass   bool   `json:"pass"`
	Detail string `json:"detail"` // JSON of what was seen, set on failure
}

var resultsRe = regexp.MustCompile(`(?s)<pre id="results">(.*?)</pre>`)

// parseBrowserResults pulls the JSON results out of the dumped DOM. The
// serializer escapes &, < and > in text, so the JSON is unescaped first.
func parseBrowserResults(dom string) ([]browserCheck, error) {
	m := resultsRe.FindStringSubmatch(dom)
	if m == nil || strings.TrimSpace(m[1]) == "" {
		return nil, errorString("no results in the dumped DOM: the test script did not run")
	}
	var results []browserCheck
	if err := json.Unmarshal([]byte(html.UnescapeString(m[1])), &results); err != nil {
		return nil, err
	}
	if len(results) == 0 {
		return nil, errorString("the test script recorded no checks")
	}
	return results, nil
}

type errorString string

func (e errorString) Error() string { return string(e) }

// findChrome returns a Chrome or Chromium binary: CHROME_BIN, then the usual
// names on PATH, then the macOS app bundles. "" when there is none.
func findChrome() string {
	if p := os.Getenv("CHROME_BIN"); p != "" {
		return p
	}
	for _, name := range []string{"google-chrome", "google-chrome-stable", "chromium", "chromium-browser", "chrome"} {
		if p, err := exec.LookPath(name); err == nil {
			return p
		}
	}
	for _, p := range []string{
		"/Applications/Google Chrome.app/Contents/MacOS/Google Chrome",
		"/Applications/Chromium.app/Contents/MacOS/Chromium",
	} {
		if _, err := os.Stat(p); err == nil {
			return p
		}
	}
	return ""
}

// fileURL makes a file:// URL from a local path, including a Windows one
// (C:\x → file:///C:/x).
func fileURL(path string) string {
	p := filepath.ToSlash(path)
	if !strings.HasPrefix(p, "/") {
		p = "/" + p
	}
	return "file://" + p
}

// tail returns at most the last n bytes of s, for error messages.
func tail(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return "…" + s[len(s)-n:]
}

// browserTestPage renders the instances the script drives:
//
//	a  in a <form>: Name, counts, MaxCount 10, a disabled box, one Checked
//	   box, and a long label that contains a short one ("…/acme/prod/" and
//	   "acme/prod/"), so Enter's exact match has something to beat
//	b  StorageKey over a Checked box, default wording; storage is seeded by
//	   a script just before it
//	c  starts empty with StorageKey (restore waits for setOptions); used for
//	   setNote, the normalize and total hooks, and a late data-max-count
//	d  labels that hold commas ("Smith, Jo"), for comma-separated filter
//	   terms that must not split a label typed in full
func browserTestPage() string {
	b := element.NewBuilder()
	b.Html().R(
		b.Head().R(
			b.Meta("charset", "utf-8").R(),
			b.Title().T("MultiCheckbox browser test"),
		),
		b.Body().R(
			b.Form("id", "form-a").R(
				MultiCheckbox{
					ID: "a", Name: "folder", Label: "Folders:", FilterLabel: "Filter folders:",
					Noun: "folder", CountNoun: "DAG", MaxCount: 10,
					Options: []CheckOption{
						{Value: "exports/db_to_db/acme/prod", Label: "exports/db_to_db/acme/prod/", Count: 2},
						{Value: "acme/prod", Label: "acme/prod/", Count: 1, Checked: true},
						{Value: "sales/q1", Label: "sales/q1/", Count: 12},
						{Value: "sales/q2", Label: "sales/q2/", Count: 3},
						{Value: "sales/locked", Label: "sales/locked/", Count: 5, Disabled: true},
					},
				}.Render(b),
			),
			b.Script().T(`localStorage.setItem('mcb-test-b', JSON.stringify(['y', 'zz']));`),
			MultiCheckbox{
				ID: "b", Name: "tag", StorageKey: "mcb-test-b",
				Options: []CheckOption{{Value: "x", Checked: true}, {Value: "y"}, {Value: "z"}},
			}.Render(b),
			b.Script().T(`localStorage.setItem('mcb-test-c', JSON.stringify(['gone', 'p/q']));`),
			MultiCheckbox{ID: "c", StorageKey: "mcb-test-c", Noun: "folder", CountNoun: "DAG"}.Render(b),
			MultiCheckbox{
				ID: "d",
				Options: []CheckOption{
					{Value: "Smith, Jo"}, {Value: "Smith, Al"}, {Value: "Jo"}, {Value: "Al Doe"},
				},
			}.Render(b),
			b.Pre("id", "results").R(),
			b.Script().T(mcbBrowserTestJS),
		),
	)
	return b.String()
}

// mcbBrowserTestJS drives the three instances and records one result per
// check. Each instance's checks are a function, run in order: later checks
// build on the state earlier ones leave (what is ticked, what the filter
// says), so read them top to bottom.
//
// Events are dispatched the way a browser would deliver them to the
// component's listeners (input on the filter, keydown Enter, click). A
// synthetic Enter never triggers a form's implicit submission, so "does not
// submit" is checked as defaultPrevented, which is what stops a real one.
const mcbBrowserTestJS = `(function () {
  'use strict';
  var results = [];
  function check(name, pass, detail) {
    results.push({ name: name, pass: !!pass, detail: pass ? '' : JSON.stringify(detail === undefined ? null : detail) });
  }

  // The values of the boxes on screen: not hidden by the filter or a note.
  function shownIn(root) {
    return Array.prototype.filter.call(root.querySelectorAll('.mcb-opt'), function (l) {
      return !l.hidden && getComputedStyle(l).display !== 'none';
    }).map(function (l) { return l.querySelector('.mcb-cb').value; });
  }
  function typeInto(input, v) { input.value = v; input.dispatchEvent(new Event('input', { bubbles: true })); }
  function pressEnter(input) {
    var e = new KeyboardEvent('keydown', { key: 'Enter', bubbles: true, cancelable: true });
    input.dispatchEvent(e);
    return e.defaultPrevented;
  }
  function same(a, b) { return JSON.stringify(a) === JSON.stringify(b); }
  // Tag, classes and attribute names, recursively: the markup's shape
  // without its values, to compare a server-rendered box with a built one.
  function shape(node) {
    if (node.nodeType !== 1) return '#text';
    var attrs = Array.prototype.map.call(node.attributes, function (a) {
      return a.name === 'class' ? 'class=' + a.value : a.name;
    }).sort();
    return node.tagName.toLowerCase() + '[' + attrs.join(',') + '](' +
      Array.prototype.map.call(node.childNodes, shape).join(' ') + ')';
  }

  try {
    instanceA();
    instanceB();
    instanceC();
    instanceD();
    check('` + browserTestSentinel + `', true);
  } catch (e) {
    check('script threw', false, String((e && e.stack) || e));
  }
  document.getElementById('results').textContent = JSON.stringify(results);

  function instanceA() {
    var A = document.getElementById('a'), api = A.mcb;
    var filter = A.querySelector('.mcb-filter'), summary = A.querySelector('.mcb-summary');
    var empty = A.querySelector('.mcb-empty');
    var serverOpt = A.querySelector('.mcb-cb[value="sales/q1"]').closest('.mcb-opt').cloneNode(true);
    var changes = 0;
    A.addEventListener('change', function () { changes++; });
    function shown() { return shownIn(A); }
    function type(v) { typeInto(filter, v); }

    // Start-up
    check('A: instance on root.mcb, init is idempotent', !!api && ElementMultiCheckbox.init(A) === api);
    check('A: script-only controls revealed', A.querySelectorAll('.mcb-js[hidden]').length === 0 &&
      getComputedStyle(filter).display !== 'none');
    check('A: initial summary', summary.textContent === '1 folder ticked \u00b7 1 DAG', summary.textContent);
    check('A: initial ticks from Checked', same(api.values(), ['acme/prod']), api.values());
    check('A: every box shown with no filter', shown().length === 5, shown());
    check('A: grid layout applied', getComputedStyle(A.querySelector('.mcb-grid')).display === 'grid');

    // Filter
    type('acme');
    check('A: filter narrows by label', same(shown(), ['exports/db_to_db/acme/prod', 'acme/prod']), shown());
    type('SALES/Q');
    check('A: filter is case-blind, a ticked non-match stays shown',
      same(shown(), ['acme/prod', 'sales/q1', 'sales/q2']), shown());
    type('zzz');
    check('A: no match leaves only the ticked box, note hidden', same(shown(), ['acme/prod']) && empty.hidden,
      [shown(), empty.hidden]);
    var acme = A.querySelector('.mcb-cb[value="acme/prod"]');
    acme.click();
    check('A: an unticked box stays until the next keystroke', !acme.closest('.mcb-opt').hidden && changes === 1, changes);
    check('A: summary with nothing ticked', summary.textContent === 'No folder ticked', summary.textContent);
    type('zzz');
    check('A: the note says nothing matches', !empty.hidden && empty.textContent === 'No folders match the filter.' &&
      getComputedStyle(empty).display !== 'none', [empty.hidden, empty.textContent]);

    // Tick shown
    type('sales/');
    changes = 0;
    A.querySelector('.mcb-tick-shown').click();
    check('A: Tick shown adds the matches, skips disabled, one change',
      same(api.values(), ['sales/q1', 'sales/q2']) && changes === 1, [api.values(), changes]);
    check('A: past MaxCount, warning text and class',
      summary.textContent === '2 folders ticked \u00b7 15 DAGs (at most 10)' && summary.classList.contains('mcb-over'),
      summary.textContent);
    var submitted = new FormData(document.getElementById('form-a')).getAll('folder');
    check('A: the form submits the ticked values', same(submitted, ['sales/q1', 'sales/q2']), submitted);

    // Clear
    changes = 0;
    A.querySelector('.mcb-clear').click();
    check('A: Clear unticks all, one change, warning gone',
      api.values().length === 0 && changes === 1 && !summary.classList.contains('mcb-over'), [api.values(), changes]);
    check('A: Clear re-applies the filter', same(shown(), ['sales/q1', 'sales/q2', 'sales/locked']), shown());
    changes = 0;
    A.querySelector('.mcb-clear').click();
    check('A: Clear with nothing ticked fires no change', changes === 0, changes);

    // Enter
    type('sales/');
    check('A: Enter on an ambiguous filter is prevented and ticks nothing',
      pressEnter(filter) && api.values().length === 0, api.values());
    type('q2');
    check('A: Enter ticks the only match', pressEnter(filter) && same(api.values(), ['sales/q2']), api.values());
    type('SALES/Q1/');
    pressEnter(filter);
    check('A: Enter ticks an exact label, case-blind', same(api.values(), ['sales/q1', 'sales/q2']), api.values());
    type('acme/prod');
    pressEnter(filter);
    check('A: Enter prefers an exact value over several matches',
      same(api.values(), ['acme/prod', 'sales/q1', 'sales/q2']), api.values());
    type('sales/locked/');
    pressEnter(filter);
    check('A: Enter never ticks a disabled box', api.values().indexOf('sales/locked') < 0, api.values());

    // Comma-separated terms. Cleared first, so the shown lists below are
    // the filter's doing and not ticked boxes kept on screen.
    A.querySelector('.mcb-clear').click();
    type('q1, exports');
    check('A: a comma list shows the matches of any term', same(shown(), ['exports/db_to_db/acme/prod', 'sales/q1']), shown());
    type('Q2,, LOCKED ');
    check('A: terms are trimmed and case-blind, empty terms dropped', same(shown(), ['sales/q2', 'sales/locked']), shown());
    type(' , ,');
    check('A: only commas is no filter', shown().length === 5, shown());
    type('q1, q2, locked');
    changes = 0;
    A.querySelector('.mcb-tick-shown').click();
    check('A: Tick shown ticks the matches of every term, skips disabled, one change',
      same(api.values(), ['sales/q1', 'sales/q2']) && changes === 1, [api.values(), changes]);
    A.querySelector('.mcb-clear').click();
    type('q1, acme/prod, sales/');
    changes = 0;
    check('A: Enter on a comma list is prevented', pressEnter(filter), null);
    check('A: Enter ticks the exact or only match of each term, skips an ambiguous one, one change',
      same(api.values(), ['acme/prod', 'sales/q1']) && changes === 1, [api.values(), changes]);
    type('zzz, Q2');
    pressEnter(filter);
    check('A: Enter adds to the ticks, a term matching nothing is skipped',
      same(api.values(), ['acme/prod', 'sales/q1', 'sales/q2']), api.values());
    changes = 0;
    pressEnter(filter);
    check('A: Enter that ticks nothing new fires no change', changes === 0, changes);

    // Script API
    type('');
    api.setValues(['sales/q2', 'not-an-option']);
    check('A: setValues ticks exactly those listed', same(api.values(), ['sales/q2']), api.values());
    changes = 0;
    api.setOptions([
      { value: 'sales/q2', label: 'sales/q2/', count: 3 },
      { value: 'new', label: '<b>new</b>', count: 7, checked: true },
      { value: 'd', disabled: true, checked: true, count: 1 }
    ]);
    check('A: setOptions keeps present ticks, applies checked, one change',
      same(api.values(), ['sales/q2', 'new', 'd']) && changes === 1, [api.values(), changes]);
    var built = A.querySelector('.mcb-cb[value="new"]').closest('.mcb-opt');
    check('A: setOptions writes labels as text',
      built.querySelector('.mcb-name').textContent === '<b>new</b>' && !built.querySelector('b'), built.outerHTML);
    check('A: setOptions builds the same markup shape as the server', shape(built) === shape(serverOpt),
      [shape(built), shape(serverOpt)]);
    check('A: summary after setOptions', summary.textContent === '3 folders ticked \u00b7 11 DAGs (at most 10)',
      summary.textContent);
    api.setOptions([]);
    check('A: an empty list says there are none', !empty.hidden && empty.textContent === 'No folders.',
      [empty.hidden, empty.textContent]);
  }

  function instanceB() {
    var B = document.getElementById('b'), summary = B.querySelector('.mcb-summary');
    check('B: remembered ticks win over Checked, unknown values ignored', same(B.mcb.values(), ['y']), B.mcb.values());
    check('B: default wording', summary.textContent === '1 item ticked', summary.textContent);
    B.querySelector('.mcb-cb[value="z"]').click();
    check('B: ticks are saved on change', localStorage.getItem('mcb-test-b') === '["y","z"]',
      localStorage.getItem('mcb-test-b'));
    check('B: one shared definition, separate instances',
      typeof ElementMultiCheckbox.init === 'function' && B.mcb !== document.getElementById('a').mcb);
  }

  function instanceC() {
    var C = document.getElementById('c'), api = C.mcb;
    var filter = C.querySelector('.mcb-filter'), summary = C.querySelector('.mcb-summary');
    var empty = C.querySelector('.mcb-empty');
    function shown() { return shownIn(C); }
    function type(v) { typeInto(filter, v); }
    var pq = { value: 'p/q', label: 'p/q/', count: 4 }, pr = { value: 'p/r', label: 'p/r/', count: 2 };
    var s = { value: 's', label: 's/', count: 1 };

    // An empty grid filled later
    check('C: an empty grid starts with the no-options note', !empty.hidden && empty.textContent === 'No folders.',
      empty.textContent);
    api.setNote('Loading folders\u2026');
    check('C: setNote on an empty grid', !empty.hidden && empty.textContent === 'Loading folders\u2026', empty.textContent);
    api.setOptions([pq, pr, s]);
    check('C: the first setOptions restores remembered ticks and clears the note',
      same(api.values(), ['p/q']) && empty.hidden, [api.values(), empty.hidden]);
    check('C: storage is rewritten after the restore', localStorage.getItem('mcb-test-c') === '["p/q"]',
      localStorage.getItem('mcb-test-c'));
    api.setOptions([pq, s]);
    check('C: a later setOptions does not re-apply storage', same(api.values(), ['p/q']), api.values());

    // normalize hook: strip a pasted "repo/" prefix and a trailing slash
    api.normalize = function (t) {
      t = String(t).trim().replace(/\/+$/, '');
      return t.indexOf('repo/') === 0 ? t.slice(5) : t;
    };
    type('repo/s/');
    check('C: normalize shapes the filter', same(shown(), ['p/q', 's']), shown());
    check('C: normalize shapes the Enter exact match', pressEnter(filter) && same(api.values(), ['p/q', 's']),
      api.values());
    type('');

    // total hook and a data-max-count set after load
    check('C: the summary sums Counts by default', summary.textContent === '2 folders ticked \u00b7 5 DAGs',
      summary.textContent);
    var per = 100;
    api.total = function (values) { return values.length * per; };
    C.setAttribute('data-max-count', '150');
    api.refresh();
    check('C: total hook and a later data-max-count',
      summary.textContent === '2 folders ticked \u00b7 200 DAGs (at most 150)' && summary.classList.contains('mcb-over'),
      summary.textContent);
    per = 1;
    api.refresh();
    check('C: refresh picks up a change in the hook inputs',
      summary.textContent === '2 folders ticked \u00b7 2 DAGs' && !summary.classList.contains('mcb-over'),
      summary.textContent);
    api.total = null;
    api.refresh();

    // Notes: p/r is back and unticked, so Tick shown or Enter would tick
    // it if they ignored the note.
    api.setOptions([pq, pr, s]);
    api.setNote('Loading folders\u2026');
    check('C: a note hides every box and keeps the ticks',
      shown().length === 0 && !empty.hidden && empty.textContent === 'Loading folders\u2026' && same(api.values(), ['p/q', 's']),
      [shown(), api.values()]);
    type('p');
    check('C: typing keeps the note up', shown().length === 0 && empty.textContent === 'Loading folders\u2026', shown());
    C.querySelector('.mcb-tick-shown').click();
    type('p/r');
    pressEnter(filter);
    type('');
    check('C: Tick shown and Enter do nothing under a note', same(api.values(), ['p/q', 's']), api.values());
    api.setOptions([pq, pr, { value: 't', label: 't/', count: 9 }]);
    check('C: setOptions drops the ticks of values no longer listed', same(api.values(), ['p/q']), api.values());
    check('C: setOptions clears the note', empty.hidden && shown().length === 3, [empty.hidden, shown()]);
    api.setNote('x');
    api.setNote(null);
    check('C: setNote(null) brings the boxes back', empty.hidden && shown().length === 3, shown());
  }

  // D: labels that hold commas, which a comma list could otherwise never
  // name in full, and the normalize hook applied to each term.
  function instanceD() {
    var D = document.getElementById('d'), api = D.mcb;
    var filter = D.querySelector('.mcb-filter');
    function shown() { return shownIn(D); }
    function type(v) { typeInto(filter, v); }

    // Split, "smith" would be ambiguous and "jo" would name the box Jo.
    type('Smith, Jo');
    pressEnter(filter);
    check('D: Enter on the whole text ticks a label holding a comma, not each part',
      same(api.values(), ['Smith, Jo']), api.values());
    D.querySelector('.mcb-clear').click();
    type('smith, doe');
    check('D: the display still splits the text into terms', same(shown(), ['Smith, Jo', 'Smith, Al', 'Al Doe']), shown());

    // A hook that strips a pasted "x:" prefix, as a page would strip a path prefix.
    api.normalize = function (t) {
      t = String(t).trim();
      return t.indexOf('x:') === 0 ? t.slice(2) : t;
    };
    type('x:jo, x:doe');
    check('D: normalize applies to every term', same(shown(), ['Smith, Jo', 'Jo', 'Al Doe']), shown());
    pressEnter(filter);
    check('D: Enter ticks each normalized term (exact, then only match)', same(api.values(), ['Jo', 'Al Doe']), api.values());
    api.normalize = null;
  }
})();
`
