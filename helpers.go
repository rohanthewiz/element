package element

import (
	"encoding/base64"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Seed once, e.g., globally or pass the *rand.Rand instance
var seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))

// seededRandLock protects seededRand which is not safe for concurrent use
var seededRandLock sync.Mutex

// Satisfy the Stringer interface
func (el Element) String() string {
	return el.sb.String()
}

// Shortcut for String()
func (el Element) S() string {
	return el.sb.String()
}

// Is this a plaintext element
func (el Element) IsText() bool {
	return el.name == "t"
}

func (el Element) IsSingleTag() bool {
	if _, ok := singleTags[el.name]; ok {
		return true
	}
	return false
}

// lowerName lowercases an element name, avoiding an allocation
// when the name is already lowercase (the common case)
func lowerName(s string) string {
	for i := 0; i < len(s); i++ {
		if s[i] >= 'A' && s[i] <= 'Z' {
			return strings.ToLower(s)
		}
	}
	return s
}

// normalizeAttrPairs validates a list of key values into attribute pairs,
// preserving insertion order so rendered output is deterministic.
// Duplicate keys keep their first position with the last value winning.
// Even number of items required or the last item is dropped
func normalizeAttrPairs(e Element, items []string) []string {
	if len(items)%2 != 0 {
		issue := fmt.Sprintf(`Even number of arguments required for Element attributes.
Args: %q dropping %q`, items, items[len(items)-1])

		if IsDebugMode() {
			fmt.Printf("![%s] %s\n", e.id, issue)
			e.issues = append(e.issues, issue)
			concerns.UpsertConcern(concernOther, e)
		}

		items = items[:len(items)-1] // drop the last item
	}

	pairs := make([]string, 0, len(items)+2)
	for i := 0; i+1 < len(items); i += 2 {
		pairs = upsertPair(pairs, items[i], items[i+1])
	}

	if IsDebugMode() {
		pairs = upsertPair(pairs, "data-ele-id", e.id) // Add the data-ele-id attribute
	}
	return pairs
}

// upsertPair appends a key/value pair, or updates the value in place
// if the key is already present
func upsertPair(pairs []string, key, value string) []string {
	for i := 0; i+1 < len(pairs); i += 2 {
		if pairs[i] == key {
			pairs[i+1] = value
			return pairs
		}
	}
	return append(pairs, key, value)
}

/*// genRandString is a fast generator of a random string
// based on math.Rand (non-cryptographic)
func genRandString(length int) string {
	byts := make([]byte, length)
	seededRand.Read(byts) // Use the seeded generator
	return hex.EncodeToString(byts)
}
*/

// genRandomId is a fast generator of a base64 encoded random string
// based on math.Rand (non-cryptographic)
// Length specifies the number of bytes to generate, however output is trimmed of any "=" padding
func genRandomId(length int) string {
	byts := make([]byte, length)
	seededRandLock.Lock()
	seededRand.Read(byts) // Use the seeded generator
	seededRandLock.Unlock()
	b64 := base64.URLEncoding.EncodeToString(byts)
	return strings.TrimRight(b64, "=") // Remove any padding
}

// ToOrdinal converts a non-negative integer to its English ordinal string representation.
// Example: 1 -> "1st", 2 -> "2nd", 11 -> "11th", 21 -> "21st"
// It handles 0 as "0th". Behaviour for negative numbers is undefined by this function.
// Credit: Google's Gemini 2.5 pro
func ToOrdinal(n int) string {
	if n < 0 {
		// Or return an error, or handle negatives if needed.
		// For simplicity, we assume non-negative input based on typical usage.
		// Let's return the number with "th" for consistency, although unusual.
		// Alternatively, panic("ToOrdinal does not support negative numbers")
		return strconv.Itoa(n) + "th"
	}

	numStr := strconv.Itoa(n)

	// Check for 11, 12, 13 suffix ("th")
	lastTwoDigits := n % 100
	if lastTwoDigits >= 11 && lastTwoDigits <= 13 {
		return numStr + "th"
	}

	// Check for 1, 2, 3 suffix ("st", "nd", "rd")
	lastDigit := n % 10
	switch lastDigit {
	case 1:
		return numStr + "st"
	case 2:
		return numStr + "nd"
	case 3:
		return numStr + "rd"
	default:
		// Includes 0, 4, 5, 6, 7, 8, 9
		return numStr + "th"
	}
}

// inlineElements are HTML elements that should stay on the same line as their content
var inlineElements = map[string]bool{
	"a": true, "abbr": true, "b": true, "bdi": true, "bdo": true,
	"cite": true, "code": true, "data": true, "dfn": true, "em": true,
	"i": true, "kbd": true, "mark": true, "q": true, "s": true,
	"samp": true, "small": true, "span": true, "strong": true, "sub": true,
	"sup": true, "time": true, "u": true, "var": true,
}

// PrettyHTML formats HTML with proper indentation and line breaks
// It parses the HTML and adds newlines and indentation for readability
func PrettyHTML(html string) string {
	if html == "" {
		return ""
	}

	var result strings.Builder
	var depth int
	indent := "  " // 2 spaces per level
	lastWasText := false
	lastWasNewline := false // Track if we just wrote a newline

	i := 0
	for i < len(html) {
		// Handle DOCTYPE, comments, and other special declarations
		if strings.HasPrefix(html[i:], "<!") {
			end := strings.Index(html[i:], ">")
			if end != -1 {
				result.WriteString(html[i : i+end+1])
				result.WriteString("\n")
				i += end + 1
				lastWasText = false
				lastWasNewline = true // We just wrote a newline
				continue
			}
		}

		// Handle opening tags
		if html[i] == '<' && i+1 < len(html) && html[i+1] != '/' {
			// Find the end of the tag
			tagEnd := strings.Index(html[i:], ">")
			if tagEnd == -1 {
				result.WriteString(html[i:])
				break
			}

			// Extract tag name
			tagContent := html[i+1 : i+tagEnd]
			spaceIdx := strings.IndexAny(tagContent, " \t\n")
			tagName := tagContent
			if spaceIdx != -1 {
				tagName = tagContent[:spaceIdx]
			}

			isInline := inlineElements[strings.ToLower(tagName)]
			isSelfClosing := singleTags[strings.ToLower(tagName)]

			// Add newline and indentation for non-inline elements
			if !isInline && i > 0 && !lastWasText {
				// Only add newline if we didn't just write one
				if !lastWasNewline {
					result.WriteString("\n")
				}
				result.WriteString(strings.Repeat(indent, depth))
			}

			// Write the opening tag
			result.WriteString(html[i : i+tagEnd+1])

			// Increase depth for non-self-closing tags
			if !isSelfClosing && !isInline {
				depth++
			}

			i += tagEnd + 1
			lastWasText = false
			lastWasNewline = false
			continue
		}

		// Handle closing tags
		if html[i] == '<' && i+1 < len(html) && html[i+1] == '/' {
			// Find the end of the tag
			tagEnd := strings.Index(html[i:], ">")
			if tagEnd == -1 {
				result.WriteString(html[i:])
				break
			}

			// Extract tag name
			tagName := html[i+2 : i+tagEnd]

			isInline := inlineElements[strings.ToLower(tagName)]

			// Decrease depth for non-inline elements
			if !isInline {
				depth--
			}

			// Add indentation for non-inline elements (but not if last was text)
			if !isInline && !lastWasText {
				result.WriteString("\n")
				result.WriteString(strings.Repeat(indent, depth))
			}

			// Write the closing tag
			result.WriteString(html[i : i+tagEnd+1])

			i += tagEnd + 1
			lastWasText = false
			lastWasNewline = false
			continue
		}

		// Handle text content
		textStart := i
		for i < len(html) && html[i] != '<' {
			i++
		}

		text := html[textStart:i]
		// Skip text that is purely whitespace, but preserve actual text exactly as-is
		if strings.TrimSpace(text) != "" {
			// Write original text to preserve meaningful spaces
			result.WriteString(text)
			lastWasText = true
			lastWasNewline = false
		}
	}

	// Add final newline
	if result.Len() > 0 {
		result.WriteString("\n")
	}

	return result.String()
}
