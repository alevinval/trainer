package adapters

import (
	"encoding/xml"
	"time"

	"github.com/alevinval/trainer/internal/trainer"
)

type (
	// GpxAdapter maps to a GPX activity file
	GpxAdapter struct {
		Time        time.Time       `xml:"metadata>time"`
		Name        string          `xml:"trk>name"`
		TrackPoints []GpxTrackPoint `xml:"trk>trkseg>trkpt"`
	}

	// GpxTrackPoint maps to a GPX file trackpoint
	GpxTrackPoint struct {
		Time time.Time `xml:"time"`
		Lat  float64   `xml:"lat,attr"`
		Lon  float64   `xml:"lon,attr"`
		Ele  float64   `xml:"ele"`
		Hr   int       `xml:"extensions>TrackPointExtension>hr"`
		Cad  int       `xml:"extensions>TrackPointExtension>cad"`
	}
)

// NewGpxAdapter returns an adapter that converts gpx files
// into trainer primitives.
func NewGpxAdapter(b []byte) (g *GpxAdapter, err error) {
	g = &GpxAdapter{}
	err = xml.Unmarshal(b, g)
	return
}

// Metadata implements trainer.MetadataProvider interface.
// It creates a metadata object with known information of the gpx file:
// The activity name and time it was carried on.
func (g *GpxAdapter) Metadata() (meta *trainer.Metadata) {
	meta = &trainer.Metadata{
		Name: g.Name,
		Time: g.Time,
	}
	return
}

// DataPoints implements DatapointProvider interface.
// It converts a list of gpxTrackPoints to a list of datapoints.
func (g *GpxAdapter) DataPoints() (list trainer.DataPointList) {
	list = make(trainer.DataPointList, len(g.TrackPoints))
	for i, trackpoint := range g.TrackPoints {
		list[i] = trackpoint.toDataPoint()
	}
	list.Process()
	return list
}

func (tp *GpxTrackPoint) toDataPoint() (dp *trainer.DataPoint) {
	dp = &trainer.DataPoint{
		Time: tp.Time,
		Coords: trainer.Point{
			Lat:       tp.Lat,
			Lon:       tp.Lon,
			Elevation: trainer.Elevation(tp.Ele),
		},
		Hr: trainer.HeartRate(tp.Hr),

		// Count both feet for cadence.
		Cad: trainer.Cadence(tp.Cad * 2),

		N: 1,
	}
	return dp
}
