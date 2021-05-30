package element

import (
	"fmt"
	"strings"
	"testing"
)

// Warning: since Go will randomize maps, multiple attributes may come out in different order
// Maybe just to use a single attribute pair or test only the length of the output if multiple attributes are required

func TestRender(t *testing.T) {
	s := &strings.Builder{}

	New(s, "span").R()
	if s.String() != `<span></span>` {
		t.Error("Failed to render an empty span", s.String())
	}

	// Span with inner text
	s.Reset()
	New(s, "span").R(
		Text(s, "This is some inner text"),
	)
	if s.String() != `<span>This is some inner text</span>` {
		t.Error("Failed to render a span with inner text", "\nGot:", s.String())
	}

	// Span with multiple inner text
	s.Reset()
	New(s, "span", "id", "special", "class", "normal").R(
		Text(s, "This is some inner text", " and some more by the way"), // we can use a list of texts
	)
	str := s.String()
	expected := `<span id="special" class="normal">This is some inner text and some more by the way</span>`
	if len(str) != len(expected) { // Go's map order is random, so have to rely on length match
		t.Error("Failed to render a span with multiple inner text",
			"\nExpected:", expected, "\nGot:", str)
	} else {
		fmt.Println("good ->", str)
	}

	// Div with text and element children
	s.Reset()
	New(s, "div", "id", "container", "class", "active").R(
		Text(s, "This is some inner text", " and some more by the way"),
		New(s, "form", "method", "post").R(),
		Text(s, "Some ending text"),
	)
	str = s.String()
	expected = `<div id="container" class="active">This is some inner text and some more by the way<form method="post"></form>Some ending text</div>`
	if len(str) != len(expected) { // Go's map order is random, so have to rely on length match
		t.Error("Failed to render div with text and element children",
			"\nExpected:", expected, "\nGot:", str)
	} else {
		fmt.Println("good ->", str)
	}

	// Deep nesting
	s.Reset()
	moreText := " - more text"
	New(s, "div", "id", "container", "class", "active").R(
		Text(s, "some text", moreText),
		New(s, "form", "method", "post").R(
			New(s, "input", "value", "some input").R(),
			New(s, "button").R(
				New(s, "span", "style", "background-color:wheat").R(Text(s, "My nice button")),
			),
		),
		Text(s, "Some ending text"),
	)
	str = s.String()
	expected = `<div id="container" class="active">some text - more text<form method="post"><input value="some input"><button><span style="background-color:wheat">My nice button</span></button></form>Some ending text</div>`
	if len(str) != len(expected) { // Go's map order is random, so have to rely on length match
		t.Error("Failed to render div with deep nesting",
			"\nExpected:", expected, "\nGot:", str)
	} else {
		fmt.Println("good ->", str)
	}
}

func TestFor(t *testing.T) {
	var testAnimals = []string{"cat", "mouse", "dog"}
	s := &strings.Builder{}

	// Div with text and element children
	New(s, "div", "id", "container", "class", "active").R(
		New(s, "ul", "class", "list").For(testAnimals, "li", "class", "animal"), // build a list
	)
	str := s.String()
	expected := `<div id="container" class="active"><ul class="list"><li class="animal">cat</li><li class="animal">mouse</li><li class="animal">dog</li></ul></div>`
	if len(str) != len(expected) { // Go's map order is random, so have to rely on length match
		t.Error("Failed to render a html list with For",
			"\nExpected:", expected, "\nGot:", str)
	} else {
		fmt.Println("good ->", str)
	}
}
