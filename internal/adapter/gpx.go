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

		metadata   *trainer.Metadata
		datapoints trainer.DataPointList
		b          []byte
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
	if err != nil {
		return nil, err
	}

	g.metadata = g.makeMetadata()
	g.datapoints = g.makeDataPoints()
	g.b = b
	return g, err
}

// Metadata implements trainer.MetadataProvider interface.
// It creates a metadata object with known information of the gpx file:
// The activity name and time it was carried on.
func (g *gpxAdapter) Metadata() (meta *trainer.Metadata) {
	return g.metadata
}

// DataPoints implements trainer.DatapointProvider interface.
// It converts a list of gpxTrackPoints to a list of datapoints.
func (g *gpxAdapter) DataPoints() (list trainer.DataPointList) {
	return g.datapoints
}

// Bytes implements trainer.BytesProvider interface.
// Returns the raw bytes of the original activity.
func (g *gpxAdapter) Bytes() []byte {
	return g.b
}

func (g *gpxAdapter) makeMetadata() (metadata *trainer.Metadata) {
	return &trainer.Metadata{
		Name: g.Name,
		Time: g.Time,
	}
}

func (g *gpxAdapter) makeDataPoints() (list trainer.DataPointList) {
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
