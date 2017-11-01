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
		&DataPoint{Hr: bpm130, Speed: speed10, Cad: cad180, n: 1},
		&DataPoint{Hr: bpm150, Speed: speed10, Cad: cad180, n: 1},
		&DataPoint{Hr: bpm150, Speed: speed10, Cad: cad180, n: 1},
		&DataPoint{Hr: bpm180, Speed: speed10, Cad: cad180, n: 1},
		&DataPoint{Hr: bpm180, Speed: speed10, Cad: cad180, n: 1},
		&DataPoint{Hr: bpm180, Speed: speed10, Cad: cad180, n: 1},
	}
	input.process()
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
	histogram := input.GetHistogram()
	assert.Equal(t, 3, len(histogram.Data()[bpm180]))

	flat := histogram.Flatten()
	assert.Equal(t, len(flat.Data()), 3)
	assert.Equal(t, 1, len(flat.Data()[bpm180]))

	for hr, flatDataPoint := range flat.Data() {
		matched := input.filterBy(func(dp *DataPoint) bool {
			return dp.Hr == hr
		})
		assert.Equal(t, flatDataPoint.AvgCad(), matched.AvgCad())
		assert.Equal(t, flatDataPoint.AvgSpeed(), matched.AvgSpeed())
		assert.Equal(t, flatDataPoint.AvgPerf(), matched.AvgPerf())
	}

	assert.Equal(t, flat.GetAvgPerf(), input.AvgPerf())
}
