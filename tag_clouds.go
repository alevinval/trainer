package trainer

import (
	"fmt"
	"sort"
	"strings"
)

// TagCloud represents a collection of strings that appear in
// a group of activities.
type TagCloud struct {
	tags map[string]int
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

var tagCloudIgnoredWords = map[string]struct{}{
	"and": {},
	"run": {},
	"de":  {},
}

// TagCloudFromActivities returns the cloud of rellevant words represented within
// list of activities.
func TagCloudFromActivities(activities ActivityList) *TagCloud {
	cloud := &TagCloud{
		tags: map[string]int{},
	}
	for _, activity := range activities {
		metadata := activity.Metadata()
		words := strings.FieldsFunc(metadata.Name, func(r rune) bool {
			_, ok := tagCloudSplitters[r]
			return ok
		})
		for _, word := range words {
			if len(word) == 1 {
				continue
			}
			if _, ok := tagCloudIgnoredWords[word]; ok {
				continue
			}
			cloud.tags[cloud.normalize(word)]++
		}
	}
	return cloud
}

// Contains tells whether a tag is present in the tag cloud or not.
func (tg *TagCloud) Contains(tag string) bool {
	_, ok := tg.tags[tg.normalize(tag)]
	return ok
}

func (tg *TagCloud) String() string {
	list := ""
	for _, word := range tg.asList() {
		list += fmt.Sprintf("%s[%d] ", word, tg.tags[word])
	}
	return list
}

func (tg *TagCloud) normalize(tag string) string {
	return strings.ToLower(tag)
}

func (tg *TagCloud) asList() []string {
	words := make([]string, 0, len(tg.tags))
	for word := range tg.tags {
		words = append(words, word)
	}
	sort.Slice(words, func(i, j int) bool {
		return tg.tags[words[i]] > tg.tags[words[j]]
	})
	return words
}
