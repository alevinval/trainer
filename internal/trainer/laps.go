package trainer

import "fmt"

// Lap represents a segment composed of a sequence of data points.
// For example, a lap of datapoints up until one kilometer.
type Lap struct {
	ElevationUp   Elevation
	ElevationDown Elevation
	Distance      Distance
	Speed         Speed
	Pace          Pace
	Performance   Performance

	dataPoints DataPointList
}

type lapBuilder struct {
	curr    *DataPoint
	prev    *DataPoint
	i       int
	offset  int
	dist    Distance
	eleUp   Elevation
	eleDown Elevation
}

// DataPoints implements datapointProvider interface.
func (l *Lap) DataPoints() DataPointList {
	return l.dataPoints
}

func (l Lap) String() string {
	return fmt.Sprintf(
		"%s - %s - %s",
		l.Distance,
		l.Pace,
		l.ElevationUp,
	)
}

func (b *lapBuilder) fromList(list DataPointList, threshold Distance) (laps []*Lap) {
	laps = make([]*Lap, 0)
	if len(list) <= 1 {
		return
	}
	b.prev = list[0]
	for b.i, b.curr = range list[1:] {
		b.accumulate()
		b.prev = b.curr
		if b.dist >= threshold {
			laps = append(laps, b.buildLap(list))
			b.reset()
			b.offset = b.i
		}
	}
	if b.i < len(list) {
		b.i = len(list)
		laps = append(laps, b.buildLap(list))
	}
	return laps
}

func (b *lapBuilder) accumulate() {
	b.dist += b.curr.Coords.DistanceTo(b.prev.Coords)
	eleChange := b.curr.Coords.Elevation - b.prev.Coords.Elevation
	if eleChange > 0 {
		b.eleUp += eleChange
	} else if eleChange < 0 {
		b.eleDown -= eleChange
	}

}

func (b *lapBuilder) buildLap(list DataPointList) *Lap {
	dataPoints := list[b.offset:b.i]
	speed := dataPoints.AvgSpeed()
	perf := dataPoints.AvgPerf()
	return &Lap{
		dataPoints:  dataPoints,
		Distance:    b.dist,
		Speed:       speed,
		Pace:        Pace(speed),
		Performance: perf,
	}
}

func (b *lapBuilder) reset() {
	b.dist = Distance(0)
	b.eleUp = Elevation(0)
	b.eleDown = Elevation(0)
}
