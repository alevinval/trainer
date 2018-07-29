package adapter

import (
	"encoding/xml"
	"time"

	"github.com/alevinval/trainer/internal/trainer"
)

type (
	// gpxAdapter maps to a GPX activity file
	gpxAdapter struct {
		Time        time.Time       `xml:"metadata>time"`
		Name        string          `xml:"trk>name"`
		TrackPoints []gpxTrackPoint `xml:"trk>trkseg>trkpt"`
	}

	// gpxTrackPoint maps to a GPX file trackpoint
	gpxTrackPoint struct {
		Time time.Time `xml:"time"`
		Lat  float64   `xml:"lat,attr"`
		Lon  float64   `xml:"lon,attr"`
		Ele  float64   `xml:"ele"`
		Hr   int       `xml:"extensions>TrackPointExtension>hr"`
		Cad  int       `xml:"extensions>TrackPointExtension>cad"`
	}
)

// Gpx returns an adapter that converts gpx files
// into trainer primitives.
func Gpx(b []byte) (provider trainer.ActivityProvider, err error) {
	g := &gpxAdapter{}
	err = xml.Unmarshal(b, g)
	return g, err
}

// Metadata implements trainer.MetadataProvider interface.
// It creates a metadata object with known information of the gpx file:
// The activity name and time it was carried on.
func (g *gpxAdapter) Metadata() (meta *trainer.Metadata) {
	meta = &trainer.Metadata{
		Name: g.Name,
		Time: g.Time,
	}
	return meta
}

// DataPoints implements DatapointProvider interface.
// It converts a list of gpxTrackPoints to a list of datapoints.
func (g *gpxAdapter) DataPoints() (list trainer.DataPointList) {
	list = make(trainer.DataPointList, len(g.TrackPoints))
	for i, trackpoint := range g.TrackPoints {
		list[i] = trackpoint.toDataPoint()
	}
	list.Process()
	return list
}

func (tp *gpxTrackPoint) toDataPoint() (dp *trainer.DataPoint) {
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
