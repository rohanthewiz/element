package element

// Builder Elements are just convenience methods for creating specific elements

func (b *Builder) Html(attrPairs ...string) Element {
	b.Text("<!DOCTYPE html>")
	return New(b.s, "html", attrPairs...)
}

func (b *Builder) Head(attrPairs ...string) Element {
	return New(b.s, "head", attrPairs...)
}

func (b *Builder) Title(attrPairs ...string) Element {
	return New(b.s, "title", attrPairs...)
}

func (b *Builder) Meta(attrPairs ...string) Element {
	return New(b.s, "meta", attrPairs...)
}

func (b *Builder) Style(attrPairs ...string) Element {
	return New(b.s, "style", attrPairs...)
}

func (b *Builder) Link(attrPairs ...string) Element {
	return New(b.s, "link", attrPairs...)
}

func (b *Builder) A(attrPairs ...string) Element {
	return New(b.s, "a", attrPairs...)
}

// B is for bold text
func (b *Builder) B(attrPairs ...string) Element {
	return New(b.s, "b", attrPairs...)
}

// I signifies italic text
func (b *Builder) I(attrPairs ...string) Element {
	return New(b.s, "i", attrPairs...)
}

// U is for underline text
func (b *Builder) U(attrPairs ...string) Element {
	return New(b.s, "u", attrPairs...)
}

func (b *Builder) Form(attrPairs ...string) Element {
	return New(b.s, "form", attrPairs...)
}

func (b *Builder) Input(attrPairs ...string) Element {
	return New(b.s, "input", attrPairs...)
}

func (b *Builder) Select(attrPairs ...string) Element {
	return New(b.s, "select", attrPairs...)
}

func (b *Builder) Option(attrPairs ...string) Element {
	return New(b.s, "option", attrPairs...)
}

func (b *Builder) Dd(attrPairs ...string) Element {
	return New(b.s, "dd", attrPairs...)
}

func (b *Builder) Dt(attrPairs ...string) Element {
	return New(b.s, "dt", attrPairs...)
}

func (b *Builder) Div(attrPairs ...string) Element {
	return New(b.s, "div", attrPairs...)
}

func (b *Builder) Body(attrPairs ...string) Element {
	return New(b.s, "body", attrPairs...)
}

func (b *Builder) P(attrPairs ...string) Element {
	return New(b.s, "p", attrPairs...)
}

func (b *Builder) Span(attrPairs ...string) Element {
	return New(b.s, "span", attrPairs...)
}

func (b *Builder) Table(attrPairs ...string) Element {
	return New(b.s, "table", attrPairs...)
}

func (b *Builder) THead(attrPairs ...string) Element {
	return New(b.s, "thead", attrPairs...)
}

func (b *Builder) TBody(attrPairs ...string) Element {
	return New(b.s, "tbody", attrPairs...)
}

func (b *Builder) Tr(attrPairs ...string) Element {
	return New(b.s, "tr", attrPairs...)
}

func (b *Builder) Th(attrPairs ...string) Element {
	return New(b.s, "th", attrPairs...)
}

func (b *Builder) Td(attrPairs ...string) Element {
	return New(b.s, "td", attrPairs...)
}

func (b *Builder) Section(attrPairs ...string) Element {
	return New(b.s, "section", attrPairs...)
}

func (b *Builder) H1(attrPairs ...string) Element {
	return New(b.s, "h1", attrPairs...)
}

func (b *Builder) H2(attrPairs ...string) Element {
	return New(b.s, "h2", attrPairs...)
}

func (b *Builder) H3(attrPairs ...string) Element {
	return New(b.s, "h3", attrPairs...)
}

func (b *Builder) H4(attrPairs ...string) Element {
	return New(b.s, "h4", attrPairs...)
}

func (b *Builder) Hr(attrPairs ...string) Element {
	return New(b.s, "hr", attrPairs...)
}

func (b *Builder) Ol(attrPairs ...string) Element {
	return New(b.s, "ol", attrPairs...)
}

func (b *Builder) Ul(attrPairs ...string) Element {
	return New(b.s, "ul", attrPairs...)
}

func (b *Builder) Li(attrPairs ...string) Element {
	return New(b.s, "li", attrPairs...)
}

func (b *Builder) Img(attrPairs ...string) Element {
	return New(b.s, "img", attrPairs...)
}
func (b *Builder) Article(attrPairs ...string) Element {
	return New(b.s, "article", attrPairs...)
}

func (b *Builder) Aside(attrPairs ...string) Element {
	return New(b.s, "aside", attrPairs...)
}

func (b *Builder) Audio(attrPairs ...string) Element {
	return New(b.s, "audio", attrPairs...)
}

func (b *Builder) BlockQuote(attrPairs ...string) Element {
	return New(b.s, "blockquote", attrPairs...)
}

func (b *Builder) Button(attrPairs ...string) Element {
	return New(b.s, "button", attrPairs...)
}

func (b *Builder) Canvas(attrPairs ...string) Element {
	return New(b.s, "canvas", attrPairs...)
}

func (b *Builder) Code(attrPairs ...string) Element {
	return New(b.s, "code", attrPairs...)
}

func (b *Builder) DataList(attrPairs ...string) Element {
	return New(b.s, "datalist", attrPairs...)
}

func (b *Builder) Details(attrPairs ...string) Element {
	return New(b.s, "details", attrPairs...)
}

func (b *Builder) Dialog(attrPairs ...string) Element {
	return New(b.s, "dialog", attrPairs...)
}

func (b *Builder) Em(attrPairs ...string) Element {
	return New(b.s, "em", attrPairs...)
}

func (b *Builder) Fieldset(attrPairs ...string) Element {
	return New(b.s, "fieldset", attrPairs...)
}

func (b *Builder) FigCaption(attrPairs ...string) Element {
	return New(b.s, "figcaption", attrPairs...)
}

func (b *Builder) Figure(attrPairs ...string) Element {
	return New(b.s, "figure", attrPairs...)
}

func (b *Builder) Footer(attrPairs ...string) Element {
	return New(b.s, "footer", attrPairs...)
}

func (b *Builder) Header(attrPairs ...string) Element {
	return New(b.s, "header", attrPairs...)
}

func (b *Builder) Iframe(attrPairs ...string) Element {
	return New(b.s, "iframe", attrPairs...)
}

func (b *Builder) Label(attrPairs ...string) Element {
	return New(b.s, "label", attrPairs...)
}

func (b *Builder) Legend(attrPairs ...string) Element {
	return New(b.s, "legend", attrPairs...)
}

func (b *Builder) Main(attrPairs ...string) Element {
	return New(b.s, "main", attrPairs...)
}

func (b *Builder) Nav(attrPairs ...string) Element {
	return New(b.s, "nav", attrPairs...)
}

func (b *Builder) Noscript(attrPairs ...string) Element {
	return New(b.s, "noscript", attrPairs...)
}

func (b *Builder) Object(attrPairs ...string) Element {
	return New(b.s, "object", attrPairs...)
}

func (b *Builder) Pre(attrPairs ...string) Element {
	return New(b.s, "pre", attrPairs...)
}

func (b *Builder) Progress(attrPairs ...string) Element {
	return New(b.s, "progress", attrPairs...)
}

func (b *Builder) Script(attrPairs ...string) Element {
	return New(b.s, "script", attrPairs...)
}

func (b *Builder) Strong(attrPairs ...string) Element {
	return New(b.s, "strong", attrPairs...)
}

func (b *Builder) Summary(attrPairs ...string) Element {
	return New(b.s, "summary", attrPairs...)
}

func (b *Builder) Svg(attrPairs ...string) Element {
	return New(b.s, "svg", attrPairs...)
}

func (b *Builder) Textarea(attrPairs ...string) Element {
	return New(b.s, "textarea", attrPairs...)
}

func (b *Builder) Time(attrPairs ...string) Element {
	return New(b.s, "time", attrPairs...)
}

func (b *Builder) Video(attrPairs ...string) Element {
	return New(b.s, "video", attrPairs...)
}
