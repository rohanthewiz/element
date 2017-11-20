package element

import (
	"testing"
)
var testAnimals = []string{"cat", "mouse", "dog"}

func TestRender(t *testing.T) {
	str := New("span").R()
	if str != `<span></span>` {
		t.Error("Failed to render an empty span")
	}
	str = New("span").R("This is some inner text")
	if str != `<span>This is some inner text</span>` {
		t.Error("Failed to render a span with inner text")
	}
	str = New("span", "id", "special", "class", "normal").R("This is some inner text")
	expected := `<span id="special" class="normal">This is some inner text</span>`
	if len(str) != len(expected) {  // Go's map order is random, so have to rely on length match
		t.Error("Failed to render a span with attributes and inner text",
			"\nExpected:", expected, "\nGot:", str)
	}
}

func TestAddAttributes(t *testing.T) {
	el := New("span")
	el.AddAttributes("id", "my-span", "class", "special-span")
	str := el.R()
	expected := `<span id="my-span" class="special-span"></span>`
	if str != expected {
		t.Error("Failed to add attributes to an element", "\nExpected:", expected, "\nGot:", str)
	}
}

func TestRenderIf(t *testing.T) {
	str := New("span").RIf(false, "This is some inner text")
	if str != "" {
		t.Error("False condition to RIf should yield empty string")
	}
	str = New("span").RIf(true, "This is some inner text")
	if str != `<span>This is some inner text</span>` {
		t.Error("True condition to RIf should yield the rendered element")
	}
}

func TestRenderOddNumOfAttributes(t *testing.T) {
	str := New("img", "src", "http://example.com/test.png", "enable").R()
	expected := `<img src="http://example.com/test.png">`
	if str != expected {
		t.Error("With odd number of attributes, element should drop the last:",
		"\nExpected:", expected, "\nGot:", str)
	}
}

func TestFor(t *testing.T) {
	str := New("div").R(
		New("ul", "class", "list").For(testAnimals, "li"), // build a list
	)
	expect := `<div><ul class="list"><li>cat</li><li>mouse</li><li>dog</li></ul></div>`
	if str != expect {
		t.Error("Iteration failed. Expected:", expect, "\nGot:", str)
	}
}

// Todo - Single tag elements
func TestForWithInterpolation(t *testing.T) {
	str := New("div", "class", "form-group").For(testAnimals, "input", "class", "text-input", "value", "{{$1}}")
	expect := `<div class="form-group"><input class="text-input" value="cat"><input class="text-input" value="mouse"><input class="text-input" value="dog"></div>`
	if len(str) != len(expect) {  // map ordering is random so test lengths as a compromise
		t.Error("Iteration failed. Expected similar to:", expect, "\nGot:", str)
	}
}

func TestForIf(t *testing.T) {
	animals := []string{"cat", "mouse", "dog"}

	str := New("div").R(
		New("ul", "class", "list").ForIf(false, animals, "li"),
	)
	expect := `<div></div>`
	if str != expect {
		t.Error("Iteration failed. Expected:", expect, "\nGot:", str)
	}

	str = New("div").R(
		New("ul", "class", "list").ForIf(true, animals, "li"),
	)
	expect = `<div><ul class="list"><li>cat</li><li>mouse</li><li>dog</li></ul></div>`
	if str != expect {
		t.Error("Iteration failed. Expected:", expect, "\nGot:", str)
	}
}
