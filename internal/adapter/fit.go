package adapter

import (
	"bytes"
	"log"
	"math"

	"github.com/alevinval/trainer/internal/trainer"
	"github.com/tormoder/fit"
)

// FitAdapter maps a fit.File into trainer primitives.
type FitAdapter struct {
	file *fit.File
}

// Fit returns an adapter that converts fit files
// into trainer primitives.
func Fit(b []byte) (adapter *FitAdapter, err error) {
	buffer := bytes.NewBuffer(b)
	adapter = &FitAdapter{}
	adapter.file, err = fit.Decode(buffer)
	return
}

// DataPoints implements trainer.DatapointProvider interface.
// It converts a list of activity records to a list of datapoints.
func (adapter *FitAdapter) DataPoints() trainer.DataPointList {
	activity, err := adapter.file.Activity()
	if err != nil {
		log.Fatalf("unexpected error processing fit file: %s", err)
	}

	list := make(trainer.DataPointList, 0)
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

// Metadata implements trainer.MetadataProvider interface.
// It creates a metadata object with known information of the fit file:
// The activity name and time it was carried on.
func (adapter *FitAdapter) Metadata() (meta *trainer.Metadata) {
	activity, _ := adapter.file.Activity()
	meta = &trainer.Metadata{
		Time: activity.Activity.Timestamp,
		Name: "", // Not available on fit files, enrichers can fill that up.
	}
	return
}
