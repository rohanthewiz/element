package element

import (
	"testing"
)

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
	if str != `<span id="special" class="normal">This is some inner text</span>` {
		t.Error("Failed to render a span with attributes and inner text")
	}
}

func TestIterate(t *testing.T) {
	animals := []string{"cat", "mouse", "dog"}

	str := New("div").R(
		New("ul", "class", "list").I(animals, New("li")), // build a list
	)
	expect := `<div><ul class="list"><li>cat</li><li>mouse</li><li>dog</li></ul></div>`
	if str != expect {
		t.Error("Iteration failed. Expected:", expect, "\nGot:", str)
	}
}
