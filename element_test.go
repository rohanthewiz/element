package element

import (
	"regexp"
	"strings"
	"testing"
)

// Warning: since Go will randomize maps, multiple attributes may come out in different order
// Maybe just to use a single attribute pair or test only the length of the output if multiple attributes are required

func TestRender(t *testing.T) {
	s := &strings.Builder{}

	New(s, "span").R()
	want := `<span.*></span>`
	got := s.String()

	if !regexp.MustCompile(want).MatchString(got) {
		t.Error("Failed to render an empty span", got)
	}

	// Span with inner text
	s.Reset()
	New(s, "span").R(
		Text(s, "This is some inner text"),
	)

	if !regexp.MustCompile(`<span.*>This is some inner text</span>`).MatchString(s.String()) {
		t.Error("Failed to render a span with inner text", "\nGot:", s.String())
	}

	// Span with multiple inner text
	s.Reset()
	New(s, "span", "id", "special", "class", "normal").R(
		Text(s, "This is some inner text", " and some more by the way"), // we can use a list of texts
	)

	got = s.String()
	// want = `<span id="special" class="normal">This is some inner text and some more by the way</span>`

	if !strings.Contains(got, "This is some inner text") ||
		!strings.Contains(got, "and some more by the way") {
		t.Error("Failed to render a span with multiple inner text. \nGot:", got)
	}

	// We can't properly test these because of the automatic el id data attribute
	/*	// Div with text and element children
		s.Reset()
		New(s, "div", "id", "container", "class", "active").R(
			Text(s, "This is some inner text", " and some more by the way"),
			New(s, "form", "method", "post").R(),
			Text(s, "Some ending text"),
		)
		got = s.String()
		want = `<div id="container" class="active">This is some inner text and some more by the way<form method="post"></form>Some ending text</div>`
		if len(got) != len(want) { // Go's map order is random, so have to rely on length match
			t.Error("Failed to render div with text and element children",
				"\nExpected:", want, "\nGot:", got)
		} else {
			fmt.Println("good ->", got)
		}

		// Deep nesting
		s.Reset()
		moreText := " - more text"
		New(s, "div", "id", "container", "class", "active").R(
			Text(s, "some text", moreText),
			New(s, "form", "method", "post").R(
				New(s, "input", "value", "some input").R(),
				New(s, "button").R(
					New(s, "span", "style", "background-color:wheat").R(Text(s, "My nice button")),
				),
			),
			Text(s, "Some ending text"),
		)
		got = s.String()
		want = `<div id="container" class="active">some text - more text<form method="post"><input value="some input"><button><span style="background-color:wheat">My nice button</span></button></form>Some ending text</div>`
		if len(got) != len(want) { // Go's map order is random, so have to rely on length match
			t.Error("Failed to render div with deep nesting",
				"\nExpected:", want, "\nGot:", got)
		} else {
			fmt.Println("good ->", got)
		}
	*/
}

func TestElement_F(t *testing.T) {
	tests := []struct {
		name     string
		format   string
		args     []any
		expected string
	}{
		{
			name:     "simple string",
			format:   "Hello",
			args:     nil,
			expected: "Hello",
		},
		{
			name:     "with formatting",
			format:   "Count: %d",
			args:     []any{42},
			expected: "Count: 42",
		},
		{
			name:     "multiple args",
			format:   "%s: %d",
			args:     []any{"Total", 100},
			expected: "Total: 100",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sb := &strings.Builder{}
			el := New(sb, "div")
			el.F(tt.format, tt.args...)

			got := sb.String()
			want := "<div.*>" + tt.expected + "</div>"
			if !regexp.MustCompile(want).MatchString(got) {
				t.Errorf("F() = %q, want %q", got, want)
			}
		})
	}
}

// We have deprecated For -- use builder.Wrap or element.ForEach instead
/*func TestFor(t *testing.T) {
	var testAnimals = []string{"cat", "mouse", "dog"}

	s := &strings.Builder{}

	// Div with text and element children
	New(s, "div", "id", "container", "class", "active").R(
		New(s, "ul", "class", "list").For(testAnimals, "li", "class", "animal"), // build a list
	)
	str := s.String()
	expected := `<div id="container" class="active"><ul class="list"><li class="animal">cat</li><li class="animal">mouse</li><li class="animal">dog</li></ul></div>`
	if len(str) != len(expected) { // Go's map order is random, so have to rely on length match
		t.Error("Failed to render a html list with For",
			"\nExpected:", expected, "\nGot:", str)
	} else {
		fmt.Println("good ->", str)
	}
}
*/

func TestDebugShowWithCopyIcon(t *testing.T) {
	// Enable debug mode
	DebugSet()
	defer DebugClear()

	// Create an element with an issue
	b := NewBuilder()
	b.Div().R(
		// This will create an issue - text not wrapped
		"unwrapped text",
	)

	// Get debug output
	output := DebugShow()

	// Check that output contains the copy icon SVG
	if !strings.Contains(output, `class="copy-icon"`) {
		t.Error("Debug output should contain copy icon class")
	}

	if !strings.Contains(output, `onclick="copyToClipboard(`) {
		t.Error("Debug output should contain copyToClipboard function call")
	}

	if !strings.Contains(output, `<svg`) {
		t.Error("Debug output should contain SVG element")
	}

	// Check that JavaScript is included
	if !strings.Contains(output, `function copyToClipboard(text)`) {
		t.Error("Debug output should contain copyToClipboard JavaScript function")
	}

	if !strings.Contains(output, `function showNotification(message)`) {
		t.Error("Debug output should contain showNotification JavaScript function")
	}

	// Check that CSS for copy icon is included
	if !strings.Contains(output, `.copy-icon {`) {
		t.Error("Debug output should contain CSS for copy icon")
	}

	if !strings.Contains(output, `.notification {`) {
		t.Error("Debug output should contain CSS for notification")
	}
}

func TestDebugClearIssues(t *testing.T) {
	// Enable debug mode
	DebugSet()
	defer DebugClear()

	// Create an element with an issue
	b := NewBuilder()
	b.Div().R(
		// This will create an issue - text not wrapped
		"unwrapped text",
	)

	// Verify issue exists
	output := DebugShow()
	if !strings.Contains(output, "unwrapped text") {
		t.Error("Debug output should contain the issue")
	}

	// Clear issues only
	DebugClearIssues()

	// Verify debug mode is still active
	if !IsDebugMode() {
		t.Error("Debug mode should still be active after clearing issues")
	}

	// Verify issues are cleared
	output = DebugShow()
	if strings.Contains(output, "unwrapped text") {
		t.Error("Issues should be cleared")
	}
	if !strings.Contains(output, "No element concerns found") {
		t.Error("Should show no concerns after clearing")
	}

	// Create new issue after clearing
	b.Div().R(
		"another unwrapped text",
	)

	// Verify new issue is tracked
	output = DebugShow()
	if !strings.Contains(output, "another unwrapped text") {
		t.Error("New issues should be tracked after clearing")
	}
}

func TestDebugShowClearButton(t *testing.T) {
	// Enable debug mode
	DebugSet()
	defer DebugClear()

	// Create an element with an issue
	b := NewBuilder()
	b.Div().R(
		"unwrapped text",
	)

	// Get debug output
	output := DebugShow()

	// Check that output contains the clear button
	if !strings.Contains(output, `class="clear-button"`) {
		t.Error("Debug output should contain clear button")
	}

	if !strings.Contains(output, `onclick="clearIssues()"`) {
		t.Error("Clear button should have onclick handler")
	}

	if !strings.Contains(output, "Clear Issues") {
		t.Error("Clear button should have 'Clear Issues' text")
	}

	// Check that clearIssues JavaScript function is included
	if !strings.Contains(output, `function clearIssues()`) {
		t.Error("Debug output should contain clearIssues JavaScript function")
	}

	if !strings.Contains(output, `/debug/clear-issues`) {
		t.Error("JavaScript should fetch /debug/clear-issues endpoint")
	}

	// Check CSS for clear button
	if !strings.Contains(output, `.clear-button {`) {
		t.Error("Debug output should contain CSS for clear button")
	}
}
