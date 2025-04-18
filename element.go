package element

import (
	"log"
	"strings"
)

type Element struct {
	El         string // just the base of the element e.g. td, h1
	arrayAttrs []string
	attrs      map[string]string
	sb         *strings.Builder
}

// New creates a new element
func New(s *strings.Builder, el string, attrs ...string) (e Element) {
	if s == nil {
		log.Println("Please supply a pointer to a string builder to element.New():", el)
	}
	e = Element{sb: s, El: strings.ToLower(el)}

	if e.IsText() {
		e.arrayAttrs = attrs // plain text will use the original list
	} else {
		e.attrs = stringlistToMap(attrs...)
	}

	e.writeOpeningTag() // write opening tag right away
	return e
}

// Text is an element core function which creates a new text element in the string builder
func Text(s *strings.Builder, texts ...string) (a struct{}) {
	if s == nil {
		log.Println("Please supply a pointer to a string builder to element.Text()")
	}
	e := Element{sb: s, El: "t"}
	e.arrayAttrs = texts
	e.writeOpeningTag() // write opening tag right away
	return
}

// R renders Elements - well kind of, as the language will run inner functions first
//
//	we don't have to do anything for children
//
// This element's Ancestors will be already in the tree (string builder) bc New() is called before R (Render)
// So, essentially this is just to let us know to add our ending tag if applicable
// The return is bogus - it's just to satisfy the any input of the parent .R()
func (e Element) R(_ ...any) (a struct{}) {
	e.close()
	return
}

// For renders a slice of items wrapped in the Element el
// with everything nested within the parent element e
// Attrs is a key, value list.
// Note that the use of an inline anonymous function gives more flexibility
// This function is just for convenience
// Consider this method deprecated as we can do so much more with anonymous functions, and even better, components.
func (e Element) For(items []string, ele string, attrs ...string) (a struct{}) {
	for _, item := range items {
		New(e.sb, ele, attrs...).R(
			New(e.sb, "t", item),
		)
	}
	e.close()
	return
}

func (e Element) writeOpeningTag() {
	if e.sb != nil {
		if e.El == "t" { // "t" is a pseudo element representing a list of strings
			for _, a := range e.arrayAttrs {
				e.sb.WriteString(a)
			}
		} else {
			e.sb.WriteString("<" + e.El)
			for k, v := range e.attrs {
				e.sb.WriteString(" " + k + "=" + `"` + v + `"`)
			}
			e.sb.WriteString(">")
		}
	}
}

func (e Element) close() {
	if !e.IsSingleTag() {
		e.sb.WriteString("</" + e.El + ">")
	}
}
