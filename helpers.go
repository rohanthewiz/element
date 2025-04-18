package element

import (
	"fmt"

	"github.com/rohanthewiz/serr"
)

// Satisfy the Stringer interface
func (e Element) String() string {
	return e.sb.String()
}

// Shortcut for String()
func (e Element) S() string {
	return e.sb.String()
}

// Is this a plaintext element
func (e Element) IsText() bool {
	return e.El == "t"
}

func (e Element) IsSingleTag() bool {
	if _, ok := singleTags[e.El]; ok {
		return true
	}
	return false
}

// stringlistToMap generates a map from a list of key values
// Even number of items required or results may be fatal
func stringlistToMap(items ...string) map[string]string {
	m := map[string]string{}

	if len(items)%2 != 0 {
		fmt.Printf(`WARN: Even number of strings required for Element attributes.
From: %s (%s)
Items %q dropping %q.
`,
			serr.FunctionName(serr.FrameLevels.FrameLevel4),
			serr.FunctionLoc(serr.FrameLevels.FrameLevel4),
			items, items[len(items)-1])
		items = items[:len(items)-1] // drop the last item
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
