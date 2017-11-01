package trainer

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
		Data       []byte
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

// DataPoints implements datapointProvider interface.
func (a *Activity) DataPoints() DataPointList {
	return a.datapoints
}

// DataPoints implements datapointProvider interface.
func (al ActivityList) DataPoints() DataPointList {
	list := DataPointList{}
	for _, activity := range al {
		list = append(list, activity.DataPoints()...)
	}
	return list
}

// GetClusters traverses the activity list to find clusters of activities
// by coordinates.
func (al ActivityList) GetClusters() ClusterList {
	return findClusters(al)
}

// GetHistogram generates a histogram for the list of activities
func (al ActivityList) GetHistogram() *Histogram {
	hist := new(Histogram)
	hist.Reset()
	hist.Feed(al)
	return hist
}
