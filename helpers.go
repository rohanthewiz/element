package element

import "fmt"

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

// Generate a map from a list of key values
// Even number of items required or results may be fatal
// Candidate for a general utilities collection
func stringlistToMap(list ...string) map[string]string {
	m := map[string]string{}

	if len(list)%2 != 0 {
		fmt.Println("Bad number of items to stringListToMap. Dropping:", list[len(list)-1:])
		return m
	}

	key := ""
	for i, item := range list {
		if i%2 == 0 {
			key = item
		} else {
			m[key] = item
		}
	}
	return m
}
