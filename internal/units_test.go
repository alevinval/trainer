package trainer

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	speed3m30s      = Speed(4.76)
	speedStationary = Speed(0)
)

func TestUnitsToString(t *testing.T) {
	for _, test := range []struct {
		unit     fmt.Stringer
		expected string
	}{
		{unit: Cadence(179), expected: "179 steps/s"},
		{unit: Pace(speed3m30s), expected: "3:30 min/km"},
		{unit: Pace(0), expected: "n/a"},
		{unit: Performance(1.12345), expected: "1.12"},
		{unit: Point{1, 2, 100}, expected: "lat=1.000000, lon=2.000000, ele=100.0m"},
		{unit: HeartRate(129), expected: "129 bpm"},
		{unit: speed3m30s, expected: "4.76 m/s"},
		{unit: Distance(0.123), expected: "0.12 m"},
		{unit: Elevation(123.3), expected: "123.3m"},
	} {
		assert.Equal(t, test.expected, test.unit.String())
	}
}

func TestPointDistance(t *testing.T) {
	p1, p2 := Point{1, 1, 1}, Point{2, 2, 2}
	assert.Equal(t, p1.DistanceTo(p1), Distance(0))
	assert.True(t, p1.DistanceTo(p2) > 0)
	assert.True(t, p2.DistanceTo(p1) > 0)
}
