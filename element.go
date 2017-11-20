package element

import (
	"strings"
)

type Element struct {
	El string // just the base of the element e.g. form
	Attr map[string]string
	Single bool
}

// Create a new element
func New(el string, attrs ...string ) (Element) {
	return Element{El: strings.ToLower(el), Attr: stringlistToMap(attrs...)}
}

func (e Element) IsSingleTag() bool {
	if _, ok := singleTags[e.El]; ok {
		return true
	}
	return false
}

// Add attributes after element creation
func (e Element) AddAttributes(attrs ...string) {
	m := stringlistToMap(attrs...)
	for k, v := range m {
		e.Attr[k] = v
	}
}

// Render with `inner` as the innerHTML
func (e Element) R(inner ...string) (str string) {
	str = "<" + e.El

	for k, v := range e.Attr {
		str += " " + k + "=" + `"` + v + `"`
	}

	if !e.IsSingleTag() {
		str += ">"
		for _, r := range inner {
			str += r
		}
		str += "</" + e.El
	}
	str += ">"
	return
}

// Render element only if condition true
func (e Element) RIf(condition bool, inner ...string) (str string) {
	if !condition {
		return
	}
	return e.R(inner...)
}

// Render each item in the slice wrapped in the Element el
// with everything nested within the parent element
// Attrs is a key, value list. A value may be marked as interpolatable with the iterated item with `{{$1}}`
// A value in the attrs list may be marked as interpolatable with the iterated item
func (e Element) For(arr []string, ele string, attrs ...string) string {
	// Find and save the index of the first interpolatable attr if any
	j := 0  // 0 is safe since we would never interpolate a key
	for i, a := range attrs {
		if i % 2 == 1 && a == "{{$1}}" {  // an attribute value wants to be interpolated
			j = i; break
		}
	}

	out := ""
	el := Element{ El: ele }
	for _, item := range arr {
		if j > 0 {
			attrs[j] = item  // replace
		}
		el.Attr = stringlistToMap(attrs...)
		if j == 0 {
			out += el.R(item) // render the element wrapping item
		} else {
			out += el.R()  // we already used the item in an attribute, so no wrap
		}
	}
	return e.R(out)
}


// Render as in For, but only if condition true
func (e Element) ForIf(condition bool, arr []string, el string, extra ...string) (str string) {
	if !condition {
		return
	}
	return e.For(arr, el, extra...)
}
