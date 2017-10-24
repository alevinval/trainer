package trainer

import (
	"errors"
	"path"
)

const (
	extGpx     fileExt = ".gpx"
	extUnknown         = "unknown"
)

var (
	ErrUnknownFormat = errors.New("unrecognized file format")
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
		Data       []byte
		metadata   *Metadata
		datapoints DataPointList
	}

	ActivityList []*Activity

	fileExt string
)

func Open(fname string) (a *Activity, err error) {
	ext := getFileExt(fname)
	if ext == extUnknown {
		return nil, ErrUnknownFormat
	}

	data, err := getFileContent(fname)
	if err != nil {
		return nil, err
	}

	provider, err := getActivityProvider(ext, data)
	if err != nil {
		return nil, err
	}

	datapoints := provider.DataPoints()
	datapoints.process()

	metadata := provider.Metadata()
	metadata.DataSource = newDataSource(FileDataSource, fname)

	return &Activity{
		Data:       data,
		metadata:   metadata,
		datapoints: datapoints,
	}, err
}

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

func getFileExt(fname string) fileExt {
	switch fileExt(path.Ext(fname)) {
	case extGpx:
		return extGpx
	default:
		return extUnknown
	}
}

func getActivityProvider(ext fileExt, data []byte) (p activityProvider, err error) {
	switch ext {
	case extGpx:
		p, err = newGpx(data)
	default:
		p, err = nil, ErrUnknownFormat
	}
	return
}
