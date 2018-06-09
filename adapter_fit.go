package trainer

import (
	"bytes"
	"math"

	"github.com/tormoder/fit"
)

type fitAdapter struct {
	file *fit.File
}

func newFitAdapter(b []byte) (adapter *fitAdapter, err error) {
	buffer := bytes.NewBuffer(b)
	adapter = &fitAdapter{}
	adapter.file, err = fit.Decode(buffer)
	return
}

func (adapter *fitAdapter) DataPoints() DataPointList {
	activity, err := adapter.file.Activity()
	if err != nil {
		panic(err)
	}

	list := make(DataPointList, 0)
	for _, record := range activity.Records {
		// Ignore datapoints with NaN coordinates for the moment being
		if math.IsNaN(record.PositionLat.Degrees()) {
			continue
		}
		list = append(list, &DataPoint{
			Time:   record.Timestamp,
			Coords: Point{record.PositionLat.Degrees(), record.PositionLong.Degrees()},
			Hr:     HeartRate(record.HeartRate),
			Cad:    Cadence(record.Cadence * 2),
			n:      1,
		})
	}
	return list
}

func (adapter *fitAdapter) Metadata() (meta *Metadata) {
	activity, _ := adapter.file.Activity()
	meta = &Metadata{
		Time: activity.Activity.Timestamp,
		Name: "n/a", // Name of the activity is not available in fit exports... will need to work on that
	}
	return
}
