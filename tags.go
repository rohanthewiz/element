package element

var singleTags map[string]bool

func init() {
	singleTags = map[string]bool{
		"img":true, "br": true, "hr": true,
		"input": true, "link": true,
	}
}
