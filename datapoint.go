package trainer

import (
	"math"
	"time"
)

type (
	// DataPoint holds all relevant data to describe the location and performance during
	// an activity.
	DataPoint struct {
		Time   time.Time
		Coords Point
		Hr     HeartRate
		Cad    Cadence
		Speed  Speed
		Perf   Performance

		n int
	}

	// DataPointList is a list of DataPoint elements
	DataPointList []*DataPoint

	getterFunc func(dp *DataPoint) float64
)

var undefinedCoords = Point{}

// DistanceTo returns the distance in meters between two datapoints.
func (dp *DataPoint) DistanceTo(next *DataPoint) float64 {
	return dp.Coords.DistanceTo(next.Coords)
}

func (dp *DataPoint) secondsTo(next *DataPoint) float64 {
	return next.Time.Sub(dp.Time).Seconds()
}

func (dp *DataPoint) computeSpeed(prevDataPoint *DataPoint) {
	if dp.Coords == undefinedCoords {
		return
	}
	meters := prevDataPoint.DistanceTo(dp)
	seconds := prevDataPoint.secondsTo(dp)
	dp.Speed = Speed(meters / seconds)
}

func (dp *DataPoint) computePerf() {
	c := float64(dp.Cad) / 60.0
	s := float64(dp.Speed)
	h := float64(dp.Hr) / 60.0
	dp.Perf = Performance(c * s / h)
}

// DataPoints implements datapointProvider interface.
func (dp *DataPoint) DataPoints() DataPointList {
	return DataPointList{dp}
}

func (list DataPointList) process() {
	if len(list) <= 1 {
		return
	}
	prev := list[0]
	for _, curr := range list[1:] {
		curr.computeSpeed(prev)
		curr.computePerf()
		prev = curr
	}
}

// AvgSpeed returns the average speed of all datapoints in the list.
func (list DataPointList) AvgSpeed() Speed {
	avg := list.weightedAverage(getterSpeed)
	return Speed(avg)
}

// AvgCad returns the average cadence of all datapoints in the list.
func (list DataPointList) AvgCad() Cadence {
	avg := list.weightedAverage(getterCad)
	return Cadence(avg)
}

// AvgPerf returns the average performance of all datapoints in the list.
func (list DataPointList) AvgPerf() Performance {
	avg := list.weightedAverage(getterPerf)
	return Performance(avg)
}

func (list DataPointList) weightedAverage(getter getterFunc) float64 {
	var sum float64
	var size int
	for _, datapoint := range list {
		weightedValue := getter(datapoint) * float64(datapoint.n)
		if math.IsNaN(weightedValue) {
			continue
		}
		sum += weightedValue
		size += datapoint.n
	}
	return sum / float64(size)
}

func (list DataPointList) filterBy(fn func(dp *DataPoint) bool) DataPointList {
	filteredList := DataPointList{}
	for _, dp := range list {
		if fn(dp) {
			filteredList = append(filteredList, dp)
		}
	}
	return filteredList
}

// DataPoints implements datapointProvider interface.
func (list DataPointList) DataPoints() DataPointList {
	return list
}

// GetHistogram generates a histogram for the list of datapoints
func (list DataPointList) GetHistogram() *Histogram {
	hist := new(Histogram)
	hist.Reset()
	hist.Feed(list)
	return hist
}

func getterCad(dp *DataPoint) float64 {
	return float64(dp.Cad)
}

func getterSpeed(dp *DataPoint) float64 {
	return float64(dp.Speed)
}

func getterPerf(dp *DataPoint) float64 {
	return float64(dp.Perf)
}
