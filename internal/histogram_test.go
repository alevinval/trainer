package trainer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	bpm130  = HeartRate(130)
	bpm150  = HeartRate(150)
	bpm180  = HeartRate(180)
	speed10 = Speed(10)
	cad180  = Cadence(180)
)

func getTestInputDataPoints() DataPointList {
	input := DataPointList{
		&DataPoint{Hr: bpm130, Speed: speed10, Cad: cad180, N: 1},
		&DataPoint{Hr: bpm150, Speed: speed10, Cad: cad180, N: 1},
		&DataPoint{Hr: bpm150, Speed: speed10, Cad: cad180, N: 1},
		&DataPoint{Hr: bpm180, Speed: speed10, Cad: cad180, N: 1},
		&DataPoint{Hr: bpm180, Speed: speed10, Cad: cad180, N: 1},
		&DataPoint{Hr: bpm180, Speed: speed10, Cad: cad180, N: 1},
	}
	input.Process()
	return input
}

func TestHistogramFeed(t *testing.T) {
	input := getTestInputDataPoints()
	histogram := input.GetHistogram()
	for _, test := range []struct {
		hr    HeartRate
		count int
	}{
		{bpm130, 1},
		{bpm150, 2},
		{bpm180, 3},
	} {
		list := histogram.Data()[test.hr]
		assert.Equal(t, test.count, len(list))
	}
}

func TestHistogramFlatten(t *testing.T) {
	input := getTestInputDataPoints()
	hist := input.GetHistogram()
	assert.Equal(t, 3, len(hist.Data()[bpm180]))

	flat := hist.Flatten()
	for hr, flatDataPoint := range flat.Data() {
		matched := input.filterBy(func(dp *DataPoint) bool {
			return dp.Hr == hr
		})
		assert.Equal(t, flatDataPoint.Cad, matched.AvgCad())
		assert.Equal(t, flatDataPoint.Speed, matched.AvgSpeed())
		assert.Equal(t, flatDataPoint.Perf, matched.AvgPerf())
	}
}
