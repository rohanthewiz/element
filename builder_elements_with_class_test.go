package element

import (
	"regexp"
	"strings"
	"testing"
)

// testElementWithClassCase defines a test case for elements with class methods
type testElementWithClassCase struct {
	name           string
	elementTag     string
	className      string
	additionalArgs []string
	expectedAttrs  map[string]string
}

func TestElementsWithClass(t *testing.T) {
	// Define test cases for elements with class methods
	testCases := []testElementWithClassCase{
		// A tag tests
		{
			name:       "AClass with only class",
			elementTag: "a",
			className:  "btn",
		},
		{
			name:           "AClass with class and href",
			elementTag:     "a",
			className:      "btn-primary",
			additionalArgs: []string{"href", "https://example.com"},
			expectedAttrs: map[string]string{
				"href": "https://example.com",
			},
		},
		// B tag tests
		{
			name:       "BClass with only class",
			elementTag: "b",
			className:  "highlight",
		},
		// I tag tests
		{
			name:       "IClass with only class",
			elementTag: "i",
			className:  "icon",
		},
		// U tag tests
		{
			name:       "UClass with only class",
			elementTag: "u",
			className:  "underline",
		},
		// Form tag tests
		{
			name:       "FormClass with only class",
			elementTag: "form",
			className:  "contact-form",
		},
		{
			name:           "FormClass with class and method",
			elementTag:     "form",
			className:      "login-form",
			additionalArgs: []string{"method", "post"},
			expectedAttrs: map[string]string{
				"method": "post",
			},
		},
		// Input tag tests
		{
			name:       "InputClass with only class",
			elementTag: "input",
			className:  "form-control",
		},
		// Select tag tests
		{
			name:       "SelectClass with only class",
			elementTag: "select",
			className:  "dropdown",
		},
		// Option tag tests
		{
			name:       "OptionClass with only class",
			elementTag: "option",
			className:  "option-item",
		},
		// Dd tag tests
		{
			name:       "DdClass with only class",
			elementTag: "dd",
			className:  "definition",
		},
		// Dt tag tests
		{
			name:       "DtClass with only class",
			elementTag: "dt",
			className:  "term",
		},
		// Div tag tests
		{
			name:       "DivClass with only class",
			elementTag: "div",
			className:  "container",
		},
		{
			name:           "DivClass with class and id",
			elementTag:     "div",
			className:      "wrapper",
			additionalArgs: []string{"id", "main-content"},
			expectedAttrs: map[string]string{
				"id": "main-content",
			},
		},
		// Body tag tests
		{
			name:       "BodyClass with only class",
			elementTag: "body",
			className:  "dark-theme",
		},
		// P tag tests
		{
			name:       "PClass with only class",
			elementTag: "p",
			className:  "paragraph",
		},
		// Span tag tests
		{
			name:       "SpanClass with only class",
			elementTag: "span",
			className:  "badge",
		},
		// Table tag tests
		{
			name:       "TableClass with only class",
			elementTag: "table",
			className:  "data-table",
		},
		// THead tag tests
		{
			name:       "THeadClass with only class",
			elementTag: "thead",
			className:  "table-head",
		},
		// TBody tag tests
		{
			name:       "TBodyClass with only class",
			elementTag: "tbody",
			className:  "table-body",
		},
		// Tr tag tests
		{
			name:       "TrClass with only class",
			elementTag: "tr",
			className:  "table-row",
		},
		// Th tag tests
		{
			name:       "ThClass with only class",
			elementTag: "th",
			className:  "header-cell",
		},
		// Td tag tests
		{
			name:       "TdClass with only class",
			elementTag: "td",
			className:  "data-cell",
		},
		// Section tag tests
		{
			name:       "SectionClass with only class",
			elementTag: "section",
			className:  "content-section",
		},
		// H1 tag tests
		{
			name:       "H1Class with only class",
			elementTag: "h1",
			className:  "main-heading",
		},
		// H2 tag tests
		{
			name:       "H2Class with only class",
			elementTag: "h2",
			className:  "section-heading",
		},
		// H3 tag tests
		{
			name:       "H3Class with only class",
			elementTag: "h3",
			className:  "subsection-heading",
		},
		// H4 tag tests
		{
			name:       "H4Class with only class",
			elementTag: "h4",
			className:  "minor-heading",
		},
		// Hr tag tests
		{
			name:       "HrClass with only class",
			elementTag: "hr",
			className:  "divider",
		},
		// Ol tag tests
		{
			name:       "OlClass with only class",
			elementTag: "ol",
			className:  "ordered-list",
		},
		// Ul tag tests
		{
			name:       "UlClass with only class",
			elementTag: "ul",
			className:  "unordered-list",
		},
		// Li tag tests
		{
			name:       "LiClass with only class",
			elementTag: "li",
			className:  "list-item",
		},
		// Img tag tests
		{
			name:       "ImgClass with only class",
			elementTag: "img",
			className:  "responsive-image",
		},
		{
			name:           "ImgClass with class and src",
			elementTag:     "img",
			className:      "logo",
			additionalArgs: []string{"src", "/images/logo.png"},
			expectedAttrs: map[string]string{
				"src": "/images/logo.png",
			},
		},
		// Article tag tests
		{
			name:       "ArticleClass with only class",
			elementTag: "article",
			className:  "blog-post",
		},
		// Aside tag tests
		{
			name:       "AsideClass with only class",
			elementTag: "aside",
			className:  "sidebar",
		},
		// Audio tag tests
		{
			name:       "AudioClass with only class",
			elementTag: "audio",
			className:  "sound-player",
		},
		// BlockQuote tag tests
		{
			name:       "BlockQuoteClass with only class",
			elementTag: "blockquote",
			className:  "quote",
		},
		// Button tag tests
		{
			name:       "ButtonClass with only class",
			elementTag: "button",
			className:  "btn-primary",
		},
		{
			name:           "ButtonClass with class and type",
			elementTag:     "button",
			className:      "btn-submit",
			additionalArgs: []string{"type", "submit"},
			expectedAttrs: map[string]string{
				"type": "submit",
			},
		},
		// Canvas tag tests
		{
			name:       "CanvasClass with only class",
			elementTag: "canvas",
			className:  "drawing-area",
		},
		// Code tag tests
		{
			name:       "CodeClass with only class",
			elementTag: "code",
			className:  "code-snippet",
		},
		// Dl tag tests
		{
			name:       "DlClass with only class",
			elementTag: "dl",
			className:  "description-list",
		},
		// Em tag tests
		{
			name:       "EmClass with only class",
			elementTag: "em",
			className:  "emphasis",
		},
		// Footer tag tests
		{
			name:       "FooterClass with only class",
			elementTag: "footer",
			className:  "page-footer",
		},
		// Header tag tests
		{
			name:       "HeaderClass with only class",
			elementTag: "header",
			className:  "page-header",
		},
		// Label tag tests
		{
			name:       "LabelClass with only class",
			elementTag: "label",
			className:  "form-label",
		},
		// Main tag tests
		{
			name:       "MainClass with only class",
			elementTag: "main",
			className:  "main-content",
		},
		// Nav tag tests
		{
			name:       "NavClass with only class",
			elementTag: "nav",
			className:  "navigation",
		},
		// Picture tag tests
		{
			name:       "PictureClass with only class",
			elementTag: "picture",
			className:  "image-container",
		},
		// Pre tag tests
		{
			name:       "PreClass with only class",
			elementTag: "pre",
			className:  "code-block",
		},
		// Time tag tests
		{
			name:       "TimeClass with only class",
			elementTag: "time",
			className:  "timestamp",
		},
		// Video tag tests
		{
			name:       "VideoClass with only class",
			elementTag: "video",
			className:  "video-player",
		},
		// Track tag tests
		{
			name:       "TrackClass with only class",
			elementTag: "track",
			className:  "subtitle-track",
		},
		// Abbr tag tests
		{
			name:       "AbbrClass with only class",
			elementTag: "abbr",
			className:  "abbreviation",
		},
		// Caption tag tests
		{
			name:       "CaptionClass with only class",
			elementTag: "caption",
			className:  "table-caption",
		},
		// FieldSet tag tests
		{
			name:       "FieldSetClass with only class",
			elementTag: "fieldset",
			className:  "form-group",
		},
		// Legend tag tests
		{
			name:       "LegendClass with only class",
			elementTag: "legend",
			className:  "form-legend",
		},
		// Progress tag tests
		{
			name:       "ProgressClass with only class",
			elementTag: "progress",
			className:  "progress-bar",
		},
		// Q tag tests
		{
			name:       "QClass with only class",
			elementTag: "q",
			className:  "inline-quote",
		},
		// Ruby tag tests
		{
			name:       "RubyClass with only class",
			elementTag: "ruby",
			className:  "ruby-text",
		},
		// Rt tag tests
		{
			name:       "RtClass with only class",
			elementTag: "rt",
			className:  "ruby-annotation",
		},
		// Noscript tag tests
		{
			name:       "NoscriptClass with only class",
			elementTag: "noscript",
			className:  "no-js-content",
		},
		// Small tag tests
		{
			name:       "SmallClass with only class",
			elementTag: "small",
			className:  "fine-print",
		},
		// Strong tag tests
		{
			name:       "StrongClass with only class",
			elementTag: "strong",
			className:  "important",
		},
		// Sub tag tests
		{
			name:       "SubClass with only class",
			elementTag: "sub",
			className:  "subscript",
		},
		// Summary tag tests
		{
			name:       "SummaryClass with only class",
			elementTag: "summary",
			className:  "details-summary",
		},
		// Sup tag tests
		{
			name:       "SupClass with only class",
			elementTag: "sup",
			className:  "superscript",
		},
		// TFoot tag tests
		{
			name:       "TFootClass with only class",
			elementTag: "tfoot",
			className:  "table-footer",
		},
		// Svg tag tests
		{
			name:       "SvgClass with only class",
			elementTag: "svg",
			className:  "vector-graphic",
		},
		// TextArea tag tests
		{
			name:       "TextAreaClass with only class",
			elementTag: "textarea",
			className:  "text-input",
		},
		// DataList tag tests
		{
			name:       "DataListClass with only class",
			elementTag: "datalist",
			className:  "data-options",
		},
		// Details tag tests
		{
			name:       "DetailsClass with only class",
			elementTag: "details",
			className:  "expandable",
		},
		// Dialog tag tests
		{
			name:       "DialogClass with only class",
			elementTag: "dialog",
			className:  "modal",
		},
		// Fieldset tag tests (alternative to FieldSet)
		{
			name:       "FieldsetClass with only class",
			elementTag: "fieldset",
			className:  "input-group",
		},
		// FigCaption tag tests
		{
			name:       "FigCaptionClass with only class",
			elementTag: "figcaption",
			className:  "image-caption",
		},
		// Figure tag tests
		{
			name:       "FigureClass with only class",
			elementTag: "figure",
			className:  "figure-container",
		},
	}

	// Run tests for each test case
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			runElementWithClassTest(t, tc)
		})
	}
}

// runElementWithClassTest executes a single element with class test case
func runElementWithClassTest(t *testing.T, tc testElementWithClassCase) {
	b := NewBuilder()

	// Get the element based on the tag
	var el Element
	switch tc.elementTag {
	case "a":
		if len(tc.additionalArgs) > 0 {
			el = b.AClass(tc.className, tc.additionalArgs...)
		} else {
			el = b.AClass(tc.className)
		}
	case "b":
		if len(tc.additionalArgs) > 0 {
			el = b.BClass(tc.className, tc.additionalArgs...)
		} else {
			el = b.BClass(tc.className)
		}
	case "i":
		if len(tc.additionalArgs) > 0 {
			el = b.IClass(tc.className, tc.additionalArgs...)
		} else {
			el = b.IClass(tc.className)
		}
	case "u":
		if len(tc.additionalArgs) > 0 {
			el = b.UClass(tc.className, tc.additionalArgs...)
		} else {
			el = b.UClass(tc.className)
		}
	case "form":
		if len(tc.additionalArgs) > 0 {
			el = b.FormClass(tc.className, tc.additionalArgs...)
		} else {
			el = b.FormClass(tc.className)
		}
	case "input":
		if len(tc.additionalArgs) > 0 {
			el = b.InputClass(tc.className, tc.additionalArgs...)
		} else {
			el = b.InputClass(tc.className)
		}
	case "select":
		if len(tc.additionalArgs) > 0 {
			el = b.SelectClass(tc.className, tc.additionalArgs...)
		} else {
			el = b.SelectClass(tc.className)
		}
	case "option":
		if len(tc.additionalArgs) > 0 {
			el = b.OptionClass(tc.className, tc.additionalArgs...)
		} else {
			el = b.OptionClass(tc.className)
		}
	case "dd":
		if len(tc.additionalArgs) > 0 {
			el = b.DdClass(tc.className, tc.additionalArgs...)
		} else {
			el = b.DdClass(tc.className)
		}
	case "dt":
		if len(tc.additionalArgs) > 0 {
			el = b.DtClass(tc.className, tc.additionalArgs...)
		} else {
			el = b.DtClass(tc.className)
		}
	case "div":
		if len(tc.additionalArgs) > 0 {
			el = b.DivClass(tc.className, tc.additionalArgs...)
		} else {
			el = b.DivClass(tc.className)
		}
	case "body":
		if len(tc.additionalArgs) > 0 {
			el = b.BodyClass(tc.className, tc.additionalArgs...)
		} else {
			el = b.BodyClass(tc.className)
		}
	case "p":
		if len(tc.additionalArgs) > 0 {
			el = b.PClass(tc.className, tc.additionalArgs...)
		} else {
			el = b.PClass(tc.className)
		}
	case "span":
		if len(tc.additionalArgs) > 0 {
			el = b.SpanClass(tc.className, tc.additionalArgs...)
		} else {
			el = b.SpanClass(tc.className)
		}
	case "table":
		if len(tc.additionalArgs) > 0 {
			el = b.TableClass(tc.className, tc.additionalArgs...)
		} else {
			el = b.TableClass(tc.className)
		}
	case "thead":
		if len(tc.additionalArgs) > 0 {
			el = b.THeadClass(tc.className, tc.additionalArgs...)
		} else {
			el = b.THeadClass(tc.className)
		}
	case "tbody":
		if len(tc.additionalArgs) > 0 {
			el = b.TBodyClass(tc.className, tc.additionalArgs...)
		} else {
			el = b.TBodyClass(tc.className)
		}
	case "tr":
		if len(tc.additionalArgs) > 0 {
			el = b.TrClass(tc.className, tc.additionalArgs...)
		} else {
			el = b.TrClass(tc.className)
		}
	case "th":
		if len(tc.additionalArgs) > 0 {
			el = b.ThClass(tc.className, tc.additionalArgs...)
		} else {
			el = b.ThClass(tc.className)
		}
	case "td":
		if len(tc.additionalArgs) > 0 {
			el = b.TdClass(tc.className, tc.additionalArgs...)
		} else {
			el = b.TdClass(tc.className)
		}
	case "section":
		if len(tc.additionalArgs) > 0 {
			el = b.SectionClass(tc.className, tc.additionalArgs...)
		} else {
			el = b.SectionClass(tc.className)
		}
	case "h1":
		if len(tc.additionalArgs) > 0 {
			el = b.H1Class(tc.className, tc.additionalArgs...)
		} else {
			el = b.H1Class(tc.className)
		}
	case "h2":
		if len(tc.additionalArgs) > 0 {
			el = b.H2Class(tc.className, tc.additionalArgs...)
		} else {
			el = b.H2Class(tc.className)
		}
	case "h3":
		if len(tc.additionalArgs) > 0 {
			el = b.H3Class(tc.className, tc.additionalArgs...)
		} else {
			el = b.H3Class(tc.className)
		}
	case "h4":
		if len(tc.additionalArgs) > 0 {
			el = b.H4Class(tc.className, tc.additionalArgs...)
		} else {
			el = b.H4Class(tc.className)
		}
	case "hr":
		if len(tc.additionalArgs) > 0 {
			el = b.HrClass(tc.className, tc.additionalArgs...)
		} else {
			el = b.HrClass(tc.className)
		}
	case "ol":
		if len(tc.additionalArgs) > 0 {
			el = b.OlClass(tc.className, tc.additionalArgs...)
		} else {
			el = b.OlClass(tc.className)
		}
	case "ul":
		if len(tc.additionalArgs) > 0 {
			el = b.UlClass(tc.className, tc.additionalArgs...)
		} else {
			el = b.UlClass(tc.className)
		}
	case "li":
		if len(tc.additionalArgs) > 0 {
			el = b.LiClass(tc.className, tc.additionalArgs...)
		} else {
			el = b.LiClass(tc.className)
		}
	case "img":
		if len(tc.additionalArgs) > 0 {
			el = b.ImgClass(tc.className, tc.additionalArgs...)
		} else {
			el = b.ImgClass(tc.className)
		}
	case "article":
		if len(tc.additionalArgs) > 0 {
			el = b.ArticleClass(tc.className, tc.additionalArgs...)
		} else {
			el = b.ArticleClass(tc.className)
		}
	case "aside":
		if len(tc.additionalArgs) > 0 {
			el = b.AsideClass(tc.className, tc.additionalArgs...)
		} else {
			el = b.AsideClass(tc.className)
		}
	case "audio":
		if len(tc.additionalArgs) > 0 {
			el = b.AudioClass(tc.className, tc.additionalArgs...)
		} else {
			el = b.AudioClass(tc.className)
		}
	case "blockquote":
		if len(tc.additionalArgs) > 0 {
			el = b.BlockQuoteClass(tc.className, tc.additionalArgs...)
		} else {
			el = b.BlockQuoteClass(tc.className)
		}
	case "button":
		if len(tc.additionalArgs) > 0 {
			el = b.ButtonClass(tc.className, tc.additionalArgs...)
		} else {
			el = b.ButtonClass(tc.className)
		}
	case "canvas":
		if len(tc.additionalArgs) > 0 {
			el = b.CanvasClass(tc.className, tc.additionalArgs...)
		} else {
			el = b.CanvasClass(tc.className)
		}
	case "code":
		if len(tc.additionalArgs) > 0 {
			el = b.CodeClass(tc.className, tc.additionalArgs...)
		} else {
			el = b.CodeClass(tc.className)
		}
	case "dl":
		if len(tc.additionalArgs) > 0 {
			el = b.DlClass(tc.className, tc.additionalArgs...)
		} else {
			el = b.DlClass(tc.className)
		}
	case "em":
		if len(tc.additionalArgs) > 0 {
			el = b.EmClass(tc.className, tc.additionalArgs...)
		} else {
			el = b.EmClass(tc.className)
		}
	case "footer":
		if len(tc.additionalArgs) > 0 {
			el = b.FooterClass(tc.className, tc.additionalArgs...)
		} else {
			el = b.FooterClass(tc.className)
		}
	case "header":
		if len(tc.additionalArgs) > 0 {
			el = b.HeaderClass(tc.className, tc.additionalArgs...)
		} else {
			el = b.HeaderClass(tc.className)
		}
	case "label":
		if len(tc.additionalArgs) > 0 {
			el = b.LabelClass(tc.className, tc.additionalArgs...)
		} else {
			el = b.LabelClass(tc.className)
		}
	case "main":
		if len(tc.additionalArgs) > 0 {
			el = b.MainClass(tc.className, tc.additionalArgs...)
		} else {
			el = b.MainClass(tc.className)
		}
	case "nav":
		if len(tc.additionalArgs) > 0 {
			el = b.NavClass(tc.className, tc.additionalArgs...)
		} else {
			el = b.NavClass(tc.className)
		}
	case "picture":
		if len(tc.additionalArgs) > 0 {
			el = b.PictureClass(tc.className, tc.additionalArgs...)
		} else {
			el = b.PictureClass(tc.className)
		}
	case "pre":
		if len(tc.additionalArgs) > 0 {
			el = b.PreClass(tc.className, tc.additionalArgs...)
		} else {
			el = b.PreClass(tc.className)
		}
	case "time":
		if len(tc.additionalArgs) > 0 {
			el = b.TimeClass(tc.className, tc.additionalArgs...)
		} else {
			el = b.TimeClass(tc.className)
		}
	case "video":
		if len(tc.additionalArgs) > 0 {
			el = b.VideoClass(tc.className, tc.additionalArgs...)
		} else {
			el = b.VideoClass(tc.className)
		}
	case "track":
		if len(tc.additionalArgs) > 0 {
			el = b.TrackClass(tc.className, tc.additionalArgs...)
		} else {
			el = b.TrackClass(tc.className)
		}
	case "abbr":
		if len(tc.additionalArgs) > 0 {
			el = b.AbbrClass(tc.className, tc.additionalArgs...)
		} else {
			el = b.AbbrClass(tc.className)
		}
	case "caption":
		if len(tc.additionalArgs) > 0 {
			el = b.CaptionClass(tc.className, tc.additionalArgs...)
		} else {
			el = b.CaptionClass(tc.className)
		}
	case "fieldset":
		if len(tc.additionalArgs) > 0 {
			el = b.FieldSetClass(tc.className, tc.additionalArgs...)
		} else {
			el = b.FieldSetClass(tc.className)
		}
	case "legend":
		if len(tc.additionalArgs) > 0 {
			el = b.LegendClass(tc.className, tc.additionalArgs...)
		} else {
			el = b.LegendClass(tc.className)
		}
	case "progress":
		if len(tc.additionalArgs) > 0 {
			el = b.ProgressClass(tc.className, tc.additionalArgs...)
		} else {
			el = b.ProgressClass(tc.className)
		}
	case "q":
		if len(tc.additionalArgs) > 0 {
			el = b.QClass(tc.className, tc.additionalArgs...)
		} else {
			el = b.QClass(tc.className)
		}
	case "ruby":
		if len(tc.additionalArgs) > 0 {
			el = b.RubyClass(tc.className, tc.additionalArgs...)
		} else {
			el = b.RubyClass(tc.className)
		}
	case "rt":
		if len(tc.additionalArgs) > 0 {
			el = b.RtClass(tc.className, tc.additionalArgs...)
		} else {
			el = b.RtClass(tc.className)
		}
	case "noscript":
		if len(tc.additionalArgs) > 0 {
			el = b.NoscriptClass(tc.className, tc.additionalArgs...)
		} else {
			el = b.NoscriptClass(tc.className)
		}
	case "small":
		if len(tc.additionalArgs) > 0 {
			el = b.SmallClass(tc.className, tc.additionalArgs...)
		} else {
			el = b.SmallClass(tc.className)
		}
	case "strong":
		if len(tc.additionalArgs) > 0 {
			el = b.StrongClass(tc.className, tc.additionalArgs...)
		} else {
			el = b.StrongClass(tc.className)
		}
	case "sub":
		if len(tc.additionalArgs) > 0 {
			el = b.SubClass(tc.className, tc.additionalArgs...)
		} else {
			el = b.SubClass(tc.className)
		}
	case "summary":
		if len(tc.additionalArgs) > 0 {
			el = b.SummaryClass(tc.className, tc.additionalArgs...)
		} else {
			el = b.SummaryClass(tc.className)
		}
	case "sup":
		if len(tc.additionalArgs) > 0 {
			el = b.SupClass(tc.className, tc.additionalArgs...)
		} else {
			el = b.SupClass(tc.className)
		}
	case "tfoot":
		if len(tc.additionalArgs) > 0 {
			el = b.TFootClass(tc.className, tc.additionalArgs...)
		} else {
			el = b.TFootClass(tc.className)
		}
	case "svg":
		if len(tc.additionalArgs) > 0 {
			el = b.SvgClass(tc.className, tc.additionalArgs...)
		} else {
			el = b.SvgClass(tc.className)
		}
	case "textarea":
		if len(tc.additionalArgs) > 0 {
			el = b.TextAreaClass(tc.className, tc.additionalArgs...)
		} else {
			el = b.TextAreaClass(tc.className)
		}
	case "datalist":
		if len(tc.additionalArgs) > 0 {
			el = b.DataListClass(tc.className, tc.additionalArgs...)
		} else {
			el = b.DataListClass(tc.className)
		}
	case "details":
		if len(tc.additionalArgs) > 0 {
			el = b.DetailsClass(tc.className, tc.additionalArgs...)
		} else {
			el = b.DetailsClass(tc.className)
		}
	case "dialog":
		if len(tc.additionalArgs) > 0 {
			el = b.DialogClass(tc.className, tc.additionalArgs...)
		} else {
			el = b.DialogClass(tc.className)
		}
	case "figcaption":
		if len(tc.additionalArgs) > 0 {
			el = b.FigCaptionClass(tc.className, tc.additionalArgs...)
		} else {
			el = b.FigCaptionClass(tc.className)
		}
	case "figure":
		if len(tc.additionalArgs) > 0 {
			el = b.FigureClass(tc.className, tc.additionalArgs...)
		} else {
			el = b.FigureClass(tc.className)
		}
	default:
		t.Fatalf("Unsupported element tag: %s", tc.elementTag)
	}

	// Verify class attribute
	if !el.HasAttribute("class", tc.className) {
		t.Errorf("%s failed. Expected class attribute to be set to %q.", tc.name, tc.className)
	}

	// Verify additional attributes
	for key, expectedValue := range tc.expectedAttrs {
		if !el.HasAttribute(key, expectedValue) {
			t.Errorf("%s failed. Expected %s attribute to be set to %q.", tc.name, key, expectedValue)
		}
	}

	// Render the element
	el.R()
	result := b.String()

	// Create and verify pattern match
	pattern := "^<" + tc.elementTag + "[^>]+class=\"" + tc.className + "\"[^>]*>"
	if !el.IsSingleTag() {
		pattern += "</" + tc.elementTag + ">$"
	} else {
		pattern += "$"
	}

	matched, err := regexp.MatchString(pattern, result)
	if err != nil || !matched {
		t.Errorf("%s failed regex match.\nExpected pattern: %s\nGot: %s", tc.name, pattern, result)
	}

	// Additional attribute checks for rendered output
	for key, value := range tc.expectedAttrs {
		attrPattern := ` ` + key + `="` + value + `"`
		if !strings.Contains(result, attrPattern) {
			t.Errorf("%s rendered output missing expected attribute %s=%q. Got: %s", tc.name, key, value, result)
		}
	}
}
