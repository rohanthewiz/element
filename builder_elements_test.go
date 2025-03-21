package element

import (
	"strings"
	"testing"
)

func TestHtmlElements(t *testing.T) {
	tests := []struct {
		name     string
		function func(b *Builder) Element
		want     string
	}{
		{name: "Html", function: func(b *Builder) Element { return b.Html() }, want: "<!DOCTYPE html><html></html>"},
		{name: "Head", function: func(b *Builder) Element { return b.Head() }, want: "<head></head>"},
		{name: "Title", function: func(b *Builder) Element { return b.Title() }, want: "<title></title>"},
		{name: "Meta", function: func(b *Builder) Element { return b.Meta() }, want: "<meta>"},
		{name: "Style", function: func(b *Builder) Element { return b.Style() }, want: "<style></style>"},
		{name: "Link", function: func(b *Builder) Element { return b.Link() }, want: "<link>"},
		{name: "A", function: func(b *Builder) Element { return b.A() }, want: "<a></a>"},
		{name: "B", function: func(b *Builder) Element { return b.B() }, want: "<b></b>"},
		{name: "I", function: func(b *Builder) Element { return b.I() }, want: "<i></i>"},
		{name: "U", function: func(b *Builder) Element { return b.U() }, want: "<u></u>"},
		{name: "Form", function: func(b *Builder) Element { return b.Form() }, want: "<form></form>"},
		{name: "Input", function: func(b *Builder) Element { return b.Input() }, want: "<input>"},
		{name: "Select", function: func(b *Builder) Element { return b.Select() }, want: "<select></select>"},
		{name: "Option", function: func(b *Builder) Element { return b.Option() }, want: "<option></option>"},
		{name: "Dd", function: func(b *Builder) Element { return b.Dd() }, want: "<dd></dd>"},
		{name: "Dt", function: func(b *Builder) Element { return b.Dt() }, want: "<dt></dt>"},
		{name: "Div", function: func(b *Builder) Element { return b.Div() }, want: "<div></div>"},
		{name: "Body", function: func(b *Builder) Element { return b.Body() }, want: "<body></body>"},
		{name: "P", function: func(b *Builder) Element { return b.P() }, want: "<p></p>"},
		{name: "Span", function: func(b *Builder) Element { return b.Span() }, want: "<span></span>"},
		{name: "Table", function: func(b *Builder) Element { return b.Table() }, want: "<table></table>"},
		{name: "THead", function: func(b *Builder) Element { return b.THead() }, want: "<thead></thead>"},
		{name: "TBody", function: func(b *Builder) Element { return b.TBody() }, want: "<tbody></tbody>"},
		{name: "Tr", function: func(b *Builder) Element { return b.Tr() }, want: "<tr></tr>"},
		{name: "Th", function: func(b *Builder) Element { return b.Th() }, want: "<th></th>"},
		{name: "Td", function: func(b *Builder) Element { return b.Td() }, want: "<td></td>"},
		{name: "Section", function: func(b *Builder) Element { return b.Section() }, want: "<section></section>"},
		{name: "H1", function: func(b *Builder) Element { return b.H1() }, want: "<h1></h1>"},
		{name: "H2", function: func(b *Builder) Element { return b.H2() }, want: "<h2></h2>"},
		{name: "H3", function: func(b *Builder) Element { return b.H3() }, want: "<h3></h3>"},
		{name: "H4", function: func(b *Builder) Element { return b.H4() }, want: "<h4></h4>"},
		{name: "Hr", function: func(b *Builder) Element { return b.Hr() }, want: "<hr>"},
		{name: "Ol", function: func(b *Builder) Element { return b.Ol() }, want: "<ol></ol>"},
		{name: "Ul", function: func(b *Builder) Element { return b.Ul() }, want: "<ul></ul>"},
		{name: "Li", function: func(b *Builder) Element { return b.Li() }, want: "<li></li>"},
		{name: "Img", function: func(b *Builder) Element { return b.Img() }, want: "<img>"},
		{name: "Article", function: func(b *Builder) Element { return b.Article() }, want: "<article></article>"},
		{name: "Aside", function: func(b *Builder) Element { return b.Aside() }, want: "<aside></aside>"},
		{name: "Audio", function: func(b *Builder) Element { return b.Audio() }, want: "<audio></audio>"},
		{name: "BlockQuote", function: func(b *Builder) Element { return b.BlockQuote() }, want: "<blockquote></blockquote>"},
		{name: "Button", function: func(b *Builder) Element { return b.Button() }, want: "<button></button>"},
		{name: "Canvas", function: func(b *Builder) Element { return b.Canvas() }, want: "<canvas></canvas>"},
		{name: "Code", function: func(b *Builder) Element { return b.Code() }, want: "<code></code>"},
		{name: "DataList", function: func(b *Builder) Element { return b.DataList() }, want: "<datalist></datalist>"},
		{name: "Details", function: func(b *Builder) Element { return b.Details() }, want: "<details></details>"},
		{name: "Dialog", function: func(b *Builder) Element { return b.Dialog() }, want: "<dialog></dialog>"},
		{name: "Em", function: func(b *Builder) Element { return b.Em() }, want: "<em></em>"},
		{name: "Fieldset", function: func(b *Builder) Element { return b.Fieldset() }, want: "<fieldset></fieldset>"},
		{name: "FigCaption", function: func(b *Builder) Element { return b.FigCaption() }, want: "<figcaption></figcaption>"},
		{name: "Figure", function: func(b *Builder) Element { return b.Figure() }, want: "<figure></figure>"},
		{name: "Footer", function: func(b *Builder) Element { return b.Footer() }, want: "<footer></footer>"},
		{name: "Header", function: func(b *Builder) Element { return b.Header() }, want: "<header></header>"},
		{name: "Iframe", function: func(b *Builder) Element { return b.Iframe() }, want: "<iframe></iframe>"},
		{name: "Label", function: func(b *Builder) Element { return b.Label() }, want: "<label></label>"},
		{name: "Legend", function: func(b *Builder) Element { return b.Legend() }, want: "<legend></legend>"},
		{name: "Main", function: func(b *Builder) Element { return b.Main() }, want: "<main></main>"},
		{name: "Nav", function: func(b *Builder) Element { return b.Nav() }, want: "<nav></nav>"},
		{name: "Noscript", function: func(b *Builder) Element { return b.Noscript() }, want: "<noscript></noscript>"},
		{name: "Object", function: func(b *Builder) Element { return b.Object() }, want: "<object></object>"},
		{name: "Pre", function: func(b *Builder) Element { return b.Pre() }, want: "<pre></pre>"},
		{name: "Progress", function: func(b *Builder) Element { return b.Progress() }, want: "<progress></progress>"},
		{name: "Script", function: func(b *Builder) Element { return b.Script() }, want: "<script></script>"},
		{name: "Strong", function: func(b *Builder) Element { return b.Strong() }, want: "<strong></strong>"},
		{name: "Summary", function: func(b *Builder) Element { return b.Summary() }, want: "<summary></summary>"},
		{name: "Svg", function: func(b *Builder) Element { return b.Svg() }, want: "<svg></svg>"},
		{name: "Textarea", function: func(b *Builder) Element { return b.Textarea() }, want: "<textarea></textarea>"},
		{name: "Time", function: func(b *Builder) Element { return b.Time() }, want: "<time></time>"},
		{name: "Video", function: func(b *Builder) Element { return b.Video() }, want: "<video></video>"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := NewBuilder()
			tt.function(b).R()
			got := b.String()
			if got != tt.want {
				t.Errorf("%s() = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}

func TestHtmlElementsWithAttributes(t *testing.T) {
	tests := []struct {
		name       string
		function   func(b *Builder) Element
		attributes []string
		want       string
	}{
		{name: "Div with class", function: func(b *Builder) Element { return b.Div("class", "container") }, attributes: []string{"class", "container"}, want: `<div class="container"></div>`},
		{name: "A with href", function: func(b *Builder) Element { return b.A("href", "https://example.com") }, attributes: []string{"href", "https://example.com"}, want: `<a href="https://example.com"></a>`},
		{name: "Img with src", function: func(b *Builder) Element { return b.Img("src", "image.jpg", "alt", "An image") }, attributes: []string{"src", "image.jpg", "alt", "An image"}, want: `<img src="image.jpg" alt="An image">`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := NewBuilder()
			tt.function(b).R()
			got := b.String()
			if got != tt.want {
				t.Errorf("%s(%s) = %v, want %v", tt.name, strings.Join(tt.attributes, ", "), got, tt.want)
			}
		})
	}
}

func TestHtmlElementsWithContent(t *testing.T) {
	tests := []struct {
		name  string
		setup func(b *Builder)
		want  string
	}{
		{
			name: "P with text",
			setup: func(b *Builder) {
				b.P().R(b.Text("Hello"))
			},
			want: "<p>Hello</p>",
		},
		{
			name: "Div with nested elements",
			setup: func(b *Builder) {
				b.Div("class", "container").R(
					b.H1().R(b.Text("Title")),
					b.P().R(b.Text("Paragraph")),
				)
			},
			want: `<div class="container"><h1>Title</h1><p>Paragraph</p></div>`,
		},
		{
			name: "Table structure",
			setup: func(b *Builder) {
				b.Table().R(
					b.THead().R(
						b.Tr().R(
							b.Th().R(b.Text("Header")),
						),
					),
					b.TBody().R(
						b.Tr().R(
							b.Td().R(b.Text("Cell")),
						),
					),
				)
			},
			want: "<table><thead><tr><th>Header</th></tr></thead><tbody><tr><td>Cell</td></tr></tbody></table>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := NewBuilder()
			tt.setup(b)
			got := b.String()
			if got != tt.want {
				t.Errorf("%s = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}
