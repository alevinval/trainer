package adapter

import (
	"bytes"
	"math"

	"github.com/alevinval/trainer/internal/trainer"
	"github.com/tormoder/fit"
)

// fitAdapter maps a fit.File into trainer primitives.
type fitAdapter struct {
	b          []byte
	metadata   *trainer.Metadata
	datapoints trainer.DataPointList
}

// Fit returns an adapter that converts fit files
// into trainer primitives.
func Fit(b []byte) (provider trainer.ActivityProvider, err error) {
	buffer := bytes.NewBuffer(b)
	file, err := fit.Decode(buffer)
	if err != nil {
		return nil, err
	}
	activity, err := file.Activity()
	if err != nil {
		return nil, err
	}

	adapter := &fitAdapter{
		metadata:   makeMetadata(activity),
		datapoints: makeDatapoints(activity),
		b:          b,
	}
	return adapter, err
}

// DataPoints implements trainer.DatapointProvider interface.
// It converts a list of activity records to a list of datapoints.
func (adapter *fitAdapter) DataPoints() (list trainer.DataPointList) {
	return adapter.datapoints
}

// Metadata implements trainer.MetadataProvider interface.
// It creates a metadata object with known information of the fit file:
// The activity name and time it was carried on.
func (adapter *fitAdapter) Metadata() (metadata *trainer.Metadata) {
	return adapter.metadata
}

// Bytes implements trainer.BytesProvider interface.
// Returns the raw bytes of the original activity.
func (adapter *fitAdapter) Bytes() []byte {
	return adapter.b
}

func makeMetadata(activity *fit.ActivityFile) (metadata *trainer.Metadata) {
	return &trainer.Metadata{
		Time: activity.Activity.Timestamp,
		Name: "", // Not available on fit files, enrichers can fill that up.
	}
}

func makeDatapoints(activity *fit.ActivityFile) (list trainer.DataPointList) {
	list = make(trainer.DataPointList, 0)
	for _, record := range activity.Records {
		// Ignore datapoints with NaN coordinates for the moment being
		if math.IsNaN(record.PositionLat.Degrees()) {
			continue
		}
		list = append(list, &trainer.DataPoint{
			Time: record.Timestamp,
			Coords: trainer.Point{
				Lat:       record.PositionLat.Degrees(),
				Lon:       record.PositionLong.Degrees(),
				Elevation: trainer.Elevation(record.GetAltitudeScaled()),
			},
			Hr:  trainer.HeartRate(record.HeartRate),
			Cad: trainer.Cadence(record.Cadence * 2),
			N:   1,
		})
	}
	list.Process()
	return list
}
