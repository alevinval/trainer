package trainer

import (
	"sort"
	"time"
)

type (
	dataPointProvider interface {
		DataPoints() DataPointList
	}

	metadataProvider interface {
		Metadata() *Metadata
	}

	activityProvider interface {
		metadataProvider
		dataPointProvider
	}

	// Activity is the main domain object, it contains metadata about
	// your activity and the normalised datapoints with processed information
	// such as speed, performance, etc...
	Activity struct {
		rawData    []byte
		metadata   *Metadata
		datapoints DataPointList
	}

	// ActivityList is a list of activities.
	ActivityList []*Activity
)

// Metadata implements metadataProvider interface.
func (a *Activity) Metadata() *Metadata {
	return a.metadata
}

// SetMetadata sets the metadata object to the provided one.
func (a *Activity) SetMetadata(newMetadata *Metadata) {
	a.metadata = newMetadata
}

// DataPoints implements datapointProvider interface.
func (a *Activity) DataPoints() DataPointList {
	return a.datapoints
}

// DataPoints implements datapointProvider interface.
func (list ActivityList) DataPoints() DataPointList {
	l := make(DataPointList, 0, len(list))
	for _, activity := range list {
		l = append(l, activity.DataPoints()...)
	}
	return l
}

// Filter returns a list of activities that pass the filter function.
func (list ActivityList) Filter(filterFn func(a *Activity) bool) ActivityList {
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
		return list[i].metadata.Time.Before(list[j].metadata.Time)
	})
}

// ChunkByDuration chunks the list of activities in time windows of a certain duration
func (list ActivityList) ChunkByDuration(d time.Duration) []ActivityList {
	chunks := []ActivityList{}
	if len(list) == 0 {
		return chunks
	}

	chunk := ActivityList{list[0]}
	chunkTime := list[0].metadata.Time
	for _, activity := range list[1:] {
		if activity.metadata.Time.Sub(chunkTime) > d {
			chunks = append(chunks, chunk)
			chunk = ActivityList{activity}
			chunkTime = activity.metadata.Time
			continue
		}
		chunk = append(chunk, activity)
	}
	if len(chunk) > 0 {
		chunks = append(chunks, chunk)
	}
	return chunks
}
