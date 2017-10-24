package trainer

import (
	"math"
	"time"
)

type (
	DataPoint struct {
		Time   time.Time
		Coords Point
		Hr     HeartRate
		Cad    Cadence
		Speed  Speed
		Perf   Performance

		n int
	}
	DataPointList []*DataPoint
)

// DistanceTo returns the distance in meters between two datapoints.
func (dp *DataPoint) DistanceTo(next *DataPoint) float64 {
	return dp.Coords.DistanceTo(next.Coords)
}

func (dp *DataPoint) secondsTo(next *DataPoint) float64 {
	return next.Time.Sub(dp.Time).Seconds()
}

func (dp *DataPoint) computeSpeed(prevDataPoint *DataPoint) {
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

// DataPoints implements datapointProvider interface
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

func (list DataPointList) AvgSpeed() Speed {
	var sum Speed
	for _, dp := range list {
		sum += dp.Speed
	}
	value := float64(sum) / float64(len(list))
	return Speed(value)
}

func (list DataPointList) AvgCad() Cadence {
	getter := func(dp *DataPoint) float64 {
		return float64(dp.Cad)
	}
	avg := list.weightedAverage(getter)
	return Cadence(avg)
}

func (list DataPointList) AvgPerf() Performance {
	getter := func(dp *DataPoint) float64 {
		return float64(dp.Perf)
	}
	avg := list.weightedAverage(getter)
	return Performance(avg)
}

func (list DataPointList) weightedAverage(getter func(dp *DataPoint) float64) float64 {
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

// DataPoints implements datapointProvider interface
func (list DataPointList) DataPoints() DataPointList {
	return list
}

func (list DataPointList) Histogram() *Histogram {
	hist := new(Histogram)
	hist.Reset()
	hist.Feed(list)
	return hist
}
