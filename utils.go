package element

import "fmt"

// Generate a map from a list of key values
// Even number of items required or results may be fatal
// Candidate for a general utilities collection
func stringlistToMap(list ...string) map[string]string {
	m := map[string]string{}
	if len(list) % 2 != 0 {
		fmt.Println("Bad number of items to stringListToMap. Dropping:", list[len(list) - 1:])
	}
	key := ""
	for i, item := range list {
		if i % 2 == 0 {
			key = item
		} else {
			m[key] = item
		}
	}
	return m
}
