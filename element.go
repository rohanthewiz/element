package element

import (
	"fmt"
	"strings"

	"github.com/rohanthewiz/serr"
)

type Element struct {
	name string // just the base of the element e.g. td, h1
	id   string // id is the unique element id
	// seq        int    // seq holds the order of the element - order is not guaranteed, but it is useful for debugging
	arrayAttrs []string
	attrs      map[string]string
	function   string // function in which the element is created
	location   string // file:line_nbr of the element creation
	sb         *strings.Builder
	issues     []string // issues can hold any issues with the element
}

func (el Element) Name() string {
	return el.name
}

// New creates a new element
func New(s *strings.Builder, el string, attrs ...string) (e Element) {
	if s == nil {
		fmt.Println("Please supply a pointer to a string builder to element.New():", el)
	}

	e = Element{sb: s, name: strings.ToLower(el)}
	e.id = e.name
	if IsDebugMode() {
		e.id += "-" + genRandomId(6) // generate a random id for the element
	}

	e.function = serr.FunctionName(serr.FrameLevels.FrameLevel3)
	e.location = serr.FunctionLoc(serr.FrameLevels.FrameLevel3)

	if e.IsText() {
		e.arrayAttrs = attrs // plain text will use the original list
	} else {
		e.attrs = stringlistToMap(e, attrs...)
	}

	e.writeOpeningTag() // write opening tag right away

	if IsDebugMode() && !e.IsSingleTag() {
		// Temporarily element id into the issues map as an open tag
		concerns.UpsertConcern(concernOpenTag, e)
	}

	return e
}

func (el Element) HasAttribute(key, value string) bool {
	if el.attrs == nil {
		return false
	}
	if val, ok := el.attrs[key]; ok && val == value {
		return true
	}
	return false
}

// Text is an element core function which creates a new text element in the string builder
func Text(s *strings.Builder, texts ...string) (x any) {
	if s == nil {
		fmt.Println("Please supply a pointer to a string builder to element.Text()")
	}
	e := Element{sb: s, name: "t"}
	e.arrayAttrs = texts
	e.writeOpeningTag() // write opening tag right away
	return
}

// R renders children of the element and any closing tag if applicable
// The element's opening tag will be already in the render tree (string builder)
// because New() [element] is called before R (Render)
// So, essentially this is just to allow any children to render and let us add our ending tag if applicable
// The return is just to have some value to pass back as an argument of the parent R()
func (el Element) R(args ...any) (x any) {
	if el.IsSingleTag() {
		if IsDebugMode() && len(args) > 0 { // We tried to render children on a single tag -- we shouldn't do that.
			issue := fmt.Sprintf(`The element is a single tag, but yet has %d child(ren). It should have none.`, len(args))
			fmt.Printf("![%s] %s\n", el.id, issue)
			el.issues = append(el.issues, issue)
			concerns.UpsertConcern(concernOther, el)
		}
		return // Single tags should not have children, so just return
	}

	if IsDebugMode() {
		for i, arg := range args {
			// Not sure, but we may want to deprecate this
			_, isStruct := arg.(struct{})

			// Check if the argument is a single element tag as we will allow single tags to not be rendered
			argIsSingleElement := false
			argEle, isElement := arg.(Element)
			if isElement && argEle.IsSingleTag() {
				argIsSingleElement = true
			}

			// Our standard now is to return an empty value of type any (i.e. nil) to the parent R().
			// Anything other than nil from the children elements (or struct{}), should be considered an issue
			// Most likely some literal text was not wrapped in t().
			if arg != nil && !isStruct && !argIsSingleElement {
				strArg := ""
				if isElement {
					strArg = arg.(Element).detailsHtml() // Get details of the element
				} else {
					strArg = fmt.Sprintf("%v", arg) // Convert to string representation
				}

				issue := fmt.Sprintf(`The %s child is not properly rendered.
		Did you forget to wrap with builder.Text()? Child: %q`,
					ToOrdinal(i+1), strArg)

				fmt.Printf("![%s] %s\n", el.id, issue)
				el.issues = append(el.issues, issue)
			}
		}

		if len(el.issues) > 0 {
			// Add / Replace element in the concerns map
			concerns.UpsertConcern(concernOther, el)
		}
	}

	el.close()

	if IsDebugMode() {
		// Remove the open tag from concerns
		concerns.UpsertConcern(concernClosedTag, el)
	}
	return
}

// T renders a list of text-only children on an Element
// Use this when an element has only text children
func (el Element) T(texts ...string) (x any) {
	el.R(Text(el.sb, texts...))
	return
}

// F renders a formatted text-only child on an Element
// This eliminates the need to use fmt.Sprintf()
// Use this when an element has only a single text child
func (el Element) F(format string, args ...any) (x any) {
	el.R(Text(el.sb, fmt.Sprintf(format, args...)))
	return
}

func (el Element) writeOpeningTag() {
	if el.sb != nil {
		if el.name == "t" { // "t" is a pseudo element representing a list of strings
			for _, a := range el.arrayAttrs {
				el.sb.WriteString(a)
			}
		} else {
			el.sb.WriteString("<" + el.name)
			for k, v := range el.attrs {
				el.sb.WriteString(fmt.Sprintf(` %s="%s"`, k, v))
			}
			el.sb.WriteString(">")
		}
	}
}

func (el Element) close() {
	if !el.IsSingleTag() {
		el.sb.WriteString("</" + el.name + ">")
	}
}

func (el Element) details() string {
	return fmt.Sprintf("Element %s in %s (%s)", el.id, el.function, el.location)
}

func (el Element) detailsHtml() string {
	copyIcon := `<svg class="copy-icon" onclick="copyToClipboard('` + el.location + `')" title="Copy location" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="9" y="9" width="13" height="13" rx="2" ry="2"></rect><path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"></path></svg>`
	return fmt.Sprintf("<strong>%s</strong> tag %s<br>%s (<strong>%s</strong> %s)", el.name, el.id, el.function, el.location, copyIcon)
}

/* Deprecating For - use b.Wrap and ForEach instead // For renders a slice of items wrapped in the Element el
// with everything nested within the parent element e
// Attrs is a key, value list.
// Note that the builder convenience functions gives more flexibility
// Consider this method deprecated as we can do so much more with builder.Wrap and components.
func (el Element) For(items []string, ele string, attrs ...string) (x any) {
	for _, item := range items {
		New(el.sb, ele, attrs...).R(
			New(el.sb, "t", item),
		)
	}
	el.close()
	return
}
*/
