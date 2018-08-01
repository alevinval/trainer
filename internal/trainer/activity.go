package trainer

import (
	"sort"
	"time"
)

type (
	DataPointProvider interface {
		DataPoints() DataPointList
	}

	MetadataProvider interface {
		Metadata() *Metadata
	}

	// ActivityProvider is the main domain object, it contains metadata about
	// your activity and the normalised datapoints with processed information
	// such as speed, performance, etc...
	ActivityProvider interface {
		MetadataProvider
		DataPointProvider
	}

	// ActivityList is a list of activities.
	ActivityList []ActivityProvider
)

// DataPoints implements datapointProvider interface.
func (list ActivityList) DataPoints() DataPointList {
	l := make(DataPointList, 0, len(list))
	for _, activity := range list {
		l = append(l, activity.DataPoints()...)
	}
	return l
}

// Filter returns a list of activities that pass the filter function.
func (list ActivityList) Filter(filterFn func(a ActivityProvider) bool) ActivityList {
	var filtered ActivityList
	for _, activity := range list {
		if filterFn(activity) {
			filtered = append(filtered, activity)
		}
	}
	return filtered
}

// SortByTime ensures the list of activities are sorted by date of the activity,
// from oldest to newest.
func (list ActivityList) SortByTime() {
	sort.Slice(list, func(i, j int) bool {
		return list[i].Metadata().Time.Before(list[j].Metadata().Time)
	})
}

// ChunkByDuration chunks the list of activities in time windows of a certain duration
func (list ActivityList) ChunkByDuration(d time.Duration) []ActivityList {
	chunks := []ActivityList{}
	if len(list) == 0 {
		return chunks
	}

	first := list[0]
	chunk, chunkTime := ActivityList{first}, first.Metadata().Time
	for _, a := range list[1:] {
		if a.Metadata().Time.Sub(chunkTime) >= d {
			chunks = append(chunks, chunk)
			chunk, chunkTime = ActivityList{a}, a.Metadata().Time
			continue
		}
		chunk = append(chunk, a)
	}
	chunks = append(chunks, chunk)
	return chunks
}
