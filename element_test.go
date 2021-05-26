package element

import (
	"strings"
	"testing"
)

var testAnimals = []string{"cat", "mouse", "dog"}

func TestRender(t *testing.T) {
	sb := &strings.Builder{}
	str := New(sb, "span").R()
	if str != `<span></span>` {
		t.Error("Failed to render an empty span")
	}

	sb.Reset()
	str = New(sb, "span").R("This is some inner text")
	if str != `<span>This is some inner text</span>` {
		t.Error("Failed to render a span with inner text")
	}

	sb.Reset()
	str = New(sb, "span", "id", "special", "class", "normal").R("This is some inner text")
	expected := `<span id="special" class="normal">This is some inner text</span>`
	if len(str) != len(expected) { // Go's map order is random, so have to rely on length match
		t.Error("Failed to render a span with attributes and inner text",
			"\nExpected:", expected, "\nGot:", str)
	}
}

func TestAddAttributes(t *testing.T) {
	sb := &strings.Builder{}
	el := New(sb, "span")
	el.AddAttributes("id", "my-span", "class", "special-span")
	str := el.R()
	expected := `<span id="my-span" class="special-span"></span>`
	if str != expected {
		t.Error("Failed to add attributes to an element", "\nExpected:", expected, "\nGot:", str)
	}
}

func TestRenderIf(t *testing.T) {
	sb := &strings.Builder{}
	str := New(sb, "span").RIf(false, "This is some inner text")
	if str != "" {
		t.Error("False condition to RIf should yield empty string")
	}

	sb.Reset()
	str = New(sb, "span").RIf(true, "This is some inner text")
	if str != `<span>This is some inner text</span>` {
		t.Error("True condition to RIf should yield the rendered element")
	}
}

func TestRenderOddNumOfAttributes(t *testing.T) {
	sb := &strings.Builder{}
	str := New(sb, "img", "src", "http://example.com/test.png", "enable").R()
	expected := `<img src="http://example.com/test.png">`
	if str != expected {
		t.Error("With odd number of attributes, element should drop the last:",
			"\nExpected:", expected, "\nGot:", str)
	}
}

func TestFor(t *testing.T) {
	sb := &strings.Builder{}

	str := New(sb, "div").R(
		New(&strings.Builder{}, "ul", "class", "list").For(testAnimals, "li"), // build a list
	)

	expect := `<div><ul class="list"><li>cat</li><li>mouse</li><li>dog</li></ul></div>`
	if str != expect {
		t.Error("Iteration failed. Expected:", expect, "\nGot:", str)
	}
}

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
func TestForIf(t *testing.T) {
	animals := []string{"cat", "mouse", "dog"}
	sb := &strings.Builder{}

	str := New(sb, "div").R(
		New(sb, "ul", "class", "list").ForIf(false, animals, "li"),
	)
	expect := `<div></div>`
	if str != expect {
		t.Error("Iteration failed. Expected:", expect, "\nGot:", str)
	}

	sb.Reset()
	str = New(sb, "div").R(
		New(sb, "ul", "class", "list").ForIf(true, animals, "li"),
	)
	expect = `<div><ul class="list"><li>cat</li><li>mouse</li><li>dog</li></ul></div>`
	if str != expect {
		t.Error("Iteration failed. Expected:", expect, "\nGot:", str)
	}
}
