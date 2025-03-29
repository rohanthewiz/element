package element

// Builder Elements are just convenience methods for creating specific elements

// Html creates an HTML root element with doctype declaration
//
// Example:
//
//	b.Html("lang", "en")
func (b *Builder) Html(attrPairs ...string) Element {
	b.Text("<!DOCTYPE html>")
	return New(b.s, "html", attrPairs...)
}

// Head creates a head element for metadata content
//
// Example:
//
//	b.Head()
func (b *Builder) Head(attrPairs ...string) Element {
	return New(b.s, "head", attrPairs...)
}

// Title creates a title element for document title
//
// Example:
//
//	b.Title()
func (b *Builder) Title(attrPairs ...string) Element {
	return New(b.s, "title", attrPairs...)
}

// Meta creates a meta element for document metadata
//
// Example:
//
//	b.Meta("charset", "UTF-8")
//	b.Meta("name", "viewport", "content", "width=device-width, initial-scale=1.0")
func (b *Builder) Meta(attrPairs ...string) Element {
	return New(b.s, "meta", attrPairs...)
}

// Style creates a style element for CSS rules
//
// Example:
//
//	b.Style("type", "text/css")
func (b *Builder) Style(attrPairs ...string) Element {
	return New(b.s, "style", attrPairs...)
}

// Link creates a link element for external resources
//
// Example:
//
//	b.Link("rel", "stylesheet", "href", "styles/main.css", "type", "text/css")
//	b.Link("rel", "icon", "href", "images/favicon.ico", "type", "image/x-icon")
//	b.Link("rel", "preload", "href", "fonts/roboto.woff2", "as", "font", "crossorigin", "anonymous")
func (b *Builder) Link(attrPairs ...string) Element {
	return New(b.s, "link", attrPairs...)
}

// A creates an anchor element for hyperlinks
//
// Example:
//
//	b.A("href", "https://example.com", "target", "_blank")
//	b.A("href", "#section1", "class", "internal-link")
func (b *Builder) A(attrPairs ...string) Element {
	return New(b.s, "a", attrPairs...)
}

// B is for bold text
//
// Example:
//
//	b.B("class", "highlight")
func (b *Builder) B(attrPairs ...string) Element {
	return New(b.s, "b", attrPairs...)
}

// I signifies italic text
//
// Example:
//
//	b.I("class", "emphasis")
func (b *Builder) I(attrPairs ...string) Element {
	return New(b.s, "i", attrPairs...)
}

// U is for underline text
//
// Example:
//
//	b.U("class", "underlined")
func (b *Builder) U(attrPairs ...string) Element {
	return New(b.s, "u", attrPairs...)
}

// Form creates a form container for user input
//
// Example:
//
//	b.Form("action", "/submit", "method", "post", "id", "contactForm")
func (b *Builder) Form(attrPairs ...string) Element {
	return New(b.s, "form", attrPairs...)
}

// Input creates an input control element
//
// Example:
//
//	b.Input("type", "text", "name", "username", "placeholder", "Enter username")
//	b.Input("type", "email", "name", "email", "required", "true")
//	b.Input("type", "checkbox", "name", "subscribe", "checked", "checked")
func (b *Builder) Input(attrPairs ...string) Element {
	return New(b.s, "input", attrPairs...)
}

// Select creates a dropdown selection control
//
// Example:
//
//	b.Select("name", "country", "id", "country-select")
func (b *Builder) Select(attrPairs ...string) Element {
	return New(b.s, "select", attrPairs...)
}

// Option creates an option item for a select element
//
// Example:
//
//	b.Option("value", "us", "selected", "selected")
func (b *Builder) Option(attrPairs ...string) Element {
	return New(b.s, "option", attrPairs...)
}

// Dd creates a description details element in a description list
//
// Example:
//
//	b.Dd("class", "definition")
func (b *Builder) Dd(attrPairs ...string) Element {
	return New(b.s, "dd", attrPairs...)
}

// Dt creates a description term element in a description list
//
// Example:
//
//	b.Dt("class", "term")
func (b *Builder) Dt(attrPairs ...string) Element {
	return New(b.s, "dt", attrPairs...)
}

// Div creates a division container element
//
// Example:
//
//	b.Div("class", "container", "id", "main-content")
//	b.Div("data-role", "panel", "class", "sidebar")
func (b *Builder) Div(attrPairs ...string) Element {
	return New(b.s, "div", attrPairs...)
}

// Body creates the document body element
//
// Example:
//
//	b.Body("class", "light-theme")
func (b *Builder) Body(attrPairs ...string) Element {
	return New(b.s, "body", attrPairs...)
}

// P creates a paragraph element
//
// Example:
//
//	b.P("class", "intro", "id", "first-paragraph")
func (b *Builder) P(attrPairs ...string) Element {
	return New(b.s, "p", attrPairs...)
}

// Span creates an inline container element
//
// Example:
//
//	b.Span("class", "highlight", "data-tooltip", "Important information")
func (b *Builder) Span(attrPairs ...string) Element {
	return New(b.s, "span", attrPairs...)
}

// Table creates a table element for tabular data
//
// Example:
//
//	b.Table("class", "data-table", "id", "results")
func (b *Builder) Table(attrPairs ...string) Element {
	return New(b.s, "table", attrPairs...)
}

// THead creates a table header group element
//
// Example:
//
//	b.THead("class", "header-group")
func (b *Builder) THead(attrPairs ...string) Element {
	return New(b.s, "thead", attrPairs...)
}

// TBody creates a table body group element
//
// Example:
//
//	b.TBody("class", "table-data")
func (b *Builder) TBody(attrPairs ...string) Element {
	return New(b.s, "tbody", attrPairs...)
}

// Tr creates a table row element
//
// Example:
//
//	b.Tr("class", "data-row", "id", "row-1")
func (b *Builder) Tr(attrPairs ...string) Element {
	return New(b.s, "tr", attrPairs...)
}

// Th creates a table header cell element
//
// Example:
//
//	b.Th("scope", "col", "class", "header-cell")
func (b *Builder) Th(attrPairs ...string) Element {
	return New(b.s, "th", attrPairs...)
}

// Td creates a table data cell element
//
// Example:
//
//	b.Td("class", "data-cell", "rowspan", "2")
func (b *Builder) Td(attrPairs ...string) Element {
	return New(b.s, "td", attrPairs...)
}

// Section creates a standalone section element
//
// Example:
//
//	b.Section("class", "content-section", "id", "about")
func (b *Builder) Section(attrPairs ...string) Element {
	return New(b.s, "section", attrPairs...)
}

// H1 creates a top-level heading element
//
// Example:
//
//	b.H1("class", "page-title", "id", "main-heading")
func (b *Builder) H1(attrPairs ...string) Element {
	return New(b.s, "h1", attrPairs...)
}

// H2 creates a second-level heading element
//
// Example:
//
//	b.H2("class", "section-title", "id", "section-2")
func (b *Builder) H2(attrPairs ...string) Element {
	return New(b.s, "h2", attrPairs...)
}

// H3 creates a third-level heading element
//
// Example:
//
//	b.H3("class", "subsection-title")
func (b *Builder) H3(attrPairs ...string) Element {
	return New(b.s, "h3", attrPairs...)
}

// H4 creates a fourth-level heading element
//
// Example:
//
//	b.H4("class", "minor-title")
func (b *Builder) H4(attrPairs ...string) Element {
	return New(b.s, "h4", attrPairs...)
}

// Hr creates a horizontal rule (line) element
//
// Example:
//
//	b.Hr("class", "divider")
func (b *Builder) Hr(attrPairs ...string) Element {
	return New(b.s, "hr", attrPairs...)
}

// Ol creates an ordered list element
//
// Example:
//
//	b.Ol("class", "numbered-list", "start", "3")
func (b *Builder) Ol(attrPairs ...string) Element {
	return New(b.s, "ol", attrPairs...)
}

// Ul creates an unordered list element
//
// Example:
//
//	b.Ul("class", "bullet-list", "id", "main-menu")
func (b *Builder) Ul(attrPairs ...string) Element {
	return New(b.s, "ul", attrPairs...)
}

// Li creates a list item element
//
// Example:
//
//	b.Li("class", "list-item", "data-index", "1")
func (b *Builder) Li(attrPairs ...string) Element {
	return New(b.s, "li", attrPairs...)
}

// Img creates an image element
//
// Example:
//
//	b.Img("src", "images/logo.png", "alt", "Company Logo", "width", "200", "height", "100")
func (b *Builder) Img(attrPairs ...string) Element {
	return New(b.s, "img", attrPairs...)
}

// Article creates an article element
//
// Example:
//
//	b.Article("class", "blog-post", "id", "post-123")
func (b *Builder) Article(attrPairs ...string) Element {
	return New(b.s, "article", attrPairs...)
}

// Aside creates an aside element
//
// Example:
//
//	b.Aside("class", "sidebar", "aria-label", "Related content")
func (b *Builder) Aside(attrPairs ...string) Element {
	return New(b.s, "aside", attrPairs...)
}

func (b *Builder) Attr(attrPairs ...string) Element {
	return New(b.s, "attr", attrPairs...)
}

// Audio creates an audio element
//
// Example:
//
//	b.Audio("src", "audio/sound.mp3", "controls", "true", "autoplay", "false")
func (b *Builder) Audio(attrPairs ...string) Element {
	return New(b.s, "audio", attrPairs...)
}

// BlockQuote creates a blockquote element
//
// Example:
//
//	b.BlockQuote("cite", "https://example.com/source", "class", "quote")
func (b *Builder) BlockQuote(attrPairs ...string) Element {
	return New(b.s, "blockquote", attrPairs...)
}

// Button creates a clickable button element
//
// Example:
//
//	b.Button("type", "submit", "class", "btn-primary")
//	b.Button("type", "button", "onclick", "doSomething()", "disabled", "disabled")
func (b *Builder) Button(attrPairs ...string) Element {
	return New(b.s, "button", attrPairs...)
}

// Canvas creates a drawable canvas element
//
// Example:
//
//	b.Canvas("width", "800", "height", "600", "id", "gameCanvas")
func (b *Builder) Canvas(attrPairs ...string) Element {
	return New(b.s, "canvas", attrPairs...)
}

// Code creates an element for code snippets
//
// Example:
//
//	b.Code("class", "language-go")
func (b *Builder) Code(attrPairs ...string) Element {
	return New(b.s, "code", attrPairs...)
}

// Dl creates a description list element
//
// Example:
//
//	b.Dl("class", "glossary")
func (b *Builder) Dl(attrPairs ...string) Element {
	return New(b.s, "dl", attrPairs...)
}

// Em creates an emphasis element
//
// Example:
//
//	b.Em("class", "highlighted")
func (b *Builder) Em(attrPairs ...string) Element {
	return New(b.s, "em", attrPairs...)
}

// Footer creates a footer element
//
// Example:
//
//	b.Footer("class", "site-footer", "id", "main-footer")
func (b *Builder) Footer(attrPairs ...string) Element {
	return New(b.s, "footer", attrPairs...)
}

// Header creates a header element
//
// Example:
//
//	b.Header("class", "page-header", "id", "main-header")
func (b *Builder) Header(attrPairs ...string) Element {
	return New(b.s, "header", attrPairs...)
}

// IFrame creates an inline frame element
//
// Example:
//
//	b.IFrame("src", "https://example.com/embed", "width", "560", "height", "315", "frameborder", "0")
func (b *Builder) IFrame(attrPairs ...string) Element {
	return New(b.s, "iframe", attrPairs...)
}

// Label creates a label for form controls
//
// Example:
//
//	b.Label("for", "username", "class", "form-label")
func (b *Builder) Label(attrPairs ...string) Element {
	return New(b.s, "label", attrPairs...)
}

// Main creates the main content element
//
// Example:
//
//	b.Main("id", "main-content", "role", "main")
func (b *Builder) Main(attrPairs ...string) Element {
	return New(b.s, "main", attrPairs...)
}

// Nav creates a navigation element
//
// Example:
//
//	b.Nav("class", "main-navigation", "aria-label", "Main Navigation")
func (b *Builder) Nav(attrPairs ...string) Element {
	return New(b.s, "nav", attrPairs...)
}

// Picture creates a container for multiple image sources
//
// Example:
//
//	b.Picture("class", "responsive-image")
func (b *Builder) Picture(attrPairs ...string) Element {
	return New(b.s, "picture", attrPairs...)
}

// Pre creates a preformatted text element
//
// Example:
//
//	b.Pre("class", "code-block")
func (b *Builder) Pre(attrPairs ...string) Element {
	return New(b.s, "pre", attrPairs...)
}

// Script creates a script element
//
// Example:
//
//	b.Script("src", "js/main.js", "defer", "true")
//	b.Script("type", "module", "src", "js/app.js")
func (b *Builder) Script(attrPairs ...string) Element {
	return New(b.s, "script", attrPairs...)
}

// Source creates a media source element
//
// Example:
//
//	b.Source("src", "video/intro.mp4", "type", "video/mp4")
func (b *Builder) Source(attrPairs ...string) Element {
	return New(b.s, "source", attrPairs...)
}

// Time creates a time element
//
// Example:
//
//	b.Time("datetime", "2023-12-25T20:00:00", "class", "event-time")
func (b *Builder) Time(attrPairs ...string) Element {
	return New(b.s, "time", attrPairs...)
}

// Video creates a video element
//
// Example:
//
//	b.Video("src", "videos/tutorial.mp4", "controls", "true", "width", "640", "height", "360")
func (b *Builder) Video(attrPairs ...string) Element {
	return New(b.s, "video", attrPairs...)
}

// Track creates a text track element
//
// Example:
//
//	b.Track("src", "subtitles-en.vtt", "kind", "subtitles", "srclang", "en", "label", "English")
//	b.Track("src", "captions.vtt", "kind", "captions", "default", "true")
func (b *Builder) Track(attrPairs ...string) Element {
	return New(b.s, "track", attrPairs...)
}

// Abbr creates an abbreviation element
//
// Example:
//
//	b.Abbr("title", "Hypertext Markup Language", "class", "term")
func (b *Builder) Abbr(attrPairs ...string) Element {
	return New(b.s, "abbr", attrPairs...)
}

// Caption creates a table caption element
//
// Example:
//
//	b.Caption("class", "table-title")
func (b *Builder) Caption(attrPairs ...string) Element {
	return New(b.s, "caption", attrPairs...)
}

// FieldSet creates a group of form controls
//
// Example:
//
//	b.FieldSet("name", "userInfo", "class", "form-group")
func (b *Builder) FieldSet(attrPairs ...string) Element {
	return New(b.s, "fieldset", attrPairs...)
}

// Legend creates a caption for a fieldset
//
// Example:
//
//	b.Legend("class", "form-legend")
func (b *Builder) Legend(attrPairs ...string) Element {
	return New(b.s, "legend", attrPairs...)
}

// Progress creates a progress indicator element
//
// Example:
//
//	b.Progress("value", "70", "max", "100")
//	b.Progress("class", "upload-progress", "id", "file-upload")
func (b *Builder) Progress(attrPairs ...string) Element {
	return New(b.s, "progress", attrPairs...)
}

// Q creates a short inline quotation element
//
// Example:
//
//	b.Q("cite", "https://example.com/quote-source")
//	b.Q("class", "quoted-text", "lang", "en")
func (b *Builder) Q(attrPairs ...string) Element {
	return New(b.s, "q", attrPairs...)
}

// Ruby creates a ruby annotation element
//
// Example:
//
//	b.Ruby("class", "pronunciation")
func (b *Builder) Ruby(attrPairs ...string) Element {
	return New(b.s, "ruby", attrPairs...)
}

// Rt creates a ruby text element
//
// Example:
//
//	b.Rt("class", "pronunciation-text")
func (b *Builder) Rt(attrPairs ...string) Element {
	return New(b.s, "rt", attrPairs...)
}

// Noscript creates content to display when scripts are disabled
//
// Example:
//
//	b.Noscript("class", "no-js-message")
func (b *Builder) Noscript(attrPairs ...string) Element {
	return New(b.s, "noscript", attrPairs...)
}

// Object creates an embedded external resource container
func (b *Builder) Object(attrPairs ...string) Element {
	return New(b.s, "object", attrPairs...)
}

// Small creates small text element
//
// Example:
//
//	b.Small("class", "copyright-notice")
//	b.Small("id", "disclaimer")
func (b *Builder) Small(attrPairs ...string) Element {
	return New(b.s, "small", attrPairs...)
}

// Strong creates a strong importance element
//
// Example:
//
//	b.Strong("class", "important")
//	b.Strong("id", "warning-text")
func (b *Builder) Strong(attrPairs ...string) Element {
	return New(b.s, "strong", attrPairs...)
}

// Sub creates a subscript element
//
// Example:
//
//	b.Sub("class", "chemical-formula")
func (b *Builder) Sub(attrPairs ...string) Element {
	return New(b.s, "sub", attrPairs...)
}

// Summary creates a disclosure summary element
//
// Example:
//
//	b.Summary("class", "accordion-header")
//	b.Summary("id", "section-title")
func (b *Builder) Summary(attrPairs ...string) Element {
	return New(b.s, "summary", attrPairs...)
}

// Sup creates a superscript element
//
// Example:
//
//	b.Sup("class", "footnote")
func (b *Builder) Sup(attrPairs ...string) Element {
	return New(b.s, "sup", attrPairs...)
}

// TFoot creates a table footer element
//
// Example:
//
//	b.TFoot("class", "table-footer")
func (b *Builder) TFoot(attrPairs ...string) Element {
	return New(b.s, "tfoot", attrPairs...)
}

// Svg creates a container for SVG graphics
func (b *Builder) Svg(attrPairs ...string) Element {
	return New(b.s, "svg", attrPairs...)
}

// TextArea creates a multiline text input control
//
// Example:
//
//	b.TextArea("name", "comments", "rows", "4", "cols", "50", "placeholder", "Enter your comments")
//	b.TextArea("id", "message", "class", "form-control", "required", "true")
func (b *Builder) TextArea(attrPairs ...string) Element {
	return New(b.s, "textarea", attrPairs...)
}

// DataList creates a datalist element for providing a list of predefined options for input elements
//
// Example:
//
//	b.DataList("id", "browsers")
func (b *Builder) DataList(attrPairs ...string) Element {
	return New(b.s, "datalist", attrPairs...)
}

// Details creates a disclosure widget that can show/hide information
//
// Example:
//
//	b.Details("open", "true", "class", "expandable-section")
func (b *Builder) Details(attrPairs ...string) Element {
	return New(b.s, "details", attrPairs...)
}

// Dialog creates a modal or non-modal dialog box
//
// Example:
//
//	b.Dialog("id", "confirmation-dialog", "open", "true")
func (b *Builder) Dialog(attrPairs ...string) Element {
	return New(b.s, "dialog", attrPairs...)
}

// Fieldset creates a container for grouping related form elements
//
// Example:
//
//	b.Fieldset("name", "personal-info", "class", "form-group")
func (b *Builder) Fieldset(attrPairs ...string) Element {
	return New(b.s, "fieldset", attrPairs...)
}

// FigCaption creates a caption for a figure element
//
// Example:
//
//	b.FigCaption("class", "image-caption")
func (b *Builder) FigCaption(attrPairs ...string) Element {
	return New(b.s, "figcaption", attrPairs...)
}

// Figure creates a container for self-contained content like images, diagrams, etc.
//
// Example:
//
//	b.Figure("class", "image-container", "id", "main-diagram")
func (b *Builder) Figure(attrPairs ...string) Element {
	return New(b.s, "figure", attrPairs...)
}
