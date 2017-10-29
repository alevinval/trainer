package trainer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const toleranceFloat = 0.0000001

var (
	testAvgDataPointList = DataPointList{
		{Cad: Cadence(100), Speed: Speed(10), Perf: Performance(5), n: 1},
		{Cad: Cadence(50), Speed: Speed(5), Perf: Performance(2.5), n: 2},
	}
)

func TestDataPointListAvg(t *testing.T) {
	assert.InDelta(t, 200.0/3.0, float64(testAvgDataPointList.AvgCad()), toleranceFloat)
	assert.InDelta(t, 20.0/3.0, float64(testAvgDataPointList.AvgSpeed()), toleranceFloat)
	assert.InDelta(t, 10.0/3.0, float64(testAvgDataPointList.AvgPerf()), toleranceFloat)
}
