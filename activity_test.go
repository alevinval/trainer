package trainer

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestActivitySortByTime(t *testing.T) {
	earliest := time.Unix(1, 0)
	oldest := time.Unix(10, 0)
	list := ActivityList{
		&Activity{
			metadata: &Metadata{
				Time: oldest,
			},
		},
		&Activity{
			metadata: &Metadata{
				Time: earliest,
			},
		},
	}

	assert.Equal(t, oldest, list[0].metadata.Time)
	list.SortByTime()
	assert.Equal(t, earliest, list[0].metadata.Time)
}
