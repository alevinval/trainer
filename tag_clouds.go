package trainer

import "strings"

// TagCloud represents a collection of strings that appear in
// a group of activities.
type TagCloud struct {
	tags []string
}

var tagCloudSplitters = map[rune]struct{}{
	' ':  {},
	'.':  {},
	',':  {},
	'(':  {},
	')':  {},
	'/':  {},
	'\\': {},
	'#':  {},
	'-':  {},
}

func tagCloudFromActivities(activities ActivityList) *TagCloud {
	wordMap := make(map[string]struct{})

	for _, activity := range activities {
		metadata := activity.Metadata()
		words := strings.FieldsFunc(metadata.Name, func(r rune) bool {
			_, ok := tagCloudSplitters[r]
			return ok
		})
		for _, word := range words {
			wordMap[word] = struct{}{}
		}
	}

	words := make([]string, 0, len(wordMap))
	for word := range wordMap {
		words = append(words, word)
	}

	return &TagCloud{tags: words}
}
