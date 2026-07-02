package element

import (
	"strings"
	"testing"
)

// Attributes are now stored as ordered pairs, so output is deterministic
// and matches the order attributes were passed in.
func TestAttributeOrderDeterministic(t *testing.T) {
	want := `<div id="main" class="row" data-x="1"></div>`

	for i := 0; i < 50; i++ { // map ordering bugs only show up across runs
		b := NewBuilder()
		b.Div("id", "main", "class", "row", "data-x", "1").R()
		if got := b.String(); got != want {
			t.Fatalf("run %d: got %q, want %q", i, got, want)
		}
	}
}

func TestAttributeDuplicateLastWins(t *testing.T) {
	b := NewBuilder()
	b.Div("class", "one", "id", "main", "class", "two").R()

	want := `<div class="two" id="main"></div>`
	if got := b.String(); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestAttributeOddCountDropsLast(t *testing.T) {
	b := NewBuilder()
	b.Div("class", "row", "dangling").R()

	want := `<div class="row"></div>`
	if got := b.String(); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestHasAttribute(t *testing.T) {
	b := NewBuilder()
	el := b.Div("class", "row", "id", "main")
	defer el.R()

	if !el.HasAttribute("class", "row") {
		t.Error("expected HasAttribute(class, row) to be true")
	}
	if !el.HasAttribute("id", "main") {
		t.Error("expected HasAttribute(id, main) to be true")
	}
	if el.HasAttribute("class", "nope") {
		t.Error("expected HasAttribute(class, nope) to be false")
	}
	if el.HasAttribute("missing", "row") {
		t.Error("expected HasAttribute(missing, row) to be false")
	}
}

func TestClassMethodsClassComesFirst(t *testing.T) {
	b := NewBuilder()
	b.DivClass("container", "id", "app").R()

	want := `<div class="container" id="app"></div>`
	if got := b.String(); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestUpperCaseElementNameLowered(t *testing.T) {
	b := NewBuilder()
	b.Ele("DIV", "class", "x").R()

	want := `<div class="x"></div>`
	if got := b.String(); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

// Debug mode should still append the data-ele-id attribute
func TestDebugModeAddsEleId(t *testing.T) {
	DebugSet()
	defer DebugClear()

	b := NewBuilder()
	b.Div("class", "row").R()

	got := b.String()
	if !strings.Contains(got, `class="row"`) || !strings.Contains(got, `data-ele-id="div-`) {
		t.Errorf("expected class and data-ele-id attributes, got %q", got)
	}
}
