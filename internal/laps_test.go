package trainer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func generateDataPointsWithDistanceBetween(n int, distance float64) DataPointList {
	deg2m := 111300.0
	list := DataPointList{}
	for n > 0 {
		n--
		dp := &DataPoint{
			Coords: Point{
				Lat: 0,
				Lon: float64(n) * distance / deg2m,
			},
		}
		list = append(list, dp)
	}
	return list
}

func TestLapsCount(t *testing.T) {
	list := generateDataPointsWithDistanceBetween(10, 1000)

	laps := list.Laps(1000)

	assert.Equal(t, 10, len(laps))
}

func TestLapsCountWithMiles(t *testing.T) {
	list := generateDataPointsWithDistanceBetween(1000, 16)

	laps := list.Laps(1600)

	assert.Equal(t, 10, len(laps))
}
