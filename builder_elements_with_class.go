package element

// AClass creates an anchor element with a class
//
// Example:
//
//	b.AClass("btn", "href", "https://example.com")
func (b *Builder) AClass(class string, attrPairs ...string) Element {
	return New(b.s, "a", append([]string{"class", class}, attrPairs...)...)
}

// BClass creates a bold text element with a class
//
// Example:
//
//	b.BClass("highlight")
func (b *Builder) BClass(class string, attrPairs ...string) Element {
	return New(b.s, "b", append([]string{"class", class}, attrPairs...)...)
}

// IClass creates an italic text element with a class
//
// Example:
//
//	b.IClass("emphasis")
func (b *Builder) IClass(class string, attrPairs ...string) Element {
	return New(b.s, "i", append([]string{"class", class}, attrPairs...)...)
}

// UClass creates an underline text element with a class
//
// Example:
//
//	b.UClass("underlined")
func (b *Builder) UClass(class string, attrPairs ...string) Element {
	return New(b.s, "u", append([]string{"class", class}, attrPairs...)...)
}

// FormClass creates a form container with a class
//
// Example:
//
//	b.FormClass("contact-form", "action", "/submit", "method", "post")
func (b *Builder) FormClass(class string, attrPairs ...string) Element {
	return New(b.s, "form", append([]string{"class", class}, attrPairs...)...)
}

// InputClass creates an input control element with a class
//
// Example:
//
//	b.InputClass("form-control", "type", "text", "name", "username")
func (b *Builder) InputClass(class string, attrPairs ...string) Element {
	return New(b.s, "input", append([]string{"class", class}, attrPairs...)...)
}

// SelectClass creates a dropdown selection control with a class
//
// Example:
//
//	b.SelectClass("form-select", "name", "country")
func (b *Builder) SelectClass(class string, attrPairs ...string) Element {
	return New(b.s, "select", append([]string{"class", class}, attrPairs...)...)
}

// OptionClass creates an option item with a class
//
// Example:
//
//	b.OptionClass("selected-option", "value", "us")
func (b *Builder) OptionClass(class string, attrPairs ...string) Element {
	return New(b.s, "option", append([]string{"class", class}, attrPairs...)...)
}

// DdClass creates a description details element with a class
//
// Example:
//
//	b.DdClass("definition")
func (b *Builder) DdClass(class string, attrPairs ...string) Element {
	return New(b.s, "dd", append([]string{"class", class}, attrPairs...)...)
}

// DtClass creates a description term element with a class
//
// Example:
//
//	b.DtClass("term")
func (b *Builder) DtClass(class string, attrPairs ...string) Element {
	return New(b.s, "dt", append([]string{"class", class}, attrPairs...)...)
}

// DivClass creates a division container element with a class
//
// Example:
//
//	b.DivClass("container", "id", "main-content")
func (b *Builder) DivClass(class string, attrPairs ...string) Element {
	return New(b.s, "div", append([]string{"class", class}, attrPairs...)...)
}

// BodyClass creates the document body element with a class
//
// Example:
//
//	b.BodyClass("light-theme")
func (b *Builder) BodyClass(class string, attrPairs ...string) Element {
	return New(b.s, "body", append([]string{"class", class}, attrPairs...)...)
}

// PClass creates a paragraph element with a class
//
// Example:
//
//	b.PClass("intro", "id", "first-paragraph")
func (b *Builder) PClass(class string, attrPairs ...string) Element {
	return New(b.s, "p", append([]string{"class", class}, attrPairs...)...)
}

// BrClass creates a line break element with a class
//
// Example:
//
//	b.BrClass("visible-break")
func (b *Builder) BrClass(class string, attrPairs ...string) Element {
	return New(b.s, "br", append([]string{"class", class}, attrPairs...)...)
}

// SpanClass creates an inline container element with a class
//
// Example:
//
//	b.SpanClass("highlight", "data-tooltip", "Important information")
func (b *Builder) SpanClass(class string, attrPairs ...string) Element {
	return New(b.s, "span", append([]string{"class", class}, attrPairs...)...)
}

// TableClass creates a table element with a class
//
// Example:
//
//	b.TableClass("data-table", "id", "results")
func (b *Builder) TableClass(class string, attrPairs ...string) Element {
	return New(b.s, "table", append([]string{"class", class}, attrPairs...)...)
}

// THeadClass creates a table header group element with a class
//
// Example:
//
//	b.THeadClass("header-group")
func (b *Builder) THeadClass(class string, attrPairs ...string) Element {
	return New(b.s, "thead", append([]string{"class", class}, attrPairs...)...)
}

// TBodyClass creates a table body group element with a class
//
// Example:
//
//	b.TBodyClass("table-data")
func (b *Builder) TBodyClass(class string, attrPairs ...string) Element {
	return New(b.s, "tbody", append([]string{"class", class}, attrPairs...)...)
}

// TrClass creates a table row element with a class
//
// Example:
//
//	b.TrClass("data-row", "id", "row-1")
func (b *Builder) TrClass(class string, attrPairs ...string) Element {
	return New(b.s, "tr", append([]string{"class", class}, attrPairs...)...)
}

// ThClass creates a table header cell element with a class
//
// Example:
//
//	b.ThClass("header-cell", "scope", "col")
func (b *Builder) ThClass(class string, attrPairs ...string) Element {
	return New(b.s, "th", append([]string{"class", class}, attrPairs...)...)
}

// TdClass creates a table data cell element with a class
//
// Example:
//
//	b.TdClass("data-cell", "rowspan", "2")
func (b *Builder) TdClass(class string, attrPairs ...string) Element {
	return New(b.s, "td", append([]string{"class", class}, attrPairs...)...)
}

// SectionClass creates a standalone section element with a class
//
// Example:
//
//	b.SectionClass("content-section", "id", "about")
func (b *Builder) SectionClass(class string, attrPairs ...string) Element {
	return New(b.s, "section", append([]string{"class", class}, attrPairs...)...)
}

// H1Class creates a top-level heading element with a class
//
// Example:
//
//	b.H1Class("page-title", "id", "main-heading")
func (b *Builder) H1Class(class string, attrPairs ...string) Element {
	return New(b.s, "h1", append([]string{"class", class}, attrPairs...)...)
}

// H2Class creates a second-level heading element with a class
//
// Example:
//
//	b.H2Class("section-title", "id", "section-2")
func (b *Builder) H2Class(class string, attrPairs ...string) Element {
	return New(b.s, "h2", append([]string{"class", class}, attrPairs...)...)
}

// H3Class creates a third-level heading element with a class
//
// Example:
//
//	b.H3Class("subsection-title")
func (b *Builder) H3Class(class string, attrPairs ...string) Element {
	return New(b.s, "h3", append([]string{"class", class}, attrPairs...)...)
}

// H4Class creates a fourth-level heading element with a class
//
// Example:
//
//	b.H4Class("minor-title")
func (b *Builder) H4Class(class string, attrPairs ...string) Element {
	return New(b.s, "h4", append([]string{"class", class}, attrPairs...)...)
}

// HrClass creates a horizontal rule element with a class
//
// Example:
//
//	b.HrClass("divider")
func (b *Builder) HrClass(class string, attrPairs ...string) Element {
	return New(b.s, "hr", append([]string{"class", class}, attrPairs...)...)
}

// OlClass creates an ordered list element with a class
//
// Example:
//
//	b.OlClass("numbered-list", "start", "3")
func (b *Builder) OlClass(class string, attrPairs ...string) Element {
	return New(b.s, "ol", append([]string{"class", class}, attrPairs...)...)
}

// UlClass creates an unordered list element with a class
//
// Example:
//
//	b.UlClass("bullet-list", "id", "main-menu")
func (b *Builder) UlClass(class string, attrPairs ...string) Element {
	return New(b.s, "ul", append([]string{"class", class}, attrPairs...)...)
}

// LiClass creates a list item element with a class
//
// Example:
//
//	b.LiClass("list-item", "data-index", "1")
func (b *Builder) LiClass(class string, attrPairs ...string) Element {
	return New(b.s, "li", append([]string{"class", class}, attrPairs...)...)
}

// ImgClass creates an image element with a class
//
// Example:
//
//	b.ImgClass("logo", "src", "images/logo.png", "alt", "Company Logo")
func (b *Builder) ImgClass(class string, attrPairs ...string) Element {
	return New(b.s, "img", append([]string{"class", class}, attrPairs...)...)
}

// ArticleClass creates an article element with a class
//
// Example:
//
//	b.ArticleClass("blog-post", "id", "post-123")
func (b *Builder) ArticleClass(class string, attrPairs ...string) Element {
	return New(b.s, "article", append([]string{"class", class}, attrPairs...)...)
}

// AsideClass creates an aside element with a class
//
// Example:
//
//	b.AsideClass("sidebar", "aria-label", "Related content")
func (b *Builder) AsideClass(class string, attrPairs ...string) Element {
	return New(b.s, "aside", append([]string{"class", class}, attrPairs...)...)
}

// AudioClass creates an audio element with a class
//
// Example:
//
//	b.AudioClass("player", "src", "audio/sound.mp3", "controls", "true")
func (b *Builder) AudioClass(class string, attrPairs ...string) Element {
	return New(b.s, "audio", append([]string{"class", class}, attrPairs...)...)
}

// BlockQuoteClass creates a blockquote element with a class
//
// Example:
//
//	b.BlockQuoteClass("quote", "cite", "https://example.com/source")
func (b *Builder) BlockQuoteClass(class string, attrPairs ...string) Element {
	return New(b.s, "blockquote", append([]string{"class", class}, attrPairs...)...)
}

// ButtonClass creates a clickable button element with a class
//
// Example:
//
//	b.ButtonClass("btn-primary", "type", "submit")
func (b *Builder) ButtonClass(class string, attrPairs ...string) Element {
	return New(b.s, "button", append([]string{"class", class}, attrPairs...)...)
}

// CanvasClass creates a drawable canvas element with a class
//
// Example:
//
//	b.CanvasClass("game-canvas", "width", "800", "height", "600")
func (b *Builder) CanvasClass(class string, attrPairs ...string) Element {
	return New(b.s, "canvas", append([]string{"class", class}, attrPairs...)...)
}

// CodeClass creates an element for code snippets with a class
//
// Example:
//
//	b.CodeClass("language-go")
func (b *Builder) CodeClass(class string, attrPairs ...string) Element {
	return New(b.s, "code", append([]string{"class", class}, attrPairs...)...)
}

// DlClass creates a description list element with a class
//
// Example:
//
//	b.DlClass("glossary")
func (b *Builder) DlClass(class string, attrPairs ...string) Element {
	return New(b.s, "dl", append([]string{"class", class}, attrPairs...)...)
}

// EmClass creates an emphasis element with a class
//
// Example:
//
//	b.EmClass("highlighted")
func (b *Builder) EmClass(class string, attrPairs ...string) Element {
	return New(b.s, "em", append([]string{"class", class}, attrPairs...)...)
}

// FooterClass creates a footer element with a class
//
// Example:
//
//	b.FooterClass("site-footer", "id", "main-footer")
func (b *Builder) FooterClass(class string, attrPairs ...string) Element {
	return New(b.s, "footer", append([]string{"class", class}, attrPairs...)...)
}

// HeaderClass creates a header element with a class
//
// Example:
//
//	b.HeaderClass("page-header", "id", "main-header")
func (b *Builder) HeaderClass(class string, attrPairs ...string) Element {
	return New(b.s, "header", append([]string{"class", class}, attrPairs...)...)
}

// LabelClass creates a label for form controls with a class
//
// Example:
//
//	b.LabelClass("form-label", "for", "username")
func (b *Builder) LabelClass(class string, attrPairs ...string) Element {
	return New(b.s, "label", append([]string{"class", class}, attrPairs...)...)
}

// MainClass creates the main content element with a class
//
// Example:
//
//	b.MainClass("content", "id", "main-content", "role", "main")
func (b *Builder) MainClass(class string, attrPairs ...string) Element {
	return New(b.s, "main", append([]string{"class", class}, attrPairs...)...)
}

// NavClass creates a navigation element with a class
//
// Example:
//
//	b.NavClass("main-navigation", "aria-label", "Main Navigation")
func (b *Builder) NavClass(class string, attrPairs ...string) Element {
	return New(b.s, "nav", append([]string{"class", class}, attrPairs...)...)
}

// PictureClass creates a container for multiple image sources with a class
//
// Example:
//
//	b.PictureClass("responsive-image")
func (b *Builder) PictureClass(class string, attrPairs ...string) Element {
	return New(b.s, "picture", append([]string{"class", class}, attrPairs...)...)
}

// PreClass creates a preformatted text element with a class
//
// Example:
//
//	b.PreClass("code-block")
func (b *Builder) PreClass(class string, attrPairs ...string) Element {
	return New(b.s, "pre", append([]string{"class", class}, attrPairs...)...)
}

// TimeClass creates a time element with a class
//
// Example:
//
//	b.TimeClass("event-time", "datetime", "2023-12-25T20:00:00")
func (b *Builder) TimeClass(class string, attrPairs ...string) Element {
	return New(b.s, "time", append([]string{"class", class}, attrPairs...)...)
}

// VideoClass creates a video element with a class
//
// Example:
//
//	b.VideoClass("video-player", "src", "videos/tutorial.mp4", "controls", "true")
func (b *Builder) VideoClass(class string, attrPairs ...string) Element {
	return New(b.s, "video", append([]string{"class", class}, attrPairs...)...)
}

// TrackClass creates a text track element with a class
//
// Example:
//
//	b.TrackClass("subtitle", "src", "subtitles-en.vtt", "kind", "subtitles")
func (b *Builder) TrackClass(class string, attrPairs ...string) Element {
	return New(b.s, "track", append([]string{"class", class}, attrPairs...)...)
}

// AbbrClass creates an abbreviation element with a class
//
// Example:
//
//	b.AbbrClass("term", "title", "Hypertext Markup Language")
func (b *Builder) AbbrClass(class string, attrPairs ...string) Element {
	return New(b.s, "abbr", append([]string{"class", class}, attrPairs...)...)
}

// CaptionClass creates a table caption element with a class
//
// Example:
//
//	b.CaptionClass("table-title")
func (b *Builder) CaptionClass(class string, attrPairs ...string) Element {
	return New(b.s, "caption", append([]string{"class", class}, attrPairs...)...)
}

// FieldSetClass creates a group of form controls with a class
//
// Example:
//
//	b.FieldSetClass("form-group", "name", "userInfo")
func (b *Builder) FieldSetClass(class string, attrPairs ...string) Element {
	return New(b.s, "fieldset", append([]string{"class", class}, attrPairs...)...)
}

// LegendClass creates a caption for a fieldset with a class
//
// Example:
//
//	b.LegendClass("form-legend")
func (b *Builder) LegendClass(class string, attrPairs ...string) Element {
	return New(b.s, "legend", append([]string{"class", class}, attrPairs...)...)
}

// ProgressClass creates a progress indicator element with a class
//
// Example:
//
//	b.ProgressClass("upload-progress", "value", "70", "max", "100")
func (b *Builder) ProgressClass(class string, attrPairs ...string) Element {
	return New(b.s, "progress", append([]string{"class", class}, attrPairs...)...)
}

// QClass creates a short inline quotation element with a class
//
// Example:
//
//	b.QClass("quoted-text", "cite", "https://example.com/quote-source")
func (b *Builder) QClass(class string, attrPairs ...string) Element {
	return New(b.s, "q", append([]string{"class", class}, attrPairs...)...)
}

// RubyClass creates a ruby annotation element with a class
//
// Example:
//
//	b.RubyClass("pronunciation")
func (b *Builder) RubyClass(class string, attrPairs ...string) Element {
	return New(b.s, "ruby", append([]string{"class", class}, attrPairs...)...)
}

// RtClass creates a ruby text element with a class
//
// Example:
//
//	b.RtClass("pronunciation-text")
func (b *Builder) RtClass(class string, attrPairs ...string) Element {
	return New(b.s, "rt", append([]string{"class", class}, attrPairs...)...)
}

// NoscriptClass creates content for when scripts are disabled with a class
//
// Example:
//
//	b.NoscriptClass("no-js-message")
func (b *Builder) NoscriptClass(class string, attrPairs ...string) Element {
	return New(b.s, "noscript", append([]string{"class", class}, attrPairs...)...)
}

// SmallClass creates small text element with a class
//
// Example:
//
//	b.SmallClass("copyright-notice")
func (b *Builder) SmallClass(class string, attrPairs ...string) Element {
	return New(b.s, "small", append([]string{"class", class}, attrPairs...)...)
}

// StrongClass creates a strong importance element with a class
//
// Example:
//
//	b.StrongClass("important")
func (b *Builder) StrongClass(class string, attrPairs ...string) Element {
	return New(b.s, "strong", append([]string{"class", class}, attrPairs...)...)
}

// SubClass creates a subscript element with a class
//
// Example:
//
//	b.SubClass("chemical-formula")
func (b *Builder) SubClass(class string, attrPairs ...string) Element {
	return New(b.s, "sub", append([]string{"class", class}, attrPairs...)...)
}

// SummaryClass creates a disclosure summary element with a class
//
// Example:
//
//	b.SummaryClass("accordion-header")
func (b *Builder) SummaryClass(class string, attrPairs ...string) Element {
	return New(b.s, "summary", append([]string{"class", class}, attrPairs...)...)
}

// SupClass creates a superscript element with a class
//
// Example:
//
//	b.SupClass("footnote")
func (b *Builder) SupClass(class string, attrPairs ...string) Element {
	return New(b.s, "sup", append([]string{"class", class}, attrPairs...)...)
}

// TFootClass creates a table footer element with a class
//
// Example:
//
//	b.TFootClass("table-footer")
func (b *Builder) TFootClass(class string, attrPairs ...string) Element {
	return New(b.s, "tfoot", append([]string{"class", class}, attrPairs...)...)
}

// SvgClass creates a container for SVG graphics with a class
//
// Example:
//
//	b.SvgClass("icon")
func (b *Builder) SvgClass(class string, attrPairs ...string) Element {
	return New(b.s, "svg", append([]string{"class", class}, attrPairs...)...)
}

// TextAreaClass creates a multiline text input control with a class
//
// Example:
//
//	b.TextAreaClass("form-control", "name", "comments", "rows", "4")
func (b *Builder) TextAreaClass(class string, attrPairs ...string) Element {
	return New(b.s, "textarea", append([]string{"class", class}, attrPairs...)...)
}

// DataListClass creates a datalist element with a class
//
// Example:
//
//	b.DataListClass("autocomplete-options", "id", "browsers")
func (b *Builder) DataListClass(class string, attrPairs ...string) Element {
	return New(b.s, "datalist", append([]string{"class", class}, attrPairs...)...)
}

// DetailsClass creates a disclosure widget with a class
//
// Example:
//
//	b.DetailsClass("expandable-section", "open", "true")
func (b *Builder) DetailsClass(class string, attrPairs ...string) Element {
	return New(b.s, "details", append([]string{"class", class}, attrPairs...)...)
}

// DialogClass creates a dialog box with a class
//
// Example:
//
//	b.DialogClass("modal", "id", "confirmation-dialog", "open", "true")
func (b *Builder) DialogClass(class string, attrPairs ...string) Element {
	return New(b.s, "dialog", append([]string{"class", class}, attrPairs...)...)
}

// FieldsetClass creates a container for related form elements with a class
//
// Example:
//
//	b.FieldsetClass("form-group", "name", "personal-info")
func (b *Builder) FieldsetClass(class string, attrPairs ...string) Element {
	return New(b.s, "fieldset", append([]string{"class", class}, attrPairs...)...)
}

// FigCaptionClass creates a caption for a figure element with a class
//
// Example:
//
//	b.FigCaptionClass("image-caption")
func (b *Builder) FigCaptionClass(class string, attrPairs ...string) Element {
	return New(b.s, "figcaption", append([]string{"class", class}, attrPairs...)...)
}

// FigureClass creates a container for self-contained content with a class
//
// Example:
//
//	b.FigureClass("image-container", "id", "main-diagram")
func (b *Builder) FigureClass(class string, attrPairs ...string) Element {
	return New(b.s, "figure", append([]string{"class", class}, attrPairs...)...)
}
