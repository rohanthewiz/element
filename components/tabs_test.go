package components

import (
	"strings"
	"testing"

	"github.com/rohanthewiz/element"
)

func TestTabs_Render(t *testing.T) {
	b := element.NewBuilder()
	Tabs{
		ID:     "demo",
		Active: 1,
		Items: []Tab{
			{Label: "Overview", Content: "Overview content"},
			{Label: "Details", Body: Badge{Text: "New"}},
		},
	}.Render(b)
	got := b.String()

	for _, want := range []string{
		`id="demo"`,
		`type="radio"`, `name="demo"`, `id="demo-t0"`, `id="demo-t1"`,
		`<label class="tab-label" for="demo-t0">Overview</label>`,
		`<div class="tab-panels">`,
		"Overview content",
		`<span class="badge badge-neutral">New</span>`,
		// Generated style rules
		"#demo > .tab-radio{display:none}",
		"#demo-t0:checked ~ .tab-panels > .tab-panel:nth-child(1){display:block}",
		"#demo-t1:checked ~ .tab-panels > .tab-panel:nth-child(2){display:block}",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("Tabs.Render() missing %q\ngot: %s", want, got)
		}
	}

	// Active tab (index 1) gets checked, and only that one
	if strings.Count(got, `checked="checked"`) != 1 {
		t.Errorf("expected exactly one checked radio, got: %s", got)
	}
	checkedPos := strings.Index(got, `checked="checked"`)
	t1Pos := strings.Index(got, `id="demo-t1"`)
	if checkedPos < t1Pos {
		t.Errorf("expected the second tab to be checked, got: %s", got)
	}
}

func TestTabs_EmptyRendersNothing(t *testing.T) {
	b := element.NewBuilder()
	Tabs{ID: "empty"}.Render(b)
	if got := b.String(); got != "" {
		t.Errorf("empty Tabs should render nothing, got: %s", got)
	}
}

func TestTabs_ActiveOutOfRangeDefaultsToFirst(t *testing.T) {
	b := element.NewBuilder()
	Tabs{Items: []Tab{{Label: "A"}, {Label: "B"}}, Active: 5}.Render(b)
	got := b.String()

	checkedPos := strings.Index(got, `checked="checked"`)
	t1Pos := strings.Index(got, `id="tabs-t1"`)
	if checkedPos == -1 || checkedPos > t1Pos {
		t.Errorf("out-of-range Active should check the first tab, got: %s", got)
	}
}
