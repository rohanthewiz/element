package element

import (
	"fmt"
	"strings"
	"testing"
)

// var testAnimals = []string{"cat", "mouse", "dog"}

// Warning: since Go will randomize maps, multiple attributes may come out in different order
// Maybe just to use a single attribute pair or test only the length of the output if multiple attributes are required

func TestRender(t *testing.T) {
	s := &strings.Builder{}

	el := New(s, "span").R()
	if el.String() != `<span></span>` {
		t.Error("Failed to render an empty span", el.String())
	}

	// Span with inner text
	s.Reset()
	el = New(s, "span").R(
		New(s, "t", "This is some inner text").R(), // "t" is a pseudo tag for plain text
	)
	if el.String() != `<span>This is some inner text</span>` {
		t.Error("Failed to render a span with inner text", "\nGot:", el.String())
	}

	// Span with multiple inner text
	s.Reset()
	el = New(s, "span", "id", "special", "class", "normal").R(
		New(s, "t", "This is some inner text", " and some more by the way"), // we can use a list of texts
	)
	str := el.String()
	expected := `<span id="special" class="normal">This is some inner text and some more by the way</span>`
	if len(str) != len(expected) { // Go's map order is random, so have to rely on length match
		t.Error("Failed to render a span with attributes and inner text",
			"\nExpected:", expected, "\nGot:", str)
	} else {
		fmt.Println("good ->", str)
	}

	// Div with text and element children
	s.Reset()
	el = New(s, "div", "id", "container", "class", "active").R(
		New(s, "t", "This is some inner text", " and some more by the way").R(),
		New(s, "form", "method", "post").R(),
		New(s, "t", "Some ending text").R(),
	)
	str = el.String()
	expected = `<div id="container" class="active">This is some inner text and some more by the way<form method="post"></form>Some ending text</div>`
	if len(str) != len(expected) { // Go's map order is random, so have to rely on length match
		t.Error("Failed to render a span with attributes and inner text",
			"\nExpected:", expected, "\nGot:", str)
	} else {
		fmt.Println("good ->", str)
	}

	// Deep nesting
	s.Reset()
	moreText := " - more text"
	el = New(s, "div", "id", "container", "class", "active").R(
		New(s, "t", "some text", moreText).R(),
		New(s, "form", "method", "post").R(
			New(s, "input", "value", "some input").R(),
			New(s, "button").R(
				New(s, "span", "style", "background-color:wheat").R(New(s, "t", "My nice button")),
			),
		),
		New(s, "t", "Some ending text").R(),
	)
	str = el.String()
	expected = `<div id="container" class="active">some text - more text<form method="post"><input value="some input"><button><span style="background-color:wheat">My nice button</span></button></form>Some ending text</div>`
	if len(str) != len(expected) { // Go's map order is random, so have to rely on length match
		t.Error("Failed to render a span with attributes and inner text",
			"\nExpected:", expected, "\nGot:", str)
	} else {
		fmt.Println("good ->", str)
	}
}

// func TestAddAttributes(t *testing.T) {
// 	sb := &strings.Builder{}
// 	el := New(sb, "span")
// 	el.AddAttributes("id", "my-span", "class", "special-span")
// 	str := el.R()
// 	expected := `<span id="my-span" class="special-span"></span>`
// 	if str != expected {
// 		t.Error("Failed to add attributes to an element", "\nExpected:", expected, "\nGot:", str)
// 	}
// }
//
// func TestRenderIf(t *testing.T) {
// 	sb := &strings.Builder{}
// 	str := New(sb, "span").RIf(false, "This is some inner text")
// 	if str != "" {
// 		t.Error("False condition to RIf should yield empty string")
// 	}
//
// 	sb.Reset()
// 	str = New(sb, "span").RIf(true, "This is some inner text")
// 	if str != `<span>This is some inner text</span>` {
// 		t.Error("True condition to RIf should yield the rendered element")
// 	}
// }
//
// func TestRenderOddNumOfAttributes(t *testing.T) {
// 	sb := &strings.Builder{}
// 	str := New(sb, "img", "src", "http://example.com/test.png", "enable").R()
// 	expected := `<img src="http://example.com/test.png">`
// 	if str != expected {
// 		t.Error("With odd number of attributes, element should drop the last:",
// 			"\nExpected:", expected, "\nGot:", str)
// 	}
// }
//
// func TestFor(t *testing.T) {
// 	sb := &strings.Builder{}
//
// 	str := New(sb, "div").R(
// 		New(&strings.Builder{}, "ul", "class", "list").For(testAnimals, "li"), // build a list
// 	)
//
// 	expect := `<div><ul class="list"><li>cat</li><li>mouse</li><li>dog</li></ul></div>`
// 	if str != expect {
// 		t.Error("Iteration failed. Expected:", expect, "\nGot:", str)
// 	}
// }

// Todo - Single tag elements
// func TestForWithInterpolation(t *testing.T) {
// 	sb := &strings.Builder{}
// 	str := New(sb, "div", "class", "form-group").For(testAnimals, "input", "class", "text-input", "value", "{{$1}}")
// 	expect := `<div class="form-group"><input class="text-input" value="cat"><input class="text-input" value="mouse"><input class="text-input" value="dog"></div>`
// 	if len(str) != len(expect) {  // map ordering is random so test lengths as a compromise
// 		t.Error("Iteration failed. Expected similar to:", expect, "\nGot:", str)
// 	}
// }
//
// func TestForIf(t *testing.T) {
// 	animals := []string{"cat", "mouse", "dog"}
// 	sb := &strings.Builder{}
//
// 	str := New(sb, "div").R(
// 		New(sb, "ul", "class", "list").ForIf(false, animals, "li"),
// 	)
// 	expect := `<div></div>`
// 	if str != expect {
// 		t.Error("Iteration failed. Expected:", expect, "\nGot:", str)
// 	}
//
// 	sb.Reset()
// 	str = New(sb, "div").R(
// 		New(sb, "ul", "class", "list").ForIf(true, animals, "li"),
// 	)
// 	expect = `<div><ul class="list"><li>cat</li><li>mouse</li><li>dog</li></ul></div>`
// 	if str != expect {
// 		t.Error("Iteration failed. Expected:", expect, "\nGot:", str)
// 	}
// }
