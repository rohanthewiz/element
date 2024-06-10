package element

import (
	"log"
	"strings"
)

type Element struct {
	El                string // just the base of the element e.g. td, h1
	arrayAttrs        []string
	attrsMap          map[string]string
	sb                *strings.Builder
	openingTagWritten bool // attributes can be manipulated until the opening tag is written
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
		e.attrsMap = stringListToMap(attrs...)
	}

	// Render opening tag bc most times that's  what we want next
	e.RenderOpeningTag()
	return e
}

// AA is an alias method for AddAttributes
func (e *Element) AA(attrs ...string) {
	e.AddAttributes(attrs...)
}

// AddAttributes adds attributes (key val pairs) to the Element's attributes map.
func (e *Element) AddAttributes(attrs ...string) {
	if e.sb == nil {
		log.Println("Please create a new element with NewNoRender() before adding attributes")
		return
	}

	if e.IsText() {
		log.Println("Cannot add additional attributes to a simple text element")
		return
	}

	addAttributePairsToMap(e.attrsMap, attrs...)
}

// NewNoRender is for creating a new element without immediately rendering the opening tag
// We would do this if we want to further manipulate attributes,  as in the case of conditional attributes.
func NewNoRender(s *strings.Builder, el string, attrs ...string) (e Element) {
	if s == nil {
		log.Println("Please supply a pointer to a string builder to element.NewNoRender():", el)
	}
	e = Element{sb: s, El: strings.ToLower(el)}

	if e.IsText() {
		e.arrayAttrs = attrs // simple text will use the original list
	} else {
		e.attrsMap = stringListToMap(attrs...)
	}
	return e
}

// Text creates a new text element
func Text(s *strings.Builder, texts ...string) (a struct{}) {
	if s == nil {
		log.Println("Please supply a pointer to a string builder to element.Text()")
	}
	e := Element{sb: s, El: "t"}
	e.arrayAttrs = texts
	e.RenderOpeningTag() // write opening tag right away
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
// The return is just to satisfy the any interface{} input param of the parent .R()
func (e Element) For(items []string, ele string, attrs ...string) (a struct{}) {
	for _, item := range items {
		New(e.sb, ele, attrs...).R(
			New(e.sb, "t", item),
		)
	}
	e.close()
	return
}

// OT is an alias method for RenderOpeningTag
func (e Element) OT() Element {
	return e.RenderOpeningTag()
}

func (e Element) RenderOpeningTag() (self Element) {
	if e.sb != nil {
		if e.IsText() {
			for _, a := range e.arrayAttrs {
				e.sb.WriteString(a)
			}
		} else {
			e.sb.WriteString("<" + e.El)
			for k, v := range e.attrsMap {
				e.sb.WriteString(" " + k + "=" + `"` + v + `"`)
			}
			e.sb.WriteString(">")
		}

		e.openingTagWritten = true
	}
	return e
}

func (e Element) close() {
	if !e.IsSingleTag() {
		e.sb.WriteString("</" + e.El + ">")
	}
}
