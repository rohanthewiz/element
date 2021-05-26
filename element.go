package element

import (
	"fmt"
	"log"
	"strings"
)

type Element struct {
	El     string // just the base of the element e.g. form
	attr   map[string]string
	Single bool
	sb     *strings.Builder
}

// Create a new element
func New(sb *strings.Builder, el string, attrs ...string) (ele Element) {
	if sb == nil {
		log.Println("Please supply a pointer to a string builder to element.New():", el)
	}
	return Element{sb: sb, El: strings.ToLower(el), attr: stringlistToMap(attrs...)}
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
		e.attr[k] = v
	}
}

// Render the element with `inner` as the innerHTML
func (e Element) R(inner ...string) (out string) {
	if e.sb != nil {
		e.W(inner...)
		out = e.sb.String()
	}
	return
}

// Write the element to the internal string builder but don't render a string
func (e Element) W(inner ...string) {
	if e.sb != nil {
		e.sb.WriteString("<" + e.El)

		for k, v := range e.attr {
			e.sb.WriteString(" " + k + "=" + `"` + v + `"`)
		}

		if !e.IsSingleTag() {
			e.sb.WriteString(">")
			for _, r := range inner {
				e.sb.WriteString(r)
			}
			e.sb.WriteString("</" + e.El)
		}
		e.sb.WriteString(">")
	}
	return
}

// Render element only if condition true
// Note that it is better to do conditionals in an anonymous function
// therefore we are deprecating this function. This function is deprecated.
func (e Element) RIf(condition bool, inner ...string) (str string) {
	if !condition {
		return
	}
	return e.R(inner...)
}

// Normally render each item in the slice wrapped in the Element el
// with everything nested within the parent element e
// Attrs is a key, value list.
// We will do interpolation in a separate function.
// Note that it is better to do conditionals in an anonymous function
func (e Element) For(items []string, ele string, attrs ...string) string {
	sb := &strings.Builder{}
	for _, item := range items {
		el := Element{sb: sb, El: ele}
		el.attr = stringlistToMap(attrs...)
		el.W(item)
	}

	fmt.Println("**->", sb.String())
	return e.R(sb.String())
}

// TODO: Make another version of this function that only does interpolation
// // Normally render each item in the slice wrapped in the Element el
// // with everything nested within the parent element e
// // Attrs is a key, value list.
// // We will do interpolation in a separate function.
// // strike thru: A value may be marked as interpolatable
// // strike thru: with the iterated item with this expression `{{$x}}`
// // strike thru: -- in this case it makes sense to have only one item in `items`
// // Note that it is better to do conditionals in an anonymous function
// func (e Element) For(items []string, ele string, attrs ...string) string {
// 	// Find and save the index of the first interpolatable attr if any
// 	j := 0  // 0 is safe since we would never interpolate a key
// 	for i, a := range attrs {
// 		if i % 2 == 1 && a == "{{$x}}" {  // an attribute value wants to be interpolated
// 			j = i; break
// 		}
// 	}
//
// 	sb := &strings.Builder{}
// 	for _, item := range items {
// 		el := Element{sb: sb, El: ele}
// 		if j > 0 {
// 			attrs[j] = item  // replace - the replaceable attribute with
// 			break
// 		}
// 		el.attr = stringlistToMap(attrs...)
//
// 		if j == 0 {
// 			el.W(item)
// 		} else {
// 			el.W()  // we already used the item in an attribute, so no wrap
// 		}
// 	}
//
// 	// fmt.Println("**->", sb.String())
// 	return e.R(sb.String())
// }

// Render as in For, but only if condition true
// Note that it is better to do conditionals in an anonymous function
// therefore we are deprecating this function. This function is deprecated.
func (e Element) ForIf(condition bool, arr []string, el string, extra ...string) (str string) {
	if !condition {
		return
	}
	return e.For(arr, el, extra...)
}
