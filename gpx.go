package trainer

import (
	"encoding/xml"
	"time"
)

type (
	gpx struct {
		Time        time.Time       `xml:"metadata>time"`
		Name        string          `xml:"trk>name"`
		TrackPoints []gpxTrackPoint `xml:"trk>trkseg>trkpt"`
	}

	gpxTrackPoint struct {
		Time time.Time `xml:"time"`
		Lat  float64   `xml:"lat,attr"`
		Lon  float64   `xml:"lon,attr"`
		Ele  float64   `xml:"ele"`
		Hr   int       `xml:"extensions>TrackPointExtension>hr"`
		Cad  int       `xml:"extensions>TrackPointExtension>cad"`
	}
)

func (tp *gpxTrackPoint) toDataPoint() (dp *DataPoint) {
	dp = &DataPoint{
		Time:   tp.Time,
		Coords: Point{tp.Lat, tp.Lon},
		Hr:     HeartRate(tp.Hr),

		// Count both feet for cadence.
		Cad: Cadence(tp.Cad * 2),

		n: 1,
	}
	return dp
}

func newGpx(b []byte) (g *gpx, err error) {
	g = new(gpx)
	err = xml.Unmarshal(b, g)
	return
}

// Metadata implements metadataProvider interface.
// It creates a metadata object with known information of the gpx file:
// The activity name and time it was carried on.
func (g *gpx) Metadata() (meta *Metadata) {
	meta = &Metadata{
		Name: g.Name,
		Time: g.Time,
	}
	return
}

// DataPoints implements datapointProvider interface.
// It converts a list of gpxTrackPoints to a list of datapoints.
func (g *gpx) DataPoints() (list DataPointList) {
	list = make(DataPointList, len(g.TrackPoints))
	for i, trackpoint := range g.TrackPoints {
		list[i] = trackpoint.toDataPoint()
	}
	return list
}
