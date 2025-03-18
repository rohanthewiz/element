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

func (b *Builder) Style(attrPairs ...string) Element {
	return New(b.s, "style", attrPairs...)
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

func (b *Builder) Image(attrPairs ...string) Element {
	return New(b.s, "img", attrPairs...)
}
