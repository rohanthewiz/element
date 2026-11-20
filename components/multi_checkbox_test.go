package components

import (
	"strings"
	"testing"

	"github.com/rohanthewiz/element"
)

// renderMC renders a MultiCheckbox and returns the HTML.
func renderMC(m MultiCheckbox) string {
	b := element.NewBuilder()
	m.Render(b)
	return b.String()
}

// markupOf cuts the rendered HTML at the scoped <style>, leaving the markup
// alone. The inline script names the same classes and attributes the markup
// uses (.mcb-js, data-count, …), so counting or ruling them out has to look
// at the markup only.
func markupOf(html string) string {
	if i := strings.Index(html, "<style>"); i >= 0 {
		return html[:i]
	}
	return html
}

func TestMultiCheckbox_Render(t *testing.T) {
	got := renderMC(MultiCheckbox{
		ID:          "folders",
		Name:        "folder",
		Label:       "Folders:",
		FilterLabel: "Filter folders:",
		Noun:        "folder",
		CountNoun:   "DAG",
		MaxCount:    300,
		StorageKey:  "pickedFolders",
		Columns:     8,
		Options: []CheckOption{
			{Value: "exports/db_to_db/acme/emea/prod", Label: "exports/db_to_db/acme/emea/prod/", Count: 2},
			{Value: "acme/prod", Label: "acme/prod/", Count: 1, Checked: true},
			{Value: "sales/q1", Label: "sales/q1/", Count: 12},
			{Value: "locked", Count: 0, Disabled: true},
		},
		Selected: []string{"sales/q1"},
	})

	for _, want := range []string{
		// Root carries the script's settings
		`class="form-field multi-checkbox"`, `id="folders"`,
		`data-noun="folder"`, `data-nouns="folders"`, `data-ticked="ticked"`,
		`data-count-noun="DAG"`, `data-count-nouns="DAGs"`,
		`data-max-count="300"`, `data-storage-key="pickedFolders"`, `data-name="folder"`,
		// Toolbar: the script-only parts start hidden
		`<label class="mcb-filter-label mcb-js" for="folders-filter" hidden="hidden">Filter folders:</label>`,
		`type="search"`, `id="folders-filter"`, `aria-controls="folders-grid"`,
		`class="mcb-summary mcb-js"`, `aria-live="polite"`,
		`title="Tick every folder the filter shows"`, `>Tick shown</button>`,
		`title="Untick every folder"`, `>Clear</button>`,
		// Grid: labelled group of wrapped boxes
		`<span class="mcb-label" id="folders-label">Folders:</span>`,
		`role="group"`, `aria-labelledby="folders-label"`,
		`<label class="mcb-opt" title="acme/prod/">`,
		`value="acme/prod"`, `name="folder"`, `data-count="1"`,
		`<span class="mcb-name">acme/prod/</span><span class="mcb-count">(1)</span>`,
		// Label defaults to the value
		`<span class="mcb-name">locked</span><span class="mcb-count">(0)</span>`,
		// Empty note is hidden while there are options, worded for the filter
		`<span class="mcb-empty" hidden="hidden">No folders match the filter.</span>`,
		// Scoped layout rules and the script
		"#folders [hidden]{display:none!important}",
		":where(#folders .mcb-grid){display:grid;grid-template-columns:repeat(8,minmax(0,max-content))",
		"max-height:16rem",
		"window.ElementMultiCheckbox",
		"document.currentScript",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("MultiCheckbox.Render() missing %q\ngot: %s", want, got)
		}
	}

	// Ticked: one by Option.Checked, one by Selected. Not the others.
	if n := strings.Count(got, `checked="checked"`); n != 2 {
		t.Errorf("want 2 ticked boxes, got %d", n)
	}
	if !strings.Contains(got, `disabled="disabled"`) {
		t.Error("disabled option should render disabled")
	}

	// Every script-only control starts hidden: label, input, summary, two buttons.
	if n := strings.Count(markupOf(got), `mcb-js`); n != 5 {
		t.Errorf("want 5 script-only controls, got %d", n)
	}

	// The script must come last in the root, so currentScript.parentNode is
	// the instance.
	if !strings.HasSuffix(got, "})(document.currentScript);\n</script></div>") {
		t.Errorf("script should be the root's last child; tail: %q", got[len(got)-80:])
	}
}

// Option labels and values are data: a label is escaped as text, and a value
// or title survives the attribute round trip unchanged.
func TestMultiCheckbox_EscapesOptionData(t *testing.T) {
	got := renderMC(MultiCheckbox{
		Name: "x",
		Options: []CheckOption{
			{Value: `a"b&amp`, Label: `<img src=x onerror=alert(1)>`},
		},
	})

	if strings.Contains(got, "<img") {
		t.Errorf("label markup should be escaped, got: %s", got)
	}
	for _, want := range []string{
		`<span class="mcb-name">&lt;img src=x onerror=alert(1)&gt;</span>`,
		`value="a&#34;b&amp;amp"`,
		`title="&lt;img src=x onerror=alert(1)&gt;"`,
	} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\ngot: %s", want, got)
		}
	}
}

func TestMultiCheckbox_Defaults(t *testing.T) {
	got := renderMC(MultiCheckbox{
		Options: []CheckOption{{Value: "a"}, {Value: "b"}},
	})

	for _, want := range []string{
		`id="multi-checkbox"`,
		`data-noun="item"`, `data-nouns="items"`,
		`>Filter:</label>`,
		`title="Tick every item the filter shows"`,
		`<span class="mcb-name">a</span></label>`, // no count span without CountNoun
	} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\ngot: %s", want, got)
		}
	}
	for _, unwanted := range []string{
		"data-count", "mcb-count", "data-max-count", "data-storage-key", "data-name",
		` name=`,          // no Name, so nothing is submitted
		"aria-labelledby", // no Label to point at
		"placeholder",
	} {
		if strings.Contains(markupOf(got), unwanted) {
			t.Errorf("did not expect %q\ngot: %s", unwanted, got)
		}
	}
}

// With no ID the component takes its Name, so one-per-form needs no extra field.
func TestMultiCheckbox_IDDefaultsToName(t *testing.T) {
	got := renderMC(MultiCheckbox{Name: "tags"})
	if !strings.Contains(got, `id="tags"`) || !strings.Contains(got, "#tags [hidden]") {
		t.Errorf("ID should default to Name, got: %s", got)
	}
}

// No options at all: the note is shown, and says so rather than blaming the filter.
func TestMultiCheckbox_NoOptions(t *testing.T) {
	got := renderMC(MultiCheckbox{Noun: "folder"})
	if !strings.Contains(got, `<span class="mcb-empty">No folders.</span>`) {
		t.Errorf("want a visible empty note, got: %s", got)
	}
}

// Custom button text drops the default tooltip, which would no longer match it.
func TestMultiCheckbox_CustomWording(t *testing.T) {
	got := renderMC(MultiCheckbox{
		Noun: "person", NounPlural: "people", TickedWord: "selected",
		CountNoun: "entry", CountNounPlural: "entries",
		TickShownText: "Select shown", ClearText: "None",
		FilterPlaceholder: "type to filter",
	})
	for _, want := range []string{
		`data-nouns="people"`, `data-ticked="selected"`, `data-count-nouns="entries"`,
		`>Select shown</button>`, `>None</button>`, `placeholder="type to filter"`,
	} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\ngot: %s", want, got)
		}
	}
	if strings.Contains(got, "title=\"Tick every") || strings.Contains(got, "title=\"Untick every") {
		t.Errorf("custom button text should drop the default tooltips, got: %s", got)
	}
}

// Only the column steps below Columns are emitted.
func TestMultiCheckbox_ColumnSteps(t *testing.T) {
	tests := []struct {
		cols  int
		steps []string // container widths expected
		none  []string // container widths not expected
	}{
		{8, []string{"1050px", "800px", "500px"}, nil},
		{0, []string{"1050px", "800px", "500px"}, nil}, // default 8
		{5, []string{"800px", "500px"}, []string{"1050px"}},
		{3, []string{"500px"}, []string{"1050px", "800px"}},
		{2, nil, []string{"1050px", "800px", "500px"}},
	}
	for _, tt := range tests {
		css := MultiCheckbox{Columns: tt.cols}.styleRules("m")
		for _, w := range tt.steps {
			if !strings.Contains(css, "@container (max-width:"+w+")") {
				t.Errorf("Columns %d: missing step at %s\n%s", tt.cols, w, css)
			}
		}
		for _, w := range tt.none {
			if strings.Contains(css, "@container (max-width:"+w+")") {
				t.Errorf("Columns %d: unexpected step at %s\n%s", tt.cols, w, css)
			}
		}
	}

	if css := (MultiCheckbox{Columns: 3, MaxHeight: "10rem"}).styleRules("m"); !strings.Contains(css, "repeat(3,") ||
		!strings.Contains(css, "max-height:10rem") {
		t.Errorf("Columns/MaxHeight not applied:\n%s", css)
	}
}

// Extra renders inside the toolbar, after the filter and before the summary;
// help and error text wire up aria-describedby on the grid.
func TestMultiCheckbox_ExtraAndMessages(t *testing.T) {
	extra := element.CompFunc(func(b *element.Builder) (x any) {
		b.LabelClass("nested").R(
			b.Input("type", "checkbox", "name", "nested"),
			b.T(" include subfolders"),
		)
		return
	})

	got := renderMC(MultiCheckbox{ID: "f", Extra: extra, HelpText: "Pick any number."})
	filterAt := strings.Index(got, `id="f-filter"`)
	extraAt := strings.Index(got, `class="nested"`)
	summaryAt := strings.Index(got, `id="f-summary"`)
	if filterAt < 0 || extraAt < filterAt || summaryAt < extraAt {
		t.Errorf("want filter < extra < summary, got %d %d %d\n%s", filterAt, extraAt, summaryAt, got)
	}
	if !strings.Contains(got, `aria-describedby="f-help"`) || !strings.Contains(got, `id="f-help">Pick any number.`) {
		t.Errorf("help text not wired, got: %s", got)
	}

	got = renderMC(MultiCheckbox{ID: "f", Error: "Pick at least one."})
	for _, want := range []string{`class="form-field has-error multi-checkbox"`, `aria-invalid="true"`, `aria-describedby="f-error"`} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\ngot: %s", want, got)
		}
	}
}
