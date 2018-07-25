package trainer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func provideActivityWithName(name string) *Activity {
	return &Activity{
		metadata: &Metadata{
			Name: name,
		},
	}
}

func TestTagCloudsContains(t *testing.T) {
	list := ActivityList{
		provideActivityWithName("long run"),
		provideActivityWithName("long runs"),
		provideActivityWithName("tempo run"),
		provideActivityWithName("race"),
	}
	cloud := TagCloudFromActivities(list)
	assert.True(t, cloud.Contains("long"))
	assert.True(t, cloud.Contains("tempo"))
	assert.True(t, cloud.Contains("race"))

	assert.False(t, cloud.Contains("workout"))
	assert.False(t, cloud.Contains("easy"))
}
