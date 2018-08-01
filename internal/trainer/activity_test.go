package trainer

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type MockActivity struct {
	MockMetadata   *Metadata
	MockDataPoints DataPointList
}

func (ma *MockActivity) Metadata() *Metadata {
	if ma.MockMetadata != nil {
		return ma.MockMetadata
	}
	return &Metadata{}
}

func (ma *MockActivity) DataPoints() DataPointList {
	if ma.MockDataPoints != nil {
		return ma.MockDataPoints
	}
	return DataPointList{}
}

func createActivityListWithTimes(times ...time.Time) ActivityList {
	list := ActivityList{}
	for _, t := range times {
		activity := &MockActivity{
			MockMetadata: &Metadata{
				Time: t,
			},
		}
		list = append(list, activity)
	}
	return list
}

func TestActivityListSortByTime(t *testing.T) {
	earliest := time.Unix(10, 0)
	oldest := time.Unix(1, 0)
	list := createActivityListWithTimes(earliest, oldest)

	assert.Equal(t, earliest, list[0].Metadata().Time)
	list.SortByTime()
	assert.Equal(t, oldest, list[0].Metadata().Time)
}

func TestActivityListChunkByDurationSplit(t *testing.T) {
	t1, t2, t3 := time.Unix(1, 0), time.Unix(11, 0), time.Unix(21, 0)
	list := createActivityListWithTimes(t1, t2, t3)

	chunks := list.ChunkByDuration(10 * time.Second)

	assert.Equal(t, 3, len(chunks))
}

func TestActivityListChunkByDurationMerge(t *testing.T) {
	t1, t2, t3 := time.Unix(1, 0), time.Unix(10, 0), time.Unix(11, 0)
	list := createActivityListWithTimes(t1, t2, t3)

	chunks := list.ChunkByDuration(10 * time.Second)

	assert.Equal(t, 2, len(chunks))
	assert.Equal(t, 2, len(chunks[0]))
}
