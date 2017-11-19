package element

import "strings"

type Element struct {
	El string // just the base of the element e.g. form
	Attr map[string]string
	Single bool
}

// Create a new element
func New(el string, attrs ...string ) (Element) {
	m := map[string]string{}
	key := ""
	for i, item := range attrs {
		if i % 2 == 0 {
			key = item
		} else {
			m[key] = item
		}
	}
	return Element{El: strings.ToLower(el), Attr: m}
}

// Render with `inner` as the innerHTML
func (e Element) R(inner ...string) (str string) {
	str = "<" + e.El

	for k, v := range e.Attr {
		str += " " + k + "=" + `"` + v + `"`
	}

	if !e.Single{ // todo keep an internal map of single elements
		str += ">"
		for _, r := range inner {
			str += r
		}
		str += "</" + e.El
	}
	str += ">"
	return
}

// Iterate
func (e Element) I(arr []string, el Element) string {
	str := ""
	for _, item := range arr {
		str += el.R(item)
	}
	return e.R(str)
}
