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

// Create a new element
func New(t *strings.Builder, el string, attrs ...string) (ele Element) {
	if t == nil {
		log.Println("Please supply a pointer to a strings builder to element.New():", el)
	}

	ele = Element{sb: t, El: strings.ToLower(el)}

	if ele.IsText() {
		ele.arrayAttrs = attrs // plain text will use the original list
	} else {
		ele.attrs = stringlistToMap(attrs...)
	}

	ele.writeOpenTag() // write opening tag right away

	return ele
}

// Render Elements - well kind of, as the language will run inner functions first
// 	we don't have to do anything for children
// This element Ancestors will be already in the tree (string builder) as New() is called before R (Render)
// So, essentially this is just to let us know to add our ending tag if applicable
func (e Element) R(children ...Element) Element {
	if !e.IsSingleTag() {
		e.sb.WriteString("</" + e.El + ">")
	}
	return e
}

func (e Element) writeOpenTag() {
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
