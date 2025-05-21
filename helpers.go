package element

import (
	"encoding/base64"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

// Seed once, e.g., globally or pass the *rand.Rand instance
var seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))

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

// stringlistToMap generates a map from a list of key values
// Even number of items required or results may be fatal
func stringlistToMap(e Element, items ...string) map[string]string {
	m := map[string]string{}

	if len(items)%2 != 0 {
		issue := fmt.Sprintf(`Even number of arguments required for Element attributes.
Args: %q dropping %q`, items, items[len(items)-1])
		// fmt.Println(issue)

		if debugMode {
			e.issues = append(e.issues, issue)
			concerns.UpsertConcern(concernOther, e)
		}

		items = items[:len(items)-1] // drop the last item
	}

	if debugMode {
		items = append(items, "data-ele-id", e.id) // Add the data-ele-id attribute
	}

	key := ""
	for i, item := range items {
		if i%2 == 0 {
			key = item
		} else {
			m[key] = item
		}
	}
	return m
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
	seededRand.Read(byts) // Use the seeded generator
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
