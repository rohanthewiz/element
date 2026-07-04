package components

import (
	"strings"
	"testing"

	"github.com/rohanthewiz/element"
)

func TestSelectField_Render(t *testing.T) {
	b := element.NewBuilder()
	SelectField{
		Label:    "Country",
		Name:     "country",
		Prompt:   "Choose one…",
		Selected: "ca",
		Required: true,
		HelpText: "Where you live",
		Options: []SelectOption{
			{Value: "us", Label: "United States"},
			{Value: "ca", Label: "Canada"},
			{Value: "xx", Label: "Atlantis", Disabled: true},
		},
	}.Render(b)
	got := b.String()

	for _, want := range []string{
		`<label class="form-label" for="country">`,
		`name="country"`, `required="required"`,
		`aria-describedby="country-help"`,
		">Choose one…</option>",
		`<option value="us">United States</option>`,
		`<option value="ca" selected="selected">Canada</option>`,
		`<option value="xx" disabled="disabled">Atlantis</option>`,
		`<small class="form-help" id="country-help">Where you live</small>`,
	} {
		if !strings.Contains(got, want) {
			t.Errorf("SelectField.Render() missing %q\ngot: %s", want, got)
		}
	}

	// Prompt should not be selected since a value is selected
	if strings.Count(got, `selected="selected"`) != 1 {
		t.Errorf("expected exactly one selected option, got: %s", got)
	}
}

func TestSelectField_PromptSelectedWhenNoValue(t *testing.T) {
	b := element.NewBuilder()
	SelectField{
		Label: "Size", Name: "size", Prompt: "Pick a size",
		Options: []SelectOption{{Value: "s"}, {Value: "m"}},
	}.Render(b)
	got := b.String()

	if !strings.Contains(got, `value="" disabled="disabled" selected="selected"`) {
		t.Errorf("prompt should be selected when nothing else is, got: %s", got)
	}
	// Options with no Label fall back to Value
	if !strings.Contains(got, `<option value="s">s</option>`) {
		t.Errorf("option label should default to value, got: %s", got)
	}
}

func TestTextAreaField_Render(t *testing.T) {
	b := element.NewBuilder()
	TextAreaField{
		Label: "Bio", Name: "bio", Rows: 6,
		Value: "Hello there", Error: "Too short",
	}.Render(b)
	got := b.String()

	for _, want := range []string{
		`class="form-field has-error"`,
		`rows="6"`,
		">Hello there</textarea>",
		`aria-invalid="true"`, `aria-describedby="bio-error"`,
		`<span class="form-error" id="bio-error">Too short</span>`,
	} {
		if !strings.Contains(got, want) {
			t.Errorf("TextAreaField.Render() missing %q\ngot: %s", want, got)
		}
	}
}

func TestCheckboxField_Render(t *testing.T) {
	b := element.NewBuilder()
	CheckboxField{
		Label: "Subscribe to newsletter", Name: "subscribe",
		Checked: true, Value: "yes",
	}.Render(b)
	got := b.String()

	for _, want := range []string{
		`type="checkbox"`, `name="subscribe"`,
		`checked="checked"`, `value="yes"`,
		`<label class="form-check-label" for="subscribe">`,
		"Subscribe to newsletter",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("CheckboxField.Render() missing %q\ngot: %s", want, got)
		}
	}

	// Input comes before the label (checkbox layout)
	if strings.Index(got, "<input") > strings.Index(got, "form-check-label") {
		t.Errorf("checkbox input should precede its label, got: %s", got)
	}
}

func TestRadioGroup_Render(t *testing.T) {
	b := element.NewBuilder()
	RadioGroup{
		Legend: "Plan", Name: "plan", Selected: "pro", Inline: true,
		Options: []SelectOption{
			{Value: "free", Label: "Free"},
			{Value: "pro", Label: "Pro"},
		},
	}.Render(b)
	got := b.String()

	for _, want := range []string{
		"<fieldset", "radio-group-inline",
		"<legend", "Plan",
		`type="radio"`, `id="plan-0"`, `id="plan-1"`,
		`value="pro" class="form-check-input" checked="checked"`,
		`<label class="form-check-label" for="plan-1">Pro</label>`,
		"</fieldset>",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("RadioGroup.Render() missing %q\ngot: %s", want, got)
		}
	}

	if strings.Count(got, `checked="checked"`) != 1 {
		t.Errorf("expected exactly one checked radio, got: %s", got)
	}
}
