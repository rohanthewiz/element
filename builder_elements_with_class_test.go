package element

import (
	"regexp"
	"strings"
	"testing"
)

func TestBuilderAClass(t *testing.T) {
	b := NewBuilder()

	// Test AClass with only class
	el := b.AClass("btn")
	if !el.HasAttribute("class", "btn") {
		t.Errorf("AClass with only class failed. Expected class attribute to be set to \"btn\".")
	}
	el.R()
	result := b.String()

	// Test regex pattern matching for anchor tag with class
	anchorPattern := `<a[^>]+class="btn"[^>]*></a>`
	matched, err := regexp.MatchString(anchorPattern, result)
	if err != nil || !matched {
		t.Errorf("AClass with only class failed regex match.\nExpected pattern: %s\nGot: %s", anchorPattern, result)
	}

	// Reset the builder for next test
	b.Reset()

	// Test AClass with class and href attribute
	b.AClass("btn-primary", "href", "https://example.com").R()
	result = b.String()
	if !strings.Contains(result, `<a `) ||
		!strings.Contains(result, ` class="btn-primary"`) ||
		!strings.Contains(result, ` href="https://example.com"`) ||
		!strings.Contains(result, `</a>`) {
		t.Errorf("AClass with class and href failed.\nGot: %s", result)
	}
}
