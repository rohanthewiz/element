package element

import (
	"bytes"
	"fmt"
	"html"
	"strings"

	"github.com/rohanthewiz/serr"
)

var singleTags = map[string]bool{
	"img": true, "br": true, "hr": true, "meta": true,
	"input": true, "link": true, "t": true,
	"area": true, "base": true, "col": true, "embed": true,
	"keygen": true, "param": true, "source": true, "track": true, "wbr": true,
}

type Element struct {
	name string // just the base of the element e.g. td, h1
	id   string // id is the unique element id
	// seq        int    // seq holds the order of the element - order is not guaranteed, but it is useful for debugging
	arrayAttrs []string
	attrPairs  []string // key/value attribute pairs in insertion order (deterministic output)
	function   string   // function in which the element is created
	location   string   // file:line_nbr of the element creation
	sb         *bytes.Buffer
	issues     []string // issues can hold any issues with the element
}

func (el Element) Name() string {
	return el.name
}

// New creates a new element
func New(s *bytes.Buffer, el string, attrs ...string) (e Element) {
	if s == nil {
		fmt.Println("Please supply a pointer to a string builder to element.New():", el)
	}

	e = Element{sb: s, name: lowerName(el)}
	e.id = e.name
	if IsDebugMode() {
		e.id += "-" + genRandomId(6) // generate a random id for the element

		// Stack walks are expensive, so only capture caller info when debugging
		e.function = serr.FunctionName(serr.FrameLevels.FrameLevel3)
		e.location = serr.FunctionLoc(serr.FrameLevels.FrameLevel3)
	}

	if e.IsText() {
		e.arrayAttrs = attrs // plain text will use the original list
	} else {
		e.attrPairs = normalizeAttrPairs(e, attrs)
	}

	e.writeOpeningTag() // write opening tag right away

	if IsDebugMode() && !e.IsSingleTag() {
		// Temporarily element id into the issues map as an open tag
		concerns.UpsertConcern(concernOpenTag, e)
	}

	return e
}

func (el Element) HasAttribute(key, value string) bool {
	for i := 0; i+1 < len(el.attrPairs); i += 2 {
		if el.attrPairs[i] == key && el.attrPairs[i+1] == value {
			return true
		}
	}
	return false
}

// writeAttrValue writes an attribute value into the double quotes the renderer
// has already opened, escaping the one character that cannot appear there: a
// double quote, which would close the attribute early and let everything after
// it be read as further attributes on the element.
//
// This is the only escaping element performs, and it is deliberate on both
// counts — that it happens at all, and that it stops here.
//
// Why it happens: element is otherwise a raw writer, and that is the right
// default for element *content*, where callers legitimately emit markup —
// b.Style().T(css) and b.Script().T(js) both depend on it. Attribute values have
// no equivalent case. A raw double quote inside one is not a stylistic choice,
// it is broken output, every time. Escaping it therefore cannot change what any
// correct program renders, only what an incorrect one does.
//
// Why only the double quote:
//
//   - The angle brackets do not terminate a quoted attribute value. The
//     tokenizer is looking for the closing quote and nothing else, so a "<" or
//     ">" is inert here.
//   - A single quote cannot close a value the renderer opened with a double one.
//   - The ampersand is the interesting omission. Encoding it would be more
//     correct in the strict sense — a bare "&" ought to be "&amp;" — but it
//     would also silently double-encode every caller who already passes
//     character references, which is exactly what a careful caller does when
//     building a JS string literal for an inline handler. Their "&#39;" would
//     render as "&amp;#39;" and appear on screen as text. That is a real
//     regression traded for a cosmetic fix, and a bare "&" cannot break out of
//     an attribute value.
//
// Callers wanting full entity encoding still have html.EscapeString, and
// applying it stays correct under this change: an already escaped value has no
// double quotes left for this function to find.
//
// The fast path matters — this runs for every attribute of every element, and
// almost no value contains a quote — so a value needing no work is written
// straight through after a single scan, with no allocation.
func writeAttrValue(sb *bytes.Buffer, val string) {
	idx := strings.IndexByte(val, '"')
	if idx < 0 { // the common case
		sb.WriteString(val)
		return
	}

	for {
		sb.WriteString(val[:idx])
		sb.WriteString("&#34;")
		val = val[idx+1:]

		idx = strings.IndexByte(val, '"')
		if idx < 0 {
			sb.WriteString(val)
			return
		}
	}
}

// Text is an element core function which creates a new text element in the string builder
func Text(s *bytes.Buffer, texts ...string) (x any) {
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
		Did you forget to wrap with builder.Text()? Child: %s`,
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
//
// T writes its arguments verbatim. That is what makes b.Style().T(css) and
// b.Script().T(js) work, and it is why T stays this way. Use TE for text that
// came from outside the program.
func (el Element) T(texts ...string) (x any) {
	el.R(Text(el.sb, texts...))
	return
}

// TE renders text-only children with HTML escaping applied — T, for text the
// program did not author.
//
// The difference between T and TE is one character on purpose, because the cost
// of picking wrong is not symmetric. Reaching for T on a database value, an API
// response or a request body is how markup in that data becomes markup in the
// page, and nothing about the output looks wrong until someone puts a tag in it.
// Reaching for TE on markup you meant to inline fails loudly and immediately —
// the tags appear on screen as text.
//
// There is deliberately no FE counterpart to F. Escaping a format string is
// ambiguous: the caller means "escape the arguments, not the template", and an
// API that guesses would be worse than one that makes the intent explicit.
//
//	b.Td().TE(fmt.Sprintf("%s (%s)", name, email))
func (el Element) TE(texts ...string) (x any) {
	escaped := make([]string, len(texts))
	for i, t := range texts {
		escaped[i] = html.EscapeString(t)
	}
	el.R(Text(el.sb, escaped...))
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
			el.sb.WriteByte('<')
			el.sb.WriteString(el.name)
			for i := 0; i+1 < len(el.attrPairs); i += 2 {
				el.sb.WriteByte(' ')
				el.sb.WriteString(el.attrPairs[i])
				el.sb.WriteString(`="`)
				writeAttrValue(el.sb, el.attrPairs[i+1])
				el.sb.WriteByte('"')
			}
			el.sb.WriteByte('>')
		}
	}
}

func (el Element) close() {
	if !el.IsSingleTag() {
		el.sb.WriteString("</")
		el.sb.WriteString(el.name)
		el.sb.WriteByte('>')
	}
}

func (el Element) details() string {
	return fmt.Sprintf("Element %s in %s (%s)", el.id, el.function, el.location)
}

func (el Element) detailsHtml() string {
	copyIcon := `<svg class="copy-icon" onclick="copyToClipboard('` + el.location + `')" title="Copy location" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="9" y="9" width="13" height="13" rx="2" ry="2"></rect><path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"></path></svg>`
	return fmt.Sprintf(`<strong>%s</strong> tag %s<br>%s (<strong>%s</strong> %s)`, el.name, el.id, el.function, el.location, copyIcon)
}
