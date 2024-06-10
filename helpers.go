package element

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

// stringListToMap generates a map from a list of key values
// Even number of items required or results may be fatal
func stringListToMap(items ...string) map[string]string {
	m := map[string]string{}

	if len(items)%2 != 0 {
		// fmt.Println("Bad number of items to stringListToMap. Dropping:", items[len(items)-1:])
		items = items[:len(items)-1]
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

// addAttributePairsToMap also discards the last item in the list if the list count is odd
func addAttributePairsToMap(m map[string]string, items ...string) {
	if len(items)%2 != 0 {
		items = items[:len(items)-1] // drop last
	}

	key := ""
	for i, item := range items {
		if i%2 == 0 {
			key = item
		} else {
			m[key] = item
		}
	}
}
